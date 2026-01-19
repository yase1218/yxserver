package mysql

import (
	"context"
	"github.com/v587-zyf/gc/db/mysql"
	"testing"
	"time"
)

func TestInsertOne(t *testing.T) {
	uri := "root:root@tcp(127.0.0.1)/test?&parseTime=true&charset=utf8mb4"
	err := mysql.Init(context.Background(), mysql.WithUri(uri))
	if err != nil {
		t.Error(err)
		return
	}

	d := &TestModel{
		Time: time.Now(),
	}
	if err = mysql.CreateModel[TestModel](d, nil); err != nil {
		t.Error(err)
		return
	}

	data := &TestModel{
		ModelBase: ModelBase{
			ID: d.ID,
		},
	}
	if err = mysql.LoadModel[TestModel](data); err != nil {
		t.Error(err)
		return
	}
	t.Log(data)
}
