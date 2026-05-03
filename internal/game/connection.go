package game

import (
	"encoding/json"
	"log"
	"runtime"

	"myebiten/internal/models"
	"myebiten/internal/models/character"
	modelitem "myebiten/internal/models/item"
)

var (
	CONNECTION_MODE_OFFLINE = "offline"
	CONNECTION_MODE_SERVER  = "server"
	CONNECTION_MODE_CLIENT  = "client"
)

type MazeDTO struct {
	H, W  int
	Walls []models.Wall
}

type connectionClient interface {
	ReadMessage() []byte
	WriteMessage(message []byte) error
}

type connectionServer interface {
	GetInput() models.Input
	WriteThingsMessage(message []byte) error
	WriteMapMessage(message []byte) error
}

func (mainScene *MainScene) UpdateGameFromServer(client connectionClient) {
	msg := client.ReadMessage()

	var newGame MainScene
	err := json.Unmarshal(msg, &newGame)
	if err != nil {
		log.Println(err)
		return
	}

	if len(newGame.Characters) == 0 {
		return
	}

	mainScene.Characters = copyCharacters(mainScene.Characters, newGame.Characters)
	mainScene.Bullets = copyBullets(mainScene.Bullets, newGame.Bullets)
	mainScene.Items = mainScene.copyItems(newGame.Items)
	mainScene.CharactersScores = newGame.CharactersScores
}

// func (c *Character) SendInputToServer(ws *websocket.Conn) {
// 	msg, err := json.Marshal(c.input)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	err = ws.WriteMessage(websocket.TextMessage, msg)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

func copyCharacters(dst []*character.Character, src []*character.Character) []*character.Character {
	if len(dst) == 0 {
		dst = append(dst, src...)
	} else {
		for i, char := range src {
			c := dst[i]
			if char == nil {
				c.Position.X = 99999
				c.Position.Y = 99999
				continue
			}

			c.Position.X = char.Position.X
			c.Position.Y = char.Position.Y
			c.Rotation = char.Rotation

			c.SetActive(char.IsActive())
		}
	}

	return dst
}

func copyBullets(dst []*models.Bullet, src []*models.Bullet) []*models.Bullet {
	if len(dst) == 0 {
		dst = append(dst, src...)
	} else {
		for i, b := range src {
			bullet := dst[i]

			bullet.Position.X = b.Position.X
			bullet.Position.Y = b.Position.Y
			bullet.Rotation = b.Rotation
			bullet.Sprite.R = b.Sprite.R

			bullet.SetActive(b.IsActive())
		}
	}

	return dst
}

func (mainScene *MainScene) copyItems(src []*modelitem.Item) []*modelitem.Item {
	for i := len(src); i < len(mainScene.Items); i++ {
		if mainScene.Items[i] != nil {
			mainScene.Items[i].SetActive(false)
		}
	}

	dst := mainScene.Items
	if len(dst) < len(src) {
		dst = append(dst, make([]*modelitem.Item, len(src)-len(dst))...)
	}
	dst = dst[:len(src)]

	for i, srcItem := range src {
		if srcItem == nil {
			if dst[i] != nil {
				dst[i].SetActive(false)
			}
			continue
		}

		if dst[i] == nil {
			dst[i] = &modelitem.Item{}
			mainScene.AddObject(dst[i], MAZE_AREA_ID)
		}

		dstItem := dst[i]
		dstItem.GameObject = srcItem.GameObject
		dstItem.Type = srcItem.Type

		if dstItem.IsActive() && dstItem.IconSprite.Image == nil {
			if sprite := getItemSprite(dstItem.Type); sprite != nil {
				dstItem.IconSprite = *sprite
			}
		}
	}

	return dst
}

func (mainScene *MainScene) ReceiveMazeUpdates() {
	client := mainScene.getGameClient()
	for {
		_, message, err := client.MapUpdateConn.ReadMessage()
		if err != nil {
			log.Println(runtime.Caller(1))
			log.Println(err)
		}

		// if *DEBUG_MODE {
		// 	log.Printf("Received map message: %s\n", message)
		// }

		var maze MazeDTO
		err = json.Unmarshal(message, &maze)
		if err != nil {
			log.Fatal(err)
		}

		mainScene.Walls = maze.Walls

		for i := range mainScene.Walls {
			mainScene.Walls[i].Hitbox = models.RectangleHitbox{W: WALL_WIDTH, H: WALL_HEIGHT}
			mainScene.Walls[i].Sprite = models.RectangleSprite{W: WALL_WIDTH, H: WALL_HEIGHT}
		}

		mainScene.Reset()
		mainScene.SetDrawingSettings(maze.H, maze.W)
	}
}

func SendMazeToClient(server connectionServer, h, w int, walls []models.Wall) {
	maze := MazeDTO{h, w, walls}

	msg, err := json.Marshal(maze)
	if err != nil {
		log.Fatal(err)
	}

	err = server.WriteMapMessage(msg)

	if err != nil {
		log.Fatal(err)
	}
}
