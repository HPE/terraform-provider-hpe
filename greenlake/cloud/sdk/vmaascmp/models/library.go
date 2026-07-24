// (C) Copyright 2021-2024 Hewlett Packard Enterprise Development LP

package models

type LayoutRespBody struct {
	ID           int `json:"id"`
	Instancetype struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Code string `json:"code"`
	} `json:"instanceType"`
	Account struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"account"`
	Code                     string      `json:"code"`
	Name                     string      `json:"name"`
	Instanceversion          string      `json:"instanceVersion"`
	Description              interface{} `json:"description"`
	Creatable                bool        `json:"creatable"`
	Memoryrequirement        int         `json:"memoryRequirement"`
	Sortorder                int         `json:"sortOrder"`
	Supportsconverttomanaged bool        `json:"supportsConvertToManaged"`
	Provisiontype            struct {
		ID                    int         `json:"id"`
		Name                  string      `json:"name"`
		Description           interface{} `json:"description"`
		Code                  string      `json:"code"`
		Aclenabled            bool        `json:"aclEnabled"`
		Multitenant           bool        `json:"multiTenant"`
		Managed               bool        `json:"managed"`
		Hostnetwork           bool        `json:"hostNetwork"`
		Customsupported       bool        `json:"customSupported"`
		Mapports              bool        `json:"mapPorts"`
		Exportserver          bool        `json:"exportServer"`
		Viewset               interface{} `json:"viewSet"`
		Servertype            string      `json:"serverType"`
		Hosttype              string      `json:"hostType"`
		Addvolumes            bool        `json:"addVolumes"`
		Hasvolumes            bool        `json:"hasVolumes"`
		Hasdatastore          bool        `json:"hasDatastore"`
		Hasnetworks           bool        `json:"hasNetworks"`
		Maxnetworks           int         `json:"maxNetworks"`
		Customizevolume       bool        `json:"customizeVolume"`
		Rootdiskcustomizable  bool        `json:"rootDiskCustomizable"`
		Rootdisksizeknown     bool        `json:"rootDiskSizeKnown"`
		Rootdiskresizable     bool        `json:"rootDiskResizable"`
		Lvmsupported          bool        `json:"lvmSupported"`
		Hostdiskmode          string      `json:"hostDiskMode"`
		Mindisk               int         `json:"minDisk"`
		Maxdisk               interface{} `json:"maxDisk"`
		Resizecopiesvolumes   bool        `json:"resizeCopiesVolumes"`
		Supportsautodatastore bool        `json:"supportsAutoDatastore"`
		Haszonepools          bool        `json:"hasZonePools"`
		Hassecuritygroups     bool        `json:"hasSecurityGroups"`
		Hasparameters         interface{} `json:"hasParameters"`
		Canenforcetags        bool        `json:"canEnforceTags"`
		Disablerootdatastore  bool        `json:"disableRootDatastore"`
		Hassnapshots          bool        `json:"hasSnapshots"`
		Optiontypes           []struct {
			ID                    int         `json:"id"`
			Name                  string      `json:"name"`
			Description           interface{} `json:"description"`
			Code                  string      `json:"code"`
			Fieldname             string      `json:"fieldName"`
			Fieldlabel            string      `json:"fieldLabel"`
			Fieldcode             string      `json:"fieldCode"`
			Fieldcontext          string      `json:"fieldContext"`
			Fieldgroup            string      `json:"fieldGroup"`
			Fieldclass            interface{} `json:"fieldClass"`
			Fieldaddon            interface{} `json:"fieldAddOn"`
			Fieldcomponent        interface{} `json:"fieldComponent"`
			Fieldinput            interface{} `json:"fieldInput"`
			Placeholder           interface{} `json:"placeHolder"`
			Verifypattern         interface{} `json:"verifyPattern"`
			Helpblock             string      `json:"helpBlock"`
			Helpblockfieldcode    interface{} `json:"helpBlockFieldCode"`
			Defaultvalue          interface{} `json:"defaultValue"`
			Optionsource          interface{} `json:"optionSource"`
			Optionlist            interface{} `json:"optionList"`
			Type                  string      `json:"type"`
			Advanced              bool        `json:"advanced"`
			Required              bool        `json:"required"`
			Exportmeta            bool        `json:"exportMeta"`
			Editable              bool        `json:"editable"`
			Creatable             bool        `json:"creatable"`
			Config                struct{}    `json:"config"`
			Displayorder          int         `json:"displayOrder"`
			Wrapperclass          interface{} `json:"wrapperClass"`
			Enabled               bool        `json:"enabled"`
			Noblank               bool        `json:"noBlank"`
			Dependsoncode         interface{} `json:"dependsOnCode"`
			Visibleoncode         interface{} `json:"visibleOnCode"`
			Requireoncode         interface{} `json:"requireOnCode"`
			Contextualdefault     bool        `json:"contextualDefault"`
			Displayvalueondetails interface{} `json:"displayValueOnDetails"`
		} `json:"optionTypes"`
		Customoptiontypes []struct {
			ID                    int         `json:"id"`
			Name                  string      `json:"name"`
			Description           interface{} `json:"description"`
			Code                  string      `json:"code"`
			Fieldname             string      `json:"fieldName"`
			Fieldlabel            string      `json:"fieldLabel"`
			Fieldcode             string      `json:"fieldCode"`
			Fieldcontext          string      `json:"fieldContext"`
			Fieldgroup            string      `json:"fieldGroup"`
			Fieldclass            interface{} `json:"fieldClass"`
			Fieldaddon            interface{} `json:"fieldAddOn"`
			Fieldcomponent        interface{} `json:"fieldComponent"`
			Fieldinput            interface{} `json:"fieldInput"`
			Placeholder           interface{} `json:"placeHolder"`
			Verifypattern         interface{} `json:"verifyPattern"`
			Helpblock             string      `json:"helpBlock"`
			Helpblockfieldcode    interface{} `json:"helpBlockFieldCode"`
			Defaultvalue          string      `json:"defaultValue"`
			Optionsource          interface{} `json:"optionSource"`
			Optionlist            interface{} `json:"optionList"`
			Type                  string      `json:"type"`
			Advanced              bool        `json:"advanced"`
			Required              bool        `json:"required"`
			Exportmeta            bool        `json:"exportMeta"`
			Editable              bool        `json:"editable"`
			Creatable             bool        `json:"creatable"`
			Config                struct{}    `json:"config"`
			Displayorder          int         `json:"displayOrder"`
			Wrapperclass          interface{} `json:"wrapperClass"`
			Enabled               bool        `json:"enabled"`
			Noblank               bool        `json:"noBlank"`
			Dependsoncode         interface{} `json:"dependsOnCode"`
			Visibleoncode         interface{} `json:"visibleOnCode"`
			Requireoncode         interface{} `json:"requireOnCode"`
			Contextualdefault     bool        `json:"contextualDefault"`
			Displayvalueondetails interface{} `json:"displayValueOnDetails"`
		} `json:"customOptionTypes"`
		Networktypes []struct {
			ID           int    `json:"id"`
			Name         string `json:"name"`
			Displayorder int    `json:"displayOrder"`
			Enabled      bool   `json:"enabled"`
			Defaulttype  bool   `json:"defaultType"`
			Externalid   string `json:"externalId"`
			Code         string `json:"code"`
		} `json:"networkTypes"`
		Storagetypes []struct {
			ID                int         `json:"id"`
			Code              string      `json:"code"`
			Name              string      `json:"name"`
			Displayorder      int         `json:"displayOrder"`
			Defaulttype       bool        `json:"defaultType"`
			Customlabel       bool        `json:"customLabel"`
			Customsize        bool        `json:"customSize"`
			Customsizeoptions interface{} `json:"customSizeOptions"`
		} `json:"storageTypes"`
		Rootstoragetypes []struct {
			ID                int         `json:"id"`
			Code              string      `json:"code"`
			Name              string      `json:"name"`
			Displayorder      int         `json:"displayOrder"`
			Defaulttype       bool        `json:"defaultType"`
			Customlabel       bool        `json:"customLabel"`
			Customsize        bool        `json:"customSize"`
			Customsizeoptions interface{} `json:"customSizeOptions"`
		} `json:"rootStorageTypes"`
		Controllertypes []struct {
			ID           int    `json:"id"`
			Name         string `json:"name"`
			Displayorder int    `json:"displayOrder"`
			Category     string `json:"category"`
			Enabled      bool   `json:"enabled"`
			Creatable    bool   `json:"creatable"`
			Maxdevices   int    `json:"maxDevices"`
		} `json:"controllerTypes"`
	} `json:"provisionType"`
	Tasksets       []interface{} `json:"taskSets"`
	Containertypes []struct {
		ID      int `json:"id"`
		Account struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"account"`
		Name             string `json:"name"`
		Shortname        string `json:"shortName"`
		Code             string `json:"code"`
		Containerversion string `json:"containerVersion"`
		Provisiontype    struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Code string `json:"code"`
		} `json:"provisionType"`
		Virtualimage struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"virtualImage"`
		Category interface{} `json:"category"`
		Config   struct {
			Extraoptions struct{} `json:"extraOptions"`
			Migrated     bool     `json:"migrated"`
		} `json:"config"`
		Containerports       []interface{} `json:"containerPorts"`
		Containerscripts     []interface{} `json:"containerScripts"`
		Containertemplates   []interface{} `json:"containerTemplates"`
		Environmentvariables []interface{} `json:"environmentVariables"`
	} `json:"containerTypes"`
	Mounts               []interface{} `json:"mounts"`
	Ports                []interface{} `json:"ports"`
	Optiontypes          []interface{} `json:"optionTypes"`
	Environmentvariables []interface{} `json:"environmentVariables"`
	Spectemplates        []interface{} `json:"specTemplates"`
	Permissions          struct {
		Resourcepermissions struct {
			Defaultstore  bool `json:"defaultStore"`
			Allplans      bool `json:"allPlans"`
			Defaulttarget bool `json:"defaultTarget"`
			Canmanage     bool `json:"canManage"`
			All           bool `json:"all"`
			Account       struct {
				ID int `json:"id"`
			} `json:"account"`
			Sites []interface{} `json:"sites"`
			Plans []interface{} `json:"plans"`
		} `json:"resourcePermissions"`
	} `json:"permissions"`
}

type LayoutsResp struct {
	InstanceTypeLayouts []LayoutRespBody `json:"instanceTypeLayouts"`
}

type GetSpecificLayout struct {
	InstanceTypeLayouts LayoutRespBody `json:"instanceTypeLayout"`
}

type InstanceTypeRespBody struct {
	ID      int `json:"id"`
	Account struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"account"`
	Name                string      `json:"name"`
	Code                string      `json:"code"`
	Description         interface{} `json:"description"`
	Provisiontypecode   string      `json:"provisionTypeCode"`
	Category            string      `json:"category"`
	Active              bool        `json:"active"`
	Environmentprefix   string      `json:"environmentPrefix"`
	Visibility          string      `json:"visibility"`
	Featured            bool        `json:"featured"`
	Versions            []string    `json:"versions"`
	Instancetypelayouts []struct {
		ID                int           `json:"id"`
		Name              string        `json:"name"`
		Provisiontypecode string        `json:"provisionTypeCode"`
		ProvisionType     ProvisionType `json:"provisionType"`
	} `json:"instanceTypeLayouts"`
}

type ProvisionType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}
type InstanceTypesResp struct {
	InstanceTypes []InstanceTypeRespBody `json:"instanceTypes"`
}
