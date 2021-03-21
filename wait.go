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
	"context"
	"os"
)

// lockWait will attempt to get a lock, and if that fails, spawn the provided
// goroutine, and cancel the context when the lock is aquired.
func lockWait(file *os.File, fLockNB, fLock func(*os.File) error, f Waiter) error {
	err := fLockNB(file)
	if err == nil {
		return nil
	}

	if err != errWouldBlock {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	go f(ctx, file)

	err = fLock(file)
	cancel()
	return err
}

// Waiter is the function type expected by the flock.*LockWait functions.
//
// This function is started during a blocking operation, and the context will be
// cancelled when the lock is finally aquired.
type Waiter func(context.Context, *os.File)

// LockWaiter will attempt to hold the *Exclusive* process flock for a specific
// file. If this call blocks, the Water will be started, with a context that is
// cancelled when the lock is *aquired*.
//
// This is used when a goroutine needs to be running whilst the lock is being
// aquired.
func LockWaiter(file *os.File, waiter Waiter) error {
	return lockWait(file, lockNB, Lock, waiter)
}

// LockSharedWaiter will attempt to hold the *Shared* process flock for a specific
// file. If this call blocks, the Water will be started, with a context that is
// cancelled when the lock is *aquired*.
//
// This is used when a goroutine needs to be running whilst the lock is being
// aquired.
func LockSharedWaiter(file *os.File, waiter Waiter) error {
	return lockWait(file, lockSharedNB, LockShared, waiter)
}

// vim: foldmethod=marker
