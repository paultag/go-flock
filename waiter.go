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
	"fmt"
	"os"
	"time"
)

// StdoutWaiter is a flock.Waiter that will write a message to os.Stdout
// when blocked and waiting for a lock.
//
// This will print out a message like "Waiting for a filesystem lock on
// 'filename' (time)", and update it by rewritting the line and clearing
// the line using ANSI escape codes.
//
// After the lock is acquired, this will write out a message saying the
// lock has been acquired, and return.
func StdoutWaiter(ctx context.Context, file *os.File) {
	start := time.Now()
	tick := time.NewTicker(time.Second)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			fmt.Fprintf(
				os.Stderr,
				"\x1b[2K\rWaiting for filesystem lock on %s (%s)",
				file.Name(),
				time.Now().Sub(start).Round(time.Second),
			)
			continue
		case <-ctx.Done():
			fmt.Fprintf(
				os.Stderr,
				"\x1b[2K\rLocked ringbuffer at %s\n",
				file.Name(),
			)
			return
		}
	}
}

// vim: foldmethod=marker
