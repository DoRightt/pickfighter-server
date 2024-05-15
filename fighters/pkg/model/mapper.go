package model

import "fightbettr.com/gen"

// FighterToProto converts a single Fighter struct into a generated proto counterpart.
func FighterToProto(f *Fighter) *gen.Fighter {
	return &gen.Fighter{
		FighterId:      f.FighterId,
		Name:           f.Name,
		NickName:       f.NickName,
		Division:       int32(f.Division),
		Status:         f.Status,
		Hometown:       f.Hometown,
		TrainsAt:       f.TrainsAt,
		FightingStyle:  f.FightingStyle,
		Age:            int32(f.Age),
		OctagonDebut:   f.OctagonDebut,
		DebutTimestamp: int32(f.DebutTimestamp),
		Reach:          f.Reach,
		LegReach:       f.LegReach,
		Wins:           int32(f.Wins),
		Loses:          int32(f.Loses),
		Draw:           int32(f.Draw),
		FighterUrl:     f.FighterUrl,
		ImageUrl:       f.ImageUrl,
		Stats: FighterStatsrToProto(&f.Stats),
	}
}

// FightersToProto converts a slice of Fighter structs into a slice of generated proto counterparts.
func FightersToProto(fs []*Fighter) []*gen.Fighter {
	var protoFighters = make([]*gen.Fighter, len(fs))

	for i, f := range fs {
		protoFighters[i] = FighterToProto(f)
	}

	return protoFighters
}

// FighterFromProto converts a generated proto counterpart into a single Fighter struct.
func FighterFromProto(f *gen.Fighter) *Fighter {
	return &Fighter{
		FighterId:      f.FighterId,
		Name:           f.Name,
		NickName:       f.NickName,
		Division:       Division(f.Division),
		Status:         f.Status,
		Hometown:       f.Hometown,
		TrainsAt:       f.TrainsAt,
		FightingStyle:  f.FightingStyle,
		Age:            int8(f.Age),
		OctagonDebut:   f.OctagonDebut,
		DebutTimestamp: int(f.DebutTimestamp),
		Reach:          f.Reach,
		LegReach:       f.LegReach,
		Wins:           int(f.Wins),
		Loses:          int(f.Loses),
		Draw:           int(f.Draw),
		FighterUrl:     f.FighterUrl,
		ImageUrl:       f.ImageUrl,
		Stats: *FighterStatsrFromProto(f.Stats),
	}
}

// FightersFromProto converts a slice of generated proto counterparts into a slice of Fighter structs.
func FightersFromProto(fs []*gen.Fighter) []*Fighter {
	var fighters = make([]*Fighter, len(fs))

	for i, f := range fs {
		fighters[i] = FighterFromProto(f)
	}

	return fighters
}

// FighterStatsrToProto converts a single Fighter stats struct into a generated proto counterpart.
func FighterStatsrToProto(f *FighterStats) *gen.FighterStats {
	return &gen.FighterStats{
		StatId:               f.StatId,
		FighterId:            f.FighterId,
		TotalSigStrLanded:    int32(f.TotalSigStrLanded),
		TotalSigStrAttempted: int32(f.TotalSigStrAttempted),
		StrAccuracy:          int32(f.StrAccuracy),
		TotalTkdLanded:       int32(f.TotalTkdLanded),
		TotalTkdAttempted:    int32(f.TotalTkdAttempted),
		TkdAccuracy:          int32(f.TkdAccuracy),
		SigStrLanded:         f.SigStrLanded,
		SigStrAbs:            f.SigStrAbs,
		SigStrDefense:        int32(f.SigStrDefense),
		TakedownDefense:      int32(f.TakedownDefense),
		TakedownAvg:          f.TakedownAvg,
		SubmissionAvg:        f.SubmissionAvg,
		KnockdownAvg:         f.KnockdownAvg,
		AvgFightTime:         f.AvgFightTime,
		WinByKO:              int32(f.WinByKO),
		WinBySub:             int32(f.WinBySub),
		WinByDec:             int32(f.WinByDec),
	}
}

// FighterStatsrFromProto converts a generated proto counterpart into a single Fighter stats struct..
func FighterStatsrFromProto(f *gen.FighterStats) *FighterStats {
	return &FighterStats{
		StatId:               f.StatId,
		FighterId:            f.FighterId,
		TotalSigStrLanded:    int(f.TotalSigStrLanded),
		TotalSigStrAttempted: int(f.TotalSigStrAttempted),
		StrAccuracy:          int(f.StrAccuracy),
		TotalTkdLanded:       int(f.TotalTkdLanded),
		TotalTkdAttempted:    int(f.TotalTkdAttempted),
		TkdAccuracy:          int(f.TkdAccuracy),
		SigStrLanded:         f.SigStrLanded,
		SigStrAbs:            f.SigStrAbs,
		SigStrDefense:        int8(f.SigStrDefense),
		TakedownDefense:      int8(f.TakedownDefense),
		TakedownAvg:          f.TakedownAvg,
		SubmissionAvg:        f.SubmissionAvg,
		KnockdownAvg:         f.KnockdownAvg,
		AvgFightTime:         f.AvgFightTime,
		WinByKO:              int(f.WinByKO),
		WinBySub:             int(f.WinBySub),
		WinByDec:             int(f.WinByDec),
	}
}
