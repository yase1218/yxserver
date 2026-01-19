package condition

import "gameserver/internal/game/player"

func checkMissionPass(p *player.Player, args []uint32) ([]uint32, bool) {
	if len(args) == 0 {
		return nil, true
	}

	missionId := int(args[0])
	missions := p.UserData.Mission.Missions

	for i := 0; i < len(missions); i++ {
		if missions[i].MissionId == missionId {
			return nil, missions[i].IsPass
		}
	}

	return nil, false
}
