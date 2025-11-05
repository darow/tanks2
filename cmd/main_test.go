package main

import (
	"testing"
)

func TestGetMazeCoordinates(t *testing.T) {
	type TestCase struct {
		input     Vector2D
		iExpected int
		jExpected int
	}

	testCases := []TestCase{
		{
			input:     Vector2D{0.0, 0.0},
			iExpected: 0,
			jExpected: 0,
		},
		{
			input:     Vector2D{10.0, 20.0},
			iExpected: 1,
			jExpected: 1,
		},
		{
			input:     Vector2D{200.0, 20.0},
			iExpected: 1,
			jExpected: 2,
		},
		{
			input:     Vector2D{5.0, 5.0},
			iExpected: 1,
			jExpected: 1,
		},
		{
			input:     Vector2D{205.0, 5.0},
			iExpected: 1,
			jExpected: 2,
		},
		{
			input:     Vector2D{205.0, 50.0},
			iExpected: 1,
			jExpected: 2,
		},
		{
			input:     Vector2D{195.0, 195.0},
			iExpected: 2,
			jExpected: 2,
		},
		{
			input:     Vector2D{195.0, 55.0},
			iExpected: 1,
			jExpected: 2,
		},
		{},
		{},
	}

	for i, testCase := range testCases {
		iActual, jActual := getMazeCoordinates(testCase.input)
		if (iActual != testCase.iExpected) || (jActual != testCase.jExpected) {
			t.Errorf("Test %d failed on input %v with expected result %d %d and actual result %d %d", i, testCase.input, testCase.iExpected, testCase.jExpected, iActual, jActual)
		}
	}
}
