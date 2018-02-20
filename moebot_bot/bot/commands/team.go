package commands

import (
	"github.com/camd67/moebot/moebot_bot/util/db"
)

type TeamCommand struct{}

func (tc *TeamCommand) Execute(pack *CommPackage) {
	processGuildRole([]string{"Nanachi", "Ozen", "Bondrewd", "Reg", "Riko", "Maruruk"}, pack.session, pack.params, pack.channel, pack.guild, pack.message, false)
}

func (tc *TeamCommand) GetPermLevel() db.Permission {
	return db.PermAll
}

func (sc *TeamCommand) GetCommandKeys() []string {
	return []string{"TEAM"}
}