package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type UserService struct {
	// not need to implement
	NotEmptyStruct bool
}
type MessageService struct {
	// not need to implement
	NotEmptyStruct bool
}

type Storage struct {
	Count int
}

var NewStorage func() *Storage

func init() {
	var count = 0
	NewStorage = func() *Storage {
		count++
		return &Storage{Count: count}
	}
}

type cons struct {
	fn     func() any
	isOnce bool // need to call fn only once
	cache  any
}

type Container struct {
	constructors map[string]cons
}

func NewContainer() *Container {
	return &Container{
		constructors: make(map[string]cons),
	}
}

func (c *Container) RegisterType(name string, constructor func() any) {
	c.constructors[name] = cons{fn: constructor}
}

func (c *Container) RegisterSingletonType(name string, constructor func() any) {
	c.constructors[name] = cons{
		fn:     constructor,
		isOnce: true,
	}
}

func (c *Container) Resolve(name string) (any, error) {
	cons, ok := c.constructors[name]
	if !ok {
		return nil, errors.New("constructor `%s` hasn't been registered")
	}

	if cons.isOnce && cons.cache != nil {
		return cons.cache, nil
	}

	result := cons.fn()
	if cons.isOnce {
		cons.cache = result
		c.constructors[name] = cons
	}

	return result, nil
}

func TestDIContainer(t *testing.T) {
	container := NewContainer()
	container.RegisterType("UserService", func() any {
		return &UserService{}
	})
	container.RegisterType("MessageService", func() any {
		return &MessageService{}
	})

	container.RegisterSingletonType("Storage", func() any {
		return NewStorage()
	})

	userService1, err := container.Resolve("UserService")
	assert.NoError(t, err)
	userService2, err := container.Resolve("UserService")
	assert.NoError(t, err)

	u1 := userService1.(*UserService)
	u2 := userService2.(*UserService)
	assert.False(t, u1 == u2)

	messageService, err := container.Resolve("MessageService")
	assert.NoError(t, err)
	assert.NotNil(t, messageService)

	paymentService, err := container.Resolve("PaymentService")
	assert.Error(t, err)
	assert.Nil(t, paymentService)

	maybeStorage, err := container.Resolve("Storage")
	assert.NoError(t, err)
	assert.NotNil(t, maybeStorage)
	storage, ok := maybeStorage.(*Storage)
	assert.True(t, ok)
	assert.Equal(t, 1, storage.Count)

	// call second time - got the same instance
	maybeStorage, err = container.Resolve("Storage")
	assert.NoError(t, err)
	assert.NotNil(t, maybeStorage)
	storage, ok = maybeStorage.(*Storage)
	assert.True(t, ok)
	assert.Equal(t, 1, storage.Count)
}
