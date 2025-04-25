package main

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// go test -race -v -count=1 .

type RWMutex struct {
	mx       sync.Mutex
	isWriter atomic.Bool
	readers  atomic.Int32
}

func (m *RWMutex) Lock() {
	// если какой-то писатель захватил замок, то дождемся отпускания замка
	m.mx.Lock()

	// теперь мы единственный писатель, замок наш
	m.isWriter.Store(true)

	// но нам надо точно удостовериться, что все читатели отпустили свои read замки
	// потому что согласно семантике rwlock надо дождаться пока читатели доделают свою работу
	// вот собственно и ждем
	m.waitForReaders()

	// вот тут у нас гарантирован инвариант, что все читатели дали зеленый свет (отпустили свои замки)
	// и писатель только один - текущий и он закрыл замок
}

func (m *RWMutex) Unlock() {
	m.isWriter.Store(false)
	m.mx.Unlock()
}

func (m *RWMutex) RLock() {
	// если есть писатель, то надо его дождаться
	m.waitForWriter()

	m.readers.Add(1)
}

func (m *RWMutex) RUnlock() {
	m.readers.Add(-1)
}

const spinLockCounter = 100

func (m *RWMutex) waitForReaders() {
	// подождем активно немного времени, вдруг успеем дождаться
	for i := 0; i < spinLockCounter; i++ {
		if m.readers.Load() == 0 {
			break
		}
	}

	// иначе будем отдавать управление планировщику пока не дождемся инварианта
	for m.readers.Load() != 0 {
		runtime.Gosched()
	}
}

func (m *RWMutex) waitForWriter() {
	// подождем активно немного времени, вдруг успеем дождаться
	for i := 0; i < spinLockCounter; i++ {
		if !m.isWriter.Load() {
			break
		}
	}

	// иначе будем отдавать управление планировщику пока не дождемся инварианта
	for m.isWriter.Load() {
		runtime.Gosched()
	}
}

func TestRWMutexWithWriter(t *testing.T) {
	var mutex RWMutex
	mutex.Lock() // writer

	var mutualExlusionWithWriter atomic.Bool
	mutualExlusionWithWriter.Store(true)
	var mutualExlusionWithReader atomic.Bool
	mutualExlusionWithReader.Store(true)

	go func() {
		mutex.Lock() // another writer
		mutualExlusionWithWriter.Store(false)
	}()

	go func() {
		mutex.RLock() // another reader
		mutualExlusionWithReader.Store(false)
	}()

	time.Sleep(time.Second)
	assert.True(t, mutualExlusionWithWriter.Load())
	assert.True(t, mutualExlusionWithReader.Load())
}

func TestRWMutexWithReaders(t *testing.T) {
	var mutex RWMutex
	mutex.RLock() // reader

	var mutualExlusionWithWriter atomic.Bool
	mutualExlusionWithWriter.Store(true)

	go func() {
		mutex.Lock() // another writer
		mutualExlusionWithWriter.Store(false)
	}()

	time.Sleep(time.Second)
	assert.True(t, mutualExlusionWithWriter.Load())
}

func TestRWMutexMultipleReaders(t *testing.T) {
	var mutex RWMutex
	mutex.RLock() // reader

	var readersCount atomic.Int32
	readersCount.Add(1)

	go func() {
		mutex.RLock() // another reader
		readersCount.Add(1)
	}()

	go func() {
		mutex.RLock() // another reader
		readersCount.Add(1)
	}()

	time.Sleep(time.Second)
	assert.Equal(t, int32(3), readersCount.Load())
}

func TestRWMutexWithWriterPriority(t *testing.T) {
	var mutex RWMutex
	mutex.RLock() // reader

	var mutualExlusionWithWriter atomic.Bool
	mutualExlusionWithWriter.Store(true)
	var readersCount atomic.Int32
	readersCount.Add(1)

	go func() {
		mutex.Lock() // another writer is waiting for reader
		mutualExlusionWithWriter.Store(false)
	}()

	time.Sleep(time.Second)

	go func() {
		mutex.RLock() // another reader is waiting for a higher priority writer
		readersCount.Add(1)
	}()

	go func() {
		mutex.RLock() // another reader is waiting for a higher priority writer
		readersCount.Add(1)
	}()

	time.Sleep(time.Second)

	assert.True(t, mutualExlusionWithWriter.Load())
	assert.Equal(t, int32(1), readersCount.Load())
}
