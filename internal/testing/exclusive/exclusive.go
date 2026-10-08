package exclusive

import (
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"testing"
	"time"
)

const (
	acquiringWait = 20 * time.Minute
	acquiringPoll = 20 * time.Millisecond
)

var held sync.Map

func Machine(t *testing.T, resource string) {
	t.Helper()

	if _, already := held.LoadOrStore(t, struct{}{}); already {
		return
	}

	path := filepath.Join(os.TempDir(), "socrates-exclusive-"+resource+".lock")
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o666)
	if err != nil {
		t.Fatalf("opening the lock for %s: %v", resource, err)
	}

	deadline := time.Now().Add(acquiringWait)
	for {
		err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
		if err == nil {
			break
		}
		if err != syscall.EWOULDBLOCK {
			_ = file.Close()
			t.Fatalf("locking %s: %v", path, err)
		}
		if time.Now().After(deadline) {
			_ = file.Close()
			t.Fatalf("%s is still held by another test after %s", resource, acquiringWait)
		}
		time.Sleep(acquiringPoll)
	}

	t.Cleanup(func() {
		held.Delete(t)
		_ = syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
		_ = file.Close()
	})
}
