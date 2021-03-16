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
	"fmt"
	"os"
	"syscall"
)

var (
	// errWouldBlock will be returned if an operation would have blocked waiting
	// for another process to release a lock.
	//
	// This will only be returned if LOCK_NB is set when requesting a lock.
	errWouldBlock error = fmt.Errorf("flock: operation would block")
)

// Lock will request an *Exclusive* flock on a specific file. Only a single
// process may hold the exclusive flock at a time.
func Lock(file *os.File) error {
	return flock(file, syscall.LOCK_EX)
}

// lockNB is unexported, since this function is hard to use correctly.
//
// In most cases, a user wants to do something while waiting for a lock (e.g.
// print a message to stdout, etc), and for that, we've exposed a more safe
// LockWaiter helper.
//
// lockNB will attempt to get an *Exclusive* lock, and return an EWOULDBLOCK
// if the operation would become blocking.
func lockNB(file *os.File) error {
	return flock(file, syscall.LOCK_EX|syscall.LOCK_NB)
}

// LockShared will request a *Shared* flock on a specific file. Many processes
// may hold a shared lock at the same time.
func LockShared(file *os.File) error {
	return flock(file, syscall.LOCK_SH)
}

// lockSharedNB is unexported, since this function is hard to use correctly.
//
// In most cases, a user wants to do something while waiting for a lock (e.g.
// print a message to stdout, etc), and for that, we've exposed a more safe
// LockSharedWaiter helper.
//
// lockNB will attempt to get an *Exclusive* lock, and return an EWOULDBLOCK
// if the operation would become blocking.
func lockSharedNB(file *os.File) error {
	return flock(file, syscall.LOCK_SH|syscall.LOCK_NB)
}

// Unlock will release either an Exclusive or Shared lock held by this process.
func Unlock(file *os.File) error {
	return flock(file, syscall.LOCK_UN)
}

func flock(file *os.File, op int) error {
	for {
		err := syscall.Flock(int(file.Fd()), op)
		if err == nil {
			return nil
		}
		switch err {
		case syscall.EBADF:
			return fmt.Errorf("flock: bad file descriptor")
		case syscall.EINTR:
			continue
		case syscall.EINVAL:
			return fmt.Errorf("flock: invalid operation provided")
		case syscall.ENOLCK:
			return fmt.Errorf("flock: kernel is out of memory for locks")
		case syscall.EWOULDBLOCK:
			return errWouldBlock
		default:
			return err
		}
	}
}

// vim: foldmethod=marker
