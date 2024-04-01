package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/dashjay/gohlc/api/hlcv1"
	"github.com/dashjay/gohlc/pkg/persistence"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

const InvalidClock = math.MaxInt64

type HLCService struct {
	persistenceHandler       persistence.Interface
	persistentIntervalSecond int64

	allocatedTimestamp atomic.Int64

	lastSavedSuccess   atomic.Bool
	lastSavedTimestamp atomic.Int64

	LockOSThread bool

	mu   sync.Mutex
	cond *sync.Cond
	hlcv1.UnimplementedHCLServiceServer
}

func (h *HLCService) loadClockFromPersistence(ctx context.Context) error {
	clock, err := h.persistenceHandler.Load(ctx)
	if err != nil {
		if errors.Is(err, persistence.ClockNotFound) {
			return h.persistenceHandler.Persist(ctx, time.Now().UnixNano())
		}
		return err
	}
	h.lastSavedTimestamp.Store(clock)
	return nil
}

func (h *HLCService) doPersistentAndUpdateLastSaved(ctx context.Context) error {
	h.mu.Lock()
	base := _max(h.allocatedTimestamp.Load(), time.Now().UnixNano())
	save := base + h.persistentIntervalSecond*int64(time.Second)
	h.mu.Unlock()

	err := h.persistenceHandler.Persist(ctx, save)
	if err != nil {
		h.lastSavedSuccess.Store(false)
		return err
	}

	h.lastSavedSuccess.Store(true)
	h.lastSavedTimestamp.Store(save)
	h.cond.Broadcast()
	return nil
}

func (h *HLCService) doPersistentAndUpdateLastSavedLoop(ctx context.Context) {
	// lock this goroutine
	if h.LockOSThread {
		runtime.LockOSThread()
	}

	for {
		startTime := time.Now()
		timeout := time.Second * time.Duration(h.persistentIntervalSecond) / 3
		subCtx, cancel := context.WithTimeout(ctx, timeout)
		err := h.doPersistentAndUpdateLastSaved(subCtx)
		cancel()
		if err != nil {
			fmt.Fprintf(os.Stderr, "do persistent and update lastSaved clock error: %s", err)
		}
		intervalLeft := timeout - time.Since(startTime)
		if intervalLeft > 0 {
			time.Sleep(intervalLeft)
		}
	}
}

func NewHLCService(path string, intervalSec int64) *HLCService {
	ph, err := persistence.NewDisk(path)
	if err != nil {
		panic(err)
	}
	hlc := &HLCService{
		persistenceHandler:       ph,
		persistentIntervalSecond: intervalSec,
		allocatedTimestamp:       atomic.Int64{},
		lastSavedTimestamp:       atomic.Int64{},
		cond:                     sync.NewCond(&sync.Mutex{}),
	}
	return hlc
}

func (h *HLCService) Start() error {
	err := h.loadClockFromPersistence(context.Background())
	if err != nil {
		return err
	}
	err = h.doPersistentAndUpdateLastSaved(context.Background())
	if err != nil {
		return err
	}
	go h.doPersistentAndUpdateLastSavedLoop(context.Background())
	return nil
}

func (h *HLCService) allocateClock(count uint32) int64 {
	c64 := int64(count)
	loopCount := 0
	for {
		h.cond.L.Lock()
		expected := h.allocatedTimestamp.Load()
		desired := _max(expected+c64, time.Now().UnixNano())
		if desired >= h.lastSavedTimestamp.Load() {
			h.cond.Wait()
		}
		h.cond.L.Unlock()
		if h.allocatedTimestamp.CompareAndSwap(expected, desired) {
			return desired - c64 + 1
		}
		loopCount++
	}
}

func (h *HLCService) Get(_ context.Context, _ *emptypb.Empty) (*hlcv1.GetResp, error) {
	clk := h.allocateClock(1)
	if clk == InvalidClock {
		return nil, status.Error(codes.Internal, "get clock error")
	}
	return &hlcv1.GetResp{Clock: clk}, nil
}

func (h *HLCService) BatchGet(_ context.Context, req *hlcv1.BatchGetReq) (*hlcv1.BatchGetResp, error) {
	first := h.allocateClock(req.Count)
	if first == InvalidClock {
		return nil, status.Errorf(codes.Internal, "get clocks error")
	}
	resp := &hlcv1.BatchGetResp{
		Req:    req,
		First:  first,
		Clocks: nil,
	}
	if req.ReturnFirst {
		return resp, nil
	}
	clocks := make([]int64, req.Count)
	for i := uint32(0); i < req.Count; i++ {
		clocks[i] = first + int64(i)
	}
	resp.Clocks = clocks
	return resp, nil
}

func _max(a, b int64) int64 {
	if a >= b {
		return a
	}
	return b
}

var _ hlcv1.HCLServiceServer = (*HLCService)(nil)
