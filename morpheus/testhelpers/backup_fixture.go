// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package testhelpers

import "testing"

// The backup fixture: the objects the backup tests attach backups to, named by
// environment variables and resolved by Terraform data sources inside each
// test configuration. Nothing here is an id: ids are allocated per appliance,
// a name resolves wherever the object exists.
//
// A backup of an instance takes a backup type from the types its cloud offers
// (a KVM/HVM instance: kvm); a host backup is a file backup of a path on a
// compute server, which Morpheus accepts for any managed server. Two storage
// buckets are needed: the update tests move a backup from one to the other.
// There are no defaults: a guessed name turns a test that should skip into one
// that fails against the API.
const (
	// EnvBackupInstanceName names an instance the backup tests may attach a
	// backup to. Nothing runs unless the job the test creates is scheduled,
	// which it is not; the record is created, read, moved and destroyed.
	EnvBackupInstanceName = "TF_VAR_testacc_morpheus_backup_instance_name"
	// EnvBackupServerName names a managed compute server (a guest VM's, or a
	// host's) for the host backup tests, which back up a file path on it.
	EnvBackupServerName = "TF_VAR_testacc_morpheus_backup_server_name"
	// EnvBackupBucketName and EnvBackupOtherBucketName name two storage
	// buckets a backup can target; the update tests move a backup from the
	// first to the other.
	EnvBackupBucketName      = "TF_VAR_testacc_morpheus_backup_bucket_name"
	EnvBackupOtherBucketName = "TF_VAR_testacc_morpheus_backup_other_bucket_name"
)

// BackupInstanceName returns the name of the instance to back up, skipping
// the test when EnvBackupInstanceName is unset.
func BackupInstanceName(t *testing.T) string {
	t.Helper()

	return envOrSkip(t, EnvBackupInstanceName, "an instance to attach a backup to")
}

// BackupServerName returns the name of the compute server to back a file up
// from, skipping the test when EnvBackupServerName is unset.
func BackupServerName(t *testing.T) string {
	t.Helper()

	return envOrSkip(t, EnvBackupServerName, "a managed compute server to back a file up from")
}

// BackupBucketName returns the name of the storage bucket the backup tests
// target, skipping the test when EnvBackupBucketName is unset.
func BackupBucketName(t *testing.T) string {
	t.Helper()

	return envOrSkip(t, EnvBackupBucketName, "a storage bucket for backups")
}

// BackupBucketNames returns the names of the two storage buckets the update
// tests target, skipping the test when either variable is unset.
func BackupBucketNames(t *testing.T) (first, other string) {
	t.Helper()

	first = BackupBucketName(t)
	other = envOrSkip(t, EnvBackupOtherBucketName, "a second storage bucket for backups")

	return first, other
}

// BackupFixtureBlock renders the data sources that resolve the backup
// fixture by name: the instance (BackupInstanceIDRef), the compute server
// (BackupServerIDRef) and the two storage buckets (BackupBucketIDRef,
// BackupOtherBucketIDRef). Every argument may be empty, in which case that
// data source is left out. The names are quoted for HCL (HCLQuote), so one
// with a quote, a backslash or a `${` in it is looked up as written.
func BackupFixtureBlock(instanceName, serverName, bucketName, otherBucketName string) string {
	var out string
	if instanceName != "" {
		out += "data \"hpe_morpheus_instance\" \"backup\" {\n  name = " + HCLQuote(instanceName) + "\n}\n\n"
	}
	if serverName != "" {
		out += "data \"hpe_morpheus_compute_server\" \"backup\" {\n  name = " + HCLQuote(serverName) + "\n}\n\n"
	}
	if bucketName != "" {
		out += "data \"hpe_morpheus_storage_bucket\" \"backup\" {\n  name = " + HCLQuote(bucketName) + "\n}\n\n"
	}
	if otherBucketName != "" {
		out += "data \"hpe_morpheus_storage_bucket\" \"backup_other\" {\n  name = " + HCLQuote(otherBucketName) + "\n}\n\n"
	}

	return out
}

// The HCL expressions for the ids BackupFixtureBlock resolves.
const (
	BackupInstanceIDRef         = "data.hpe_morpheus_instance.backup.id"
	BackupServerIDRef           = "data.hpe_morpheus_compute_server.backup.id"
	BackupBucketIDRef           = "data.hpe_morpheus_storage_bucket.backup.id"
	BackupOtherBucketIDRef      = "data.hpe_morpheus_storage_bucket.backup_other.id"
	BackupInstanceDataSource    = "data.hpe_morpheus_instance.backup"
	BackupServerDataSource      = "data.hpe_morpheus_compute_server.backup"
	BackupBucketDataSource      = "data.hpe_morpheus_storage_bucket.backup"
	BackupOtherBucketDataSource = "data.hpe_morpheus_storage_bucket.backup_other"
)
