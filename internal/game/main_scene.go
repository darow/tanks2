package game

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"log"
	"math"
	"math/rand"
	"myebiten/internal/models"
	"myebiten/internal/weapons"
	images "myebiten/resources"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nfnt/resize"
)

type MainScene struct {
	models.SceneUI

	stateEndingTimer *time.Timer
	itemSpawnTicker  *time.Ticker

	state     int
	leftAlive int

	Maze             [][]MazeNode
	Bullets          []*models.Bullet
	Walls            []models.Wall
	Characters       []*models.Character
	CharactersScores []uint

	scoreUITexts []models.UIText
	pauseMenu    models.UIPanel

	g *Game
}

func CreateMainScene() *MainScene {
	bullets := make([]*models.Bullet, weapons.BULLETS_COUNT*4)
	for i := range bullets {
		bullets[i] = models.CreateBullet(models.BULLET_RADIUS)
	}

	UIScore1 := models.CreateUIText("Player 1: 0", REGULAR_FONT)
	UIScore2 := models.CreateUIText("Player 2: 0", REGULAR_FONT)
	UIScores := []models.UIText{UIScore1, UIScore2}

	mainSceneUI := buildMainSceneUI(UIScores)
	return &MainScene{
		SceneUI: mainSceneUI,
		state:   STATE_MAZE_CREATING,
		Bullets: bullets,
	}
}

func buildMainSceneUI(scores []models.UIText) models.SceneUI {
	ebitenImage := ebiten.NewImage(SCREEN_SIZE_WIDTH, SCREEN_SIZE_HEIGHT)

	scene := models.CreateScene(ebitenImage, float64(SCREEN_SIZE_HEIGHT), float64(SCREEN_SIZE_WIDTH))

	rootArea := scene.GetRootArea()

	mainArea := rootArea.NewArea(
		rootArea.Height*0.8,
		rootArea.Width,
		models.DrawingSettings{
			Offset: models.Vector2D{X: 0.0, Y: rootArea.Height / 10},
			Scale:  1.0,
		})
	scene.AddDrawingArea(MAIN_PLAYING_AREA_ID, mainArea)

	UIArea1 := rootArea.NewArea(
		rootArea.Height*0.1,
		rootArea.Width,
		models.DrawingSettings{
			Offset: models.Vector2D{X: 0.0, Y: 0.0},
			Scale:  1.0,
		})
	scene.AddDrawingArea(UI_AREA1_ID, UIArea1)

	ScoreArea := rootArea.NewArea(
		rootArea.Height*0.1,
		rootArea.Width,
		models.DrawingSettings{
			Offset: models.Vector2D{X: 0.0, Y: rootArea.Height * 0.9},
			Scale:  1.0,
		})
	scene.AddDrawingArea(SCORE_AREA_ID, ScoreArea)

	scoreArea1 := ScoreArea.NewArea(
		0.99*ScoreArea.Height,
		0.2*ScoreArea.Width,
		models.DrawingSettings{
			Offset: models.Vector2D{X: 0.2 * ScoreArea.Width, Y: 0.5 * ScoreArea.Height},
			Scale:  1.0,
		})
	scene.AddDrawingArea(SCORE_AREA_1_ID, scoreArea1)

	scoreArea2 := ScoreArea.NewArea(
		0.99*ScoreArea.Height,
		0.2*ScoreArea.Width,
		models.DrawingSettings{
			Offset: models.Vector2D{X: 0.6 * ScoreArea.Width, Y: 0.5 * ScoreArea.Height},
			Scale:  1.0,
		})
	scene.AddDrawingArea(SCORE_AREA_2_ID, scoreArea2)

	scoreAreaIDs := []string{SCORE_AREA_1_ID, SCORE_AREA_2_ID}
	for i := range scores {
		scoreUI := &scores[i]
		scoreUI.SetActive(true)
		scene.AddObject(scoreUI, scoreAreaIDs[i])
	}

	return scene
}

