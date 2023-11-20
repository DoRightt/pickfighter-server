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
	TotalSigStrLanded    int     `json:"totalSigStrLandned"`
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
