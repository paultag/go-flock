package flock_test

import (
	"io/ioutil"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"pault.ag/go/flock"
)

func BenchmarkFlockLock(b *testing.B) {
	f, err := ioutil.TempFile("", "")
	assert.NoError(b, err)
	defer os.Remove(f.Name())

	locker, err := flock.New(f)
	assert.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		locker.Lock()
		locker.Unlock()
	}
}

func BenchmarkFlockSharedLock(b *testing.B) {
	f, err := ioutil.TempFile("", "")
	assert.NoError(b, err)
	defer os.Remove(f.Name())

	locker, err := flock.New(f)
	assert.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		locker.SharedLock()
		locker.Unlock()
	}
}
