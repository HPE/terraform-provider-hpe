// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package testhelpers

import (
	"os"
	"testing"
	"time"
)

// TestUnitReleaseLockLeavesAnotherHoldersLock pins that releasing a lock
// removes it only while it still carries the releasing holder's token: a
// lock taken over by another holder after going stale is left alone.
func TestUnitReleaseLockLeavesAnotherHoldersLock(t *testing.T) {
	t.Parallel()

	path := t.TempDir() + "/pool.lock"
	if err := os.WriteFile(path, []byte("other-holder"), 0o600); err != nil {
		t.Fatal(err)
	}

	releaseLock(path, "mine")

	if _, err := os.Stat(path); err != nil {
		t.Fatal("another holder's lock was removed")
	}

	if err := os.WriteFile(path, []byte("mine"), 0o600); err != nil {
		t.Fatal(err)
	}

	releaseLock(path, "mine")

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("own lock was not removed")
	}
}

// TestUnitKeepLockFreshBumpsModTime pins that a held lock's modification time
// is refreshed, so a long holder is never taken for a stale one.
func TestUnitKeepLockFreshBumpsModTime(t *testing.T) {
	t.Parallel()

	path := t.TempDir() + "/pool.lock"
	if err := os.WriteFile(path, []byte("mine"), 0o600); err != nil {
		t.Fatal(err)
	}

	old := time.Now().Add(-time.Hour)
	if err := os.Chtimes(path, old, old); err != nil {
		t.Fatal(err)
	}

	stop := make(chan struct{})
	go keepLockFresh(path, 20*time.Millisecond, stop)

	deadline := time.Now().Add(2 * time.Second)
	for {
		info, err := os.Stat(path)
		if err != nil {
			t.Fatal(err)
		}

		if info.ModTime().After(old.Add(30 * time.Minute)) {
			break
		}

		if time.Now().After(deadline) {
			t.Fatal("modification time was not refreshed")
		}

		time.Sleep(10 * time.Millisecond)
	}

	close(stop)
}

// TestUnitLockTokensDiffer pins that two tokens never collide, which is what
// lets a holder recognise its own lock.
func TestUnitLockTokensDiffer(t *testing.T) {
	t.Parallel()

	first, second := lockToken(), lockToken()
	if first == second {
		t.Fatal("two lock tokens were equal")
	}
}

// TestUnitTakeOverStaleLock drives the takeover through the function itself:
// a stale lock is removed; a fresh lock is left; a lock replaced with a fresh
// one while the takeover lock is held is left; a takeover lock held by another
// waiter blocks the takeover; and a stale takeover lock is cleared so the next
// attempt can proceed.
func TestUnitTakeOverStaleLock(t *testing.T) {
	t.Parallel()

	const stale = 20 * time.Minute

	write := func(t *testing.T, path, content string, age time.Duration) {
		t.Helper()

		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}

		then := time.Now().Add(-age)
		if err := os.Chtimes(path, then, then); err != nil {
			t.Fatal(err)
		}
	}

	exists := func(path string) bool {
		_, err := os.Stat(path)

		return err == nil
	}

	t.Run("stale lock is removed", func(t *testing.T) {
		t.Parallel()

		path := t.TempDir() + "/pool.lock"
		write(t, path, "old", time.Hour)

		if !takeOverStaleLock(path, stale) || exists(path) {
			t.Fatal("stale lock not taken over")
		}

		if exists(path + ".takeover") {
			t.Fatal("takeover lock left behind")
		}
	})

	t.Run("fresh lock is left", func(t *testing.T) {
		t.Parallel()

		path := t.TempDir() + "/pool.lock"
		write(t, path, "fresh", time.Minute)

		if takeOverStaleLock(path, stale) || !exists(path) {
			t.Fatal("fresh lock was taken over")
		}
	})

	// The re-check under the takeover lock is the point of the design, so it
	// is driven on its own, not in parallel: the hook is package state.
	t.Run("lock replaced under the takeover lock is left", func(t *testing.T) {
		path := t.TempDir() + "/pool.lock"
		write(t, path, "old", time.Hour)

		// Between this waiter taking the takeover lock and re-reading the
		// lock's age, another waiter's takeover completed and it created a
		// fresh lock at the path.
		afterTakeoverLock = func() {
			if err := os.Remove(path); err != nil {
				t.Fatal(err)
			}

			write(t, path, "fresh-holder", time.Second)
		}
		defer func() { afterTakeoverLock = nil }()

		if takeOverStaleLock(path, stale) {
			t.Fatal("a lock made fresh under the takeover lock was taken over")
		}

		content, err := os.ReadFile(path)
		if err != nil || string(content) != "fresh-holder" {
			t.Fatalf("fresh holder's lock disturbed: %q %v", content, err)
		}

		if exists(path + ".takeover") {
			t.Fatal("takeover lock left behind")
		}
	})

	t.Run("takeover in progress blocks", func(t *testing.T) {
		t.Parallel()

		path := t.TempDir() + "/pool.lock"
		write(t, path, "old", time.Hour)
		write(t, path+".takeover", "", time.Second)

		if takeOverStaleLock(path, stale) || !exists(path) {
			t.Fatal("took over while another waiter held the takeover lock")
		}

		if !exists(path + ".takeover") {
			t.Fatal("another waiter's fresh takeover lock was removed")
		}
	})

	t.Run("stale takeover lock is cleared", func(t *testing.T) {
		t.Parallel()

		path := t.TempDir() + "/pool.lock"
		write(t, path, "old", time.Hour)
		write(t, path+".takeover", "", 5*time.Minute)

		// First attempt clears the dead waiter's takeover lock; the second
		// takes the stale lock over.
		if takeOverStaleLock(path, stale) {
			t.Fatal("took over on the same attempt that cleared the stale takeover lock")
		}

		if exists(path + ".takeover") {
			t.Fatal("stale takeover lock not cleared")
		}

		if !takeOverStaleLock(path, stale) || exists(path) {
			t.Fatal("stale lock not taken over after the takeover lock was cleared")
		}
	})
}
