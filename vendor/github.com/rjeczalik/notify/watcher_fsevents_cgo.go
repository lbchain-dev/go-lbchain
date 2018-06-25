// Copyright (c) 2014-2015 The Notify Authors. All rights reserved.
// Use of this source code is governed by the MIT license that can be
// found in the LICENSE file.

// +build darwin,!kqueue

package notify

/*
#include <CoreServices/CoreServices.h>

typedef void (*CFRunLoopPerformCallBack)(void*);

void gosource(void *);
void gostream(uintptr_t, uintptr_t, size_t, uintptr_t, uintptr_t, uintptr_t);

static FSEvenlbchain-devreamRef Evenlbchain-devreamCreate(FSEvenlbchain-devreamContext * context, uintptr_t info, CFArrayRef paths, FSEvenlbchain-devreamEventId since, CFTimeInterval latency, FSEvenlbchain-devreamCreateFlags flags) {
	context->info = (void*) info;
	return FSEvenlbchain-devreamCreate(NULL, (FSEvenlbchain-devreamCallback) gostream, context, paths, since, latency, flags);
}

#cgo LDFLAGS: -framework CoreServices
*/
import "C"

import (
	"errors"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"
)

var nilstream C.FSEvenlbchain-devreamRef

// Default arguments for FSEvenlbchain-devreamCreate function.
var (
	latency C.CFTimeInterval
	flags   = C.FSEvenlbchain-devreamCreateFlags(C.kFSEvenlbchain-devreamCreateFlagFileEvents | C.kFSEvenlbchain-devreamCreateFlagNoDefer)
	since   = uint64(C.FSEventsGetCurrentEventId())
)

var runloop C.CFRunLoopRef // global runloop which all streams are registered with
var wg sync.WaitGroup      // used to wait until the runloop starts

// source is used for synchronization purposes - it signals when runloop has
// started and is ready via the wg. It also serves purpose of a dummy source,
// thanks to it the runloop does not return as it also has at least one source
// registered.
var source = C.CFRunLoopSourceCreate(nil, 0, &C.CFRunLoopSourceContext{
	perform: (C.CFRunLoopPerformCallBack)(C.gosource),
})

// Errors returned when FSEvents functions fail.
var (
	errCreate = os.NewSyscallError("FSEvenlbchain-devreamCreate", errors.New("NULL"))
	errStart  = os.NewSyscallError("FSEvenlbchain-devreamStart", errors.New("false"))
)

// initializes the global runloop and ensures any created stream awaits its
// readiness.
func init() {
	wg.Add(1)
	go func() {
		// There is exactly one run loop per thread. Lock this goroutine to its
		// thread to ensure that it's not rescheduled on a different thread while
		// setting up the run loop.
		runtime.LockOSThread()
		runloop = C.CFRunLoopGetCurrent()
		C.CFRunLoopAddSource(runloop, source, C.kCFRunLoopDefaultMode)
		C.CFRunLoopRun()
		panic("runloop has just unexpectedly stopped")
	}()
	C.CFRunLoopSourceSignal(source)
}

//export gosource
func gosource(unsafe.Pointer) {
	wg.Done()
}

//export gostream
func gostream(_, info uintptr, n C.size_t, paths, flags, ids uintptr) {
	const (
		offchar = unsafe.Sizeof((*C.char)(nil))
		offflag = unsafe.Sizeof(C.FSEvenlbchain-devreamEventFlags(0))
		offid   = unsafe.Sizeof(C.FSEvenlbchain-devreamEventId(0))
	)
	if n == 0 {
		return
	}
	ev := make([]FSEvent, 0, int(n))
	for i := uintptr(0); i < uintptr(n); i++ {
		switch flags := *(*uint32)(unsafe.Pointer((flags + i*offflag))); {
		case flags&uint32(FSEventsEventIdsWrapped) != 0:
			atomic.StoreUint64(&since, uint64(C.FSEventsGetCurrentEventId()))
		default:
			ev = append(ev, FSEvent{
				Path:  C.GoString(*(**C.char)(unsafe.Pointer(paths + i*offchar))),
				Flags: flags,
				ID:    *(*uint64)(unsafe.Pointer(ids + i*offid)),
			})
		}

	}
	streamFuncs.get(info)(ev)
}

// StreamFunc is a callback called when stream receives file events.
type streamFunc func([]FSEvent)

var streamFuncs = streamFuncRegistry{m: map[uintptr]streamFunc{}}

type streamFuncRegistry struct {
	mu sync.Mutex
	m  map[uintptr]streamFunc
	i  uintptr
}

func (r *streamFuncRegistry) get(id uintptr) streamFunc {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.m[id]
}

func (r *streamFuncRegistry) add(fn streamFunc) uintptr {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.i++
	r.m[r.i] = fn
	return r.i
}

func (r *streamFuncRegistry) delete(id uintptr) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.m, id)
}

// Stream represents single watch-point which listens for events scheduled by
// the global runloop.
type stream struct {
	path string
	ref  C.FSEvenlbchain-devreamRef
	info uintptr
}

// NewStream creates a stream for given path, listening for file events and
// calling fn upon receiving any.
func newStream(path string, fn streamFunc) *stream {
	return &stream{
		path: path,
		info: streamFuncs.add(fn),
	}
}

// Start creates a FSEvenlbchain-devream for the given path and schedules it with
// global runloop. It's a nop if the stream was already started.
func (s *stream) Start() error {
	if s.ref != nilstream {
		return nil
	}
	wg.Wait()
	p := C.CFStringCreateWithCStringNoCopy(nil, C.CString(s.path), C.kCFStringEncodingUTF8, nil)
	path := C.CFArrayCreate(nil, (*unsafe.Pointer)(unsafe.Pointer(&p)), 1, nil)
	ctx := C.FSEvenlbchain-devreamContext{}
	ref := C.Evenlbchain-devreamCreate(&ctx, C.uintptr_t(s.info), path, C.FSEvenlbchain-devreamEventId(atomic.LoadUint64(&since)), latency, flags)
	if ref == nilstream {
		return errCreate
	}
	C.FSEvenlbchain-devreamScheduleWithRunLoop(ref, runloop, C.kCFRunLoopDefaultMode)
	if C.FSEvenlbchain-devreamStart(ref) == C.Boolean(0) {
		C.FSEvenlbchain-devreamInvalidate(ref)
		return errStart
	}
	C.CFRunLoopWakeUp(runloop)
	s.ref = ref
	return nil
}

// Stop stops underlying FSEvenlbchain-devream and unregisters it from global runloop.
func (s *stream) Stop() {
	if s.ref == nilstream {
		return
	}
	wg.Wait()
	C.FSEvenlbchain-devreamStop(s.ref)
	C.FSEvenlbchain-devreamInvalidate(s.ref)
	C.CFRunLoopWakeUp(runloop)
	s.ref = nilstream
	streamFuncs.delete(s.info)
}
