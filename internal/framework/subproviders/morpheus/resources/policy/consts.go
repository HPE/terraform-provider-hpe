// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package policy

const (
	// AssociatedResourceTypeGlobal is the resource type for global policies
	AssociatedResourceTypeGlobal = "Global"
	// AssociatedResourceTypeCloud is the resource type for cloud-scoped policies
	AssociatedResourceTypeCloud = "Cloud"
	// AssociatedResourceTypeGroup is the resource type for group-scoped policies
	AssociatedResourceTypeGroup = "Group"
	// AssociatedResourceTypeUser is the resource type for user-scoped policies
	AssociatedResourceTypeUser = "User"
	// AssociatedResourceTypeRole is the resource type for role-scoped policies
	AssociatedResourceTypeRole = "Role"
	// AssociatedResourceTypeNetwork is the resource type for network-scoped policies
	AssociatedResourceTypeNetwork = "Network"
	// AssociatedResourceTypePlan is the resource type for plan-scoped policies
	AssociatedResourceTypePlan = "Plan"
	// AssociatedResourceTypeLabel is the resource type for label-scoped policies
	AssociatedResourceTypeLabel = "Label"
)

const (
	// WorkflowTypeFlow indicates a workflow of type "flow"
	WorkflowTypeFlow = "flow"
	// WorkflowTypeWorkflow indicates a workflow of type "workflow"
	WorkflowTypeWorkflow = "workflow"
)
