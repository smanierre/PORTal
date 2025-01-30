package types

type RankMap map[Grade]string

type RankTuple struct {
	Grade Grade
	Rank  string
}

func (r RankMap) SortedRanks() []RankTuple {
	// Not the prettiest, but it won't need to change and isn't worth over complicating
	return []RankTuple{
		{Grade: E1, Rank: r[E1]},
		{Grade: E2, Rank: r[E2]},
		{Grade: E3, Rank: r[E3]},
		{Grade: E4, Rank: r[E4]},
		{Grade: E5, Rank: r[E5]},
		{Grade: E6, Rank: r[E6]},
		{Grade: E7, Rank: r[E7]},
		{Grade: E8, Rank: r[E8]},
		{Grade: E9, Rank: r[E9]},
	}
}

func GetRanks(service string) RankMap {
	switch service {
	case "f":
		return AfRankMap
	case "a":
		return AfRankMap
	case "n":
		return NavyRankMap
	case "m":
		return MarineRankMap
	}
	return map[Grade]string{}
}

var AfRankMap RankMap = map[Grade]string{
	E1: "AB",
	E2: "Amn",
	E3: "A1C",
	E4: "SrA",
	E5: "SSgt",
	E6: "TSgt",
	E7: "MSgt",
	E8: "SMSgt",
	E9: "CMSgt",
}

var ArmyRankMap RankMap = map[Grade]string{
	E1: "PVT",
	E2: "PV2",
	E3: "PFC",
	E4: "CPL/SFC",
	E5: "SGT",
	E6: "SSG",
	E7: "SFC",
	E8: "MSG",
	E9: "SGM",
}

var MarineRankMap RankMap = map[Grade]string{
	E1: "PVT",
	E2: "PFC",
	E3: "LCpl",
	E4: "Cpl",
	E5: "Sgt",
	E6: "SSgt",
	E7: "GySgt",
	E8: "MSgt",
	E9: "MgySgt",
}

var NavyRankMap RankMap = map[Grade]string{
	E1: "SR",
	E2: "SA",
	E3: "SN",
	E4: "PO3",
	E5: "PO2",
	E6: "PO1",
	E7: "CPO",
	E8: "SCPO",
	E9: "MCPO",
}
