// {{{ Copyright (c) Paul R. Tagliamonte <paultag@gmail.com> 2020-2021
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE. }}}

package flock

import (
	"log"
	"os"
	"syscall"
)

// New will create a new sync.Locker using `flock(2)` on a UNIX system.
func New(file *os.File) (*Locker, error) {
	return &Locker{
		fd: file,
	}, nil
}

// Locker will allow for the locking or unlocking of a flock, or file lock.
type Locker struct {
	fd *os.File
}

// Lock will attempt to aquire an exclusive lock on the file, which
// only a single process may hold this lock at a time.
func (l *Locker) Lock() {
	l.flock(syscall.LOCK_EX)
}

// SharedLock will attempt to aquire a shared lock on the file, which many
// processes may hold at the same time.
func (l *Locker) SharedLock() {
	l.flock(syscall.LOCK_SH)
}

// Unlock will release an flock held by this process.
func (l *Locker) Unlock() {
	l.flock(syscall.LOCK_UN)
}

func (l *Locker) flock(op int) {
	for {
		err := syscall.Flock(int(l.fd.Fd()), op)
		if err == nil {
			return
		}
		log.Printf("flock: error during operation: %s", err)
	}
}

// vim: foldmethod=marker
