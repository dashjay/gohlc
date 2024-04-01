package persistence_test

import (
	"context"
	"github.com/dashjay/gohlc/pkg/persistence"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
	"time"
)

func TestDisk(t *testing.T) {
	fPath := filepath.Join(t.TempDir(), "temp-file")
	d, err := persistence.NewDisk(fPath)
	assert.Nil(t, err)

	t.Run("rw", func(t *testing.T) {
		_, err = d.Load(context.TODO())
		assert.ErrorIs(t, err, persistence.ClockNotFound)
		clock := time.Now().UnixNano()
		err = d.Persist(context.TODO(), clock)
		assert.Nil(t, err)
		getClock, err := d.Load(context.TODO())
		assert.Nil(t, err)
		assert.Equal(t, clock, getClock)
	})

	t.Run("r-with-deadline", func(t *testing.T) {
		ctx, cancel := context.WithDeadline(context.Background(), time.Now())
		defer cancel()
		_, err = d.Load(ctx)
		assert.Nil(t, err)
	})
}
