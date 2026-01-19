package mongo

type MongoOption struct {
	uri string
	db  string
}

type Option func(o *MongoOption)

func NewMongoOption() *MongoOption {
	return &MongoOption{}
}

func WithUri(uri string) Option {
	return func(o *MongoOption) {
		o.uri = uri
	}
}

func WithDb(db string) Option {
	return func(o *MongoOption) {
		o.db = db
	}
}
