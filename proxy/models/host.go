package models

// Host is an IPv4 + its associated port.
type Host struct {
	IP   [4]byte
	Port uint16
}
