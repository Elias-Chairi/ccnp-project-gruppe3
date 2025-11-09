package util

import (
	"fmt"
	"net"
	"strings"
)

// IPList is a custom flag type that holds a list of IP addresses
type IPList []net.IP

// Set parses a comma-separated string of IP addresses and populates the IPList
func (ipl *IPList) Set(value string) error {
	// Split the input string by commas
	ipStrs := strings.Split(value, ",")
	// Trim spaces from each element
	for i := range ipStrs {
		ipStrs[i] = strings.TrimSpace(ipStrs[i])
	}
	// Parse each element into net.IP and append to the list
	for _, ipStr := range ipStrs {
		if ipStr == "" {
			continue
		}
		ip := net.ParseIP(ipStr)
		if ip == nil {
			return fmt.Errorf("invalid IP address: %s", ipStr)
		}
		*ipl = append(*ipl, ip)
	}
	return nil
}

// String returns a comma-separated string representation of the IPList
func (ipl *IPList) String() string {
	var ipStrs []string
	for _, ip := range *ipl {
		ipStrs = append(ipStrs, ip.String())
	}
	return strings.Join(ipStrs, ", ")
}
