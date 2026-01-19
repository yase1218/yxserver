package db_model

import (
	"context"
	"github.com/v587-zyf/gc/db/db_model"
	"github.com/v587-zyf/gc/db/mysql"
	"testing"
	"time"
)

func TestInsertOne(t *testing.T) {
	uri := "root:root@tcp(127.0.0.1)/test?&parseTime=true&charset=utf8mb4"
	err := db_model.Init(context.Background(), mysql.WithUri(uri))
	if err != nil {
		t.Error(err)
		return
	}

	err = db_model.Init(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	id, err := GetTestModel().GetSeqId()
	if err != nil {
		t.Error(err)
		return
	}

	d := &Test{
		Id:   id,
		Time: time.Now(),
	}
	if err = GetTestModel().Create(d); err != nil {
		t.Error(err)
		return
	}

	data, _ := GetTestModel().GetOne(id)
	t.Log(data)
}
