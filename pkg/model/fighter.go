package model

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
	TotalSigStrLanded    int     `json:"totalSigStrLandned,omitempty"`
	TotalSigStrAttempted int     `json:"totalSigStrAttempted,omitempty"`
	StrAccuracy          int     `json:"strAccuracy,omitempty"`
	TotalTkdLanded       int     `json:"totalTkdLanded,omitempty"`
	TotalTkdAttempted    int     `json:"totalTkdAttempted,omitempty"`
	TkdAccuracy          int     `json:"tkdAccuracy,omitempty"`
	SigStrLanded         float32 `json:"sigStrLanded,omitempty"`
	SigStrAbs            float32 `json:"sigStrAbs,omitempty"`
	SigStrDefense        int8    `json:"sigStrDefense,omitempty"`
	TakedownDefense      int8    `json:"takedownDefense,omitempty"`
	TakedownAvg          float32 `json:"takedownAvg,omitempty"`
	SubmissionAvg        float32 `json:"submissionAvg,omitempty"`
	KnockdownAvg         float32 `json:"knockdownAvg,omitempty"`
	AvgFightTime         string  `json:"avgFightTime,omitempty"`
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

type FighterReq struct {
	FighterId int32 `json:"fighter_id"`
	*Fighter
}

type FighterStatsReq struct {
	FighterId int32 `json:"fighter_id"`
	*FighterStats
}

type FightersRequest struct {
	Status string `json:"status"`
}
