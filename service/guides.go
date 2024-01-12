package service

import (
	"context"
	"github.com/bwmarrin/discordgo"
)

func (s *Service) updateTobGuide(ctx context.Context, session *discordgo.Session) {
	mageSetup := []string{
		"# Mage Setup 5s & 4s",
		"https://i.imgur.com/KtpIQBH.png",
		"Runes: Blood, Soul, Water, Death\nSpellbook: Ancients\n\n- When entering Verzik, this setup will allow for 1 brew, 1 restore, and either 1 stam/shark/free inv for staff. \n- Because you are only using 1 restore for the entirety of Verzik, you can opt to flick piety during P1. Additionally, only putting on protect from ranged during nylo spawn on P2 will also help with prayer usage.\n- Heart is a requirement for the shadow to be useful. If you don't have a saturated heart, drop the magus ring and shadow for an Elidinis' ward (f) and Bellator ring.\n- Shadow uses:\n - Initial hit on Sotetseg\n - 0 hammers land on Sotetseg\n - P2 Verzik",
		"# Mfrz (4s)",
		"https://i.imgur.com/R1fwTy5.png",
		"Runes: Blood, Soul, Water, Death\nSpellbook: Ancients\n\n- When entering Verzik, this setup will allow for 4 inv for potions. This can be allocated to brew/restore/shark/stam or however you see fit\n- Ring Usages:\n - Maiden: Ultor\n - Bloat: Ultor\n - Nylo: Ultor\n - Sotetseg: 0-1 hammer Bellator, otherwise Ultor\n - Xarpus: Ultor\n - Verzik: P1 & P2 Bellator, P3 Ultor",
		"# Mage Trio",
		"https://i.imgur.com/cycZgde.png",
		"Runes: Blood, Soul, Water, Death\nSpellbook: Ancients\n\n- When specing down Maiden and Xarpus, everyone will do 1 1\n- When entering Verzik, this setup will allow for 2 brews, 2 restores, and either 1 stam/shark/free inv for staff\n- Ring Usages:\n - Maiden: Ultor\n - Bloat: Ultor\n - Nylo: Ultor\n - Sotetseg: 0-1 hammer Bellator (can use lightbearer for a phase if needed to backup), otherwise Ultor\n - Xarpus: Ultor\n - Verzik: P1 Lightbearer, P2 Bellator, P3 Ultor",
		"# Mage Duo",
		"https://i.imgur.com/EZmfQZD.png",
		"- The Virtus Robe Bottoms offer the same max hit when using trident on the nylos while also allowing a max hit increase on your barrages\n- Make sure to tbow in mage top/bottom & torm for better range bonus\n- Spellbook swap using the magic cape to arceuus after nylo (can drop after)\n- Dinny B spec during nylo on the south doubles wave 22 to ensure maximum usage\n- Ring Usages:\n - Maiden: Ultor\n - Bloat: Ultor\n - Nylo: Ultor\n - Sotetseg: Both you and your ranger will spec with Ultor on wait until you hit 60% spec. You will then camp lb until 80%. For the rest of the fight: 0-1 hammer Bellator, otherwise Ultor\n - Xarpus: Ultor\n - Verzik: P1 Lightbearer, P2 Bellator, P3 Ultor",
		"# BACK TO THE TOP: https://discord.com/channels/1172535371905646612/1184610126670336080/1184610126670336080",
	}

	rangeSetup := []string{
		"# Range Setup 5s & 4s",
		"https://i.imgur.com/JRs4rxx.png",
		"Runes: Law, Cosmic, Astral, Death\nSpellbook: Lunars\n\n- The runes allow for veng/sbs thralls\n- Since the ranger is not responsible for pot share, a 1 dose divine ranging pot is used for Nylo\n- Ring Usages:\n - Maiden: Ultor\n - Bloat: Ultor\n - Nylo: Ultor\n - Sotetseg: 0-1 hammer Bellator, otherwise Ultor\n - Xarpus: Ultor\n - Verzik: P1 & P2 Bellator, P3 Ultor",
		"# Range Setup Trio",
		"https://i.imgur.com/Ld44MgT.png",
		"Runes: Law, Cosmic, Astral, Death\nSpellbook: Lunars\n\n- The runes allow for veng/sbs thralls\n- Since the ranger is not responsible for pot share, a 1 dose divine ranging pot is used for Nylo\n- Ring Usages:\n - Maiden: Ultor\n - Bloat: Ultor\n - Nylo: Ultor\n - Sotetseg: 0-1 hammer Bellator (can use lightbearer for a phase if needed to backup), otherwise Ultor\n - Xarpus: Ultor\n - Verzik: P1 Lightbearer, P2 Bellator, P3 Ultor",
		"# Range Setup Duo",
		"https://i.imgur.com/PdbfjnO.png",
		"Runes: Lava, Blood, Cosmic, Death\nSpellbook: Lunars\n\n- The runes allow for veng/sbs thralls\n- There is no pot share in this either, so bring a divine ranging\n- You are the only BGS - here is the priority:\n - Maiden: BGS after two hammers hit - if your mager misses 1, you need to backup\n - Nylo: BGS the boss down to 0\n - Xarpus: Same as maiden\n- Share leftover divine sbc doses to mager after Sotetseg\n- Ring Usages:\n - Maiden: Ultor\n - Bloat: Ultor\n - Nylo: Ultor\n - Sotetseg: Both you and your mager will spec with Ultor on wait until you hit 60% spec. You will then camp lb until 80%. For the rest of the fight: 0-1 hammer Bellator, otherwise Ultor\n - Xarpus: Ultor\n - Verzik: P1 Lightbearer, P2 Bellator, P3 Ultor",
		"# BACK TO THE TOP: https://discord.com/channels/1172535371905646612/1184622959466389514/1184622959466389514",
	}

	meleeSetup := []string{
		"# Melee Setup 5s & 4s",
		"https://i.imgur.com/3nAocWw.png",
		"Runes: Law, Cosmic, Astral, Death\nSpellbook: Lunars\n\n- The runes allow for veng/pot share/sbs thralls\n- Pot share ranging potion at the end of nylo waves\n- The Dihn's is for stacking at maiden in 5s and can also be used during nylo waves during south wave 22 double melee crabs in 5s and 4s\n- Ring Usages:\n - Maiden: Ultor\n - Bloat: Ultor\n - Nylo: Ultor\n - Sotetseg: 0-1 hammer Bellator, otherwise Ultor\n - Xarpus: Ultor\n - Verzik: P1 & P2 Bellator, P3 Ultor",
		"# Melee Trio Setup",
		"https://i.imgur.com/3jgoepr.png",
		"Runes: Law, Cosmic, Astral, Death\nSpellbook: Lunars\n\n- The runes allow for veng/pot share/sbs thralls/sbs death charge\n- Pot share ranging potion at the end of nylo waves\n- Ring Usages:\n - Maiden: Ultor\n - Bloat: Ultor\n - Nylo: Ultor\n - Sotetseg: 0-1 hammer Bellator (can use lightbearer for a phase if needed to backup), otherwise Ultor\n - Xarpus: Ultor\n - Verzik: P1 Lightbearer, P2 Bellator, P3 Ultor",
		"# BACK TO THE TOP: https://discord.com/channels/1172535371905646612/1184829923479797820/1184829923479797820",
	}

	s.deleteAllMessages(session, "tob mage", s.config.DiscTobMageGuideChan, "1184610126670336080")
	s.deleteAllMessages(session, "tob range", s.config.DiscTobRangeGuideChan, "1184622959466389514")
	s.deleteAllMessages(session, "tob melee", s.config.DiscTobMeleeGuideChan, "1184829923479797820")

	s.log.Debug("Updating mage tob guide...")
	for _, line := range mageSetup {
		_, err := session.ChannelMessageSend(s.config.DiscTobMageGuideChan, line)
		if err != nil {
			s.log.Error("ERROR SENDING MESSAGES TO DISCORD TOB MAGE CHANNEL")
			return
		}
	}

	s.log.Debug("Updating range tob guide...")
	for _, line := range rangeSetup {
		_, err := session.ChannelMessageSend(s.config.DiscTobRangeGuideChan, line)
		if err != nil {
			s.log.Error("ERROR SENDING MESSAGES TO DISCORD TOB RANGE CHANNEL")
			return
		}
	}

	s.log.Debug("Updating melee tob guide...")
	for _, line := range meleeSetup {
		_, err := session.ChannelMessageSend(s.config.DiscTobMeleeGuideChan, line)
		if err != nil {
			s.log.Error("ERROR SENDING MESSAGES TO DISCORD TOB MELEE CHANNEL")
			return
		}
	}
}

