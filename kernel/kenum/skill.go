package kenum

const (
	Skill_Type_Active  = iota + 1 // 主动
	Skill_Type_Passive            // 被动
)
const (
	Skill_CDR_All       = iota // 享受全冷却减免
	Skill_CDR_skillType        // 只享受技能类型冷却
	Skill_CDR_General          // 只享受通用冷却
	Skill_CDR_None             // 不享受冷却减免
)
const (
	Skill_Type_Add_General = iota // 根据技能类型享受加成
	Skill_Type_Add_Active         // 只享受主动技伤害加成
	Skill_Type_Add_Passive        // 只享受被动技伤害加成
	Skill_Type_Add_All            // 享受所有加成
	Skill_Type_Add_None           // 不享受加成
)
