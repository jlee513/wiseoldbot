package util

import "github.com/bwmarrin/discordgo"

// ******************************** KC HOF ********************************

var HofSlayerBosses = []BossInfo{
	{BossName: "sire", ImageLink: "https://i.imgur.com/GhbmqEB.png"},
	{BossName: "hydra", ImageLink: "https://i.imgur.com/25GU0Ph.png"},
	{BossName: "cerberus", ImageLink: "https://i.imgur.com/UoxGuQi.png"},
	{BossName: "grotesqueguardians", ImageLink: "https://i.imgur.com/M7ylVBZ.png"},
	{BossName: "kraken", ImageLink: "https://i.imgur.com/Q6EbJb1.png"},
	{BossName: "smokedevil", ImageLink: "https://i.imgur.com/2AYntQ5.png"},
}

var HofGWDBosses = []BossInfo{
	{BossName: "commanderzilyana", ImageLink: "https://i.imgur.com/aNm4Ydd.png"},
	{BossName: "kreearra", ImageLink: "https://i.imgur.com/lX8SfgN.png"},
	{BossName: "kriltsutsaroth", ImageLink: "https://i.imgur.com/hh8cMvp.png"},
	{BossName: "nex", ImageLink: "https://i.imgur.com/pqiVQBC.png"},
	{BossName: "generalgraardor", ImageLink: "https://i.imgur.com/hljv9ZW.png"},
}

var HofWildyBosses = []BossInfo{
	{BossName: "artio", ImageLink: "https://i.imgur.com/bw6zLpU.png"},
	{BossName: "callisto", ImageLink: "https://i.imgur.com/bw6zLpU.png"},
	{BossName: "calvarion", ImageLink: "https://i.imgur.com/v3KX75y.png"},
	{BossName: "vetion", ImageLink: "https://i.imgur.com/v3KX75y.png"},
	{BossName: "spindel", ImageLink: "https://i.imgur.com/4zknWSX.png"},
	{BossName: "venenatis", ImageLink: "https://i.imgur.com/4zknWSX.png"},
	{BossName: "chaoselemental", ImageLink: "https://i.imgur.com/YAvIpbm.png"},
	{BossName: "chaosfanatic", ImageLink: "https://i.imgur.com/azV2sD1.png"},
	{BossName: "crazyarchaeologist", ImageLink: "https://i.imgur.com/23LXv53.png"},
	{BossName: "scorpia", ImageLink: "https://i.imgur.com/9aaguxB.png"},
}

var HofOtherBosses = []BossInfo{
	{BossName: "corporealbeast", ImageLink: "https://i.imgur.com/zEDN4Pf.png"},
	{BossName: "prime", ImageLink: "https://i.imgur.com/kJBtqHB.png"},
	{BossName: "rexbro", ImageLink: "https://i.imgur.com/PvlGWFZ.png"},
	{BossName: "supreme", ImageLink: "https://i.imgur.com/BOgkBuD.png"},
	{BossName: "gauntlet", ImageLink: "https://i.imgur.com/weiHWnz.png"},
	{BossName: "gauntlethard", ImageLink: "https://i.imgur.com/xzW4TGR.png"},
	{BossName: "giantmole", ImageLink: "https://i.imgur.com/coKk2pr.png"},
	{BossName: "jad", ImageLink: "https://i.imgur.com/H9aO1Ot.png"},
	{BossName: "zuk", ImageLink: "https://i.imgur.com/mKstHza.png"},
	{BossName: "kq", ImageLink: "https://i.imgur.com/ZuaFoBR.png"},
	{BossName: "kbd", ImageLink: "https://i.imgur.com/r5vkw1s.png"},
	{BossName: "sarachnis", ImageLink: "https://i.imgur.com/98THH8O.png"},
	{BossName: "skotizo", ImageLink: "https://i.imgur.com/YUcQu4d.png"},
	{BossName: "muspah", ImageLink: "https://i.imgur.com/sW2cLQ2.png"},
	{BossName: "vorkath", ImageLink: "https://i.imgur.com/6biF3P2.png"},
	{BossName: "nightmare", ImageLink: "https://i.imgur.com/4aDkxms.png"},
	{BossName: "phosanis", ImageLink: "https://i.imgur.com/4aDkxms.png"},
	{BossName: "zulrah", ImageLink: "https://i.imgur.com/tPllWNF.png"},
}

