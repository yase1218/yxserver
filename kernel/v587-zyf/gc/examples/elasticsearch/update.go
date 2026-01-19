package elasticSearch

import (
	"context"
	"fmt"
	"github.com/v587-zyf/gc/db/elastic"
)

func Update(indexName, id string, doc interface{}) error {
	_, err := elastic.Get().Update().Index(indexName).Id(id).Doc(doc).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}
