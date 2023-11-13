package main

import (
	"github.com/bwmarrin/discordgo"
)

func kickOffHallOfFameUpdate(session *discordgo.Session) {
	slayerBosses := hallOfFameRequestInfo{Bosses: map[string]string{"sire": "https://i.imgur.com/GhbmqEB.png", "hydra": "https://i.imgur.com/25GU0Ph.png", "cerberus": "https://i.imgur.com/UoxGuQi.png", "grotesqueguardians": "https://i.imgur.com/M7ylVBZ.png", "kraken": "https://i.imgur.com/Q6EbJb1.png", "smokedevil": "https://i.imgur.com/2AYntQ5.png"}, DiscChan: config.DiscSlayerBossesChan}
	gwd := hallOfFameRequestInfo{Bosses: map[string]string{"commanderzilyana": "https://i.imgur.com/aNm4Ydd.png", "kreearra": "https://i.imgur.com/lX8SfgN.png", "kriltsutsaroth": "https://i.imgur.com/hh8cMvp.png", "nex": "https://i.imgur.com/pqiVQBC.png", "generalgraardor": "https://i.imgur.com/hljv9ZW.png"}, DiscChan: config.DiscGwdChan}
	wildy := hallOfFameRequestInfo{Bosses: map[string]string{"artio": "https://i.imgur.com/bw6zLpU.png", "callisto": "https://i.imgur.com/bw6zLpU.png", "calvarion": "https://i.imgur.com/v3KX75y.png", "vetion": "https://i.imgur.com/v3KX75y.png", "spindel": "https://i.imgur.com/4zknWSX.png", "venenatis": "https://i.imgur.com/4zknWSX.png", "chaoselemental": "https://i.imgur.com/YAvIpbm.png", "chaosfanatic": "https://i.imgur.com/azV2sD1.png", "crazyarchaeologist": "https://i.imgur.com/23LXv53.png", "scorpia": "https://i.imgur.com/9aaguxB.png"}, DiscChan: config.DiscWildyChan}
	other := hallOfFameRequestInfo{Bosses: map[string]string{"corporealbeast": "https://i.imgur.com/zEDN4Pf.png", "prime": "https://i.imgur.com/kJBtqHB.png", "rexbro": "https://i.imgur.com/PvlGWFZ.png", "supreme": "https://i.imgur.com/BOgkBuD.png", "gauntlet": "https://i.imgur.com/weiHWnz.png", "gauntlethard": "https://i.imgur.com/xzW4TGR.png", "giantmole": "https://i.imgur.com/coKk2pr.gif", "jad": "https://i.imgur.com/H9aO1Ot.png", "zuk": "https://i.imgur.com/mKstHza.png", "kq": "https://i.imgur.com/ZuaFoBR.png", "kbd": "https://i.imgur.com/r5vkw1s.png", "sarachnis": "https://i.imgur.com/98THH8O.png", "skotizo": "https://i.imgur.com/YUcQu4d.png", "muspah": "https://i.imgur.com/sW2cLQ2.png", "vorkath": "https://i.imgur.com/6biF3P2.png", "phosanis": "https://i.imgur.com/4aDkxms.png", "nightmare": "https://i.imgur.com/4aDkxms.png", "zulrah": "https://i.imgur.com/tPllWNF.png"}, DiscChan: config.DiscOtherChan}
	misc := hallOfFameRequestInfo{Bosses: map[string]string{"barrows": "https://i.imgur.com/ajoK20v.png", "hespori": "https://i.imgur.com/b0qYGHS.png", "mimic": "https://i.imgur.com/jC7yTC3.png", "obor": "https://i.imgur.com/dwLvSbR.png", "bryophyta": "https://i.imgur.com/3cdyp4X.png", "derangedarchaeologist": "https://i.imgur.com/cnHpevF.png", "wintertodt": "https://i.imgur.com/6oFef2Y.png", "zalcano": "https://i.imgur.com/edN11Nf.png", "tempoross": "https://i.imgur.com/fRj3JA4.png", "rift": "https://i.imgur.com/MOiyXeH.png"}, DiscChan: config.DiscMiscChan}
	dt2 := hallOfFameRequestInfo{Bosses: map[string]string{"duke": "https://i.imgur.com/RYPmrXy.png", "leviathan": "https://i.imgur.com/mEQRq5c.png", "whisperer": "https://i.imgur.com/cFGWb6Y.png", "vardorvis": "https://i.imgur.com/WMPuShZ.png"}, DiscChan: config.DiscDT2Chan}
	raids := hallOfFameRequestInfo{Bosses: map[string]string{"cox": "https://i.imgur.com/gxdWXtH.png", "coxcm": "https://i.imgur.com/gxdWXtH.png", "tob": "https://i.imgur.com/pW1sJAQ.png", "tobcm": "https://i.imgur.com/pW1sJAQ.png", "toa": "https://i.imgur.com/2GvzqGw.png", "toae": "https://i.imgur.com/2GvzqGw.png"}, DiscChan: config.DiscRaidsChan}
	pvp := hallOfFameRequestInfo{Bosses: map[string]string{"bhh": "https://i.imgur.com/zSQhlWk.png", "bhr": "https://i.imgur.com/Y3Sga7t.png", "lms": "https://i.imgur.com/rzW7ZXx.png", "arena": "https://i.imgur.com/uNP6Ggu.png", "zeal": "https://i.imgur.com/Ws7HvKL.png"}, DiscChan: config.DiscPVPChan}
	clues := hallOfFameRequestInfo{Bosses: map[string]string{"clueall": "https://i.imgur.com/wX3Ei7U.png", "cluebeginner": "https://i.imgur.com/fUmzJkW.png", "clueeasy": "https://i.imgur.com/phnSCHj.png", "cluemedium": "https://i.imgur.com/t5iH8Xa.png", "cluehard": "https://i.imgur.com/a0xwcGI.png", "clueelite": "https://i.imgur.com/ibNRk3G.png", "cluemaster": "https://i.imgur.com/12rCLVv.png"}, DiscChan: config.DiscCluesChan}

	updateHallOfFame(session, slayerBosses)
	updateHallOfFame(session, gwd)
	updateHallOfFame(session, wildy)
	updateHallOfFame(session, other)
	updateHallOfFame(session, misc)
	updateHallOfFame(session, dt2)
	updateHallOfFame(session, raids)
	updateHallOfFame(session, pvp)
	updateHallOfFame(session, clues)
	updateCollectionLog(session)
}

func updateHallOfFame(session *discordgo.Session, requestInfo hallOfFameRequestInfo) {
	// First, delete all the messages within the channel
	messages, err := session.ChannelMessages(requestInfo.DiscChan, 50, "", "", "")
	if err != nil {
		return
	}

	var messageIDs []string
	for _, message := range messages {
		messageIDs = append(messageIDs, message.ID)
	}

	err = session.ChannelMessagesBulkDelete(requestInfo.DiscChan, messageIDs)
	if err != nil {
		return
	}

	// Now add all the bosses
	for bossIdForTemple, imageURL := range requestInfo.Bosses {
		getPodiumFromTemple(session, bossIdForTemple, requestInfo.DiscChan, imageURL)
	}
	return
}
