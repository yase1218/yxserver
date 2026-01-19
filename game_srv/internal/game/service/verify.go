package service

import (
	"fmt"
	"gameserver/internal/game/model"
	"gameserver/internal/game/player"
	"kernel/tools"
	"msg"
	"time"

	"github.com/zy/game_data/template"
)

func VerifyLogin(p *player.Player, now time.Time) error {
	if p == nil {
		return fmt.Errorf("verify nil player")
	}

	if p.UserData == nil {
		return fmt.Errorf("verify nil user data")
	}

	u := p.UserData

	verify_base_info(u, now)
	verify_stage_info(u, now)
	verify_items(u, now)
	verify_task(u, now)
	verify_mission(u, now)
	verify_ships(u, now)
	verify_team(u, now)
	verify_equip(u, now)
	verify_shop(u, now)
	verify_playmethod(u, now)
	verify_weapon(u, now)
	verify_treasure(u, now)
	verify_poker(u, now)
	verify_activity(u, now)
	verify_cardpool(u, now)
	verify_appearance(u, now)
	verify_pet_data(u, now)
	verify_friend(u, now)
	verify_fight(u, now)
	verify_peak_fight(u, now)
	verify_contract(u, now)
	verify_desert(u, now)
	verify_arena(u, now)
	verify_atlas(u, now)
	verify_luck_sale(u, now)
	verify_function_preview(u, now)
	verify_equip_stage(u, now)
	verify_mails(u, now)
	verify_resources_pass(u, now)
	verify_likes(u, now)
	verify_ranks(u, now)
	verify_weekpass(u, now)
	verify_personalized(u, now)
	return nil
}

// func verify_peak_fight(u *model.UserData, now time.Time) {
// 	if u.PeakFight == nil {
// 		u.PeakFight = &model.PeakFight{}
// 	}
// }

func verify_mails(u *model.UserData, now time.Time) {
	if u.MailData == nil {
		u.MailData = &model.UserMail{
			Mails: make(map[int64]*model.Mail),
		}
	}
}

