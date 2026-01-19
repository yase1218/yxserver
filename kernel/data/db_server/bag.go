package db_server

import "sync"

type Bag struct {
	Data map[uint32]int64 `bson:"data"` // id:num
	lock sync.RWMutex     `bson:"-" json:"-"`
}

func NewBag() *Bag {
	return &Bag{
		Data: make(map[uint32]int64),
	}
}

func (b *Bag) Get(id uint32) int64 {
	b.lock.RLock()
	defer b.lock.RUnlock()

	return b.Data[id]
}

func (b *Bag) Set(id uint32, num int64) {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.Data[id] = num
}

func (b *Bag) Add(id uint32, num int64) {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.Data[id] += num
}

func (b *Bag) Del(id uint32) {
	b.lock.Lock()
	defer b.lock.Unlock()

	delete(b.Data, id)
}

func (b *Bag) GetAll() map[uint32]int64 {
	b.lock.RLock()
	defer b.lock.RUnlock()

	return b.Data
}
