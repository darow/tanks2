package game

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"myebiten/internal/models"
	"myebiten/internal/models/character"
)

func (mainScene *MainScene) Update() error {
	connectionMode := mainScene.getConnectionMode()
	client := mainScene.getGameClient()
	server := mainScene.getGameServer()

	if connectionMode == CONNECTION_MODE_CLIENT {
		return mainScene.updateClientFrame(client)
	}

	if err := mainScene.updateState(connectionMode, server); err != nil {
		return err
	}

	mainScene.updateCharacters(connectionMode, server)
	mainScene.updateBullets()

	if connectionMode == CONNECTION_MODE_SERVER {
		mainScene.syncToClient(server)
	}

	return nil
}

func (mainScene *MainScene) updateClientFrame(client connectionClient) error {
	playerID := client.GetPlayerID()
	if playerID < 0 || playerID >= len(mainScene.Characters) {
		return errors.New("client player id is outside characters list")
	}

	char := mainScene.Characters[playerID]

	char.Input.Update()

	msg, err := json.Marshal(char.Input)
	if err != nil {
		return err
	}

	if err := client.WriteMessage(msg); err != nil {
		return err
	}

	mainScene.UpdateGameFromServer(client)
	return nil
}

func (mainScene *MainScene) updateState(connectionMode string, server connectionServer) error {
	switch mainScene.state {
	case STATE_MAZE_CREATING:
		mainScene.startNewRound(connectionMode, server)
	case STATE_GAME_RUNNING:
		mainScene.updateRunningState()
	case STATE_GAME_ENDING:
		mainScene.updateEndingState()
	default:
		return errors.New("invalid state")
	}

	return nil
}

func (mainScene *MainScene) startNewRound(connectionMode string, server connectionServer) {
	mainScene.Reset()
	mainScene.itemSpawnTicker = time.NewTicker(ITEM_SPAWN_INTERVAL * time.Second)

	h, w, walls := mainScene.SetupLevel()
	if connectionMode != CONNECTION_MODE_OFFLINE {
		SendMazeToClient(server, h, w, walls)
	}

	mainScene.leftAlive = len(mainScene.Characters)
	mainScene.state = STATE_GAME_RUNNING
	mainScene.SanityCheck()
}

func (mainScene *MainScene) updateRunningState() {
	select {
	case <-mainScene.itemSpawnTicker.C:
		mainScene.SpawnItem()
	default:
		if mainScene.leftAlive <= 1 {
			mainScene.stateEndingTimer = time.NewTimer(STATE_GAME_ENDING_TIMER_SECONDS * time.Second)
			mainScene.state = STATE_GAME_ENDING
		}
	}
}

func (mainScene *MainScene) updateEndingState() {
	select {
	case <-mainScene.stateEndingTimer.C:
		for _, char := range mainScene.Characters {
			if char.IsActive() {
				mainScene.updateScores(char.ID)
				break
			}
		}
		mainScene.state = STATE_MAZE_CREATING
	default:
	}
}

func (mainScene *MainScene) updateCharacters(connectionMode string, server connectionServer) {
	for i, char := range mainScene.Characters {
		if !char.IsActive() {
			continue
		}

		if connectionMode == CONNECTION_MODE_SERVER {
			char.Input = server.GetInput(i)
		} else if i == normalizePlayerID() {
			char.Input.Update()
		}

		char.ProcessInput()
		char.Move()
		mainScene.DetectCharacterToWallCollision(char)
		mainScene.collectItems(char, i)
	}
}

func (mainScene *MainScene) updateBullets() {
	for _, bullet := range mainScene.Bullets {
		if !bullet.IsActive() {
			continue
		}

		bullet.Move()
		mainScene.DetectBulletToWallCollision(bullet)
		mainScene.detectBulletHitsCharacters(bullet)
	}
}

func (mainScene *MainScene) detectBulletHitsCharacters(bullet *models.Bullet) {
	for _, char := range mainScene.Characters {
		if !char.IsActive() {
			continue
		}

		if char.DetectBulletToCharacterCollision(bullet) {
			bullet.SetActive(false)
			char.SetActive(false)
			mainScene.leftAlive--
		}
	}
}

func (mainScene *MainScene) collectItems(char *character.Character, charIndex int) {
	for _, item := range mainScene.Items {
		if !item.DetectCharacterCollision(char) {
			continue
		}

		item.SetActive(false)
		mainScene.applyItemEffect(item, char, charIndex)
	}
}

func (mainScene *MainScene) syncToClient(server connectionServer) {
	msg, err := json.Marshal(mainScene)
	if err != nil {
		log.Fatal(err)
	}

	if err := server.WriteThingsMessage(msg); err != nil {
		log.Fatal(err)
	}
}
