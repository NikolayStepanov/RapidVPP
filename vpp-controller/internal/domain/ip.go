package domain

import "net"

// Route represents a routing table entry
type Route struct {
	// Destination network with prefix length
	Dst IPWithPrefix
	// Virtual Routing and Forwarding table ID
	// 0 represents the default VRF (global routing table)
	VRF uint32
	// Next hop information for traffic forwarding
	// Multiple next hops enable ECMP (Equal-Cost Multi-Path) routing
	// Empty slice indicates a drop route
	NextHops []NextHop
}

// NextHop defines a forwarding path for a route
type NextHop struct {
	// Destination IP address of the next hop
	// Required for regular routes, ignored for drop routes
	IP net.IP
	// Output interface index for packet forwarding
	// Set to 0xFFFFFFFF for attached/connected routes
	IfIndex uint32
	// Weight for load balancing across multiple next hops
	// Used in ECMP (Equal-Cost Multi-Path) scenarios
	// Higher values indicate preferred paths
	Weight uint8
	// When true, packets matching this route are dropped
	// In this case, IP and IfIndex fields are typically ignored
	Drop bool
}

// VRF Virtual Routing and Forwarding
type VRF struct {
	ID   uint32
	Name string
}
