package elasticSearch

import (
	"context"
	"github.com/v587-zyf/gc/db/elastic"
	"testing"
)

func TestDo(t *testing.T) {
	host := "127.0.0.1"
	port := 9200
	err := elastic.Init(context.Background(), elastic.WithHost(host), elastic.WithPort(port))
	if err != nil {
		panic(err)
	}

	id := "1"
	indexName := "student"
	Create(indexName, id)

	err = Search(indexName)
	if err != nil {
		panic(err)
	}

	err = Query(indexName)
	if err != nil {
		panic(err)
	}
}