func verify_base_info(u *model.UserData, now time.Time) {
	if u.BaseInfo == nil {
		u.BaseInfo = &model.UserBase{}
	}
	if u.BaseInfo.ClientSettings == nil {
		u.BaseInfo.ClientSettings = make(map[uint32]string)
	}
	if u.BaseInfo.ApData == nil {
		u.BaseInfo.ApData = &model.ApInfo{
			RecoverStartTime: uint32(now.Unix()),
		}
	}
	if u.BaseInfo.HookData == nil {
		u.BaseInfo.HookData = &model.OnHookData{
			StartTime: uint32(now.Unix()),
		}
	}

	if u.BaseInfo.HookData.Items == nil {
		u.BaseInfo.HookData.Items = make([]*model.FloatItem, 0)
	}

	if u.BaseInfo.QuickOnHookData == nil {
		u.BaseInfo.QuickOnHookData = &model.QuickOnHookData{}
	}

	if u.BaseInfo.QuickOnHookData.Items == nil {
		u.BaseInfo.QuickOnHookData.Items = make([]*model.FloatItem, 0)
	}

	if u.BaseInfo.SupportId == nil {
		u.BaseInfo.SupportId = make([]uint32, 0)
	}

	if u.BaseInfo.GuideData == nil {
		u.BaseInfo.GuideData = make([]*model.GuideInfo, 0)
	}
	if u.BaseInfo.PopUps == nil {
		u.BaseInfo.PopUps = make([]*model.PopUpInfo, 0)
	}
	if u.BaseInfo.MissData == nil {
		u.BaseInfo.MissData = &model.MissionData{}
	}

	if u.BaseInfo.MissData.Ads == nil {
		u.BaseInfo.MissData.Ads = make([]*model.AdInfo, 0)
	}

	if u.BaseInfo.BattleData == nil {
		u.BaseInfo.BattleData = &msg.BattleData{}
	}
	if u.BaseInfo.BattleData.SkillData == nil {
		u.BaseInfo.BattleData.SkillData = make([]uint32, 0)
	}
	if u.BaseInfo.BattleData.PokerData == nil {
		u.BaseInfo.BattleData.PokerData = make([]*msg.BattlePokerData, 0)
	}
	if u.BaseInfo.BattleData.PartsId == nil {
		u.BaseInfo.BattleData.PartsId = make([]uint32, 0)
	}
	if u.BaseInfo.BattleData.PokerAddUp == nil {
		u.BaseInfo.BattleData.PokerAddUp = make([]*msg.BattlePokerAddUp, 0)
	}
	if u.BaseInfo.BattleData.PartsShop == nil {
		u.BaseInfo.BattleData.PartsShop = make([]*msg.BattlePartsShop, 0)
	}
	if u.BaseInfo.BattleData.InteractiveData == nil {
		u.BaseInfo.BattleData.InteractiveData = make([]*msg.BattleInteractiveData, 0)
	}
	if u.BaseInfo.BattleData.SummonedData == nil {
		u.BaseInfo.BattleData.SummonedData = make([]uint32, 0)
	}
	if u.BaseInfo.BattleData.StageEvent == nil {
		u.BaseInfo.BattleData.StageEvent = make([]uint32, 0)
	}
	if u.BaseInfo.BattleData.SettleData == nil {
		u.BaseInfo.BattleData.SettleData = make([]*msg.LittleSettleData, 0)
	}

	if u.BaseInfo.QuestionIds == nil {
		u.BaseInfo.QuestionIds = make([]string, 0)
	}

	if u.BaseInfo.RewardQuestionIds == nil {
		u.BaseInfo.RewardQuestionIds = make([]string, 0)
	}

	if u.BaseInfo.Charge == nil {
		u.BaseInfo.Charge = make([]*model.ChargeInfo, 0)
	}

	if u.BaseInfo.MonthCard == nil {
		u.BaseInfo.MonthCard = make([]*model.MonthcardInfo, 0)
	}

	if u.BaseInfo.MainFund == nil {
		u.BaseInfo.MainFund = make([]*model.MainFundInfo, 0)
	}

	if u.BaseInfo.Ad == nil {
		u.BaseInfo.Ad = make([]*model.AdInfo, 0)
	}

	if u.BaseInfo.Attrs == nil {
		u.BaseInfo.Attrs = make(map[uint32]*model.Attr)
	}

	if u.BaseInfo.RankLikesMap == nil {
		u.BaseInfo.RankLikesMap = make(map[uint32]uint32)
	}

	if u.BaseInfo.DailyApData == nil {
		u.BaseInfo.DailyApData = make([]*model.DailyApInfo, 0)
	}

	if u.BaseInfo.TalentData == nil {
		u.BaseInfo.TalentData = &model.TalentInfo{}
	}

	if u.BaseInfo.TalentData.Attrs == nil {
		u.BaseInfo.TalentData.Attrs = make(map[uint32]*model.Attr)
	}

	if u.BaseInfo.TalentData.Parts == nil {
		u.BaseInfo.TalentData.Parts = make([]uint32, 0)
	}

	if u.BaseInfo.ComboSkill == nil {
		u.BaseInfo.ComboSkill = make([]uint32, 0)
	}

	if u.BaseInfo.Adventures == nil {
		u.BaseInfo.Adventures = make(map[uint32]*model.AdventureInfo)
	}

	if u.BaseInfo.RankMissionReward == nil {
		u.BaseInfo.RankMissionReward = make([]int, 0)
	}

	if u.BaseInfo.FirstChargePackage == nil {
		u.BaseInfo.FirstChargePackage = make([]*model.FirstChargePackageData, 0)
	}
}

func verify_stage_info(u *model.UserData, now time.Time) {
	if u.StageInfo == nil {
		u.StageInfo = &model.UserStage{}
	}

	if u.StageInfo.StageFirstEnter == nil {
		u.StageInfo.StageFirstEnter = make([]int, 0)
	}

	if u.StageInfo.StageFirstPass == nil {
		u.StageInfo.StageFirstPass = make([]int, 0)
	}
}

func verify_items(u *model.UserData, now time.Time) {
	if u.Items == nil {
		u.Items = &model.UserItems{}
	}

	if u.Items.Items == nil {
		u.Items.Items = make([]*model.Item, 0)
	}
}

func verify_task(u *model.UserData, now time.Time) {
	if u.Task == nil {
		u.Task = &model.AccountTask{}
	}

	//if u.Task == nil {
	//	u.Task = &model.AccountTask{}
	//}

	if u.Task.DailyTasks == nil {
		u.Task.DailyTasks = make([]*model.Task, 0)
	}

	if u.Task.WeeklyTasks == nil {
		u.Task.WeeklyTasks = make([]*model.Task, 0)
	}

	if u.Task.AchieveTasks == nil {
		u.Task.AchieveTasks = make([]*model.Task, 0)
	}

	if u.Task.MainTasks == nil {
		u.Task.MainTasks = make([]*model.Task, 0)
	}

	if u.Task.FinshedTasks == nil {
		u.Task.FinshedTasks = make([]uint32, 0)
	}

	if u.Task.HistoryData == nil {
		u.Task.HistoryData = make([]*model.TaskHistroyData, 0)
	}
}

