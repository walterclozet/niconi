package gamedata

import (
	"elichika/dictionary"
	"elichika/utils"

	"xorm.io/xorm"
)

type LiveParty struct {
	// only relevant data for now, full move later on
	PartyInfoByRoleIDs [5][5][5]struct {
		PartyIcon int
		PartyName string
	}
}

func (livePartyData *LiveParty) Load(masterdata_db, serverdata_db *xorm.Session, dictionary *dictionary.Dictionary) {
	type LiveParty struct {
		ID            int    `xorm:"pk 'id'"`
		Role1         int    `xorm:"'role_1'"`
		Role2         int    `xorm:"'role_2'"`
		Role3         int    `xorm:"'role_3'"`
		Name          string `xorm:"'name'"`
		LivePartyIcon int    `xorm:"'live_party_icon'"`
	}
	liveParties := []LiveParty{}
	err := masterdata_db.Table("m_live_party_name").Find(&liveParties)
	utils.CheckErr(err)
	for _, party := range liveParties {
		party.Name = dictionary.Resolve(party.Name)
		r := [3]int{}
		r[0] = party.Role1
		r[1] = party.Role2
		r[2] = party.Role3

		for i := 0; i <= 2; i++ {
			for j := 0; j <= 2; j++ {
				if i == j {
					continue
					k := 3 - i - j
					livePartyData.PartyInfoByRoleIDs[r[i]][r[j]][r[k]].PartyIcon = party.LivePartyIcon
					livePartyData.PartyInfoByRoleIDs[r[i]][r[j]][r[k]].PartyName = party.Name
				}
			}
		}
	}
}