var HofMiscBosses = []BossInfo{
	{BossName: "barrows", ImageLink: "https://i.imgur.com/ajoK20v.png"},
	{BossName: "hespori", ImageLink: "https://i.imgur.com/b0qYGHS.png"},
	{BossName: "mimic", ImageLink: "https://i.imgur.com/jC7yTC3.png"},
	{BossName: "obor", ImageLink: "https://i.imgur.com/dwLvSbR.png"},
	{BossName: "bryophyta", ImageLink: "https://i.imgur.com/3cdyp4X.png"},
	{BossName: "derangedarchaeologist", ImageLink: "https://i.imgur.com/cnHpevF.png"},
	{BossName: "wintertodt", ImageLink: "https://i.imgur.com/6oFef2Y.png"},
	{BossName: "zalcano", ImageLink: "https://i.imgur.com/edN11Nf.png"},
	{BossName: "rift", ImageLink: "https://i.imgur.com/MOiyXeH.png"},
}

var HofDT2Bosses = []BossInfo{
	{BossName: "duke", ImageLink: "https://i.imgur.com/RYPmrXy.png"},
	{BossName: "leviathan", ImageLink: "https://i.imgur.com/mEQRq5c.png"},
	{BossName: "whisperer", ImageLink: "https://i.imgur.com/cFGWb6Y.png"},
	{BossName: "vardorvis", ImageLink: "https://i.imgur.com/WMPuShZ.png"},
}

var HofRaidsBosses = []BossInfo{
	{BossName: "cox", ImageLink: "https://i.imgur.com/gxdWXtH.png"},
	{BossName: "coxcm", ImageLink: "https://i.imgur.com/gxdWXtH.png"},
	{BossName: "tob", ImageLink: "https://i.imgur.com/pW1sJAQ.png"},
	{BossName: "tobcm", ImageLink: "https://i.imgur.com/pW1sJAQ.png"},
	{BossName: "toa", ImageLink: "https://i.imgur.com/2GvzqGw.png"},
	{BossName: "toae", ImageLink: "https://i.imgur.com/2GvzqGw.png"},
}

var HofPVPBosses = []BossInfo{
	{BossName: "bhh", ImageLink: "https://i.imgur.com/zSQhlWk.png"},
	{BossName: "bhr", ImageLink: "https://i.imgur.com/Y3Sga7t.png"},
	{BossName: "lms", ImageLink: "https://i.imgur.com/rzW7ZXx.png"},
	{BossName: "arena", ImageLink: "https://i.imgur.com/uNP6Ggu.png"},
	{BossName: "zeal", ImageLink: "https://i.imgur.com/Ws7HvKL.png"},
}

var HofCluesBosses = []BossInfo{
	{BossName: "cluebeginner", ImageLink: "https://i.imgur.com/fUmzJkW.png"},
	{BossName: "clueeasy", ImageLink: "https://i.imgur.com/phnSCHj.png"},
	{BossName: "cluemedium", ImageLink: "https://i.imgur.com/t5iH8Xa.png"},
	{BossName: "cluehard", ImageLink: "https://i.imgur.com/a0xwcGI.png"},
	{BossName: "clueelite", ImageLink: "https://i.imgur.com/ibNRk3G.png"},
	{BossName: "cluemaster", ImageLink: "https://i.imgur.com/12rCLVv.png"},
	{BossName: "clueall", ImageLink: "https://i.imgur.com/wX3Ei7U.png"},
}

// ******************************** SPEED HOF ********************************

var HofSpeedCategories = map[string]struct{}{
	"TzHaar": {}, "Slayer": {}, "Nightmare": {}, "Nex": {}, "Solo Bosses": {}, "Chambers Of Xeric": {}, "Chambers Of Xeric Challenge Mode": {}, "Theatre Of Blood": {}, "Theatre Of Blood Hard Mode": {}, "Tombs Of Amascut": {}, "Tombs Of Amascut Expert": {}, "Agility": {}, "Desert Treasure 2": {},
}

