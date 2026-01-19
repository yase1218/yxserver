package mongo

import (
	"testing"
	"time"
)

func TestInsertOne(t *testing.T) {
	id, err := GetTestIdSeq()
	if err != nil {
		t.Error(err)
		return
	}

	d := &Test{
		ID:   id,
		Time: time.Now(),
	}
	if err = GetTestMongo().InsertOne(d); err != nil {
		t.Error(err)
		return
	}

	data, _ := GetTestMongo().LoadOne(id)
	t.Log(data)
}
