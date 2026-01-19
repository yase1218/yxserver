package tdb

import c "github.com/v587-zyf/gc/tabledb"

var fileInfos = []c.FileInfo{
	{"globals.xlsx", []c.SheetInfo{
		{"global", c.LoadGlobalConf, c.GlobalBaseCfg{}},
	}},
	
	{"chapter.xlsx", []c.SheetInfo{
			{SheetName: "chapter", Initer: c.MapLoader("ChapterChapterCfgs", "Id"), ObjPropType: ChapterChapterCfg{}},
	}},
	{"item.xlsx", []c.SheetInfo{
			{SheetName: "item", Initer: c.MapLoader("ItemItemCfgs", "Id"), ObjPropType: ItemItemCfg{}},
	}},
}


type TableBase struct {
	// NOTE 关于client的配置：
	// client:对象名,对象类型　，对象名要小写．
	// mapKey 即对应的我们的结构里的key, 要看具体的型中key是什么　，一段是大写的
	
    ChapterChapterCfgs		map[int]*ChapterChapterCfg
    ItemItemCfgs		map[int]*ItemItemCfg
}	


func GetChapterChapterCfg( Id int) *ChapterChapterCfg {
	return tdb.ChapterChapterCfgs[Id]
}

func RangChapterChapterCfgs(f func(conf *ChapterChapterCfg)bool){
	for _,v := range tdb.ChapterChapterCfgs{
		if !f(v){
			return
		}
	}
}

func GetItemItemCfg( Id int) *ItemItemCfg {
	return tdb.ItemItemCfgs[Id]
}

func RangItemItemCfgs(f func(conf *ItemItemCfg)bool){
	for _,v := range tdb.ItemItemCfgs{
		if !f(v){
			return
		}
	}
}

type ChapterChapterCfg struct {
	Id int	`col:"id" client:"id"`	// 关卡id
	Locked bool	`col:"locked" client:"locked"`	// 是否锁定
}

type ItemItemCfg struct {
	Id int	`col:"id" client:"id"`	// id
	Type int	`col:"type" client:"type"`	// 类型
	Name string	`col:"name" client:"name"`	// 名称
	Deal_type int	`col:"deal_type" client:"deal_type"`	// 交易类型
	Deal_price float64	`col:"deal_price" client:"deal_price"`	// 交易初始价格
	Deal_num int	`col:"deal_num" client:"deal_num"`	// 交易捆绑数量
	Desc string	`col:"desc" client:"desc"`	// 备注
}



type InitConf struct {
	Strength_max int `conf:"strength_max" default:"240"` //体力上限
	Add_strength_minutes int `conf:"add_strength_minutes" default:"6"` //恢复体力时间(分钟)
	Complete_chapter_less_strength int `conf:"complete_chapter_less_strength" default:"10"` //完成关卡消耗体力
	Add_strength_num int `conf:"add_strength_num" default:"1"` //每次恢复体力值
	Clear_special_coin_get_num int `conf:"clear_special_coin_get_num" default:"120"` //消除一次特殊货币获得数量
	Ore_extra float64 `conf:"ore_extra" default:"6.25"` //额外小矿石产出参数
	Add_gacha_minutes int `conf:"add_gacha_minutes" default:"480"` //高级盲盒进度上涨时间(分钟)
	Gacha_max int `conf:"gacha_max" default:"90"` //高级盲盒进度上限
	Unlock_mine_grid_ton int `conf:"unlock_mine_grid_ton" default:"1"` //挖矿槽位解锁Ton价格
	Mine_grid_max int `conf:"mine_grid_max" default:"100"` //挖矿槽位上限
	Mine_reward_seconds int `conf:"mine_reward_seconds" default:"850"` //挖矿收益时间(秒)

}
