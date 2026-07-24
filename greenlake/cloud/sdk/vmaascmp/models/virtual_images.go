// (C) Copyright 2021 Hewlett Packard Enterprise Development LP
package models

type VirtualImages struct {
	VirtualImages []VirtualImage `json:"virtualImages"`
}

type GetSpecificVirtualImage struct {
	VirtualImages VirtualImage `json:"virtualImage"`
}

type VirtualImage struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	Description interface{} `json:"description"`
	OwnerID     int         `json:"ownerId"`
	Tenant      struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"tenant"`
	ImageType            string      `json:"imageType"`
	UserUploaded         bool        `json:"userUploaded"`
	UserDefined          bool        `json:"userDefined"`
	SystemImage          bool        `json:"systemImage"`
	IsCloudInit          bool        `json:"isCloudInit"`
	SSHUsername          interface{} `json:"sshUsername"`
	SSHPassword          interface{} `json:"sshPassword"`
	SSHKey               interface{} `json:"sshKey"`
	OsType               interface{} `json:"osType"`
	MinRAM               int         `json:"minRam"`
	MinRAMGB             float64     `json:"minRamGB"`
	MinDisk              int64       `json:"minDisk"`
	MinDiskGB            int         `json:"minDiskGB"`
	RawSize              interface{} `json:"rawSize"`
	RawSizeGB            interface{} `json:"rawSizeGB"`
	TrialVersion         bool        `json:"trialVersion"`
	VirtioSupported      bool        `json:"virtioSupported"`
	IsAutoJoinDomain     bool        `json:"isAutoJoinDomain"`
	VMToolsInstalled     bool        `json:"vmToolsInstalled"`
	InstallAgent         bool        `json:"installAgent"`
	IsForceCustomization bool        `json:"isForceCustomization"`
	IsSysprep            bool        `json:"isSysprep"`
	FipsEnabled          bool        `json:"fipsEnabled"`
	UserData             interface{} `json:"userData"`
	ConsoleKeymap        interface{} `json:"consoleKeymap"`
	StorageProvider      interface{} `json:"storageProvider"`
	ExternalID           string      `json:"externalId"`
	Visibility           string      `json:"visibility"`
	Config               struct{}    `json:"config"`
	Volumes              []struct {
		Name       string `json:"name"`
		MaxStorage int64  `json:"maxStorage"`
		RawSize    int64  `json:"rawSize"`
		Size       int    `json:"size"`
		RootVolume bool   `json:"rootVolume"`
		Resizeable bool   `json:"resizeable"`
	} `json:"volumes"`
	StorageControllers []struct {
		Name string `json:"name"`
		Type struct {
			ID   int    `json:"id"`
			Code string `json:"code"`
			Name string `json:"name"`
		} `json:"type"`
		MaxDevices         int `json:"maxDevices"`
		ReservedUnitNumber int `json:"reservedUnitNumber"`
	} `json:"storageControllers"`
	NetworkInterfaces []interface{} `json:"networkInterfaces"`
	DateCreated       string        `json:"dateCreated"`
	LastUpdated       string        `json:"lastUpdated"`
}
