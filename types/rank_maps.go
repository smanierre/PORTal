package types

func GetRanks(service string) map[Grade]string {
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

var AfRankMap = map[Grade]string{
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

var ArmyRankMap = map[Grade]string{
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

var MarineRankMap = map[Grade]string{
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

var NavyRankMap = map[Grade]string{
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
