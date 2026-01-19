package worker_queue

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/v587-zyf/gc/utils"
	"sync"
	"testing"
	"time"
)

func TestWorkerQueue(t *testing.T) {
	var as = assert.New(t)

	t.Run("", func(t *testing.T) {
		var mu = &sync.Mutex{}
		var listA []int
		var listB []int

		var count = 10000
		var wg = &sync.WaitGroup{}
		wg.Add(count)
		if err := Init(context.Background(), WithMaxCount(100)); err != nil {
			t.Fatal(err)
		}
		for i := 0; i < count; i++ {
			listA = append(listA, i)

			v := i
			Push(func() {
				defer wg.Done()
				var latency = time.Duration(utils.AlphabetNumeric.Intn(100)) * time.Microsecond
				time.Sleep(latency)
				mu.Lock()
				listB = append(listB, v)
				mu.Unlock()
			})
		}
		wg.Wait()
		as.ElementsMatch(listA, listB)
	})
}
