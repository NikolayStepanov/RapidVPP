package interfaces

type SetInterfaceStateRequests struct {
	AdminUp bool `json:"admin_up"`
}

type AddIPRequest struct {
	Address string `json:"address"`
	Prefix  uint8  `json:"prefix"`
}
