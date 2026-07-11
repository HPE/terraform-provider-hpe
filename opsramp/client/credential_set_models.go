// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package client

// CredentialSet represents a credential set in OpsRamp
type CredentialSet struct {
	UniqueId                        string `json:"uniqueId,omitempty"`
	Name                            string `json:"name"`
	Secure                          bool   `json:"secure"`
	Port                            int    `json:"port,omitempty"`
	SnmpVersion                     string `json:"snmpVersion,omitempty"`
	Description                     string `json:"description,omitempty"`
	AutoEnableMode                  bool   `json:"autoEnableMode"`
	Universal                       bool   `json:"universal"`
	SpSecure                        bool   `json:"spSecure"`
	SpPort                          int    `json:"spPort,omitempty"`
	TimeoutMs                       int    `json:"timeoutMs,omitempty"`
	CredentialType                  string `json:"credentialType"`
	DomainName                      string `json:"domainName,omitempty"`
	UserName                        string `json:"userName,omitempty"`
	Password                        string `json:"password,omitempty"`
	TransportType                   string `json:"transportType,omitempty"`
	Community                       string `json:"community,omitempty"`
	SpUserName                      string `json:"spUserName,omitempty"`
	SpPassword                      string `json:"spPassword,omitempty"`
	SpAuthScope                     string `json:"spAuthScope,omitempty"`
	FileAuthScope                   string `json:"fileAuthScope,omitempty"`
	EsxUserName                     string `json:"esxUserName,omitempty"`
	EsxPassword                     string `json:"esxPassword,omitempty"`
	SpNameSpace                     string `json:"spNameSpace,omitempty"`
	AuthProtocol                    string `json:"authProtocol,omitempty"`
	EncryptPassword                 string `json:"encryptPassword,omitempty"`
	SnmpContext                     string `json:"snmpContext,omitempty"`
	SecurityLevel                   string `json:"securityLevel,omitempty"`
	SecurityName                    string `json:"securityName,omitempty"`
	ApiEndPoint                     string `json:"apiEndPoint,omitempty"`
	AccountId                       string `json:"accountId,omitempty"`
	AccountName                     string `json:"accountName,omitempty"`
	AccountKey                      string `json:"accountKey,omitempty"`
	ManagementCertificate           string `json:"managementCertificate,omitempty"`
	ManagementCertificatePassphrase string `json:"managementCertificatePassphrase,omitempty"`
	SshCredentialType               string `json:"sshCredentialType,omitempty"`
	CollectorType                   string `json:"collectorType,omitempty"`
	EnablePassword                  string `json:"enablePassword,omitempty"`
}
