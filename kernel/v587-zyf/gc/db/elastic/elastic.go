package elastic

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
	"os"
	"time"
)

type Elastic struct {
	options *ElasticOption

	client *elastic.Client

	ctx    context.Context
	cancel context.CancelFunc
}

func NewElastic() *Elastic {
	m := &Elastic{
		options: NewMysqlOption(),
	}

	return m
}

func (e *Elastic) Init(ctx context.Context, opts ...any) (err error) {
	e.ctx, e.cancel = context.WithCancel(ctx)
	if len(opts) > 0 {
		for _, opt := range opts {
			opt.(Option)(e.options)
		}
	}

	var linkUrl string
	if e.options.https {
		linkUrl = fmt.Sprintf("https://%s:%d", e.options.host, e.options.port)
	} else {
		linkUrl = fmt.Sprintf("http://%s:%d", e.options.host, e.options.port)
	}
	esopts := []elastic.ClientOptionFunc{
		// elasticsearch 服务地址，多个服务地址使用逗号分隔
		elastic.SetURL(linkUrl),
		// 启用gzip压缩
		elastic.SetGzip(true),
		// 设置监控检查时间间隔
		elastic.SetHealthcheckInterval(10 * time.Second),
		// 设置错误日志输出
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		// 设置info日志输出
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
		// 设置嗅探功能
		elastic.SetSniff(false),
	}
	if e.options.userName != "" && e.options.password != "" {
		esopts = append(esopts, elastic.SetBasicAuth(e.options.userName, e.options.password))
	}
	e.client, err = elastic.NewClient(esopts...)
	if err != nil {
		fmt.Printf("elastic connection error: %v\n", err)
		return
	}

	_, _, err = e.client.Ping(e.options.host).Do(context.Background())
	if err != nil {
		log.Printf("elastic connection error: %v\n", err)
		return
	}
	//fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)
	//
	//esVersion, err := e.client.ElasticsearchVersion(e.options.host)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("Elasticsearch version %s\n", esVersion)
	return
}

func (e *Elastic) Get() *elastic.Client {
	return e.client
}
func (e *Elastic) GetCtx() context.Context {
	return e.ctx
}