var SpeedBossNameToCategory = map[string]string{
	"TzHaar Fight Cave":                       "TzHaar",
	"Inferno":                                 "TzHaar",
	"Alchemical Hydra":                        "Slayer",
	"Grotesque Guardians":                     "Slayer",
	"The Nightmare Solo":                      "Nightmare",
	"The Nightmare 2":                         "Nightmare",
	"The Nightmare 3":                         "Nightmare",
	"The Nightmare 4":                         "Nightmare",
	"The Nightmare 5":                         "Nightmare",
	"The Nightmare 6+":                        "Nightmare",
	"Phosani's Nightmare":                     "Nightmare",
	"Nex":                                     "Nex",
	"Fragment of Seren":                       "Solo Bosses",
	"Galvek":                                  "Solo Bosses",
	"The Gauntlet":                            "Solo Bosses",
	"The Corrupted Gauntlet":                  "Solo Bosses",
	"Hespori":                                 "Solo Bosses",
	"Phantom Muspah":                          "Solo Bosses",
	"Vorkath":                                 "Solo Bosses",
	"Zulrah":                                  "Solo Bosses",
	"Tempoross":                               "Solo Bosses",
	"Chambers of Xeric Solo":                  "Chambers Of Xeric",
	"Chambers of Xeric 2":                     "Chambers Of Xeric",
	"Chambers of Xeric 3":                     "Chambers Of Xeric",
	"Chambers of Xeric 4":                     "Chambers Of Xeric",
	"Chambers of Xeric 5":                     "Chambers Of Xeric",
	"Chambers of Xeric 6":                     "Chambers Of Xeric",
	"Chambers of Xeric 7":                     "Chambers Of Xeric",
	"Chambers of Xeric 8":                     "Chambers Of Xeric",
	"Chambers of Xeric 9":                     "Chambers Of Xeric",
	"Chambers of Xeric 10":                    "Chambers Of Xeric",
	"Chambers of Xeric 11-15":                 "Chambers Of Xeric",
	"Chambers of Xeric 16-23":                 "Chambers Of Xeric",
	"Chambers of Xeric 24+":                   "Chambers Of Xeric",
	"Chambers of Xeric - Challenge mode Solo": "Chambers Of Xeric Challenge mode",
	"Chambers of Xeric - Challenge mode 2":    "Chambers Of Xeric Challenge mode",
	"Chambers of Xeric - Challenge mode 3":    "Chambers Of Xeric Challenge mode",
	"Chambers of Xeric - Challenge mode 4":    "Chambers Of Xeric Challenge mode",
	"Chambers of Xeric - Challenge mode 5":    "Chambers Of Xeric Challenge mode",
	"Chambers of Xeric - Challenge mode 6":    "Chambers Of Xeric Challenge mode",
	"Chambers of Xeric - Challenge mode 7":    "Chambers Of Xeric Challenge mode",
	"Chambers of Xeric - Challenge mode 8":    "Chambers Of Xeric Challenge mode",
	"Chambers of Xeric - Challenge mode 9":    "Chambers Of Xeric Challenge mode",
	"Chambers of Xeric - Challenge mode 10":   "Chambers Of Xeric Challenge mode",
	"Theatre of Blood Room 2":                 "Theatre Of Blood",
	"Theatre of Blood Room 3":                 "Theatre Of Blood",
	"Theatre of Blood Room 4":                 "Theatre Of Blood",
	"Theatre of Blood Room 5":                 "Theatre Of Blood",
	"Theatre of Blood Overall 2":              "Theatre Of Blood",
	"Theatre of Blood Overall 3":              "Theatre Of Blood",
	"Theatre of Blood Overall 4":              "Theatre Of Blood",
	"Theatre of Blood Overall 5":              "Theatre Of Blood",
	"Theatre of Blood - Hard Room 3":          "Theatre Of Blood Hard Mode",
	"Theatre of Blood - Hard Room 4":          "Theatre Of Blood Hard Mode",
	"Theatre of Blood - Hard Room 5":          "Theatre Of Blood Hard Mode",
	"Theatre of Blood - Hard Overall 3":       "Theatre Of Blood Hard Mode",
	"Theatre of Blood - Hard Overall 4":       "Theatre Of Blood Hard Mode",
	"Theatre of Blood - Hard Overall 5":       "Theatre Of Blood Hard Mode",
	"Tombs of Amascut Room Solo":              "Tombs Of Amascut",
	"Tombs of Amascut Room 2":                 "Tombs Of Amascut",
	"Tombs of Amascut Room 3":                 "Tombs Of Amascut",
	"Tombs of Amascut Room 4":                 "Tombs Of Amascut",
	"Tombs of Amascut Room 5":                 "Tombs Of Amascut",
	"Tombs of Amascut Room 6":                 "Tombs Of Amascut",
	"Tombs of Amascut Room 7":                 "Tombs Of Amascut",
	"Tombs of Amascut Overall Solo":           "Tombs Of Amascut",
	"Tombs of Amascut Overall 2":              "Tombs Of Amascut",
	"Tombs of Amascut Overall 3":              "Tombs Of Amascut",
	"Tombs of Amascut Overall 4":              "Tombs Of Amascut",
	"Tombs of Amascut Overall 5":              "Tombs Of Amascut",
	"Tombs of Amascut Overall 6":              "Tombs Of Amascut",
	"Tombs of Amascut Overall 7":              "Tombs Of Amascut",
	"Tombs of Amascut Expert Room Solo":       "Tombs Of Amascut Expert",
	"Tombs of Amascut Expert Room 2":          "Tombs Of Amascut Expert",
	"Tombs of Amascut Expert Room 3":          "Tombs Of Amascut Expert",
	"Tombs of Amascut Expert Room 4":          "Tombs Of Amascut Expert",
	"Tombs of Amascut Expert Room 5":          "Tombs Of Amascut Expert",
	"Tombs of Amascut Expert Room 6":          "Tombs Of Amascut Expert",
	"Tombs of Amascut Expert Room 7":          "Tombs Of Amascut Expert",
	"Tombs of Amascut Expert Room 8":          "Tombs Of Amascut Expert",
	"Tombs of Amascut Expert Overall Solo":    "Tombs Of Amascut Expert",
	"Tombs of Amascut Expert Overall 2":       "Tombs Of Amascut Expert",
	"Tombs of Amascut Expert Overall 3":       "Tombs Of Amascut Expert",
	"Tombs of Amascut Expert Overall 4":       "Tombs Of Amascut Expert",
	"Tombs of Amascut Expert Overall 5":       "Tombs Of Amascut Expert",
	"Tombs of Amascut Expert Overall 6":       "Tombs Of Amascut Expert",
	"Tombs of Amascut Expert Overall 7":       "Tombs Of Amascut Expert",
	"Tombs of Amascut Expert Overall 8":       "Tombs Of Amascut Expert",
	"Hallowed Sepulchre":                      "Agility",
	"Prifddinas Agility Course":               "Agility",
	"Duke Sucellus":                           "Desert Treasure 2",
	"The Leviathan":                           "Desert Treasure 2",
	"The Whisperer":                           "Desert Treasure 2",
	"Vardorvis":                               "Desert Treasure 2",
}

