// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package affinitylock serialises read-modify-write cycles against a single
// affinity group.
//
// An affinity group's membership can only be written as a whole: there is no
// API for adding or removing one server. Anything that changes a group must
// therefore read the current membership, adjust it, and write all of it back.
// Two such cycles interleaving will lose one of the changes.
//
// That happens in ordinary use. Terraform applies resources of the same type in
// parallel, so two membership resources created together race each other; and
// the group resource itself must echo membership back on every update, so it
// races the membership resources too. All of them take the lock for the group
// they are touching, which is why this lives in a shared package rather than in
// any one resource.
//
// The lock covers a single provider process, which is where Terraform runs
// every resource in one apply. It cannot serialise two concurrent applies in
// separate processes; callers guard that case by asserting the result of their
// write afterwards.
package affinitylock

import (
	"fmt"
	"sync"
)

// locks holds one mutex per affinity group, created on first use.
var locks sync.Map

// Scope distinguishes the two kinds of affinity group. Their ids are only
// unique within their parent, so a cloud group and a cluster group can share an
// id while being unrelated.
type Scope string

const (
	// Cloud identifies a cloud affinity group.
	Cloud Scope = "cloud"
	// Cluster identifies a cluster affinity group.
	Cluster Scope = "cluster"
)

// Acquire locks the given affinity group and returns its release function.
//
// Callers should defer the returned function immediately:
//
//	defer affinitylock.Acquire(affinitylock.Cluster, clusterID, groupID)()
func Acquire(scope Scope, parentID, groupID int64) func() {
	key := fmt.Sprintf("%s:%d:%d", scope, parentID, groupID)

	value, _ := locks.LoadOrStore(key, &sync.Mutex{})

	mutex, ok := value.(*sync.Mutex)
	if !ok {
		// Unreachable: only *sync.Mutex is ever stored under these keys.
		return func() {}
	}

	mutex.Lock()

	return mutex.Unlock
}