func verify_mission(u *model.UserData, now time.Time) {
	if u.Mission == nil {
		u.Mission = &model.AccountMission{}
	}

	if u.Mission.Missions == nil {
		u.Mission.Missions = make([]*model.Mission, 0)
	}

	if u.Mission.Challenges == nil {
		u.Mission.Challenges = make([]*model.Mission, 0)
	}

	if u.Mission.ExtraIds == nil {
		u.Mission.ExtraIds = make([]int, 0)
	}
}

func verify_ships(u *model.UserData, now time.Time) {
	if u.Ships == nil {
		u.Ships = &model.UserShips{}
	}

	if u.Ships.Ships == nil {
		u.Ships.Ships = make([]*model.Ship, 0)
	}
	for _, ship := range u.Ships.Ships {
		if ship.CoatMap == nil {
			ship.CoatMap = make(map[int]*model.CoatItem)
			coatCfg := template.GetShipTemplate().GetShip(ship.Id).GetOriginalCoat()
			ship.CoatMap[coatCfg.CoatId] = model.NewCoatItem(coatCfg)
		}
	}
}

func verify_team(u *model.UserData, now time.Time) {
	if u.Team == nil {
		u.Team = &model.AccountTeam{}
	}

	if u.Team.TeamData == nil {
		u.Team.TeamData = make([]*model.Team, 0)
	}

	if u.Team.BattleData == nil {
		u.Team.BattleData = make([]*model.BattleTeam, 0)
	}
}

func verify_equip(u *model.UserData, now time.Time) {
	if u.Equip == nil {
		u.Equip = &model.AccountEquip{}
	}

	if u.Equip.EquipData == nil {
		u.Equip.EquipData = make([]*model.Equip, 0)
	}

	if u.Equip.EquipPosData == nil {
		u.Equip.EquipPosData = make([]*model.EquipPos, 0)
	}

	if u.Equip.SuitReward == nil {
		u.Equip.SuitReward = make([]*model.SuitInfo, 0)
	}

	if u.Equip.EquipSuits == nil {
		u.Equip.EquipSuits = make([]*model.EquipSuit, 0)
	}

	if u.Equip.GemBag == nil {
		u.Equip.GemBag = make(map[uint64]*model.GemBagSlot)
	}

	if u.Equip.GemPos == nil {
		u.Equip.GemPos = make([][]uint64, msg.EquipPos_EquipPos_Max-1)
	}
}

func verify_shop(u *model.UserData, now time.Time) {
	if u.Shop == nil {
		u.Shop = &model.AccountShop{}
	}

	if u.Shop.Items == nil {
		u.Shop.Items = make([]*model.ShopItem, 0)
	}
}

func verify_playmethod(u *model.UserData, now time.Time) {
	if u.PlayMethod == nil {
		u.PlayMethod = &model.AccountPlayMethod{}
	}

	if u.PlayMethod.Data == nil {
		u.PlayMethod.Data = make([]*model.PlayMethodData, 0)
	}
}

func verify_weapon(u *model.UserData, now time.Time) {
	if u.Weapon == nil {
		u.Weapon = &model.AccountWeapon{}
	}

	if u.Weapon.Attrs == nil {
		u.Weapon.Attrs = make(map[uint32]*model.Attr)
	}

	if u.Weapon.Weapons == nil {
		u.Weapon.Weapons = make([]*model.Weapon, 0)
	}

	if u.Weapon.SecondaryWeapons == nil {
		u.Weapon.SecondaryWeapons = make([]*model.SecondaryWeapon, 0)
	}
}

func verify_treasure(u *model.UserData, now time.Time) {
	if u.Treasure == nil {
		u.Treasure = &model.AccountTreasure{}
	}

	if u.Treasure.ShipData == nil {
		u.Treasure.ShipData = make([]*model.ShipTreasure, 0)
	}

	if u.Treasure.MissData == nil {
		u.Treasure.MissData = make([]uint32, 0)
	}

	if u.Treasure.CommData == nil {
		u.Treasure.CommData = make([]uint32, 0)
	}

	if u.Treasure.WeaponData == nil {
		u.Treasure.WeaponData = make([]*model.WeaponTreasure, 0)
	}
}

