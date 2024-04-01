package service_test

import (
	"context"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/dashjay/gohlc/api/hlcv1"
	"github.com/dashjay/gohlc/pkg/service"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestHLC(t *testing.T) {
	tempFile := filepath.Join(t.TempDir(), "temp-file")
	hlc := service.NewHLCService(tempFile, 1)
	hlc.Start()
	var query int64
	var count int64

	allocateOne := func() int64 {
		resp, err := hlc.Get(context.TODO(), &emptypb.Empty{})
		assert.Nil(t, err)
		assert.NotEqual(t, service.InvalidClock, resp.Clock)
		return resp.Clock
	}

	allocateBatch := func(in uint32) int64 {
		resp, err := hlc.BatchGet(context.TODO(), &hlcv1.BatchGetReq{
			Count:       in,
			ReturnFirst: true,
		})
		assert.Nil(t, err)
		assert.NotEqual(t, service.InvalidClock, resp.First)
		return resp.First
	}

	start := time.Now()
	wg := sync.WaitGroup{}
	for i := 0; i < 200; i++ {
		wg.Add(1)
		times := 100000
		go func() {
			defer wg.Done()
			for j := 0; j < times; j++ {
				_ = allocateOne()
			}
			atomic.AddInt64(&query, int64(times))
			atomic.AddInt64(&count, int64(times))
		}()
	}
	wg.Wait()
	t.Logf("finished %d query & get %d clocks in %.2f sec", query, count, time.Since(start).Seconds())

	query = 0
	count = 0

	start = time.Now()
	wg = sync.WaitGroup{}
	for i := 0; i < 200; i++ {
		wg.Add(1)
		times := 100000
		per := 50
		go func() {
			defer wg.Done()
			for j := 0; j < times; j++ {
				_ = allocateBatch(uint32(per))
			}
			atomic.AddInt64(&query, int64(times))
			atomic.AddInt64(&count, int64(times*per))
		}()
	}
	wg.Wait()
	t.Logf("finished %d query & get %d clocks in %.2f sec", query, count, time.Since(start).Seconds())
}