var HofSpeedDt2 = []BossInfo{
	{BossName: "Duke Sucellus"},
	{BossName: "The Leviathan"},
	{BossName: "The Whisperer"},
	{BossName: "Vardorvis"},
}

var HofSpeedTzhaar = []BossInfo{
	{BossName: "TzHaar Fight Cave"},
	{BossName: "Inferno"},
}

var HofSpeedSlayer = []BossInfo{
	{BossName: "Alchemical Hydra"},
	{BossName: "Grotesque Guardians"},
}

var HofSpeedNightmare = []BossInfo{
	{BossName: "The Nightmare Solo"},
	{BossName: "The Nightmare 2"},
	{BossName: "The Nightmare 3"},
	{BossName: "The Nightmare 4"},
	{BossName: "The Nightmare 5"},
	{BossName: "The Nightmare 6+"},
	{BossName: "Phosani's Nightmare"},
}

var HofSpeedNex = []BossInfo{
	{BossName: "Nex", ImageLink: "https://i.imgur.com/hnw9ki8.png"},
}

var HofSpeedSolo = []BossInfo{
	{BossName: "Fragment of Seren"},
	{BossName: "Galvek"},
	{BossName: "The Gauntlet"},
	{BossName: "The Corrupted Gauntlet"},
	{BossName: "Hespori"},
	{BossName: "Phantom Muspah"},
	{BossName: "Vorkath"},
	{BossName: "Zulrah"},
	{BossName: "Tempoross"},
}

