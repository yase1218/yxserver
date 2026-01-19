package elasticSearch

import (
	"context"
	"fmt"
	"github.com/v587-zyf/gc/db/elastic"
)

type mi = map[string]interface{}

// CreateIndex 创建索引
func CreateIndex(indexName string) {
	mapping := mi{
		"settings": mi{
			"number_of_shards":   3,
			"number_of_replicas": 2,
		},
		"mappings": mi{
			"_doc": mi{ //type名
				"properties": mi{
					"id": mi{ //整形字段, 允许精确匹配
						"type": "integer",
					},
					"name": mi{
						"type":            "text",     //字符串类型且进行分词, 允许模糊匹配
						"analyzer":        "ik_smart", //设置分词工具
						"search_analyzer": "ik_smart",
						"fields": mi{ //当需要对模糊匹配的字符串也允许进行精确匹配时假如此配置
							"keyword": mi{
								"type":         "keyword",
								"ignore_above": 256,
							},
						},
					},
					"date_field": mi{ //时间类型, 允许精确匹配
						"type": "date",
					},
					"keyword_field": mi{ //字符串类型, 允许精确匹配
						"type": "keyword",
					},
					"nested_field": mi{ //嵌套类型
						"type": "nested",
						"properties": mi{
							"id": mi{
								"type": "integer",
							},
							"start_time": mi{ //长整型, 允许精确匹配
								"type": "long",
							},
							"end_time": mi{
								"type": "long",
							},
						},
					},
				},
			},
		},
	}
	_, err := elastic.Get().CreateIndex(indexName).BodyJson(mapping).Do(context.Background())
	if err != nil {
		fmt.Println("es createIndex err:", err)
	}
}

// CheckIndexExists 索引是否存在 存在:true
func CheckIndexExists(indexName string) bool {
	isHas, err := elastic.Get().IndexExists(indexName).Do(context.Background())
	if err != nil {
		fmt.Println("check index err:", err)
	}
	return isHas
}

func UpdateIndex(indexName string) {
	mapping := mi{
		"properties": mi{
			"id": mi{
				"type": "integer",
			},
		},
	}
	_, err := elastic.Get().PutMapping().Index(indexName).BodyJson(mapping).Do(context.Background())
	if err != nil {
		fmt.Println("update index err:", err)
	}
}

// DelIndex 删除索引
func DelIndex(indexName string) {
	_, err := elastic.Get().DeleteIndex(indexName).Do(context.Background())
	if err != nil {
		fmt.Println("del index err:", err)
	}
}
