package ip

type VRFResponse struct {
	ID         uint32 `json:"id"`
	Name       string `json:"name"`
	IPv4       bool   `json:"ipv4"`
	IPv6       bool   `json:"ipv6"`
	RouteCount int    `json:"route_count"`
}
