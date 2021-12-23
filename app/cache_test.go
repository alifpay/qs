package app

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestCache(t *testing.T) {

	for i := 0; i < 25; i++ {
		ix := time.Now().Add(20 * time.Second).Unix()
		str := strconv.FormatInt(ix, 10) + strconv.Itoa(i)
		addCache(str)
	}
	time.Sleep(1 * time.Second)
	for i := 25; i < 50; i++ {
		ix := time.Now().Add(14 * time.Second).Unix()
		str := strconv.FormatInt(ix, 10) + strconv.Itoa(i)
		addCache(str)
	}
	if qCache.unsorted {
		qCache.mu.Lock()
		qCache.items.Sort()
		qCache.unsorted = false
		qCache.mu.Unlock()
	}
	for _, v := range qCache.items {
		fmt.Println(v)
	}
}