var HofSpeedCox = []BossInfo{
	{BossName: "Chambers of Xeric Solo"},
	{BossName: "Chambers of Xeric 2"},
	{BossName: "Chambers of Xeric 3"},
	{BossName: "Chambers of Xeric 4"},
	{BossName: "Chambers of Xeric 5"},
	{BossName: "Chambers of Xeric 6"},
	{BossName: "Chambers of Xeric 7"},
	{BossName: "Chambers of Xeric 8"},
	{BossName: "Chambers of Xeric 9"},
	{BossName: "Chambers of Xeric 10"},
	{BossName: "Chambers of Xeric 11-15"},
	{BossName: "Chambers of Xeric 16-23"},
	{BossName: "Chambers of Xeric 24+"},
}

var HofSpeedCoxCm = []BossInfo{
	{BossName: "Chambers of Xeric - Challenge mode Solo"},
	{BossName: "Chambers of Xeric - Challenge mode 2"},
	{BossName: "Chambers of Xeric - Challenge mode 3"},
	{BossName: "Chambers of Xeric - Challenge mode 4"},
	{BossName: "Chambers of Xeric - Challenge mode 5"},
	{BossName: "Chambers of Xeric - Challenge mode 6"},
	{BossName: "Chambers of Xeric - Challenge mode 7"},
	{BossName: "Chambers of Xeric - Challenge mode 8"},
	{BossName: "Chambers of Xeric - Challenge mode 9"},
	{BossName: "Chambers of Xeric - Challenge mode 10"},
}

var HofSpeedTob = []BossInfo{
	{BossName: "Theatre of Blood Room 2"},
	{BossName: "Theatre of Blood Room 3"},
	{BossName: "Theatre of Blood Room 4"},
	{BossName: "Theatre of Blood Room 5"},
	{BossName: "Theatre of Blood Overall 2"},
	{BossName: "Theatre of Blood Overall 3"},
	{BossName: "Theatre of Blood Overall 4"},
	{BossName: "Theatre of Blood Overall 5"},
}

var HofSpeedTobHm = []BossInfo{
	{BossName: "Theatre of Blood - Hard Room 3"},
	{BossName: "Theatre of Blood - Hard Room 4"},
	{BossName: "Theatre of Blood - Hard Room 5"},
	{BossName: "Theatre of Blood - Hard Overall 3"},
	{BossName: "Theatre of Blood - Hard Overall 4"},
	{BossName: "Theatre of Blood - Hard Overall 5"},
}

var HofSpeedToa = []BossInfo{
	{BossName: "Tombs of Amascut Room Solo"},
	{BossName: "Tombs of Amascut Room 2"},
	{BossName: "Tombs of Amascut Room 3"},
	{BossName: "Tombs of Amascut Room 4"},
	{BossName: "Tombs of Amascut Room 5"},
	{BossName: "Tombs of Amascut Room 6"},
	{BossName: "Tombs of Amascut Room 7"},
	{BossName: "Tombs of Amascut Overall Solo"},
	{BossName: "Tombs of Amascut Overall 2"},
	{BossName: "Tombs of Amascut Overall 3"},
	{BossName: "Tombs of Amascut Overall 4"},
	{BossName: "Tombs of Amascut Overall 5"},
	{BossName: "Tombs of Amascut Overall 6"},
	{BossName: "Tombs of Amascut Overall 7"},
}

var HofSpeedToae = []BossInfo{
	{BossName: "Tombs of Amascut Expert Room Solo"},
	{BossName: "Tombs of Amascut Expert Room 2"},
	{BossName: "Tombs of Amascut Expert Room 3"},
	{BossName: "Tombs of Amascut Expert Room 4"},
	{BossName: "Tombs of Amascut Expert Room 5"},
	{BossName: "Tombs of Amascut Expert Room 6"},
	{BossName: "Tombs of Amascut Expert Room 7"},
	{BossName: "Tombs of Amascut Expert Room 8"},
	{BossName: "Tombs of Amascut Expert Overall Solo"},
	{BossName: "Tombs of Amascut Expert Overall 2"},
	{BossName: "Tombs of Amascut Expert Overall 3"},
	{BossName: "Tombs of Amascut Expert Overall 4"},
	{BossName: "Tombs of Amascut Expert Overall 5"},
	{BossName: "Tombs of Amascut Expert Overall 6"},
	{BossName: "Tombs of Amascut Expert Overall 7"},
	{BossName: "Tombs of Amascut Expert Overall 8"},
}

