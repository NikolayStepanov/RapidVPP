package ip

type SetInterfaceStateRequests struct {
	AdminUp bool `json:"admin_up"`
}

type AddIPRequest struct {
	Address string `json:"address"`
	Prefix  uint8  `json:"prefix"`
}

type AddRouteRequest struct {
	Destination string           `json:"destination"`
	VRF         uint32           `json:"vrf"`
	NextHops    []NextHopRequest `json:"next_hops"`
}

type NextHopRequest struct {
	IP      string `json:"ip,omitempty"`
	IfIndex uint32 `json:"if_index"`
	Weight  uint8  `json:"weight"`
	Drop    bool   `json:"drop"`
}

type CreateVRFRequest struct {
	ID   uint32 `json:"id"`
	Name string `json:"name"`
}
