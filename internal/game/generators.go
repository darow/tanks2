package game

type Generator struct {
	sources func(int, int) []Coordinates
	next    func(Coordinates, Coordinates, int, int) (Coordinates, bool)
}

func getHorizontalSnakeSources(N, M int) []Coordinates {
	lastJ := M
	if N%2 == 0 {
		lastJ = 1
	}

	return []Coordinates{{1, 1}, {N, lastJ}}
}

func HorizontalSnakeNext(currentCoord, root Coordinates, N, M int) (Coordinates, bool) {
	if currentCoord == root {
		return Coordinates{}, false
	}

	dir := -1
	if root.i > currentCoord.i || (root.i == currentCoord.i && ((root.j > currentCoord.j && root.i%2 != 0) || (root.j < currentCoord.j && root.i%2 == 0))) {
		dir = 1
	}

	i, j := currentCoord.i, currentCoord.j

	if i%2 == 1 {
		if 1 <= j+dir && j+dir <= M {
			return Coordinates{i, j + dir}, true
		}

		return Coordinates{i + dir, j}, true
	}

	if 1 <= j-dir && j-dir <= M {
		return Coordinates{i, j - dir}, true
	}

	return Coordinates{i + dir, j}, true
}

var Generators []Generator = []Generator{
	{
		sources: getHorizontalSnakeSources,
		next:    HorizontalSnakeNext,
	},
}
