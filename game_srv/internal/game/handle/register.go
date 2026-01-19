package handle

import (
	"msg"
)

func init() {
	init_gate_handler()
	init_fight_handler()
}

func init_gate_handler() {
	// 账号
	reg_gate(msg.MsgId_ID_RequestLogout, RequestLogoutHandle)
	reg_gate(msg.MsgId_ID_RequestClientHeart, RequestClientHeartHandle)
	reg_gate(msg.MsgId_ID_RequestUpdateNick, RequestUpdateNickHandle)
	reg_gate(msg.MsgId_ID_RequestUpgrade, RequestUpgradeHandle)
	reg_gate(msg.MsgId_ID_RequestSetShip, RequestSetShipHandle)
	reg_gate(msg.MsgId_ID_RequestSetSupportShip, RequestSetSupportShipHandle)
	reg_gate(msg.MsgId_ID_InitPlayerNameAndShipReq, RequestSetPlayerNameAndShipReq)
	reg_gate(msg.MsgId_ID_RandomPlayerNameReq, RequestRandomPlayerName)

	// 系统
	reg_gate(msg.MsgId_ID_RequestBuyAp, RequestBuyApHandle)
	reg_gate(msg.MsgId_ID_RequestUseCdk, RequestUseCdkHandle)
	reg_gate(msg.MsgId_ID_RequestRedPoint, RequestRedPointHandle)
	reg_gate(msg.MsgId_ID_RequestSetGuideInfo, RequestSetGuideInfoHandle)
	reg_gate(msg.MsgId_ID_RequestStaticsAction, RequestStaticsActionHandle)
	reg_gate(msg.MsgId_ID_RequestReconnectInfo, RequestReconnectInfoHandle)
	reg_gate(msg.MsgId_ID_RequestUploadBattleData, RequestUploadBattleDataHandle)
	reg_gate(msg.MsgId_ID_RequestSetVideo, RequestSetVideoHandle)
	reg_gate(msg.MsgId_ID_RequestClickUpdateNick, RequestClickUpdateNickHandle)
	reg_gate(msg.MsgId_ID_RequestSetPopUp, RequestSetPopUpHandle)
	reg_gate(msg.MsgId_ID_RequestGetQuestionReward, RequestGetQuestionRewardHandle)
	reg_gate(msg.MsgId_ID_RequestStartAd, RequestStartAdHandle)
	reg_gate(msg.MsgId_ID_RequestGetMonthCardReward, RequestGetMonthCardRewardHandle)
	reg_gate(msg.MsgId_ID_RequestGetMainFundReward, RequestGetMainFundRewardHandle)
	reg_gate(msg.MsgId_ID_RequestMonthCardDailyReward, RequestMonthCardDailyRewardHandle)
	reg_gate(msg.MsgId_ID_RequestGetAp, RequestGetApHandle)
	reg_gate(msg.MsgId_ID_RequestTalentUpgrade, RequestTalentUpgradeHandle)
	reg_gate(msg.MsgId_ID_RequestConstructionLottery, RequestConstructionLotteryHandle)
	reg_gate(msg.MsgId_ID_RequestAdventureGuideReward, RequestGetAdventureGuideRewardHandle)

	// 挂机
	reg_gate(msg.MsgId_ID_RequestOnHookData, RequestOnHookDataHandle)
	reg_gate(msg.MsgId_ID_RequestOnHookReward, RequestOnHookRewardHandle)
	reg_gate(msg.MsgId_ID_RequestQuickOnHookInfo, RequestQuickOnHookInfoHandle)
	reg_gate(msg.MsgId_ID_RequestQuickOnHookReward, RequestQuickOnHookRewardHandle)

	// 道具
	reg_gate(msg.MsgId_ID_RequestLoadItems, RequestLoadItemsHandle)
	reg_gate(msg.MsgId_ID_RequestUseItem, RequestUseItemHandle)
	reg_gate(msg.MsgId_ID_RequestComposeItem, RequestComposeItemHandle) // 道具合成

	// 战斗
	reg_gate(msg.MsgId_ID_RequestLoadMissions, RequestLoadMissionsHandle)
	reg_gate(msg.MsgId_ID_RequestStartBattle, RequestStartBattleHandle)
	reg_gate(msg.MsgId_ID_RequestEndBattle, RequestEndBattleHandle)
	reg_gate(msg.MsgId_ID_RequestGetMissionReward, RequestGetMissionRewardHandle)
	reg_gate(msg.MsgId_ID_RequestExitBattle, RequestExitBattleHandle)
	reg_gate(msg.MsgId_ID_RequestUpdateBattleStory, RequestUpdateBattleStoryHandle)
	reg_gate(msg.MsgId_ID_RequestSetBattleSpeed, RequestSetBattleSpeedHandle)
	reg_gate(msg.MsgId_ID_RequestPlayMethodData, RequestPlayMethodDataHandle)
	reg_gate(msg.MsgId_ID_RequestPlayMethodStartBattle, RequestPlayMethodStartBattleHandle)
	reg_gate(msg.MsgId_ID_RequestPlayMethodEndBattle, RequestPlayMethodEndBattleHandle)
	reg_gate(msg.MsgId_ID_RequestPlayMethodSwap, RequestPlayMethodSwapHandle)
	reg_gate(msg.MsgId_ID_RequestPlayMethodUpdateWeapon, RequestPlayMethodUpdateWeaponHandle)
	reg_gate(msg.MsgId_ID_RequestMissionExtraReward, RequestMissionExtraRewardHandle)

	// GM
	reg_gate(msg.MsgId_ID_RequestGMCommand, RequestGMCommandHandle)

	// 机甲
	reg_gate(msg.MsgId_ID_RequestLoadShips, RequestLoadShipsHandle)
	reg_gate(msg.MsgId_ID_RequestShipUpgrade, RequestShipUpgradeHandle)
	reg_gate(msg.MsgId_ID_RequestShipStarUpgrade, RequestShipStarUpgradeHandle)
	reg_gate(msg.MsgId_ID_RequestExchangeShip, RequestExchangeShipHandle)
	reg_gate(msg.MsgId_ID_RequestGlobalAttrDetail, RequestGlobalAttrDetailHandle)
	reg_gate(msg.MsgId_ID_RequestShipPreview, RequestShipPreviewHandle)

	// 编队
	reg_gate(msg.MsgId_ID_RequestLoadTeams, RequestLoadTeamsHandle)
	reg_gate(msg.MsgId_ID_RequestUpdateTeam, RequestUpdateTeamHandle)
	reg_gate(msg.MsgId_ID_RequestSetBattleTeam, RequestSetBattleTeamHandle)

	// 任务
	reg_gate(msg.MsgId_ID_RequestLoadTaskByType, RequestLoadTaskByTypeHandle)
	reg_gate(msg.MsgId_ID_RequestGetTaskReward, RequestGetTaskRewardHandle)
	reg_gate(msg.MsgId_ID_RequestGetActiveReward, RequestGetActiveRewardHandle)
	reg_gate(msg.MsgId_ID_RequestBatchGetTaskReward, RequestBatchGetTaskRewardHandle)

	// 装备系统
	reg_gate(msg.MsgId_ID_RequestLoadEquip, RequestLoadEquipHandle)
	reg_gate(msg.MsgId_ID_RequestEquipPos, RequestEquipPosHandle)
	reg_gate(msg.MsgId_ID_RequestEquipPosUpgrade, RequestEquipPosUpgradeHandle)
	reg_gate(msg.MsgId_ID_RequestEquipUpgradeStage, RequestEquipUpgradeStageHandle)
	reg_gate(msg.MsgId_ID_RequestSuitReward, RequestSuitRewardHandle)
	reg_gate(msg.MsgId_ID_RequestPutInSuit, RequestPutInSuitHandle)
	reg_gate(msg.MsgId_ID_RequestUseSuit, RequestUseSuitHandle)
	reg_gate(msg.MsgId_ID_RequestAllEquipUpgrade, RequestAllEquipUpgradeHandle)

	// 宝石
	reg_gate(msg.MsgId_ID_LoadGemReq, LoadGemReqHandle)
	reg_gate(msg.MsgId_ID_SocketGemReq, SocketGemReqHandle)
	reg_gate(msg.MsgId_ID_UnSocketGemReq, UnSocketGemReqHandle)
	reg_gate(msg.MsgId_ID_GemLockReq, GemLockReqHandle)
	reg_gate(msg.MsgId_ID_GemComposeReq, GemComposeReqHandle)
	reg_gate(msg.MsgId_ID_GemRefreshReq, GemRefreshReqHandle)

	// 装备本
	reg_gate(msg.MsgId_ID_LoadEquipStageReq, LoadEquipStageReqHandle)
	reg_gate(msg.MsgId_ID_EquipStageDropRecordReq, EquipStageDropRecordReqHandle)
	reg_gate(msg.MsgId_ID_EquipStageTeamInviteReq, EquipStageTeamInviteReqHandle)
	reg_gate(msg.MsgId_ID_EquipStageTeamAcceptReq, EquipStageTeamAcceptReqHandle)
	reg_gate(msg.MsgId_ID_EquipStageTeamLeaveReq, EquipStageTeamLeaveReq)             // 离队
	reg_gate(msg.MsgId_ID_EquipStageTeamDissolveReq, EquipStageTeamDissolveReqHandle) // 解散队伍
	reg_gate(msg.MsgId_ID_EquipStageMatchReq, EquipStageMatchReqHandle)
	reg_gate(msg.MsgId_ID_EquipStageAcceptReq, EquipStageAcceptReqHandle)
	reg_gate(msg.MsgId_ID_EquipStageMatchCancelReq, EquipStageMatchCancelReqHandle)
	reg_gate(msg.MsgId_ID_EquipStageLoadReq, EquipStageLoadReqHandle)
	reg_gate(msg.MsgId_ID_EquipStageRewardPickReq, EquipStageRewardPickReqHandle)
	reg_gate(msg.MsgId_ID_EquipStageBuyRewardReq, EquipStageBuyRewardReqHandle)
	reg_gate(msg.MsgId_ID_EquipStageLeaveReq, EquipStageLeaveReqHandle)
	reg_gate(msg.MsgId_ID_RequestCreateEquipStageTeam, CreateEquipStageTeamHandle)

	// 邮件 TODO
	reg_gate(msg.MsgId_ID_RequestLoadMail, RequestLoadMailHandle)
	reg_gate(msg.MsgId_ID_RequestReadMail, RequestReadMailHandle)
	reg_gate(msg.MsgId_ID_RequestGetMailReward, RequestGetMailRewardHandle)
	reg_gate(msg.MsgId_ID_RequestBatchGetMailReward, RequestBatchGetMailRewardHandle)
	reg_gate(msg.MsgId_ID_RequestDelMail, RequestDelMailHandle)

	// 商店
	reg_gate(msg.MsgId_ID_RequestLoadShop, RequestLoadShopHandle)
	reg_gate(msg.MsgId_ID_RequestBuyShopItem, RequestBuyShopItemHandle)
	reg_gate(msg.MsgId_ID_RequestRefreshShop, RequestRefreshShopHandle)
	reg_gate(msg.MsgId_ID_RequestGetShopItem, RequestGetShopItemHandle)
	reg_gate(msg.MsgId_ID_FirstChargePackageReq, RequestFirstChargePackageHandle)
	reg_gate(msg.MsgId_ID_FirstChargePackageGetReq, RequestFirstChargePackageGetHandler)

	// 武器系统
	reg_gate(msg.MsgId_ID_RequestLoadWeapon, RequestLoadWeaponHandle)
	reg_gate(msg.MsgId_ID_RequestWeaponUpgrade, RequestWeaponUpgradeHandle)
	reg_gate(msg.MsgId_ID_RequestSetSecondaryWeapon, RequestSetSecondaryWeaponHandle)
	reg_gate(msg.MsgId_ID_RequestActiveWeapon, RequestActiveWeaponHandle)

	// 秘宝系统
	reg_gate(msg.MsgId_ID_RequestLoadRareTreasure, RequestLoadRareTreasureHandle)

	// 世界探索 2025.9未来会改
	// reg(msg.MsgId_ID_ExploreReq, ExploreReqHandle)
	// reg(msg.MsgId_ID_ExploreStageReq, ExploreStageReqHandle)
	// reg(msg.MsgId_ID_ExploreDrawCardReq, ExploreDrawCardReqHandle)
	// reg(msg.MsgId_ID_ExploreCollectReq, ExploreCollectReqHandle)
	// reg(msg.MsgId_ID_ExploreAllCollectReq, ExploreAllCollectReqHandle)
	// reg(msg.MsgId_ID_ExploreReviveReq, ExploreReviveReqHandle)
	// reg(msg.MsgId_ID_ExploreEnterStageReq, ExploreEnterStageReqHandle)
	// reg(msg.MsgId_ID_ExploreLeaveStageReq, ExploreLeaveStageReqHandle)
	// reg(msg.MsgId_ID_ExploreUnlockBuildingReq, ExploreUnlockBuildingReqHandle)
	// reg(msg.MsgId_ID_ExploreGetAllBuildingsReq, ExploreGetAllBuildingsReqHandle)
	// reg(msg.MsgId_ID_ExploreGetFreeCardReq, ExploreGetFreeCardReqHandle)
	// reg(msg.MsgId_ID_ExploreUseFreeCardReq, ExploreUseFreeCardReqHandle)

	// 扑克系统
	reg_gate(msg.MsgId_ID_RequestLoadPoker, RequestLoadPokerHandle)

	// 抽奖系统
	reg_gate(msg.MsgId_ID_RequestLoadCardPool, RequestLoadCardPoolHandle)
	reg_gate(msg.MsgId_ID_RequestLottery, RequestLotteryHandle)
	reg_gate(msg.MsgId_ID_RequestLotteryTimesReward, RequestLotteryTimesRewardHandle)

	// 活动
	reg_gate(msg.MsgId_ID_RequestActConfig, RequestActConfigHandle)
	reg_gate(msg.MsgId_ID_RequestLoadActData, RequestLoadActDataHandle)
	reg_gate(msg.MsgId_ID_RequestGetActReward, RequestGetActRewardHandle)
	reg_gate(msg.MsgId_ID_RequestBuyPassGrade, RequestBuyPassGradeHandle)
	reg_gate(msg.MsgId_ID_RequestSuppSign, RequestSuppSignHandle)
	//reg_gate(msg.MsgId_ID_RequestEnterDesertAct, RequestEnterDesertActHandle)
	reg_gate(msg.MsgId_ID_RequestActPreviewReward, RequestActPreviewRewardHandle)

	// 排行榜
	reg_gate(msg.MsgId_ID_GetTargetRankInfoReq, RequestRankData)
	reg_gate(msg.MsgId_ID_AddLikesForRankPlayerReq, RequestAddLikesForRank)
	reg_gate(msg.MsgId_ID_GetMaxFirstPassRankRewardReq, RequestGetMaxFirstPassReward)
	reg_gate(msg.MsgId_ID_GetFirstPassRecordDataReq, RequestFirstPassRecordData)
	// reg(msg.MsgId_ID_RequestMissionRank, RequestMissionRankHandle)
	// reg(msg.MsgId_ID_RequestPetRank, RequestPetRankHandle)
	// reg(msg.MsgId_ID_RequestLikes, RequestLikesHandle)
	// reg(msg.MsgId_ID_RequestRankRewardInfo, RequestRankRewardInfoHandle)
	// reg(msg.MsgId_ID_RequestRankMissionReward, RequestRankMissionRewardHandle)
	reg_gate(msg.MsgId_ID_RequestSpecialMissionRank, RequestSpecialMissionRankHandle)
	// reg(msg.MsgId_ID_RequestDesertRank, RequestDesertRankHandle)

	// 外观
	reg_gate(msg.MsgId_ID_RequestLoadAppearance, RequestLoadAppearanceHandle)
	reg_gate(msg.MsgId_ID_RequestUseAppearance, RequestUseAppearanceHandle)
	reg_gate(msg.MsgId_ID_RequestActiveAppearance, RequestActiveAppearanceHandle)

	// 宠物系统
	reg_gate(msg.MsgId_ID_GetPlayerPetsInfoReq, RequestLoadPetHandle)
	reg_gate(msg.MsgId_ID_ActavitePetReq, RequestActPet)
	reg_gate(msg.MsgId_ID_UpdatePetStarLvReq, RequestUpdatePetStarLv)
	reg_gate(msg.MsgId_ID_UpdatePetLvReq, RequestUpdatePetLv)
	reg_gate(msg.MSGID_ID_UPDATEPETSTATEREQ, RequestUpdatePetsState)
	// reg_gate(msg.MsgId_ID_RequestLoadPet, RequestLoadPetHandle)
	// reg(msg.MsgId_ID_RequestGetPetMaterial, RequestGetPetMaterialHandle)
	// reg(msg.MsgId_ID_RequestGetAdventureReward, RequestGetAdventureRewardHandle)
	// reg(msg.MsgId_ID_RequestCompletePetEgg, RequestCompletePetEggHandle)
	// reg(msg.MsgId_ID_RequestActivePetPart, RequestActivePetPartHandle)
	// reg(msg.MsgId_ID_RequestUsePetPart, RequestUsePetPartHandle)
	// reg(msg.MsgId_ID_RequestPetCompleteGrowUp, RequestPetCompleteGrowUpHandle)
	// reg(msg.MsgId_ID_RequestPetGrowUpAccelerate, RequestPetGrowUpAccelerateHandle)
	// reg(msg.MsgId_ID_RequestUsePetSuit, RequestUsePetSuitHandle)
	// reg(msg.MsgId_ID_RequestRenamePet, RequestRenamePetHandle)
	// reg(msg.MsgId_ID_RequestSuccinct, RequestSuccinctHandle)
	// reg(msg.MsgId_ID_RequestEvolution, RequestEvolutionHandle)
	// reg(msg.MsgId_ID_RequestPetReward, RequestPetRewardHandle)

	// 社交
	reg_gate(msg.MsgId_ID_RequestPlayerDetailInfo, RequestPlayerDetailHandle)

	// 聊天
	reg_gate(msg.MsgId_ID_RequestWorldChat, RequestWorldChatHandle)
	reg_gate(msg.MsgId_ID_RequestEnterWorldChat, RequestEnterWorldChatHandle)
	reg_gate(msg.MsgId_ID_RequestLeaveWorldChat, RequestLeaveWorldChatHandle)
	reg_gate(msg.MsgId_ID_RequestWorldChatData, RequestWorldChatDataHandle)

	// // 好友
	reg_gate(msg.MsgId_ID_RequestFriList, RequestFriListHandle)
	reg_gate(msg.MsgId_ID_RequestFriApplyList, RequestFriApplyListHandle)
	reg_gate(msg.MsgId_ID_RequestAddFriend, RequestAddFriendHandle)
	reg_gate(msg.MsgId_ID_RequestFriApplyOp, RequestFriApplyOpHandle)
	reg_gate(msg.MsgId_ID_RequestDelFriend, RequestDelFriendHandle)
	reg_gate(msg.MsgId_ID_RequestRecommandFriend, RequestRecommandFriendHandle)
	reg_gate(msg.MsgId_ID_RequestSearchPlayer, RequestSearchPlayerHandle)
	reg_gate(msg.MsgId_ID_RequestFriBlackList, RequestFriBlackListHandle)
	reg_gate(msg.MsgId_ID_RequestBlackListOp, RequestBlackListOpHandle)
	reg_gate(msg.MsgId_ID_RequestPrivateChat, RequestPrivateChatHandle)
	// reg(msg.MsgId_ID_EmoteSendReq, EmoteSend)

	// 联盟 TODO 重构
	// reg(msg.MsgId_ID_SearchAllianceReq, HandleSearchAlliance)
	// reg(msg.MsgId_ID_CreateAllianceReq, HandleCreateAlliance)
	// reg(msg.MsgId_ID_JoinAllianceReq, HandleJoinAlliance)
	// reg(msg.MsgId_ID_QuickJoinAllianceReq, HandleQuickJoinAlliance)
	// reg(msg.MsgId_ID_GetAllianceMembersReq, HandleGetAllianceMembers)
	// reg(msg.MsgId_ID_QuitAllianceReq, HandleQuitAlliance)
	// reg(msg.MsgId_ID_AllianceManageReq, HandleAllianceManage)
	// reg(msg.MsgId_ID_AllianceApplyListReq, HandleAllianceApplyList)
	// reg(msg.MsgId_ID_HandleAllianceApplyReq, HandleAllianceApply)
	// reg(msg.MsgId_ID_GetRedPacketListReq, HandleGetRedPacketList)
	// reg(msg.MsgId_ID_ClaimRedPacketReq, HandleClaimRedPacket)
	// reg(msg.MsgId_ID_AllianceInfoUpdateReq, HandleAllianceInfoUpdate)
	// reg(msg.MsgId_ID_RequenstEnterUnionBattle, RequenstEnterUnionBattleHandle)
	// reg(msg.MsgId_ID_ClearAllianceApplyReq, ClearAllianceApplyReqHandle)

	// 联盟排行榜  TODO 重构
	// reg(msg.MsgID_AllianceRankReqId, AllianceRankReq)

	// 客户端战斗消息
	reg_gate(msg.MsgID_CreateFightReqId, CreateFight)
	reg_gate(msg.MsgID_StartFightReqId, StartFight)
	reg_gate(msg.MsgID_EndFightReqId, EndFight)
	reg_gate(msg.MsgID_GetWeaponReqId, GetWeapon)
	reg_gate(msg.MsgID_SelectWeaponReqId, SelectWeapon)
	reg_gate(msg.MsgID_GetAccessoryReqId, GetAccessories)
	reg_gate(msg.MsgID_SelectAccessoryReqId, SelectAccessories)
	reg_gate(msg.MsgID_FightBroadcastReqId, FightBroadcast)
	reg_gate(msg.MsgID_FightMonsterReqId, FightMonster)

	// 巅峰战场 TODO 重构
	reg_gate(msg.MsgID_PeakFightReqId, PeakFight)
	reg_gate(msg.MsgID_PeakFightMatchReqId, PeakFightMatch)
	reg_gate(msg.MsgID_PeakFightCancelMatchReqId, PeakFightCancelMatch)
	reg_gate(msg.MsgID_PeakFightRankReqId, PeakFightRank)

	// 疯狂合约
	reg_gate(msg.MsgID_ContractReqId, Contract)
	reg_gate(msg.MsgID_ContractSignReqId, ContractSign)
	reg_gate(msg.MsgID_ContractCancelReqId, ContractCancel)
	reg_gate(msg.MsgID_ContractRandReqId, ContractRand)
	reg_gate(msg.MsgID_ContractRewardReqId, ContractReward)

	// 沙漠大冒险
	reg_gate(msg.MsgID_DesertReqId, Desert)
	reg_gate(msg.MsgID_DesertLampRewardReqId, DesertLampReward)

	// 客户端自定义内容
	reg_gate(msg.MsgID_ClientSettingsReqId, ClientSettings)
	reg_gate(msg.MsgID_ClientSettingsUpdateReqId, ClientSettingsUpdate)
	reg_gate(msg.MsgID_ClientSettingsDeleteReqId, ClientSettingsDelete)

	//竞技场 TODO 不开放
	// reg(msg.MsgID_ArenaInfoReqId, RequestArenaInfo)
	// reg(msg.MsgID_ArenaRankReqId, RequestRanks)
	// reg(msg.MsgID_ArenaPkReqId, RequestArenaPk)
	// reg(msg.MsgID_ArenaRecordReqId, RequestArenaPkRecord)
	// reg(msg.MsgID_ArenaSetDefendReqId, RequestArenaSetDefend)
	// reg(msg.MsgID_ArenaPkListReqId, RequestArenPkList)
	// reg(msg.MsgID_ArenaBuyPkCntReqId, RequestArenBuyPkCnt)
	// reg(msg.MsgID_ArenaMonsterInfoReqId, RequestMonsterInfo)
	// reg(msg.MsgID_ArenaUnlockMonsterReqId, RequestArenUnlockMonster)
	// reg(msg.MsgID_ArenaUnlockPosReqId, RequestArenUnlockPos)
	// reg(msg.MsgID_ArenaGetRewardReqId, RequestArenGetRewardList)

	// 图鉴
	reg_gate(msg.MsgID_AtlasReqId, Atlas)
	reg_gate(msg.MsgID_AtlasRewardReqId, AtlasReward)

	// 幸运售货机
	reg_gate(msg.MsgID_LuckSaleReqId, LuckSale)
	reg_gate(msg.MsgID_LuckSaleExtractReqId, LuckSaleExtract)
	reg_gate(msg.MsgID_LuckSaleTaskRewardReqId, LuckSaleTaskReward)

	// 功能预览
	reg_gate(msg.MsgId_ID_FunctionPreviewReq, FunctionPreview)
	reg_gate(msg.MsgId_ID_FunctionPreviewRewardReq, FunctionPreviewReward)

	// 资源本
	reg_gate(msg.MsgId_ID_GetResourcesPassBaseDataReq, HandleGetResourcesPassBaseData)
	reg_gate(msg.MsgId_ID_BuyResoucesPassItemReq, HandleBuyResoucePassItem)
	reg_gate(msg.MsgId_ID_ResoucePassAttackReq, HandleResourcesAttack)
	reg_gate(msg.MsgId_ID_ResourcePassStateChangeReq, HandleResourceStateChange)
	reg_gate(msg.MsgId_ID_ResourcePassRankListReq, HandleGetResourcesPassRankList)

	// 周常活动
	reg_gate(msg.MsgId_ID_RequestGetWeekPassAct, RequestGetWeekPassActHandle)
	reg_gate(msg.MsgId_ID_SetWeekPassDeputyWeaponAndFactionReq, RequestSetFactionAndWeaponHandle)
	reg_gate(msg.MsgId_ID_GetContractInfoReq, RequestGetContractInfo)
	reg_gate(msg.MsgId_ID_GetSecretInfoReq, RequestSecretInfo)

	// 推送
	reg_gate(msg.MsgId_ID_RequestPersonalized, RequestPersonalizedActHandle)
	reg_gate(msg.MsgId_ID_RequestUnOutTimePersonalized, RequestUnOutTimePersonalizedHandle)

	// 充值
	reg_gate(msg.MsgId_ID_RequestCreateOrder, RequestCreateOrderHandle)

	// 皮肤
	reg_gate(msg.MsgId_ID_RequestActiveCoat, RequestActiveCoatHandle)
	reg_gate(msg.MsgId_ID_RequestPutOnCoat, RequestPutOnCoatHandle)

	// 海边派对黄历 周环境
	reg_gate(msg.MsgId_ID_RequestBattlePassWeekFields, BattlePassWeekFieldsHandle)

	// 拒绝重连
	reg_gate(msg.MsgId_ID_RejectReconnectFightReq, RejectReconnect)
}

func init_fight_handler() {
	// 战斗服战斗消息
	reg_fight(msg.MsgID_FsCreateFightAckId, FsCreateFightAck)
	reg_fight(msg.MsgID_FsFightResultNtfId, FsFightResultNtf)
	reg_fight(msg.MsgID_FsPickItemNtfId, FsPickItemNtf)
}
