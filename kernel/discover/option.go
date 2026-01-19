package discover

import "time"

type DiscoverOption struct {
	endpoints   []string
	dialTimeout time.Duration
	updateFn    UpdateFn
	removeFn    RemoveFn
}

type Option func(*DiscoverOption)

func NewDiscoverOption() *DiscoverOption {
	return &DiscoverOption{}
}

func WithEndpoints(endpoints []string) Option {
	return func(o *DiscoverOption) {
		o.endpoints = endpoints
	}
}

func WithDialTimeout(dialTimeout time.Duration) Option {
	return func(o *DiscoverOption) {
		o.dialTimeout = dialTimeout
	}
}

func WithUpdateFn(updateFn UpdateFn) Option {
	return func(o *DiscoverOption) {
		o.updateFn = updateFn
	}
}

func WithRemoveFn(removeFn RemoveFn) Option {
	return func(o *DiscoverOption) {
		o.removeFn = removeFn
	}
}
