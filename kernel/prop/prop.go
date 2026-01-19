package prop

import "msg"

const PropBaseRate = 10000.0

var PropEntryType = map[int]bool{}

type (
	PropEntry struct {
		Total       int
		Base        int
		BasePercent int
	}
	Prop struct {
		props map[int]*PropEntry
	}
)

func (pe *PropEntry) calc() {
	base := float64(pe.Base) * (1.0 + float64(pe.BasePercent)/PropBaseRate)
	pe.Base = int(base)
}

func NewProp() *Prop {
	p := &Prop{}

	p.Reset()
	return p
}

func (p *Prop) Reset() {
	p.props = make(map[int]*PropEntry)
	for t := range PropEntryType {
		p.props[t] = &PropEntry{}
	}
}

// 增加一个全属性
func (p *Prop) AddAllProp(prop *Prop) {
	p.Add(prop.baseValueMap())
}

func (p *Prop) Add(props map[int]int) {
	for k, v := range props {
		p.AddOne(k, v)
	}
}

func (p *Prop) Calc() {
	p.calcEach()
}

func (p *Prop) getBaseValue(id int) int {
	if _, ok := PropEntryType[id]; ok {
		_p := p.props[id]
		if _p == nil {
			return 0
		}
		return _p.Base
	}
	return 0
}

func (p *Prop) Get(id int) int {
	if _, ok := PropEntryType[id]; ok {
		_p := p.props[id]
		if _p == nil {
			return 0
		}
		return _p.Total
	}
	return 0
}

func (p *Prop) AddOne(id, value int) {
	if p.props[id] == nil {
		p.props[id] = &PropEntry{}
	}
	p.props[id].Base += value
}

func (p *Prop) baseValueMap() map[int]int {
	m := make(map[int]int, len(msg.Attribute_value))
	for _, pId := range msg.Attribute_value {
		v := p.getBaseValue(int(pId))
		if v > 0 {
			m[int(pId)] = v
		}
	}
	return m
}

func (p *Prop) calcEach() {

}
