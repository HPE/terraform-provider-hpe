// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package compute

import (
	"strings"
	"testing"
)

// The value of this change is the text a user reads when a lookup fails, so the
// tests pin the text: the scope explanation, the two routes to the answer, and
// the cluster hint appearing only when there is one to give.

func TestUnitResourcePoolNotFoundByNameMessage(t *testing.T) {
	t.Parallel()

	t.Run("without hint", func(t *testing.T) {
		t.Parallel()

		got := resourcePoolNotFoundByNameMessage("hvm-cluster-01", 1, "")

		for _, want := range []string{
			`No resource pool named "hvm-cluster-01" was found in cloud 1.`,
			"returns cloud-level pools only",
			"attached to the cluster, not the cloud, and is not listed",
			"hpe_morpheus_cluster data source",
			"(permissions.resource_pool.id)",
			"or look it up here by id.",
		} {
			if !strings.Contains(got, want) {
				t.Errorf("message missing %q\n--- got ---\n%s", want, got)
			}
		}

		if strings.Contains(got, "A cluster named") {
			t.Errorf("message must not contain a cluster hint when none was supplied\n--- got ---\n%s", got)
		}

		if !strings.HasSuffix(got, "by id.") {
			t.Errorf("message without a hint must end with the by-id sentence, got suffix %q", got[len(got)-20:])
		}
	})

	t.Run("with hint", func(t *testing.T) {
		t.Parallel()

		hint := clusterResourcePoolHint("hvm-cluster-01", 1, 7)
		got := resourcePoolNotFoundByNameMessage("hvm-cluster-01", 1, hint)

		if !strings.HasSuffix(got, hint) {
			t.Errorf("hint must be the final paragraph\n--- got ---\n%s", got)
		}

		if !strings.Contains(got, "by id.\n\nA cluster named") {
			t.Errorf("hint must be separated from the body by a blank line\n--- got ---\n%s", got)
		}
	})
}

func TestUnitClusterResourcePoolHint(t *testing.T) {
	t.Parallel()

	got := clusterResourcePoolHint("hvm-cluster-01", 1, 7)
	want := `A cluster named "hvm-cluster-01" exists in cloud 1; its provisioning pool has id 7. ` +
		`Set id = 7, or use the hpe_morpheus_cluster data source.`

	if got != want {
		t.Errorf("hint text\n want: %s\n  got: %s", want, got)
	}
}

func TestUnitResourcePoolMultipleMatchesMessage(t *testing.T) {
	t.Parallel()

	got := resourcePoolMultipleMatchesMessage("shared", 3, []int64{12, 40})
	want := `2 resource pools named "shared" were found in cloud 3 (ids 12, 40). Use id to select one.`

	if got != want {
		t.Errorf("message text\n want: %s\n  got: %s", want, got)
	}
}

func TestUnitResourcePoolNotFoundByIDMessage(t *testing.T) {
	t.Parallel()

	got := resourcePoolNotFoundByIDMessage(999999, 1)
	want := "Resource pool id 999999 was not found in cloud 1."

	if got != want {
		t.Errorf("message text\n want: %s\n  got: %s", want, got)
	}
}
