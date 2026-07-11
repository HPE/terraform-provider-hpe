// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package client

type CreateServicemapLink struct {
	// Mandatory
	Id     string  `json:"id"`
	Parent *Parent `json:"parent,omitempty"`
}

type ReadServicemapLink struct {
	// Mandatory
	Results []ReadServicemap `json:"results"`
}
