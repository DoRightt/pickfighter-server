package scraper

import (
	"strconv"
	"strings"
)

type Division int

func (d Division) String() string {
	switch d {
	case Flyweight:
		return "Flyweight"
	case Bantamweight:
		return "Bantamweight"
	case Featherweight:
		return "Featherweight"
	case Lightweight:
		return "Lightweight"
	case Welterweight:
		return "Welterweight"
	case Middleweight:
		return "Middleweight"
	case Lightheavyweight:
		return "Light Heavyweight"
	case Heavyweight:
		return "Heavyweight"
	case WomensStrawweight:
		return "Women's Strawweight"
	case WomensFlyweight:
		return "Women's Flyweight"
	case WomensBantamweight:
		return "Women's Bantamweight"
	case WomensFeatherweight:
		return "Women's Featherweight"
	default:
		return "Unknown"
	}
}

const (
	Flyweight Division = iota
	Bantamweight
	Featherweight
	Lightweight
	Welterweight
	Middleweight
	Lightheavyweight
	Heavyweight
	WomensStrawweight
	WomensFlyweight
	WomensBantamweight
	WomensFeatherweight
)

type FightersCollection struct {
	Fighters []Fighter
}

type FighterStats struct {
	TotalSigStrLanded   int     `json:"totalSigStrLandned"`
	TotalSigStrAttempted int     `json:"totalSigStrAttempted"`
	StrAccuracy          int     `json:"strAccuracy"`
	TotalTkdLanded       int     `json:"totalTkdLanded"`
	TotalTkdAttempted    int     `json:"totalTkdAttempted"`
	TkdAccuracy          int     `json:"tkdAccuracy"`
	SigStrLanded         float32 `json:"sigStrLanded"`
	SigStrAbs            float32 `json:"sigStrAbs"`
	SigStrDefense        int8    `json:"sigStrDefense"`
	TakedownDefense      int8    `json:"takedownDefense"`
	TakedownAvg          float32 `json:"takedownAvg"`
	SubmissionAvg        float32 `json:"submissionAvg"`
	KnockdownAvg         float32 `json:"knockdownAvg"`
	AvgFightTime         string  `json:"avgFightTime"`
	WinByKO              int     `json:"winByKO"`
	WinBySub             int     `json:"winBySub"`
	WinByDec             int     `json:"winByDec"`
}

type Fighter struct {
	Name           string       `json:"name"`
	NickName       string       `json:"nickName"`
	Division       Division     `json:"division"`
	Status         string       `json:"status"`
	Hometown       string       `json:"hometown"`
	TrainsAt       string       `json:"trainsAt"`
	FightingStyle  string       `json:"fightingStyle"`
	Age            int8         `json:"age"`
	Height         float32      `json:"height"`
	Weight         float32      `json:"weight"`
	OctagonDebut   string       `json:"octagonDebut"`
	DebutTimestamp int          `json:"debutTimestamp"`
	Reach          float32      `json:"reach"`
	LegReach       float32      `json:"legReach"`
	Wins           int          `json:"wins"`
	Loses          int          `json:"loses"`
	Draw           int          `json:"draw"`
	FighterUrl     string       `json:"fighterUrl"`
	ImageUrl       string       `json:"imageUrl"`
	Stats          FighterStats `json:"stats"`
}

func (f *Fighter) SetStatistic(stat string) {
	parts := strings.Split(strings.Split(stat, " ")[0], "-")
	var scores []int

	for _, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil {
			l.Errorf("[%s] Conversion error: %s, with part: '%s' of %s", f.Name, err, part, parts)
			scores = append(scores, 0)
			return
		}
		scores = append(scores, num)
	}

	f.Wins = scores[0]
	f.Loses = scores[1]
	f.Draw = scores[2]
}

func (f *Fighter) SetDivision(d string) {
	switch d {
	case "Flyweight Division":
		f.Division = Flyweight
	case "Bantamweight Division":
		f.Division = Bantamweight
	case "Featherweight Division":
		f.Division = Featherweight
	case "Lightweight Division":
		f.Division = Lightweight
	case "Welterweight Division":
		f.Division = Welterweight
	case "Middleweight Division":
		f.Division = Middleweight
	case "Light Heavyweight Division":
		f.Division = Lightheavyweight
	case "Heavyweight Division":
		f.Division = Heavyweight
	case "Women's Strawweight Division":
		f.Division = WomensStrawweight
	case "Women's Flyweight Division":
		f.Division = WomensFlyweight
	case "Women's Bantamweight Division":
		f.Division = WomensBantamweight
	case "Women's Featerweight Division":
		f.Division = WomensFeatherweight
	}
}
