package elasticSearch

import (
	"context"
	"fmt"
	elastics "github.com/olivere/elastic/v7"
	"github.com/v587-zyf/gc/db/elastic"
)

func Query(indexName string) error {
	// 指定字段
	q := elastics.NewQueryStringQuery("UID")
	res, err := elastic.Get().Search(indexName).Query(q).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(res)

	/*
		// 大于
		qGt := elastic.NewRangeQuery("age").Gt(5)
		res, err = client.Search(indexName).Query(qGt).Do(context.Background())
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println("2----------------------------")
		fmt.Println(res)

		// 多个条件
		boolQ := elastic.NewBoolQuery()
		boolQ.Must(elastic.NewMatchQuery("name", "小学生"))
		boolQ.Filter(elastic.NewRangeQuery("age").Gt(5))
		res, err = client.Search(indexName).Query(boolQ).Do(context.Background())
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println("3----------------------------")
		fmt.Println(res)

		// 分页
		pg := 1
		size := 1
		res, err = client.Search(indexName).Size(size).From((pg - 1) * 1).Do(context.Background())
		for _, v := range res.Hits.Hits {
			log.Println("分页: ", string(v.Source))
		}

		// range 范围匹配
		// id >= 10, id <= 100
		qLte := elastic.NewRangeQuery("age").Gte(1).Lte(25)
		res, err = client.Search(indexName).Query(qLte).Do(context.Background())
		for _, v := range res.Hits.Hits {
			log.Println("范围匹配 = ", string(v.Source))
		}

		// name字段模糊匹配
		qMatch := elastic.NewMatchQuery("name", "*li*")
		res, err = client.Search(indexName).Query(qMatch).Do(context.Background())
		log.Println(res, err)
		for _, v := range res.Hits.Hits {
			log.Println("name字段模糊匹配 = ", string(v.Source))
		}*/

	/*
		条件查询
		var query elastic.Query

		// match_all
		query = elastic.NewMatchAllQuery()

		// term
		query = elastic.NewTermQuery("field_name", "field_value")

		// terms
		query = elastic.NewTermsQuery("field_name", "field_value")

		// match
		query = elastic.NewMatchQuery("field_name", "field_value")

		// match_phrase
		query = elastic.NewMatchPhraseQuery("field_name", "field_value")

		// match_phrase_prefix
		query = elastic.NewMatchPhrasePrefixQuery("field_name", "field_value")

		//range Gt:大于; Lt:小于; Gte:大于等于; Lte:小于等
		query = elastic.NewRangeQuery("field_name").Gte(1).Lte(2)

		//regexp
		query = elastic.NewRegexpQuery("field_name", "regexp_value")

		_, err := client.Search().Index("index_name").Query(query).Do(context.Background())

		if err != nil {
			panic(err)
		}

		//排序顺序, true为降徐， false为升序
		client.Search().Index("index_name").Sort("field_name", true).Do(context.Background())

		// 还可以通过SortBy进行多个排序
		sorts := []elastic.Sorter{
			elastic.NewFieldSort("field_name01").Asc(), // 升序
			elastic.NewFieldSort("field_name02").Desc(), // 降徐
		}
		client.Search().Index("index_name").SortBy(sorts...).Do(context.Background())
	*/

	return nil
}
