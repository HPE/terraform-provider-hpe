// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package testhelpers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// LockFloatingIpPool serialises the tests that allocate from the SDN floating
// IP pool across packages (each package is its own `go test` process, so a
// package-level mutex cannot do it). Morpheus allocates a pool's lowest free
// address, so two concurrent allocations race for the same address on the
// controller; the loser's InstanceIP is left behind and poisons that address
// for every later allocation (MORPH-19206). The lock is released when the
// test finishes.
//
// The lock is a file created with O_EXCL, which is atomic on every platform
// the provider builds for (flock is not available on Windows), holding a
// token unique to the holder. While held, the holder refreshes the file's
// modification time, so a long test is never mistaken for a crashed one. A
// lock whose modification time is older than the stale age is one whose
// holder has gone; it is taken over under a second exclusive file, the
// takeover lock, so that only one waiter removes it — and that waiter checks
// the lock is still stale once it holds the takeover lock, so a lock another
// waiter has just created in its place is left alone. Nothing depends on
// inode numbers or on a rename. On release the lock is removed only if it
// still holds the releasing test's token.
func LockFloatingIpPool(t *testing.T, poolID string) {
	t.Helper()

	const (
		pollInterval = time.Second
		heartbeat    = time.Minute
		waitAtMost   = 30 * time.Minute
	)

	path := filepath.Join(os.TempDir(), "tf-acc-sdn-fip-pool-"+poolID+".lock")
	token := lockToken()
	deadline := time.Now().Add(waitAtMost)

	for {
		f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
		if err == nil {
			_, werr := f.WriteString(token)
			_ = f.Close()

			if werr != nil {
				_ = os.Remove(path)
				t.Fatalf("floating IP pool lock: %s", werr)
			}

			stop := make(chan struct{})
			go keepLockFresh(path, heartbeat, stop)

			t.Cleanup(func() {
				close(stop)
				releaseLock(path, token)
			})

			return
		}

		if !os.IsExist(err) {
			t.Fatalf("floating IP pool lock: %s", err)
		}

		if takeOverStaleLock(path, lockStaleAfter) {
			// The stale lock is gone; try the exclusive create again at once.
			continue
		}

		if time.Now().After(deadline) {
			t.Fatalf("floating IP pool lock %s: still held after %s", path, waitAtMost)
		}

		time.Sleep(pollInterval)
	}
}

// afterTakeoverLock, when set, runs once the takeover lock is held and before
// the lock's age is re-read. Tests use it to interleave another waiter's
// takeover at that point; production leaves it nil.
var afterTakeoverLock func()

const (
	// lockStaleAfter is how long a lock may go without a heartbeat before it
	// is taken for a crashed holder's.
	lockStaleAfter = 20 * time.Minute
	// takeoverStaleAfter bounds a takeover: a takeover lock older than this
	// belongs to a waiter that died mid-takeover and is removed.
	takeoverStaleAfter = time.Minute
)

// takeOverStaleLock removes the lock at path if it is stale, and reports
// whether it did. The removal happens under an exclusive takeover lock
// (path + ".takeover"), so of several waiters that find the lock stale only
// one acts; and the lock's age is re-read once the takeover lock is held, so
// a lock that was replaced in the meantime — the previous taker's fresh lock
// — is recognised as fresh and left in place. A takeover lock that is itself
// stale is removed, so a waiter that died mid-takeover does not block the
// rest for ever.
func takeOverStaleLock(path string, staleAfter time.Duration) bool {
	info, err := os.Stat(path)
	if err != nil || time.Since(info.ModTime()) <= staleAfter {
		return false
	}

	takeover := path + ".takeover"

	tf, err := os.OpenFile(takeover, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		if os.IsExist(err) {
			if ti, serr := os.Stat(takeover); serr == nil && time.Since(ti.ModTime()) > takeoverStaleAfter {
				_ = os.Remove(takeover)
			}
		}

		return false
	}

	_ = tf.Close()
	defer func() { _ = os.Remove(takeover) }()

	if afterTakeoverLock != nil {
		afterTakeoverLock()
	}

	// Under the takeover lock: is the lock still the stale one?
	info, err = os.Stat(path)
	if err != nil {
		// Already gone: the path is free for the caller's create.
		return true
	}

	if time.Since(info.ModTime()) <= staleAfter {
		return false
	}

	return os.Remove(path) == nil
}

// keepLockFresh bumps the lock file's modification time every interval until
// stopped, so a held lock never looks stale.
func keepLockFresh(path string, interval time.Duration, stop <-chan struct{}) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			return
		case now := <-ticker.C:
			_ = os.Chtimes(path, now, now)
		}
	}
}

// releaseLock removes the lock file if it still carries the token, i.e. it is
// still this holder's; a lock taken over after going stale is left alone.
func releaseLock(path, token string) {
	content, err := os.ReadFile(path)
	if err != nil || string(content) != token {
		return
	}

	_ = os.Remove(path)
}

func lockToken() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%d-%d", os.Getpid(), time.Now().UnixNano())
	}

	return fmt.Sprintf("%d-%s", os.Getpid(), hex.EncodeToString(b))
}
