package main

import (
	"strconv"
	"testing"

	"myebiten/internal/models"
)

func TestGetMazeCoordinates(t *testing.T) {
	type TestCase struct {
		input     models.Vector2D
		iExpected int
		jExpected int
	}

	testCases := []TestCase{
		{
			input:     models.Vector2D{X: 0.0, Y: 0.0},
			iExpected: 0,
			jExpected: 0,
		},
		{
			input:     models.Vector2D{X: 10.0, Y: 20.0},
			iExpected: 1,
			jExpected: 1,
		},
		{
			input:     models.Vector2D{X: 200.0, Y: 20.0},
			iExpected: 1,
			jExpected: 2,
		},
		{
			input:     models.Vector2D{X: 5.0, Y: 5.0},
			iExpected: 1,
			jExpected: 1,
		},
		{
			input:     models.Vector2D{X: 205.0, Y: 5.0},
			iExpected: 1,
			jExpected: 2,
		},
		{
			input:     models.Vector2D{X: 205.0, Y: 50.0},
			iExpected: 1,
			jExpected: 2,
		},
		{
			input:     models.Vector2D{X: 195.0, Y: 195.0},
			iExpected: 2,
			jExpected: 2,
		},
		{
			input:     models.Vector2D{X: 195.0, Y: 55.0},
			iExpected: 1,
			jExpected: 2,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			iActual, jActual := getMazeCoordinates(tc.input)
			if (iActual != tc.iExpected) || (jActual != tc.jExpected) {
				t.Errorf("Test %d failed on input %v with expected result %d %d and actual result %d %d", i, tc.input, tc.iExpected, tc.jExpected, iActual, jActual)
			}
		})
	}
}