var HofSpeedAgility = []BossInfo{
	{BossName: "Hallowed Sepulchre"},
	{BossName: "Prifddinas Agility Course"},
}

var SpeedAutocompleteCategories = []*discordgo.ApplicationCommandOptionChoice{
	{
		Name:  "TzHaar",
		Value: "TzHaar",
	},
	{
		Name:  "Chambers Of Xeric",
		Value: "Chambers Of Xeric",
	},
	{
		Name:  "Chambers Of Xeric Challenge Mode",
		Value: "Chambers Of Xeric Challenge Mode",
	},
	{
		Name:  "Nightmare",
		Value: "Nightmare",
	},
	{
		Name:  "Theatre Of Blood Hard Mode",
		Value: "Theatre Of Blood Hard Mode",
	},
	{
		Name:  "Agility",
		Value: "Agility",
	},
	{
		Name:  "Tombs Of Amascut Expert",
		Value: "Tombs Of Amascut Expert",
	},
	{
		Name:  "Solo Bosses",
		Value: "Solo Bosses",
	},
	{
		Name:  "Nex",
		Value: "Nex",
	},
	{
		Name:  "Slayer",
		Value: "Slayer",
	},
	{
		Name:  "Desert Treasure 2",
		Value: "Desert Treasure 2",
	},
}

// ******************************** LOOT LOG INFO ********************************

