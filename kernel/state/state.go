package state

import (
	"fmt"
	"github.com/v587-zyf/gc/iface"
)

type State struct {
	curState    iface.IState
	transitions map[iface.StateTrigger]map[iface.IState]iface.IState
	isRun       bool
}

func NewState() *State {
	return &State{
		transitions: make(map[iface.StateTrigger]map[iface.IState]iface.IState),
		isRun:       false,
	}
}

func (f *State) Start(initState iface.IState) {
	f.curState = initState
	initState.Enter()
	initState.Execute()
}

func (f *State) Run() {
	if f.curState != nil && f.isRun {
		f.curState.ExecuteBefore()
		f.curState.Execute()
	}
}

func (f *State) Stop() {
	if f.curState != nil {
		f.curState.End()
	}
}

func (f *State) Pause() {
	f.isRun = false
}

func (f *State) Recover() {
	f.isRun = true
}

func (f *State) GetCurState() iface.IState {
	return f.curState
}

func (f *State) GetRunState() bool {
	return f.isRun
}

func (f *State) Register(trigger iface.StateTrigger, fromStates []iface.IState, toState iface.IState) {
	var stateMap map[iface.IState]iface.IState
	var ok bool
	for _, fromState := range fromStates {
		if stateMap, ok = f.transitions[trigger]; !ok {
			stateMap = make(map[iface.IState]iface.IState)
			f.transitions[trigger] = stateMap
		}
		stateMap[fromState] = toState
	}
}

func (f *State) SetInitState(state iface.IState) {
	f.curState = state
}

func (f *State) Trigger(trigger iface.StateTrigger) {
	if stateMap, ok := f.transitions[trigger]; ok {
		if state, ok := stateMap[f.curState]; ok {
			f.curState.End()
			f.curState = state
			state.Enter()
			return
		}
	}
}

func (f *State) String() string {
	str := fmt.Sprintf("isRun:%v curState:%t", f.isRun, f.curState != nil)
	return str
}