func verify_poker(u *model.UserData, now time.Time) {
	if u.Poker == nil {
		u.Poker = &model.AccountPoker{}
	}

	if u.Poker.ShipData == nil {
		u.Poker.ShipData = make([]*model.ShipPoker, 0)
		u.Poker.MissData = make([]int, 0)
		u.Poker.CommData = make([]uint32, 0)
		u.Poker.WeaponData = make([]*model.WeaponPoker, 0)
	}
}

func verify_activity(u *model.UserData, now time.Time) {
	if u.AccountActivity == nil {
		u.AccountActivity = &model.AccountActivity{}
	}

	if u.AccountActivity.HisData == nil {
		u.AccountActivity.HisData = make([]uint32, 0)
	}

	if u.AccountActivity.Activities == nil {
		u.AccountActivity.Activities = make([]*model.Activity, 0)
	}

	if u.AccountActivity.PreRewardTps == nil {
		u.AccountActivity.PreRewardTps = make([]uint32, 0)
	}
}

func verify_cardpool(u *model.UserData, now time.Time) {
	if u.CardPool == nil {
		u.CardPool = &model.AccountCardPool{}
	}

	if u.CardPool.CardPools == nil {
		u.CardPool.CardPools = make([]*model.CardPool, 0)
	} else {
		cardList := template.GetLotteryShipTemplate().GetInitPool()
		if len(u.CardPool.CardPools) != len(cardList) {
			pools := u.CardPool.CardPools

			for i := range cardList {
				isFind := false
				for _, v := range pools {
					if v.CardId == cardList[i].Id {
						isFind = true
						break
					}
				}

				if !isFind {
					var freeTimes uint32 = 0
					var nextResetTime uint32 = 0
					if cardList[i].XDayFreeTime > 0 {
						freeTimes = 1
					}

					if cardList[i].CardType != 2 {
						nextResetTime = tools.GetDailyRefreshTime()
					}
					pools = append(pools,
						model.NewCardPool(cardList[i].Id, freeTimes,
							tools.GetDailyXRefreshTime(cardList[i].XDayFreeTime, template.GetSystemItemTemplate().RefreshHour),
							nextResetTime, 0, 0, 0))
				}
			}

			u.CardPool.CardPools = pools
		}
	}
}

func verify_appearance(u *model.UserData, now time.Time) {
	if u.Appearance == nil {
		u.Appearance = &model.AccountAppearance{}
	}

	if u.Appearance.Appearances == nil {
		u.Appearance.Appearances = make([]*model.Appearance, 0)
	}

	if u.Appearance.Attrs == nil {
		u.Appearance.Attrs = make(map[uint32]*model.Attr)
	}
}

func verify_pet_data(u *model.UserData, now time.Time) {
	if u.PetData == nil {
		u.PetData = model.NewAccountPet()
	}

	// if u.PetData.PetMaterials == nil {
	// 	u.PetData.PetMaterials = make([]*model.PetMaterial, 0)
	// }

	// if u.PetData.PetEggs == nil {
	// 	u.PetData.PetEggs = make([]*model.PetEgg, 0)
	// }

	// if u.PetData.AdRewardItems == nil {
	// 	u.PetData.AdRewardItems = make([]*model.SimpleItem, 0)
	// }

	// if u.PetData.PetParts == nil {
	// 	u.PetData.PetParts = make([]uint32, 0)
	// }

	// if u.PetData.Attrs == nil {
	// 	u.PetData.Attrs = make(map[uint32]*model.Attr)
	// }

	if u.PetData.Pets == nil {
		u.PetData.Pets = make(map[uint32]*model.Pet, 0)
	}
}

func verify_friend(u *model.UserData, now time.Time) {
	if u.FriendData == nil {
		u.FriendData = &model.UserFriend{}
	}

	if u.FriendData.Friend == nil {
		u.FriendData.Friend = make(map[uint64]struct{})
	}

	if u.FriendData.Black == nil {
		u.FriendData.Black = make(map[uint64]struct{})
	}
}

func verify_fight(u *model.UserData, now time.Time) {
	if u.Fight == nil {
		u.Fight = &model.Fight{}
	}

	if u.Fight.Weapons == nil {
		u.Fight.Weapons = make([]uint32, 0)
	}
}

