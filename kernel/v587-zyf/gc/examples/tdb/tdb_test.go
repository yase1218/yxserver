package tdb

import (
	"testing"
)

func TestInit(t *testing.T) {
	path := "./"
	_ = Init(path)

	t.Log(GetItemItemCfg(1))
}
