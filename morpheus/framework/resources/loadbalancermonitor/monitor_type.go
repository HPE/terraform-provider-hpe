// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancermonitor

// Load balancer type codes as defined in the Morpheus seed data.
const (
	LBTypeNsxT = "nsx-t"
	LBTypeNsxV = "nsx-v"
)

// MonitorType constants are the canonical lowercase values accepted by the
// provider's monitor_type attribute. They are the union of the NSX-T display
// names (lowercased) and the NSX-V option source values.
const (
	MonitorTypeDNS     = "dns"
	MonitorTypeHTTP    = "http"
	MonitorTypeHTTPS   = "https"
	MonitorTypeICMP    = "icmp"
	MonitorTypeLDAP    = "ldap"
	MonitorTypeMSSQL   = "mssql"
	MonitorTypePassive = "passive"
	MonitorTypeTCP     = "tcp"
	MonitorTypeUDP     = "udp"
)

// nsxtMonitorTypes maps canonical lowercase monitor type names to the NSX-T
// profile values expected by the Morpheus API.
var nsxtMonitorTypes = map[string]string{
	MonitorTypeHTTP:    "LBHttpMonitorProfile",
	MonitorTypeHTTPS:   "LBHttpsMonitorProfile",
	MonitorTypeICMP:    "LBIcmpMonitorProfile",
	MonitorTypePassive: "LBPassiveMonitorProfile",
	MonitorTypeTCP:     "LBTcpMonitorProfile",
	MonitorTypeUDP:     "LBUdpMonitorProfile",
}

// nsxvMonitorTypes maps canonical lowercase monitor type names to the NSX-V
// values expected by the Morpheus API (lowercase, same as the canonical form).
var nsxvMonitorTypes = map[string]string{
	MonitorTypeDNS:   "dns",
	MonitorTypeHTTP:  "http",
	MonitorTypeHTTPS: "https",
	MonitorTypeLDAP:  "ldap",
	MonitorTypeMSSQL: "mssql",
	MonitorTypeTCP:   "tcp",
	MonitorTypeUDP:   "udp",
}

// nsxtMonitorTypesReverse maps NSX-T profile values back to canonical
// lowercase monitor type names. Built from nsxtMonitorTypes.
var nsxtMonitorTypesReverse = func() map[string]string {
	m := make(map[string]string, len(nsxtMonitorTypes))
	for canonical, profile := range nsxtMonitorTypes {
		m[profile] = canonical
	}

	return m
}()

// canonicalMonitorType converts an API monitor type value to the canonical
// lowercase form based on the load balancer type code. For NSX-T the profile
// name (e.g. LBHttpMonitorProfile) is mapped back to the lowercase name
// (e.g. http). For NSX-V and other types the value is returned as-is since
// it is already in canonical form.
func canonicalMonitorType(apiValue, lbTypeCode string) string {
	if lbTypeCode == LBTypeNsxT {
		if canonical, ok := nsxtMonitorTypesReverse[apiValue]; ok {
			return canonical
		}
	}

	return apiValue
}
