package game

import (
	"bytes"
	"fmt"
	"image"
	"log"
	"math"
	"math/rand"
	"time"

	"myebiten/internal/models"
	"myebiten/internal/models/character"
	"myebiten/internal/models/item"
	"myebiten/internal/weapons"
	wsClient "myebiten/internal/websocket/client"
	wsServer "myebiten/internal/websocket/server"
	images "myebiten/resources"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nfnt/resize"
)

type MainScene struct {
	models.SceneUI `json:"-"`

	stateEndingTimer *time.Timer
	itemSpawnTicker  *time.Ticker

	PlayersCount int
	state        int
	leftAlive    int

	Maze             [][]MazeNode
	Bullets          []*models.Bullet
	Items            []*item.Item
	Walls            []models.Wall
	Characters       []*character.Character
	defaultWeapons   []character.Weapon
	CharactersScores []uint

	ScoreUITexts []models.UIText
	pauseMenu    models.UIPanel `json:"-"`

	getConnectionMode func() string
	getGameClient     func() *wsClient.Client
	getGameServer     func() *wsServer.Server
}

func CreateMainScene(playersCount int) *MainScene {
	bullets := make([]*models.Bullet, weapons.DEFAULT_GUN_BULLETS_COUNT*playersCount+weapons.MINIGUN_BULLETS_COUNT*playersCount)
	for i := range bullets {
		bullets[i] = models.CreateBullet(weapons.DEFAULT_GUN_BULLET_RADIUS)
	}

	UIScores := make([]models.UIText, playersCount)
	for i := range UIScores {
		UIScores[i] = models.CreateUIText(fmt.Sprintf("Player %d: 0", i+1), REGULAR_FONT)
	}

	mainSceneUI := buildMainSceneUI(UIScores)
	return &MainScene{
		SceneUI:          mainSceneUI,
		state:            STATE_MAZE_CREATING,
		ScoreUITexts:     UIScores,
		Bullets:          bullets,
		CharactersScores: make([]uint, playersCount),
		PlayersCount:     playersCount,
	}
}

func buildMainSceneUI(scores []models.UIText) models.SceneUI {
	ebitenImage := ebiten.NewImage(SCREEN_SIZE_WIDTH, SCREEN_SIZE_HEIGHT)

	scene := models.CreateSceneUI(ebitenImage, float64(SCREEN_SIZE_HEIGHT), float64(SCREEN_SIZE_WIDTH))

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

	for i := range scores {
		areaID := scoreAreaID(i)
		scoreArea := ScoreArea.NewArea(
			0.99*ScoreArea.Height,
			ScoreArea.Width/float64(len(scores)),
			models.DrawingSettings{
				Offset: models.Vector2D{
					X: (float64(i) + 0.5) * ScoreArea.Width / float64(len(scores)),
					Y: 0.5 * ScoreArea.Height,
				},
				Scale: 1.0,
			})
		scene.AddDrawingArea(areaID, scoreArea)

		scoreUI := &scores[i]
		scoreUI.SetActive(true)
		scene.AddObject(scoreUI, areaID)
	}

	return scene
}

func (mainScene *MainScene) updateScores(id int) {
	if id < 0 || id >= len(mainScene.CharactersScores) || id >= len(mainScene.ScoreUITexts) {
		return
	}

	mainScene.CharactersScores[id]++
	mainScene.ScoreUITexts[id].SetText(fmt.Sprintf("Player %d: %d", id+1, mainScene.CharactersScores[id]))
}

func (mainScene *MainScene) syncScoreUITexts() {
	for i := range mainScene.ScoreUITexts {
		if i >= len(mainScene.CharactersScores) {
			break
		}
		mainScene.ScoreUITexts[i].SetText(fmt.Sprintf("Player %d: %d", i+1, mainScene.CharactersScores[i]))
	}
}

func scoreAreaID(index int) string {
	return fmt.Sprintf("%s_%d", SCORE_AREA_ID, index+1)
}

