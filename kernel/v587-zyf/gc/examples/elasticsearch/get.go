package elasticSearch

import (
	"context"
	"fmt"
	"github.com/v587-zyf/gc/db/elastic"
)

func Get(indexName, id string) {
	res, err := elastic.Get().Get().Index(indexName).Id(id).Do(context.Background())
	if err != nil {
		fmt.Println("search err:", err)
	}
	fmt.Println("search source:", string(res.Source))
}