func (mainScene *MainScene) Update() error {
	if mainScene.g.connMode == CONNECTION_MODE_CLIENT {
		char := mainScene.Characters[0]

		char.Input.Update()

		msg, err := json.Marshal(char.Input)
		if err != nil {
			log.Fatal(err)
		}

		err = mainScene.g.client.WriteMessage(msg)
		if err != nil {
			log.Fatal(err)
		}

		if char.Input.Shoot {
			char.Input.Shoot = false
		}

		mainScene.g.UpdateGameFromServer()

		return nil
	}

	switch mainScene.state {
	case STATE_MAZE_CREATING:
		mainScene.Reset()
		mainScene.itemSpawnTicker = time.NewTicker(ITEM_SPAWN_INTERVAL * time.Second)
		h, w, walls := mainScene.SetupLevel()
		if mainScene.g.connMode != CONNECTION_MODE_OFFLINE {
			mainScene.g.SendMazeToClient(h, w, walls)
		}
		mainScene.leftAlive = 2
		mainScene.state = STATE_GAME_RUNNING

		mainScene.SanityCheck()

	case STATE_GAME_RUNNING:
		select {
		case <-mainScene.itemSpawnTicker.C:
			mainScene.SpawnItem()
		default:
			if mainScene.leftAlive <= 1 {
				mainScene.stateEndingTimer = time.NewTimer(STATE_GAME_ENDING_TIMER_SECONDS * time.Second)
				mainScene.state = STATE_GAME_ENDING
			}
		}

	case STATE_GAME_ENDING:
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

	default:
		return errors.New("invalid state")
	}

	for i, char := range mainScene.Characters {
		if !char.IsActive() {
			continue
		}

		if i == 1 && mainScene.g.connMode == CONNECTION_MODE_SERVER {
			// process client's character's input
			msg := mainScene.g.server.ReadMessage()

			var input models.Input
			err := json.Unmarshal(msg, &input)
			if err != nil {
				continue
			}

			char.Input = input
		} else {
			char.Input.Update()
		}

		char.ProcessInput()

		char.Move()

		mainScene.DetectCharacterToWallCollision(char)
	}

	for _, bullet := range mainScene.Bullets {
		if !bullet.IsActive() {
			continue
		}

		bullet.Move()

		mainScene.DetectBulletToWallCollision(bullet)

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

	if mainScene.g.connMode == CONNECTION_MODE_SERVER {
		msg, err := json.Marshal(mainScene.g)
		if err != nil {
			log.Fatal(err)
		}
		err = mainScene.g.server.WriteThingsMessage(msg)
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

func (mainScene *MainScene) updateScores(id int) {
	mainScene.CharactersScores[id]++
	mainScene.scoreUITexts[id].SetText(fmt.Sprintf("Player %d: %d", id+1, mainScene.CharactersScores[id]))
}

func (mainScene *MainScene) getClosestWalls(c *models.Character) []*models.Wall {
	// yes, this is shit, I see it too, dw it will all change
	i, j := getMazeCoordinates(c.Position)
	k := 0
	if mainScene.Maze[i][j].topWall != nil {
		wallsToCheck[k] = mainScene.Maze[i][j].topWall
		k++
	}
	if mainScene.Maze[i][j].bottomWall != nil {
		wallsToCheck[k] = mainScene.Maze[i][j].bottomWall
		k++
	}
	if mainScene.Maze[i][j].leftWall != nil {
		wallsToCheck[k] = mainScene.Maze[i][j].leftWall
		k++
	}
	if mainScene.Maze[i][j].rightWall != nil {
		wallsToCheck[k] = mainScene.Maze[i][j].rightWall
		k++
	}

	if mainScene.Maze[i-1][j].leftWall != nil {
		wallsToCheck[k] = mainScene.Maze[i-1][j].leftWall
		k++
	}
	if mainScene.Maze[i-1][j].rightWall != nil {
		wallsToCheck[k] = mainScene.Maze[i-1][j].rightWall
		k++
	}

	if mainScene.Maze[i+1][j].leftWall != nil {
		wallsToCheck[k] = mainScene.Maze[i+1][j].leftWall
		k++
	}
	if mainScene.Maze[i+1][j].rightWall != nil {
		wallsToCheck[k] = mainScene.Maze[i+1][j].rightWall
		k++
	}

	if mainScene.Maze[i][j-1].topWall != nil {
		wallsToCheck[k] = mainScene.Maze[i][j-1].topWall
		k++
	}
	if mainScene.Maze[i][j-1].bottomWall != nil {
		wallsToCheck[k] = mainScene.Maze[i][j-1].bottomWall
		k++
	}

	if mainScene.Maze[i][j+1].topWall != nil {
		wallsToCheck[k] = mainScene.Maze[i][j+1].topWall
		k++
	}
	if mainScene.Maze[i][j+1].bottomWall != nil {
		wallsToCheck[k] = mainScene.Maze[i][j+1].bottomWall
		k++
	}

	return wallsToCheck[:k]
}

func (mainScene *MainScene) DetectCharacterToWallCollision(c *models.Character) {
	closestWalls := mainScene.getClosestWalls(c)
	for _, w := range closestWalls {
		isCollide := c.DetectWallCollision(*w)
		if isCollide {
			c.MoveBack()
		}
	}
}

func (mainScene *MainScene) DetectBulletToWallCollision(b *models.Bullet) {
	wh := float64(WALL_HEIGHT)
	ww := float64(WALL_WIDTH)

	i, j := getMazeCoordinates(b.Position)
	nodeCenter := getSceneCoordinates(i, j)

	if i >= len(mainScene.Maze)-1 || i < 0 || j >= len(mainScene.Maze[0])-1 || j < 0 {
		return
	}

	// Here top and bottom mean how these directions appear on the screen
	// meaning, that distToTop actually measures the distance to the wall
	// that is stored as MazeNode.bottomWall
	distToTop := b.Position.Y - (nodeCenter.Y - wh/2 + ww)
	distToRight := (nodeCenter.X + wh/2 - ww) - b.Position.X
	distToLeft := b.Position.X - (nodeCenter.X - wh/2 + ww)
	distToBottom := (nodeCenter.Y + wh/2 - ww) - b.Position.Y

	minDist := min(distToBottom, distToLeft, distToRight, distToTop)

	if minDist > b.R {
		return
	}

	var horizontalReflection, verticalReflection bool = false, false

	if minDist == distToBottom {
		// again, because of mismatch between how directions are logically stored
		// and how they are presented on screen bottom reflection requires MazeNode.topWall

		// check if the top wall is present
		horizontalReflection = !mainScene.Maze[i][j].up

		// if close to the left check 3 corner walls, same if close to the right
		// if at least one corner wall is present, perform reflection
		verticalReflection = (distToLeft < b.R) && !(mainScene.Maze[i][j].up && mainScene.Maze[i+1][j].left && mainScene.Maze[i][j-1].up) ||
			(distToRight < b.R) && !(mainScene.Maze[i][j].up && mainScene.Maze[i+1][j].right && mainScene.Maze[i][j+1].up)

		// prioritize reflection of the main wall
		verticalReflection = verticalReflection && !horizontalReflection
	}

	if minDist == distToTop {
		// check comments in minDist == distToBottom block
		horizontalReflection = !mainScene.Maze[i][j].down

		verticalReflection = (distToLeft < b.R) && !(mainScene.Maze[i][j].down && mainScene.Maze[i-1][j].left && mainScene.Maze[i][j-1].down) ||
			(distToRight < b.R) && !(mainScene.Maze[i][j].down && mainScene.Maze[i-1][j].right && mainScene.Maze[i][j+1].down)

		verticalReflection = verticalReflection && !horizontalReflection
	}

	if minDist == distToLeft {
		// check comments in minDist == distToBottom block
		verticalReflection = !mainScene.Maze[i][j].left

		horizontalReflection = (distToTop < b.R) && !(mainScene.Maze[i][j].down && mainScene.Maze[i-1][j].left && mainScene.Maze[i][j-1].down) ||
			(distToBottom < b.R) && !(mainScene.Maze[i][j].up && mainScene.Maze[i+1][j].left && mainScene.Maze[i][j-1].up)

		horizontalReflection = horizontalReflection && !verticalReflection
	}

	if minDist == distToRight {
		// check comments in minDist == distToBottom block
		verticalReflection = !mainScene.Maze[i][j].right

		horizontalReflection = (distToTop < b.R) && !(mainScene.Maze[i][j].down && mainScene.Maze[i-1][j].right && mainScene.Maze[i][j+1].down) ||
			(distToBottom < b.R) && !(mainScene.Maze[i][j].up && mainScene.Maze[i+1][j].right && mainScene.Maze[i][j+1].up)

		horizontalReflection = horizontalReflection && !verticalReflection
	}

	if verticalReflection {
		cosine := math.Abs(b.Speed.X) / b.Speed.Length()
		l := minDist / cosine
		L := b.R / cosine
		t := (L - l) / b.Speed.Length()

		b.Position.X -= t * b.Speed.X
		b.Position.Y -= t * b.Speed.Y

		b.Speed.X = -b.Speed.X

	} else if horizontalReflection {
		cosine := math.Abs(b.Speed.Y) / b.Speed.Length()
		l := minDist / cosine
		L := b.R / cosine
		t := (L - l) / b.Speed.Length()

		b.Position.X -= t * b.Speed.X
		b.Position.Y -= t * b.Speed.Y

		b.Speed.Y = -b.Speed.Y
	}
}

func (mainScene *MainScene) SetupLevel() (int, int, []models.Wall) {
	h := rand.Intn(MAX_BOARD_HEIGHT-MIN_BOARD_HEIGHT) + MIN_BOARD_HEIGHT
	w := rand.Intn(MAX_BOARD_WIDTH-MIN_BOARD_WIDTH) + MIN_BOARD_WIDTH

	walls := mainScene.CreateMaze(h, w)
	mainScene.SetDrawingSettings(h, w)
	mainScene.SetCharacters(h, w)

	return h, w, walls
}

func (mainScene *MainScene) SetCharacters(h, w int) {
	spawnPlaces := []models.Vector2D{}
	for range mainScene.Characters {
		i := rand.Intn(h) + 1
		j := rand.Intn(w) + 1
		spawnPlace := getSceneCoordinates(i, j)
		spawnPlaces = append(spawnPlaces, spawnPlace)
	}

	i := 0
	for _, char := range mainScene.Characters {
		if !char.Active {
			continue
		}

		char.Position.X = spawnPlaces[i].X
		char.Position.Y = spawnPlaces[i].Y

		char.Rotation = math.Pi / 2

		char.Speed.X = 0
		char.Speed.Y = 0

		i++
	}
}

func (mainScene *MainScene) CreateMaze(h, w int) []models.Wall {
	mainScene.Walls = make([]models.Wall, 0)

	mainScene.Maze = createMaze(h, w)
	mainScene.Walls = buildMaze(mainScene.Maze, mainScene.Walls)

	return mainScene.Walls
}

func (mainScene *MainScene) SetDrawingSettings(h, w int) {
	mainArea := mainScene.GetArea(MAIN_PLAYING_AREA_ID)
	if mainArea == nil {
		log.Fatal("Main playing area is not set")
	}

	areaHeight := mainArea.Height
	areaWidth := mainArea.Width

	mazeHeight := float64(h*(WALL_HEIGHT-WALL_WIDTH) + WALL_WIDTH)
	mazeWidth := float64(w*(WALL_HEIGHT-WALL_WIDTH) + WALL_WIDTH)

	scalingFactor := min(areaHeight/mazeHeight, areaWidth/mazeWidth)

	mazeHeight *= scalingFactor
	mazeWidth *= scalingFactor

	newDrawingSettings := models.DrawingSettings{
		Offset: models.Vector2D{X: (areaWidth - mazeWidth) / 2, Y: (areaHeight - mazeHeight) / 2},
		Scale:  scalingFactor,
	}

	mazeArea := mainArea.NewArea(mazeHeight, mazeWidth, newDrawingSettings)
	mainScene.AddDrawingArea(MAZE_AREA_ID, mazeArea)

	for _, bullet := range mainScene.Bullets {
		mainScene.AddObject(bullet, MAZE_AREA_ID)
	}

	for _, char := range mainScene.Characters {
		mainScene.AddObject(char, MAZE_AREA_ID)
	}

	for _, wall := range mainScene.Walls {
		mainScene.AddObject(&wall, MAZE_AREA_ID)
	}
}

func (mainScene *MainScene) Reset() {
	for _, bullet := range mainScene.Bullets {
		bullet.SetActive(false)
	}

	for _, char := range mainScene.Characters {
		char.SetActive(true)
		char.Input.Reset()
	}

	// This needs to be remade, quick solution
	mainScene.Objects = mainScene.Objects[:2] //len(g.activeScene.Objects)-len(g.Walls)]
	am := mainScene.AreaIDs
	for obj, id := range am {
		if id == MAZE_AREA_ID {
			delete(am, obj)
		}
	}
	mainArea := mainScene.GetArea(MAIN_PLAYING_AREA_ID)
	mainArea.Children = nil
}

func (mainScene *MainScene) SpawnItem() {

}

func (mainScene *MainScene) CreateCharacter(id int) {
	CHARACTER_IMAGE_TO_RESIZE, _, err := image.Decode(bytes.NewReader(images.TankV2png))
	if err != nil {
		log.Fatal(err)
	}
	resizedCharacterImage := resize.Resize(models.CHARACTER_WIDTH, 0, CHARACTER_IMAGE_TO_RESIZE, resize.Lanczos3)
	charImage := ebiten.NewImageFromImage(resizedCharacterImage)

	cs1 := models.ControlSettings{
		RotateRightButton:  ebiten.KeyF,
		RotateLeftButton:   ebiten.KeyS,
		MoveForwardButton:  ebiten.KeyE,
		MoveBackwardButton: ebiten.KeyD,
		ShootButton:        ebiten.KeySpace,
	}

	defaultWeapon := weapons.DefaultWeapon{
		Clip:     mainScene.Bullets[:weapons.BULLETS_COUNT],
		Cooldown: 5,
	}

	char := models.CreateCharacter(id, charImage, &defaultWeapon, cs1)
	mainScene.Characters = append(mainScene.Characters, &char)
}

// debug function
func (mainScene *MainScene) SanityCheck() {
	if len(mainScene.Objects) != len(mainScene.Bullets)+len(mainScene.Characters)+len(mainScene.Walls)+2 {
		log.Println("discrepancy between the expected number of objects on the scene and actual number")
	}

	log.Println(len(mainScene.Areas))
	log.Println(len(mainScene.AreaIDs))
}
