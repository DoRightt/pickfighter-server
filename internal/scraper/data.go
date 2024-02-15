package scraper

import (
	fighterModel "projects/fb-server/pkg/model"
	"strconv"
	"strings"
)

// SetStatistic sets the statistical data for a Fighter based on the provided string 'stat'.
// The function splits the input string, extracts individual parts, converts them to integers,
// and sets the Wins, Loses, and Draw fields of the Fighter accordingly. If conversion errors occur,
// it logs an error and sets the corresponding value to 0.
func SetStatistic(f *fighterModel.Fighter, stat string) {
	parts := strings.Split(strings.Split(stat, " ")[0], "-")
	var scores []int

	for _, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil {
			l.Errorf("[%s] Conversion error: %s, with part: '%s' of %s", f.Name, err, part, parts)
			scores = append(scores, 0)
		} else {
			scores = append(scores, num)
		}
		
	}

	f.Wins = scores[0]
	f.Loses = scores[1]
	f.Draw = scores[2]
}

// SetDivision sets division based on Division type.
func SetDivision(f *fighterModel.Fighter, d string) {
	switch d {
	case "Flyweight Division":
		f.Division = fighterModel.Flyweight
	case "Bantamweight Division":
		f.Division = fighterModel.Bantamweight
	case "Featherweight Division":
		f.Division = fighterModel.Featherweight
	case "Lightweight Division":
		f.Division = fighterModel.Lightweight
	case "Welterweight Division":
		f.Division = fighterModel.Welterweight
	case "Middleweight Division":
		f.Division = fighterModel.Middleweight
	case "Light Heavyweight Division":
		f.Division = fighterModel.Lightheavyweight
	case "Heavyweight Division":
		f.Division = fighterModel.Heavyweight
	case "Women's Strawweight Division":
		f.Division = fighterModel.WomensStrawweight
	case "Women's Flyweight Division":
		f.Division = fighterModel.WomensFlyweight
	case "Women's Bantamweight Division":
		f.Division = fighterModel.WomensBantamweight
	case "Women's Featerweight Division":
		f.Division = fighterModel.WomensFeatherweight
	}
}
