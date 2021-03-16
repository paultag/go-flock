# pault.ag/go/flock

[![Go Reference](https://pkg.go.dev/badge/pault.ag/go/flock.svg)](https://pkg.go.dev/pault.ag/go/flock)
[![Go Report Card](https://goreportcard.com/badge/pault.ag/go/flock)](https://goreportcard.com/report/pault.ag/go/flock)

flock allows processes to take out advisory locks on files using the
`flock(2)` syscalls.

The returned `flock.Locker` is a `sync.Locker`, except it also contains the
ability to take out a shared lock (like a sync.RWLock). The API diverged since
`flock.Locker.Unlock` will release either the Exclusive or Shared lock.
