package model

import template2 "github.com/zy/game_data/template"

// Attr 属性
type Attr struct {
	Id         uint32
	InitValue  float32
	LevelValue float32
	Add        float32
	FinalValue float32
	Factor     float32
}

func DeepCopyAttr(a *Attr) *Attr {
	return &Attr{
		Id:         a.Id,
		InitValue:  a.InitValue,
		LevelValue: a.LevelValue,
		Add:        a.Add,
		FinalValue: a.FinalValue,
		Factor:     a.Factor,
	}
}

func NewAttr(id uint32, initValue float32) *Attr {
	return &Attr{
		Id:         id,
		InitValue:  initValue,
		Add:        0,
		FinalValue: initValue,
		Factor:     1.0,
		LevelValue: 0.0,
	}
}

// SetFactor 设置初始因子
func (a *Attr) SetFactor(factor float32) {
	a.Factor = factor
}

func (a *Attr) AddInitValue(data float32) {
	a.InitValue += data
}

func (a *Attr) AddLevelValue(data float32) {
	a.LevelValue += data
}

func (a *Attr) SetLevelValue(data float32) {
	a.LevelValue = data
}

func (a *Attr) AddValue(data float32) {
	a.Add += data
}

func (a *Attr) DelValue(data float32) {
	a.Add -= data
	if a.Add < 0 {
		a.Add = 0
	}
}

// GetRawValue init + level + add
func (a *Attr) GetRawValue() float32 {
	return a.InitValue + a.LevelValue + a.Add
}

func (a *Attr) CalcFinalValue() float32 {
	config := template2.GetAttrListTemplate().GetAttr(a.Id)
	if config.ValueType == 2 {
		a.FinalValue = ((a.InitValue+a.LevelValue)*a.Factor + a.Add) / 100
	} else {
		a.FinalValue = (a.InitValue+a.LevelValue)*a.Factor + a.Add
	}
	return a.FinalValue
}

func (a *Attr) SetFinalValue(value float32) {
	a.FinalValue = value
}

func (a *Attr) GetFinalValue() float32 {
	return a.FinalValue
}
