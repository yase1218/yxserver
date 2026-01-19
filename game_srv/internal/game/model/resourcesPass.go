package model

type ResourcesPass struct {
	PassType     uint32 `bson:"pass_type"`
	BuyCount     uint32 `bson:"buy_count"`
	HitCount     uint32 `bson:"hit_count"`
	Total        uint32 `bson:"total"`
	ExpressState bool   `bson:"express_state"`
}
