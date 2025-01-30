package data

import (
	"strconv"
	"strings"

	"github.com/DoRightt/pickfighter-server/scraper/pkg/logger"
	"github.com/DoRightt/pickfighter-server/scraper/pkg/model"
)

const defaultStat = "0-0-0 (W-L-D)"

// SetStatistic sets the statistical data for a Fighter based on the provided string 'stat'.
// The function splits the input string, extracts individual parts, converts them to integers,
// and sets the Wins, Loses, and Draw fields of the Fighter accordingly. If conversion errors occur,
// it logs an error and sets the corresponding value to 0.
func SetStatistic(f *model.Fighter, stat string) {
	l := logger.Get()

	if len(stat) == 0 {
		stat = defaultStat
	}

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
func SetDivision(f *model.Fighter, d string) {
	switch d {
	case "Flyweight Division":
		f.DivisionId = model.Flyweight
	case "Bantamweight Division":
		f.DivisionId = model.Bantamweight
	case "Featherweight Division":
		f.DivisionId = model.Featherweight
	case "Lightweight Division":
		f.DivisionId = model.Lightweight
	case "Welterweight Division":
		f.DivisionId = model.Welterweight
	case "Middleweight Division":
		f.DivisionId = model.Middleweight
	case "Light Heavyweight Division":
		f.DivisionId = model.Lightheavyweight
	case "Heavyweight Division":
		f.DivisionId = model.Heavyweight
	case "Women's Strawweight Division":
		f.DivisionId = model.WomensStrawweight
	case "Women's Flyweight Division":
		f.DivisionId = model.WomensFlyweight
	case "Women's Bantamweight Division":
		f.DivisionId = model.WomensBantamweight
	case "Women's Featherweight Division":
		f.DivisionId = model.WomensFeatherweight
	default:
		f.DivisionId = 0
	}
}
