package tda

import (
	"testing"
)

func TestTda(t *testing.T) {
	var (
		err error
	)
	if err = Init(); err != nil {
		t.Errorf("tda init err:%v", err)
		return
	}

	accountId := "1"
	distinctId := "1"
	props := map[string]any{}
	if err = GetTda().SetUser(accountId, distinctId, props); err != nil {
		t.Errorf("tda set user err:%v", err)
		return
	}
}
