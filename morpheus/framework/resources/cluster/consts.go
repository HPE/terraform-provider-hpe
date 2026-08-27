// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cluster

const (
	// Possible cluster type codes
	clusterTypeCodeMVM                = "mvm-cluster" // i.e. HVM
	clusterTypeCodeKubernetes         = "kubernetes-cluster"
	clusterTypeCodeExternalKubernetes = "external-kubernetes-cluster"
	clusterTypeCodeAKS                = "aks-cluster"
	clusterTypeCodeGKE                = "gke-cluster"
	clusterTypeCodeEKS                = "eks-cluster"
	clusterTypeCodeDocker             = "docker-cluster"

	// Possible cluster statuses
	clusterStatusOk             = "ok"
	clusterStatusCancelled      = "cancelled"
	clusterStatusDenied         = "denied"
	clusterStatusFailed         = "failed"
	clusterStatusDeprovisioned  = "deprovisioned"
	clusterStatusDeprovisioning = "deprovisioning"
	clusterStatusPending        = "pending"
	clusterStatusPendingRemoval = "pendingRemoval"
	clusterStatusProvisioning   = "provisioning"
	clusterStatusProvisioned    = "provisioned"
	clusterStatusRemoved        = "removed"
	clusterStatusRemoving       = "removing"
	clusterStatusRunning        = "running"
	clusterStatusStarting       = "starting"
	clusterStatusStopping       = "stopping"
	clusterStatusSuspended      = "suspended"
	clusterStatusSyncing        = "syncing"
	clusterStatusWarning        = "warning"
)
