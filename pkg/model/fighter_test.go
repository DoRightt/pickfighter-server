package model

import "testing"

func TestFighter_String(t *testing.T) {
testCases := []struct {
		input    Division
		expected string
	}{
		{Flyweight, "Flyweight"},
		{Bantamweight, "Bantamweight"},
		{Featherweight, "Featherweight"},
		{Lightweight, "Lightweight"},
		{Welterweight, "Welterweight"},
		{Middleweight, "Middleweight"},
		{Lightheavyweight, "Light Heavyweight"},
		{Heavyweight, "Heavyweight"},
		{WomensStrawweight, "Women's Strawweight"},
		{WomensFlyweight, "Women's Flyweight"},
		{WomensBantamweight, "Women's Bantamweight"},
		{WomensFeatherweight, "Women's Featherweight"},
	}

	for _, testCase := range testCases {
		result := testCase.input.String()

		if result != testCase.expected {
			t.Errorf("For input %v, expected %s, but got %s", testCase.input, testCase.expected, result)
		}
	}
}
