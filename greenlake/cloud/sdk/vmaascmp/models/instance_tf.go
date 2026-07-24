// (C) Copyright 2021-2024 Hewlett Packard Enterprise Development LP

// This file contains struct models related to terraform, but not the API
package models

type TFInstance struct {
	History    []GetInstanceHistoryProcesses `tf:"history,computed"`
	Containers []GetInstanceContainer        `tf:"containers,computed"`
	Snapshot   TFInstanceSnapshot            `tf:"snapshot,computed,sub"`
	Volume     []TFInstanceVolume            `tf:"volume,computed"`
	Network    []TFInstanceNetwork           `tf:"network,computed"`
	Status     string                        `tf:"status,computed"`
}

type TFInstanceSnapshot struct {
	ID               int    `tf:"id"`
	Name             string `tf:"name"`
	Description      string `tf:"description"`
	IsSnapshotExists bool   `tf:"is_snapshot_exists"`
}

type TFInstanceVolume struct {
	Name        string `tf:"name"`
	Size        int    `tf:"size"`
	DatastoreID string `tf:"datastore_id"`
	ID          int    `tf:"id"`
	Root        bool   `tf:"root"`
	StorageType int    `tf:"storage_type"`
	Controller  string `tf:"controller"`
}

type TFInstanceNetwork struct {
	ID          int    `tf:"id"`
	InterfaceID int    `tf:"interface_id"`
	IsPrimary   bool   `tf:"is_primary"`
	InternalID  int    `tf:"internal_id"`
	Name        string `tf:"name"`
}