func (s *Service) updateTrioCMGuide(ctx context.Context, session *discordgo.Session) {
	prepGuide := []string{
		"# Gear & Overview",
		"https://i.imgur.com/eLw2Ace.png",
		"Runes: Astral, Lava, Cosmic, Water\nSpellbook: Lunars\nNote: The ring of shadows is charged with 100 soul runes, 100 death runes, 100 blood runes, and 100 law runes\n\nOverview:\n- Pick up 2 ovl and 1 enh at Tekton (drop 10 death runes for space)\n- Get 8 mushroom, 18 juice & 2 planks at scavs\n- Kill Melee Vanguard\n- Die in vangs at ~50% hp\n- Make 8 rest, 18 brews\n- Cross rope if you're fast enough\n- SBS to barrage small muta\n- Melee hand at Olm",
		"# How To Spellbook Swap Without Stopping\nYou must red  click something and then spellbook swap. This ensures that you will not be stopped running while spellbook swapping unless you get interrupted by a hit from something (a scav). Try to red click doors that allow you to move forward.",
		"# Prep Walkthrough",
		"https://i.imgur.com/3MSlpw4.png",
		"<:cm_bp:1184588592136323144> Lure Tekton. \n<:cm_bp:1184588592136323144> It's important to give enough time for the other roles to hammer twice and spec as they are swinging their first scythe (https://imgur.com/a/f9hv4fK)\n<:cm_bp:1184588592136323144> Reset threshold: 2 hammers\n<:cm_bp:1184588592136323144> Use both your anglers after tekton to recover from taking two vengs\n<:cm_bp:1184588592136323144> Once dead, drop your death runes and pickup both ovls & 1 enh",
		"https://i.imgur.com/Gg6W33B.png",
		"<:cm_bp:1184588592136323144> You will be last into the room\n<:cm_bp:1184588592136323144> Run along the south side of the room and tag the remaining crabs, smashing all but 1 and bring it in position for the yellow crystal\n<:cm_bp:1184588592136323144> Shift to get the black crystal",
		"https://i.imgur.com/noDvfRC.png",
		"<:cm_bp:1184588592136323144> Drop the 2 ovls, enh, and salve near the tinderbox\n<:cm_bp:1184588592136323144> Scav for 8 mushroom, 18 juice & 2 planks\n<:cm_bp:1184588592136323144> Note: Tweaked drops so that whenever any potion secondary would be dropped, Scavengers now drop 3-11 Stinkhorn mushroom, 3-6 Cicely and 5-14 Endarkened Juice\n<:cm_bp:1184588592136323144> Immediately make the chest, go back for your dropped items prepare for the ice demon pop\n<:cm_bp:1184588592136323144> Bank the required stuff and uncharge the ring of shadows. Drop soul runes\n### BANKING STEPS: \n<:cm_bp:1184588592136323144> Private: BGS, Salve, Pickaxe, & Voidwaker\n<:cm_bp:1184588592136323144> Hammer spec ice demon once and then inbetween tbow hits, bank hammer & avernic\n\n<:cm_bp:1184588592136323144> Ovl and enh\n<:cm_bp:1184588592136323144> SBS Thralls",
		"https://i.imgur.com/z0smyQr.png",
		"<:cm_bp:1184588592136323144> SBS Thralls\n<:cm_bp:1184588592136323144> Don't allow the shaman to jump and ZCB them when they're close to full HP\n<:cm_bp:1184588592136323144> Grab only buchu\n<:cm_bp:1184588592136323144> Grab all tools, plant the buchu, drop buchu, rake, and seed dibber. \n<:cm_bp:1184588592136323144> SBS Thralls while running to vanguards",
		"https://i.imgur.com/w1mPBZ2.png",
		"<:cm_bp:1184588592136323144>Melee vanguard\n<:cm_bp:1184588592136323144> Die ~50% hp\n\n### BANKING STEPS: \n<:cm_bp:1184588592136323144> Private: Withdraw all\n<:cm_bp:1184588592136323144> Private: Deposit all\n<:cm_bp:1184588592136323144> Private: Withdraw rune pouch & spade\n\n<:cm_bp:1184588592136323144> Get 26 vials of water and put into shared\n<:cm_bp:1184588592136323144> Pick 26 buchu\n<:cm_bp:1184588592136323144> Make 8 restores and 18 brews. \n<:cm_bp:1184588592136323144> Make while running and pick up melee vanguard loot",
		"https://i.imgur.com/RmSpqNo.png",
		"<:cm_bp:1184588592136323144> Once you get to the thieving chest, follow the banking information\n\n### BANKING STEPS: \n<:cm_bp:1184588592136323144> Private: Withdraw all\n<:cm_bp:1184588592136323144> Private: Deposit BGS, Hammer, Empty ring of shadows\n<:cm_bp:1184588592136323144> If they are still doing vespula: \n    <:cm_bp:1184588592136323144> Public: Withdraw 2 restore, 1 ovl, 1 brew for use at vespula.\n<:cm_bp:1184588592136323144> If they are done with vespula:\n    <:cm_bp:1184588592136323144> Private: Deposit your mage switches except for staff, occult, and torm\n    <:cm_bp:1184588592136323144> Public: Withdraw 2 restore, 1 ovl, and fill up on brews\n\n<:cm_bp:1184588592136323144> When running to vespula, sip your ovl once, brew 3 times, and restore once to go back to 120 stats with 115 hp",
		"https://i.imgur.com/GMFRpAE.png",
		"<:cm_bp:1184588592136323144> If alive, it will be too low to use zcb, just use mage\n<:cm_bp:1184588592136323144> If you had to attack vespula:\n    <:cm_bp:1184588592136323144> Private: Deposit your mage switches except for staff, occult, and torm\n    <:cm_bp:1184588592136323144> Public: Withdraw 1 ovl, 2 rest, and fill up on brews",
		"https://i.imgur.com/mk8tsm2.png",
		"<:cm_bp:1184588592136323144> Haul ass and cross the rope praying mage with an empty inventory spot\n<:cm_bp:1184588592136323144> Drop the crystal and instantly click the rope to cross back\n<:cm_bp:1184588592136323144> DO NOT USE THE LAST CHEST BEFORE GOING DOWN",
		"https://i.imgur.com/8lesyx3.png",
		"<:cm_bp:1184588592136323144> You have time to brew before getting all the way to guardians so this is the time to brew up as much as you can\n<:cm_bp:1184588592136323144> Share the enh with everyone on the marked tiles\n<:cm_bp:1184588592136323144> Ensure everyone is venged before vasa\n<:cm_bp:1184588592136323144> SBS Thralls after guardians",
		"https://i.imgur.com/jdY6BUj.png",
		"<:cm_bp:1184588592136323144> Go in at 118 hp with zcb equipped\n<:cm_bp:1184588592136323144> Ensure you have 1 inventory spot so you can equip the tbow\n<:cm_bp:1184588592136323144> After vasa special, spec with the zcb and switch to tbow. Make sure to not cross the line that spawns muta\n<:cm_bp:1184588592136323144> If you have room and it's not too out of the way, pickup the brews but ignore the ovl\n<:cm_bp:1184588592136323144>SBS Thralls when running to mystics",
		"https://i.imgur.com/dvMYLTw.png",
		"<:cm_bp:1184588592136323144> Zcb spec a full hp mystic\n<:cm_bp:1184588592136323144> Ensure you aren't taking too much damage by using corner safespots\n<:cm_bp:1184588592136323144> When there are 1-2 mystics left, start making your way to the exit door\n<:cm_bp:1184588592136323144> Tendrils Skip: Once the last one is killed, spam click the symbol on the ground. Once you move, click through the door and run through into muta",
		"https://i.imgur.com/AQED5ek.png",
		"<:cm_bp:1184588592136323144> SBS Thralls then SBS ancients and wait until small muta is close to 50% hp then freeze\n<:cm_bp:1184588592136323144> Mage small muta\n<:cm_bp:1184588592136323144> Determine whether you or the surger will tank big muta",
		"https://i.imgur.com/R4Feexl.png",
		"<:cm_bp:1184588592136323144> Melee Hand\n<:cm_bp:1184588592136323144> SBS Thralls on the way to Olm\n<:cm_bp:1184588592136323144> SBS Thralls every transition phase and for Olm Head\n\n### BANKING STEPS: \n<:cm_bp:1184588592136323144> Shared: Deposit all your raids potions\n<:cm_bp:1184588592136323144> Private: Deposit your pickaxe, salve, voidwaker, BGS\n<:cm_bp:1184588592136323144> Private: Withdraw your mage switches, scythe and hammer\n<:cm_bp:1184588592136323144> Shared: Withdraw until you have 1 ovl, 2 restore, and fill up with brews",
		"# BACK TO THE TOP: https://discord.com/channels/1172535371905646612/1184590018967249006/1184590018967249006",
	}

	chinGuide := []string{
		"# Gear & Overview",
		"https://i.imgur.com/yjOcpcv.png",
		"Runes: Death, Water, Astral, Lava\nSpellbook: Lunars\n\nOverview:\n* Kill Ranged Vanguard\n* Tag deadly rangers at rope with chins\n* ZGS big muta\n* Runner at Olm",
		"# Chin Walkthrough",
		"https://i.imgur.com/3MSlpw4.png",
		"<:cm_bp:1184588592136323144> Wait for the prepper to lure tekton and spec with your dwh twice (https://imgur.com/a/f9hv4fK)\n<:cm_bp:1184588592136323144> Reset threshold: 2 hammers\n<:cm_bp:1184588592136323144> Run immediately to the crab room to start it - go far for the cyan crab",
		"https://i.imgur.com/Gg6W33B.png",
		"<:cm_bp:1184588592136323144> As second in the room, run along the north side of the room and tag at least 2 crabs while setting up for the Cyan crystal\n<:cm_bp:1184588592136323144> Smash 1 of the crabs in place and bring the other one for the yellow crystal\n<:cm_bp:1184588592136323144> Shift to get the black crystal",
		"https://i.imgur.com/noDvfRC.png",
		"<:cm_bp:1184588592136323144> Drop two potions by the tinderbox and immediately start chopping\n<:cm_bp:1184588592136323144> In between chopping ticks, pick up the tinderbox and start heading over towards the braziers\n<:cm_bp:1184588592136323144> 20+ kindling for the first unlit brazier and ~14 kindling for the second unlit brazier\n<:cm_bp:1184588592136323144> Drop tinderbox and axe when finished and bank using the following banking information\n### BANKING STEPS: \n<:cm_bp:1184588592136323144> Private: Deposit Salve, Pickaxe, ZGS, Chins\n\n<:cm_bp:1184588592136323144> Ovl and enh",
		"https://i.imgur.com/z0smyQr.png",
		"<:cm_bp:1184588592136323144> Don't allow the shaman to jump and ZCB them when they're close to full HP\n<:cm_bp:1184588592136323144> Ignore drops and go to vanguards",
		"https://i.imgur.com/w1mPBZ2.png",
		"<:cm_bp:1184588592136323144> First one down will go west, the other person will go north\n<:cm_bp:1184588592136323144> Ranged vanguard\n<:cm_bp:1184588592136323144> Pickup your vanguard loot while dropping your brought in potions if not used - This will allow one space for thieving",
		"https://i.imgur.com/RmSpqNo.png",
		"<:cm_bp:1184588592136323144> Thieve on the way to feed the scav at the end of the room. Continually feed until total grub count is equal to the following:",
		"https://i.imgur.com/emqZqfk.png",
		"<:cm_bp:1184588592136323144> Once all the grubs are deposited, drop your lockpick and go to the chest. Follow the banking procedure:\n### BANKING STEPS: \n<:cm_bp:1184588592136323144> Shared: Deposit all raids potions\n<:cm_bp:1184588592136323144> Private: Withdraw all\n<:cm_bp:1184588592136323144> Private: Deposit Scythe, Hammer\n<:cm_bp:1184588592136323144> Shared: Withdraw 2 restore, 1 ovl\n\n<:cm_bp:1184588592136323144> Equip your mage gear + staff and spam click door to go in on the first tick\n<:cm_bp:1184588592136323144> When running to vespula, sip your ovl once",
		"https://i.imgur.com/GMFRpAE.png",
		"<:cm_bp:1184588592136323144> No Poke Vespula: Spam click the door wait a tick https://streamable.com/xs87l\n<:cm_bp:1184588592136323144> Ensure your bolts are equiped and your zcb is in your switches\n<:cm_bp:1184588592136323144> Go in with mage + shadow equipped - hit it once on entry and another time from the safe spot before doing your switch for zcb\n<:cm_bp:1184588592136323144> On the way to rope, equip your melee gear. Follow the banking steps\n### BANKING STEPS: \n<:cm_bp:1184588592136323144> Private: Deposit your mage switches except staff, occult, and torm\n<:cm_bp:1184588592136323144> Shared: Withdraw until you have 1 ovl, 2 restore, and fill up with brews",
		"https://i.imgur.com/mk8tsm2.png",
		"<:cm_bp:1184588592136323144> Wait until the crosser is 1 tick from the rope and tag the rangers\n<:cm_bp:1184588592136323144> After tagging, immediately walk to the safe spot and wait for the crosser to drop the crystal before tagging the ranger again\n<:cm_bp:1184588592136323144> After the room is cleared, release your chins\n<:cm_bp:1184588592136323144> DO NOT USE THE LAST CHEST BEFORE GOING DOWN",
		"https://i.imgur.com/8lesyx3.png",
		"<:cm_bp:1184588592136323144> You have time to brew before getting all the way to guardians so this is the time to brew up as much as you can\n<:cm_bp:1184588592136323144> Share the enh with everyone on the marked tiles\n<:cm_bp:1184588592136323144> Ensure everyone is venged before vasa",
		"https://i.imgur.com/jdY6BUj.png",
		"<:cm_bp:1184588592136323144> Go in at 118 hp with zcb equipped\n<:cm_bp:1184588592136323144> Ensure you have 1 inventory spot so you can equip the tbow\n<:cm_bp:1184588592136323144> After vasa special, spec with the zcb and switch to tbow. Make sure to not cross the line that spawns muta\n<:cm_bp:1184588592136323144> If you have room and it's not too out of the way, pickup the brews but ignore the ovl",
		"https://i.imgur.com/dvMYLTw.png",
		"<:cm_bp:1184588592136323144> DO NOT ZCB SPEC - SAVE DOUBLE ZGS FOR BIG MUTA\n<:cm_bp:1184588592136323144> Ensure you aren't taking too much damage by using corner safespots\n<:cm_bp:1184588592136323144> When there are 1-2 mystics left, start making your way to the exit door\n<:cm_bp:1184588592136323144> Tendrils Skip: Once the last one is killed, spam click the symbol on the ground. Once you move, click through the door and run through into muta",
		"https://i.imgur.com/AQED5ek.png",
		"<:cm_bp:1184588592136323144> Mage small muta\n<:cm_bp:1184588592136323144> When big muta spawns, stand next to the tree and tbow big muta\n<:cm_bp:1184588592136323144> Launch a ZGS spec around 60% hp. Spec again if missed",
		"https://i.imgur.com/R4Feexl.png",
		"<:cm_bp:1184588592136323144> Head Runner\n\n### BANKING STEPS: \n<:cm_bp:1184588592136323144> Shared: Deposit all your raids potions\n<:cm_bp:1184588592136323144> Private: Deposit your pickaxe, salve, voidwaker, ZGS\n<:cm_bp:1184588592136323144> Private: Withdraw your mage switches, scythe and hammer\n<:cm_bp:1184588592136323144> Shared: Withdraw until you have 1 ovl, 2 restore, and fill up with brews",
		"# BACK TO THE TOP: https://discord.com/channels/1172535371905646612/1184590304934891600/1184590304934891600",
	}

	surgeGuide := []string{
		"# Gear & Overview",
		"https://i.imgur.com/7aVooRx.png",
		"Runes: Wrath, Law, Lava, Mist\nSpellbook: Normals\n\nOverview:\n*  Surge Ice Demon\n*  Kill Mage Vanguard\n*  Telegrab at rope\n*  Skip at Olm",
		"# Surge Walkthrough",
		"https://i.imgur.com/3MSlpw4.png",
		"<:cm_bp:1184588592136323144> Wait for the prepper to lure tekton and spec with your dwh twice (https://imgur.com/a/f9hv4fK)\n<:cm_bp:1184588592136323144> Reset threshold: 2 hammers\n<:cm_bp:1184588592136323144> Run immediately to the crab room to start it - go close for the magenta crab",
		"https://i.imgur.com/Gg6W33B.png",
		"<:cm_bp:1184588592136323144> Enter the room first\n<:cm_bp:1184588592136323144> Tag a crab and manually step in front of the orb to pop it and do the magenta crystal\n<:cm_bp:1184588592136323144> Position for the cyan crystal\n<:cm_bp:1184588592136323144> Position for the yellow crystal\n<:cm_bp:1184588592136323144> Shift to get the black crystal",
		"https://i.imgur.com/noDvfRC.png",
		"<:cm_bp:1184588592136323144> Drop two potions by the tinderbox and immediately start chopping\n<:cm_bp:1184588592136323144> In between chopping ticks, pick up the tinderbox and start heading over towards the braziers\n<:cm_bp:1184588592136323144> 20+ kindling for the first unlit brazier and ~14 kindling for the second unlit brazier\n<:cm_bp:1184588592136323144> Drop tinderbox and axe when finished and bank using the following banking information\n### BANKING STEPS: \n<:cm_bp:1184588592136323144> Private: Deposit your pickaxe, salve\n\n<:cm_bp:1184588592136323144> Ovl and enh\n<:cm_bp:1184588592136323144> Surge ice demon",
		"https://i.imgur.com/z0smyQr.png",
		"<:cm_bp:1184588592136323144> Don't allow the shaman to jump and ZCB them when they're close to full HP\n<:cm_bp:1184588592136323144> Ignore drops and go to vanguards",
		"https://i.imgur.com/w1mPBZ2.png",
		"<:cm_bp:1184588592136323144> First one down will go west, the other person will go north\n<:cm_bp:1184588592136323144> Mage vanguard\n<:cm_bp:1184588592136323144> Pickup your vanguard loot while dropping your brought in potions if not used - This will allow one space for thieving",
		"https://i.imgur.com/RmSpqNo.png",
		"<:cm_bp:1184588592136323144> Thieve on the way to feed the scav at the end of the room. Continually feed until total grub count is equal to the following:",
		"https://i.imgur.com/emqZqfk.png",
		"<:cm_bp:1184588592136323144> Once all the grubs are deposited, drop your lockpick and go to the chest. Follow the banking procedure:\n### BANKING STEPS: \n<:cm_bp:1184588592136323144> Shared: Deposit all raids potions\n<:cm_bp:1184588592136323144> Private: Withdraw all\n<:cm_bp:1184588592136323144> Private: Deposit Scythe, Harm, Tome, Hammer\n<:cm_bp:1184588592136323144> Shared: Withdraw 2 restore, 1 ovl, 1 brew\n\n<:cm_bp:1184588592136323144> Equip your mage gear + staff and spam click door to go in on the first tick\n<:cm_bp:1184588592136323144> When running to vespula, sip your ovl once, brew 3 times, and restore once to go back to 120 stats with 115 hp",
		"https://i.imgur.com/GMFRpAE.png",
		"<:cm_bp:1184588592136323144> No Poke Vespula: Spam click the door wait a tick https://streamable.com/xs87l\n<:cm_bp:1184588592136323144> Ensure your bolts are equiped and your zcb is in your switches\n<:cm_bp:1184588592136323144> Go in with mage + shadow equipped - hit it once on entry and another time from the safe spot before doing your switch for zcb\n<:cm_bp:1184588592136323144> On the way to rope, equip your melee gear. Follow the banking steps\n### BANKING STEPS: \n<:cm_bp:1184588592136323144> Private: Deposit your mage switches except staff, occult, and torm\n<:cm_bp:1184588592136323144> Shared: Withdraw until you have 1 ovl, 2 restore, and fill up with brews",
		"https://i.imgur.com/mk8tsm2.png",
		"<:cm_bp:1184588592136323144> Tag one of the magers and ensure you have an empty inv for the crystal\n<:cm_bp:1184588592136323144> Telegrab the crystal when dropped and open the gate\n<:cm_bp:1184588592136323144> DO NOT USE THE LAST CHEST BEFORE GOING DOWN",
		"https://i.imgur.com/8lesyx3.png",
		"<:cm_bp:1184588592136323144> You have time to brew before getting all the way to guardians so this is the time to brew up as much as you can\n<:cm_bp:1184588592136323144> Share the enh with everyone on the marked tiles",
		"https://i.imgur.com/jdY6BUj.png",
		"<:cm_bp:1184588592136323144> Go in at 118 hp with zcb equipped\n<:cm_bp:1184588592136323144> Ensure you have 1 inventory spot so you can equip the tbow\n<:cm_bp:1184588592136323144> After vasa special, spec with the zcb and switch to tbow. Make sure to not cross the line that spawns muta\n<:cm_bp:1184588592136323144> If you have room and it's not too out of the way, pickup the brews but ignore the ovl",
		"https://i.imgur.com/dvMYLTw.png",
		"<:cm_bp:1184588592136323144> Zcb spec a full hp mystic\n<:cm_bp:1184588592136323144> Ensure you aren't taking too much damage by using corner safespots\n<:cm_bp:1184588592136323144> When there are 1-2 mystics left, start making your way to the exit door\n<:cm_bp:1184588592136323144> Tendrils Skip: Once the last one is killed, spam click the symbol on the ground. Once you move, click through the door and run through into muta",
		"https://i.imgur.com/AQED5ek.png",
		"<:cm_bp:1184588592136323144> Mage small muta\n<:cm_bp:1184588592136323144> Determine whether you or the surger will tank big muta",
		"https://i.imgur.com/R4Feexl.png",
		"<:cm_bp:1184588592136323144> Skipper\n\n### BANKING STEPS: \n<:cm_bp:1184588592136323144> Shared: Deposit all your raids potions\n<:cm_bp:1184588592136323144> Private: Deposit your pickaxe, salve, voidwaker\n<:cm_bp:1184588592136323144> Private: Withdraw your mage switches, scythe and hammer\n<:cm_bp:1184588592136323144> Shared: Withdraw until you have 1 ovl, 2 restore, and fill up with brews",
		"# BACK TO THE TOP: https://discord.com/channels/1172535371905646612/1184590344860487791/1184590344860487791",
	}

	usefulInfo := []string{
		"# Bank Tag Layouts Information\nbanktaglayoutsplugin:cm,28254:0,24780:1,28256:2,28258:3,13239:4,11212:5,24419:6,24420:7,24133:8,28035:9,27553:10,22981:11,25975:12,28328:13,25912:14,24421:15,20997:16,27238:17,27365:18,27241:19,27275:20,24666:21,24234:22,24668:23,22249:24,27235:25,26235:26,10476:27,19720:28,24664:29,23444:30,25739:31,26374:32,13441:33,21944:34,25818:35,27690:36,23681:37,23718:38,193:39,13441:40,28327:41,560:42,11804:43,12018:44,27281:45,13441:46,13441:47,16687:48,19266:49,564:50,13700:51,12919:54,12631:55,565:56,566:57,563:58,560:59,20997:64,27238:65,27365:66,27241:67,24423:68,24666:69,24234:70,24668:71,22249:72,27235:73,26235:74,10476:75,19720:76,24664:77,23444:78,20714:79,26374:80,3024:81,21944:82,1353:83,27690:84,23681:85,23718:86,193:87,6685:88,27275:89,1523:90,25739:91,12018:92,27281:93,13441:94,13441:95,21880:96,563:97,4695:98,19266:99,12919:102,12631:103,20997:112,27238:113,27365:114,27241:115,27275:116,24666:117,24234:118,24668:119,22249:120,27235:121,26235:122,10476:123,19720:124,24664:125,23444:126,25739:127,26374:128,3024:129,21944:130,1353:131,27690:132,23681:133,23718:134,193:135,6685:136,10033:137,1523:138,11808:139,12018:140,27281:141,13441:142,13441:143,560:144,16687:145,19266:146,13700:147,12919:150,12631:151,banktag:cm,22386,21944,27241,27238,27235,3024,27281,27275,23718,10476,24234,1353,27365,20714,19720,11808,11804,193,28327,1523,22249,12018,12919,21880,23680,13441,24133,27690,26374,563,560,555,12631,25912,24664,25975,10033,24666,24668,4695,4699,28258,28256,28254,22981,566,564,565,24780,23444,13239,24420,24421,24419,25739,24423,9075,26235,27553,20997,6685,28035,25818,11212",
		"## Top two rows",
		"https://i.imgur.com/DclvjUM.png",
		"* These are equipped items - except the ring of shadows.\n* The Inquisitor's set is for the surge role\n* Use slayer helm instead of masori helm if you're on a shamans task",
		"## Rows 3-8 (Prep)",
		"https://i.imgur.com/TatlDJS.png",
		"* The top two rows are the inventory setup for the prep role\n* The runes on row 7 are the runes inside the rune pouch\n* The runes on row 8 are the runes inside the ring of shadows (fill with 100 charges before the raid starts)\n* The 4 1 dose potions and two anglers on the right are for pre-pots\n* Ensure you only take in 10 death runes in your inventory setup",
		"## Rows 9-13 (Surge)",
		"https://i.imgur.com/MVLwdk9.png",
		"* The top two rows are the inventory setup for the surge role\n* The runes on row 13 are the runes inside the rune pouch\n* The 4 1 dose potions and two anglers on the right are for pre-pots",
		"## Rows 15-19 (Chin)",
		"https://i.imgur.com/3QmsFLj.png",
		"* The top two rows are the inventory setup for the chin role\n* The runes on row 19 are the runes inside the rune pouch\n* The 4 1 dose potions and two anglers on the right are for pre-pots\n* Ensure you only take in 5 chins in your inventory setup",
		"# Tile Markers for CMs",
		"https://pastebin.com/RqvZqtu8",
	}

	s.deleteAllMessages(session, "cm prep", s.config.DiscTrioCMPrepGuideChan, "1184590018967249006")
	s.deleteAllMessages(session, "cm chin", s.config.DiscTrioCMChinGuideChan, "1184590304934891600")
	s.deleteAllMessages(session, "cm surge", s.config.DiscTrioCMSurgeGuideChan, "1184590344860487791")
	s.deleteAllMessages(session, "cm useful", s.config.DiscTrioCMUsefulInfoChan, "1183782621457694762")

	s.log.Debug("Updating trio cm prep guide...")
	for _, line := range prepGuide {
		_, err := session.ChannelMessageSend(s.config.DiscTrioCMPrepGuideChan, line)
		if err != nil {
			s.log.Error("ERROR SENDING MESSAGES TO DISCORD CM TRIO PREP CHANNEL")
			return
		}
	}
	s.log.Debug("Updating trio cm chin guide...")
	for _, line := range chinGuide {
		_, err := session.ChannelMessageSend(s.config.DiscTrioCMChinGuideChan, line)
		if err != nil {
			s.log.Error("ERROR SENDING MESSAGES TO DISCORD CM TRIO CHIN CHANNEL")
			return
		}
	}
	s.log.Debug("Updating trio cm surge guide...")
	for _, line := range surgeGuide {
		_, err := session.ChannelMessageSend(s.config.DiscTrioCMSurgeGuideChan, line)
		if err != nil {
			s.log.Error("ERROR SENDING MESSAGES TO DISCORD CM TRIO SURGE CHANNEL")
			return
		}
	}
	s.log.Debug("Updating trio cm useful info...")
	for _, line := range usefulInfo {
		_, err := session.ChannelMessageSend(s.config.DiscTrioCMUsefulInfoChan, line)
		if err != nil {
			s.log.Error("ERROR SENDING MESSAGES TO DISCORD CM TRIO USEFUL INFO CHANNEL")
			return
		}
	}
}

func (s *Service) deleteAllMessages(session *discordgo.Session, guide string, channel string, initialImage string) {
	// First, delete all the messages within the channel
	messages, err := session.ChannelMessages(channel, 100, "", initialImage, "")
	if err != nil {
		s.log.Error("Failed to get all messages for deletion from the " + guide + " channel")
		return
	}
	var messageIDs []string
	for _, message := range messages {
		messageIDs = append(messageIDs, message.ID)
	}
	err = session.ChannelMessagesBulkDelete(channel, messageIDs)
	if err != nil {
		s.log.Error("Failed to delete all messages for deletion from the " + guide + " channel. Will try deleting one by one...")
		for _, message := range messageIDs {
			err = session.ChannelMessageDelete(channel, message)
			if err != nil {
				s.log.Error("Failed to delete messages one by one for deletion from the " + guide + " channel...")
				return
			}
		}
	}

	return
}
