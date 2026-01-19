package elasticSearch

import (
	"context"
	"fmt"
	"github.com/v587-zyf/gc/db/elastic"
)

func Del(indexName, id string) error {
	_, err := elastic.Get().Delete().Index(indexName).Id(id).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}