var LootLogClanPoint = map[string]string{

	// COX
	"Dexterous prayer scroll":      "COX",
	"Arcane prayer scroll":         "COX",
	"Twisted buckler":              "COX",
	"Dragon hunter crossbow":       "COX",
	"Dinh's bulwark":               "COX",
	"Ancestral hat":                "COX",
	"Ancestral robe top":           "COX",
	"Ancestral robe bottom":        "COX",
	"Dragon claws":                 "COX",
	"Elder maul":                   "COX",
	"Kodai insignia":               "COX",
	"Twisted bow":                  "COX",
	"Twisted ancestral colour kit": "COX",
	"Metamorphic dust":             "COX",

	// TOB
	"Avernic defender hilt":         "TOB",
	"Ghrazi rapier":                 "TOB",
	"Sanguinesti staff":             "TOB",
	"Sanguinesti staff (uncharged)": "TOB",
	"Justiciar faceguard":           "TOB",
	"Justiciar chestguard":          "TOB",
	"Justiciar legguards":           "TOB",
	"Scythe of vitur":               "TOB",
	"Scythe of vitur (uncharged)":   "TOB",
	"Holy ornament kit":             "TOB",
	"Sanguine ornament kit":         "TOB",
	"Sanguine dust":                 "TOB",

	// TOA
	"Osmumten's fang":              "TOA",
	"Lightbearer":                  "TOA",
	"Elidinis' ward":               "TOA",
	"Masori mask":                  "TOA",
	"Masori body":                  "TOA",
	"Masori chaps":                 "TOA",
	"Tumeken's shadow":             "TOA",
	"Tumeken's shadow (uncharged)": "TOA",
	"Thread of elidinis":           "TOA",
	"Eye of the corruptor":         "TOA",
	"Jewel of the sun":             "TOA",
	"Breach of the scarab":         "TOA",

	// Clue Mega Rares
	"Gilded full helm":             "Clue Mega Rares",
	"Gilded platebody":             "Clue Mega Rares",
	"Gilded platelegs":             "Clue Mega Rares",
	"Gilded plateskirt":            "Clue Mega Rares",
	"Gilded kiteshield":            "Clue Mega Rares",
	"Gilded med helm":              "Clue Mega Rares",
	"Gilded chainbody":             "Clue Mega Rares",
	"Gilded sq shield":             "Clue Mega Rares",
	"Gilded 2h sword":              "Clue Mega Rares",
	"Gilded spear":                 "Clue Mega Rares",
	"Gilded hasta":                 "Clue Mega Rares",
	"3rd age full helmet":          "Clue Mega Rares",
	"3rd age platebody":            "Clue Mega Rares",
	"3rd age platelegs":            "Clue Mega Rares",
	"3rd age plateskirt":           "Clue Mega Rares",
	"3rd age kiteshield":           "Clue Mega Rares",
	"3rd age range coif":           "Clue Mega Rares",
	"3rd age range top":            "Clue Mega Rares",
	"3rd age range legs":           "Clue Mega Rares",
	"3rd age vambraces":            "Clue Mega Rares",
	"3rd age mage hat":             "Clue Mega Rares",
	"3rd age robe top":             "Clue Mega Rares",
	"3rd age robe":                 "Clue Mega Rares",
	"3rd age amulet":               "Clue Mega Rares",
	"Gilded scimitar":              "Clue Mega Rares",
	"Gilded boots":                 "Clue Mega Rares",
	"Gilded coif":                  "Clue Mega Rares",
	"Gilded d'hide vambraces":      "Clue Mega Rares",
	"Gilded d'hide body":           "Clue Mega Rares",
	"Gilded d'hide chaps":          "Clue Mega Rares",
	"Gilded pickaxe":               "Clue Mega Rares",
	"Gilded axe":                   "Clue Mega Rares",
	"Gilded spade":                 "Clue Mega Rares",
	"3rd age longsword":            "Clue Mega Rares",
	"3rd age wand":                 "Clue Mega Rares",
	"3rd age cloak":                "Clue Mega Rares",
	"3rd age bow":                  "Clue Mega Rares",
	"Bucket helm (g)":              "Clue Mega Rares",
	"Ring of coins":                "Clue Mega Rares",
	"3rd age druidic robe top":     "Clue Mega Rares",
	"3rd age druidic robe bottoms": "Clue Mega Rares",
	"3rd age druidic cloak":        "Clue Mega Rares",
	"3rd age druidic staff":        "Clue Mega Rares",
	"3rd age pickaxe":              "Clue Mega Rares",
	"3rd age axe":                  "Clue Mega Rares",

	// Slayer
	"Hydra's claw":          "Slayer",
	"Hydra leather":         "Slayer",
	"Primordial crystal":    "Slayer",
	"Pegasian crystal":      "Slayer",
	"Eternal crystal":       "Slayer",
	"Occult necklace":       "Slayer",
	"Dragon harpoon":        "Slayer",
	"Black tourmaline core": "Slayer",
	"Smoke battlestaff":     "Slayer",
	"Mist battlestaff":      "Slayer",
	"Dust battlestaff":      "Slayer",
	"Imbued heart":          "Slayer",
	"Eternal gem":           "Slayer",
	"Zenyte shard":          "Slayer",
	"Monkey tail":           "Slayer",
	"Heavy frame":           "Slayer",
	"Light frame":           "Slayer",
	"Ballista spring":       "Slayer",
	"Ballista limbs":        "Slayer",

	// Wilderness
	"Claws of callisto":       "Wilderness",
	"Tyrannical ring":         "Wilderness",
	"Voidwaker hilt":          "Wilderness",
	"Dragon pickaxe":          "Wilderness",
	"Skull of vet'ion":        "Wilderness",
	"Ring of the gods":        "Wilderness",
	"Voidwaker blade":         "Wilderness",
	"Fangs of venenatis":      "Wilderness",
	"Treasonous ring":         "Wilderness",
	"Voidwaker gem":           "Wilderness",
	"Amulet of avarice":       "Wilderness",
	"Craw's bow":              "Wilderness",
	"Craw's bow (u)":          "Wilderness",
	"Thammaron's sceptre":     "Wilderness",
	"Thammaron's sceptre (u)": "Wilderness",
	"Viggora's chainmace":     "Wilderness",
	"Viggora's chainmace (u)": "Wilderness",
	"Ancient crystal":         "Wilderness",
	"Ancient relic":           "Wilderness",
	"Ancient effigy":          "Wilderness",
	"Ancient medallion":       "Wilderness",
	"Ancient statuette":       "Wilderness",
	"Ancient totem":           "Wilderness",
	"Ancient emblem":          "Wilderness",
	"Dagon'hai robe top":      "Wilderness",
	"Dagon'hai hat":           "Wilderness",
	"Dagon'hai robe bottom":   "Wilderness",

	// DT2
	"Chromium ingot":             "DT2",
	"Ultor vestige":              "DT2",
	"Executioner's axe head":     "DT2",
	"Virtus mask":                "DT2",
	"Virtus robe top":            "DT2",
	"Virtus robe bottom":         "DT2",
	"Ancient blood ornament kit": "DT2",
	"Bellator vestige":           "DT2",
	"Siren's staff":              "DT2",
	"Venator vestige":            "DT2",
	"Leviathan's lure":           "DT2",
	"Magus vestige":              "DT2",
	"Eye of the duke":            "DT2",

	// Other
	"Bryophyta's essence":     "Other",
	"Hill giant club":         "Other",
	"Dragon full helm":        "Other",
	"Infernal cape":           "Other",
	"Wyvern visage":           "Other",
	"Skeletal visage":         "Other",
	"Draconic visage":         "Other",
	"Dragon warhammer":        "Other",
	"Amulet of eternal glory": "Other",
	"Evil chicken head":       "Other",
	"Evil chicken wings":      "Other",
	"Evil chicken legs":       "Other",
	"Evil chicken feet":       "Other",
	"Abyssal lantern":         "Other",
	"Abyssal needle":          "Other",
	"Abyssal green dye":       "Other",
	"Abyssal blue dye":        "Other",
	"Abyssal red dye":         "Other",

	// Zulrah
	"Tanzanite mutagen": "Zulrah",
	"Magma mutagen":     "Zulrah",
	"Uncut onyx":        "Zulrah",
	"Serpentine visage": "Zulrah",
	"Magic fang":        "Zulrah",
	"Tanzanite fang":    "Zulrah",

	// Gauntlet
	"Enhanced crystal weapon seed": "Gauntlet",
	"Crystal armour seed":          "Gauntlet",

	// Godwars
	"Armadyl helmet":            "Godwars",
	"Armadyl chestplate":        "Godwars",
	"Armadyl chainskirt":        "Godwars",
	"Armadyl hilt":              "Godwars",
	"Bandos chestplate":         "Godwars",
	"Bandos tassets":            "Godwars",
	"Bandos boots":              "Godwars",
	"Bandos hilt":               "Godwars",
	"Saradomin sword":           "Godwars",
	"Saradomin's light":         "Godwars",
	"Armadyl crossbow":          "Godwars",
	"Saradomin hilt":            "Godwars",
	"Steam battlestaff":         "Godwars",
	"Zamorakian spear":          "Godwars",
	"Staff of the dead":         "Godwars",
	"Zamorak hilt":              "Godwars",
	"Zaryte vambraces":          "Godwars",
	"Nihil horn":                "Godwars",
	"Torva full helm (damaged)": "Godwars",
	"Torva platebody (damaged)": "Godwars",
	"Torva platelegs (damaged)": "Godwars",
	"Torva full helm":           "Godwars",
	"Torva platebody":           "Godwars",
	"Torva platelegs":           "Godwars",
	"Ancient hilt":              "Godwars",

	// Nightmare
	"Nightmare staff":         "Nightmare",
	"Inquisitor's great helm": "Nightmare",
	"Inquisitor's hauberk":    "Nightmare",
	"Inquisitor's plateskirt": "Nightmare",
	"Inquisitor's mace":       "Nightmare",
	"Eldritch orb":            "Nightmare",
	"Harmonised orb":          "Nightmare",
	"Volatile orb":            "Nightmare",

	// Muspah
	"Venator shard": "Muspah",

	// DKS
	"Warrior ring":   "DKS",
	"Berserker ring": "DKS",
	"Dragon axe":     "DKS",
	"Seers ring":     "DKS",
	"Archers ring":   "DKS",

	// Zalcano
	"Zalcano shard":     "Zalcano",
	"Crystal tool seed": "Zalcano",

	// Jars
	"Jar of dirt":      "Jars",
	"Jar of swamp":     "Jars",
	"Jar of sand":      "Jars",
	"Jar of souls":     "Jars",
	"Jar of miasma":    "Jars",
	"Jar of darkness":  "Jars",
	"Jar of stone":     "Jars",
	"Jar of decay":     "Jars",
	"Jar of chemicals": "Jars",
	"Jar of eyes":      "Jars",
	"Jar of dreams":    "Jars",
	"Jar of spirits":   "Jars",
	"Jar of smoke":     "Jars",
}
