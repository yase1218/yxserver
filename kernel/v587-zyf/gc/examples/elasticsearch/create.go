package elasticSearch

import (
	"context"
	"fmt"
	"github.com/v587-zyf/gc/db/elastic"
)

type Student struct {
	Name  string
	Age   int
	Hobby []string
}

func Create(indexName, id string) {
	s := Student{
		Name: "小学生",
		Age:  7,
		Hobby: []string{
			"soccer",
			"basketball",
			"tennis",
		},
	}
	res, err := elastic.Get().
		Index().
		Index(indexName).
		Id(id).
		BodyJson(s).
		Do(context.Background())
	if err != nil {
		fmt.Println("Create err:", err.Error())
		return
	}
	fmt.Println(res)
}