func verify_peak_fight(u *model.UserData, now time.Time) {
	if u.PeakFight == nil {
		season := uint32(0)
		season_t := template.GetBattlePassTemplate().GetCurSeason(uint32(now.Unix()))
		if season_t != nil {
			season = season_t.Season
		} else {
			season = 1
		}
		u.PeakFight = &model.PeakFight{
			BattleMatchId: 1,
			Season:        season,
		}
	}
	if u.PeakFight.BattleMatchId == 0 {
		u.PeakFight.BattleMatchId = 1
	}
}

func verify_contract(u *model.UserData, now time.Time) {
	if u.Contract == nil {
		u.Contract = &model.Contract{}
	}

	if u.Contract.TaskIds == nil {
		u.Contract.TaskIds = make([]uint32, 0)
	}
}

func verify_desert(u *model.UserData, now time.Time) {
	if u.Desert == nil {
		u.Desert = &model.DesertFight{}
	}

	if u.Desert.RewardTimes == nil {
		u.Desert.RewardTimes = make([]uint32, 0)
	}
}

func verify_arena(u *model.UserData, now time.Time) {
	if u.Arena == nil {
		u.Arena = &model.ArenaPlayerData{}
	}

	if u.Arena.TodayUseShips == nil {
		u.Arena.TodayUseShips = make([]uint32, 0)
	}

	if u.Arena.UnlockMonsterIds == nil {
		u.Arena.UnlockMonsterIds = make([]uint32, 0)
	}

	if u.Arena.DefendMonster == nil {
		u.Arena.DefendMonster = make([]int32, 0)
	}

	if u.Arena.Records == nil {
		u.Arena.Records = make([]*model.ArenaPlayerPkRecordData, 0)
	}
}

func verify_atlas(u *model.UserData, now time.Time) {
	if u.Atlas == nil {
		u.Atlas = &model.Atlas{}
	}

	if u.Atlas.Data == nil {
		u.Atlas.Data = make(map[uint32]*model.AtlasUnit)
	}
}

func verify_luck_sale(u *model.UserData, now time.Time) {
	if u.LuckSale == nil {
		u.LuckSale = &model.LuckSale{
			Jackpot: -1,
		}
	}

	if u.LuckSale.Data == nil {
		u.LuckSale.Data = make(map[int]*model.LuckSaleUnit)
	}

	if u.LuckSale.Task == nil {
		u.LuckSale.Task = make([]*model.LuckSaleTaskUnit, 0)
	}
}

func verify_function_preview(u *model.UserData, now time.Time) {
	if u.FunctionPreview == nil {
		u.FunctionPreview = &model.FunctionPreview{}
	}

	if u.FunctionPreview.Data == nil {
		u.FunctionPreview.Data = make(map[uint32]msg.TaskState)
	}
}

func verify_equip_stage(u *model.UserData, now time.Time) {
	if u.EquipStage == nil {
		u.EquipStage = &model.UserEquipStage{}
	}

	if u.EquipStage.Records == nil {
		u.EquipStage.Records = make(map[uint32]*model.EquipStageRecord)
	}
}

func verify_resources_pass(u *model.UserData, now time.Time) {
	if u.ResourcesPass == nil {
		u.ResourcesPass = &model.UserResourcesPass{
			PassList: []*model.ResourcesPass{
				{
					PassType: MoneyPass,
				},
				{
					PassType: EquipPass,
				},
				{
					PassType: SidearmPass,
				},
				{
					PassType: PetPass,
				},
			},
			LastResetTime: time.Now(),
		}
	}
}

func verify_likes(u *model.UserData, now time.Time) {
	if u.Likes == nil {
		u.Likes = &model.LikesInfo{
			LikesMap: make(map[template.RankType]bool),
		}
	}
}

func verify_ranks(u *model.UserData, now time.Time) {
	if u.Ranks == nil {
		u.Ranks = &model.RankInfo{
			NormalPassRewardInfo: make(map[uint32]bool),
			ElitePassRewardInfo:  make(map[uint32]bool),
		}
	}
}

func verify_weekpass(u *model.UserData, now time.Time) {
	if u.WeekPass == nil {
		u.WeekPass = &model.WeekPass{
			ContractInfo:   make(map[uint32]bool),
			SecretCount:    0,
			SecretBoxState: 0,
		}
	}
}

func verify_personalized(u *model.UserData, now time.Time) {
	if u.Personalized == nil {
		u.Personalized = model.NewAccountPersonalized(u.UserId)
	}
	if u.Personalized.ItemsMap == nil {
		u.Personalized.AccountId = u.UserId
		u.Personalized.ItemsMap = make(map[int]*model.PersonalizedItem)
	}
}
