package main

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Group struct {
	cancel context.CancelCauseFunc
	ctx    context.Context
	wg     sync.WaitGroup
}

func NewErrGroup(ctx context.Context) (*Group, context.Context) {
	newCtx, cancel := context.WithCancelCause(ctx)
	g := &Group{
		cancel: cancel,
		ctx:    newCtx,
	}

	return g, newCtx
}

func (g *Group) Go(action func() error) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		err := action()
		if err != nil {
			g.cancel(err)
		}
	}()
}

func (g *Group) Wait() (err error) {
	g.wg.Wait()
	return context.Cause(g.ctx)
}

// go test -v homework_test.go

func TestErrGroupWithoutError(t *testing.T) {
	var counter atomic.Int32
	group, _ := NewErrGroup(context.Background())

	for i := 0; i < 5; i++ {
		group.Go(func() error {
			time.Sleep(time.Second)
			counter.Add(1)
			return nil
		})
	}

	err := group.Wait()
	assert.Equal(t, int32(5), counter.Load())
	assert.NoError(t, err)
}

func TestErrGroupWithError(t *testing.T) {
	var counter atomic.Int32
	group, ctx := NewErrGroup(context.Background())

	for i := 0; i < 5; i++ {
		group.Go(func() error {
			timer := time.NewTimer(time.Second)
			defer timer.Stop()

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-timer.C:
				counter.Add(1)
				return nil
			}
		})
	}

	group.Go(func() error {
		return errors.New("error")
	})

	err := group.Wait()
	assert.Equal(t, int32(0), counter.Load())
	assert.Error(t, err)
}