func (mainScene *MainScene) getClosestWalls(c *character.Character) []*models.Wall {
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

func (mainScene *MainScene) DetectCharacterToWallCollision(c *character.Character) {
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
		// cosine := math.Abs(b.Speed.X) / b.Speed.Length()
		// l := minDist / cosine
		// L := b.R / cosine
		// t := (L - l) / b.Speed.Length()

		// b.Position.X -= t * b.Speed.X
		// b.Position.Y -= t * b.Speed.Y

		b.Speed.X = -b.Speed.X

	} else if horizontalReflection {
		// cosine := math.Abs(b.Speed.Y) / b.Speed.Length()
		// l := minDist / cosine
		// L := b.R / cosine
		// t := (L - l) / b.Speed.Length()

		// b.Position.X -= t * b.Speed.X
		// b.Position.Y -= t * b.Speed.Y

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
		if !char.IsActive() {
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

	for _, item := range mainScene.Items {
		mainScene.AddObject(item, MAZE_AREA_ID)
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

	for _, item := range mainScene.Items {
		item.SetActive(false)
	}
	mainScene.Items = nil

	for _, char := range mainScene.Characters {
		char.SetActive(true)
		char.Input.Reset()
		char.SwitchToDefaultWeapon()
	}

	// This needs to be remade, quick solution
	mainScene.Objects = mainScene.Objects[:len(mainScene.ScoreUITexts)] //len(g.activeScene.Objects)-len(g.Walls)]
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
	if len(mainScene.Maze) < 3 || len(mainScene.Maze[0]) < 3 {
		return
	}

	i := rand.Intn(len(mainScene.Maze)-2) + 1
	j := rand.Intn(len(mainScene.Maze[0])-2) + 1

	position := getSceneCoordinates(i, j)
	itemType, sprite := getRandomItemSprite()
	item := item.CreateItem(itemType, position, sprite)
	mainScene.Items = append(mainScene.Items, item)
	mainScene.AddObject(item, MAZE_AREA_ID)
}

func (mainScene *MainScene) CreateCharacter(id int) {
	CHARACTER_IMAGE_TO_RESIZE, _, err := image.Decode(bytes.NewReader(images.TankV2png))
	if err != nil {
		log.Fatal(err)
	}
	resizedCharacterImage := resize.Resize(character.CHARACTER_WIDTH, 0, CHARACTER_IMAGE_TO_RESIZE, resize.Lanczos3)
	charImage := ebiten.NewImageFromImage(resizedCharacterImage)

	cs := controlSettingsForPlayer(id)

	clip := models.CreatePool(mainScene.Bullets[id*weapons.DEFAULT_GUN_BULLETS_COUNT : (id+1)*weapons.DEFAULT_GUN_BULLETS_COUNT])
	defaultWeapon := weapons.NewDefaultWeapon(clip)

	mainScene.defaultWeapons = append(mainScene.defaultWeapons, defaultWeapon)

	char := character.CreateCharacter(id, charImage, defaultWeapon, cs)
	char.SetActive(true)
	mainScene.Characters = append(mainScene.Characters, &char)
	mainScene.AddObject(&char, MAZE_AREA_ID)
}

func controlSettingsForPlayer(id int) models.ControlSettings {
	controlSettings := []models.ControlSettings{
		{
			RotateRightButton:  ebiten.KeyD,
			RotateLeftButton:   ebiten.KeyA,
			MoveForwardButton:  ebiten.KeyW,
			MoveBackwardButton: ebiten.KeyS,
			ShootButton:        ebiten.KeySpace,
		},
		{
			RotateRightButton:  ebiten.KeyArrowRight,
			RotateLeftButton:   ebiten.KeyArrowLeft,
			MoveForwardButton:  ebiten.KeyArrowUp,
			MoveBackwardButton: ebiten.KeyArrowDown,
			ShootButton:        ebiten.KeySlash,
		},
		{
			RotateRightButton:  ebiten.KeyL,
			RotateLeftButton:   ebiten.KeyJ,
			MoveForwardButton:  ebiten.KeyI,
			MoveBackwardButton: ebiten.KeyK,
			ShootButton:        ebiten.KeyO,
		},
		{
			RotateRightButton:  ebiten.KeyNumpad6,
			RotateLeftButton:   ebiten.KeyNumpad4,
			MoveForwardButton:  ebiten.KeyNumpad8,
			MoveBackwardButton: ebiten.KeyNumpad5,
			ShootButton:        ebiten.KeyNumpad0,
		},
	}

	if id < len(controlSettings) {
		return controlSettings[id]
	}

	return controlSettings[id%len(controlSettings)]
}

func clientControlSettings() models.ControlSettings {
	return models.ControlSettings{
		RotateRightButton:  ebiten.KeyD,
		RotateLeftButton:   ebiten.KeyA,
		MoveForwardButton:  ebiten.KeyW,
		MoveBackwardButton: ebiten.KeyS,
		ShootButton:        ebiten.KeySpace,
	}
}

// debug function
func (mainScene *MainScene) SanityCheck() {
	if len(mainScene.Objects) != len(mainScene.Bullets)+len(mainScene.Items)+len(mainScene.Characters)+len(mainScene.Walls)+len(mainScene.ScoreUITexts) {
		log.Println("discrepancy between the expected number of objects on the scene and actual number")
	}

	log.Println(len(mainScene.Areas))
	log.Println(len(mainScene.AreaIDs))
}
