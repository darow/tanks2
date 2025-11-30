package game

const (
	ROOT_AREA_ID         = "root_area"
	MAIN_PLAYING_AREA_ID = "main_playing_area"
	MAZE_AREA_ID         = "maze_area"
	UI_AREA1_ID          = "ui_area_1"
	SCORE_AREA_ID        = "score_area"
	SCORE_AREA_1_ID      = "score_area_1"
	SCORE_AREA_2_ID      = "score_area_2"
	SCORE_AREA_3_ID      = "score_area_3"
	SCORE_AREA_4_ID      = "score_area_4"
)

const (
	MENU_SCENE_ID  = 1
	LOBBY_SCENE_ID = 2
	MAIN_SCENE_ID  = 3
)

var TILE_ID_SEQUENCE = 0

func (g *Game) CreateCharacter(id int) {
	g.scenes[MAIN_SCENE_ID].(*MainScene).CreateCharacter(id)
}

func (g *Game) SetActiveScene(sceneID int) {
	g.activeScene = g.scenes[sceneID]
}
