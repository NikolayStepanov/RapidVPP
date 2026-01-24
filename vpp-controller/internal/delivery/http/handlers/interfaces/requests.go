package interfaces

type SetInterfaceStateRequests struct {
	AdminUp bool `json:"admin_up"`
}

type AddIPRequest struct {
	Address string `json:"address"`
	Prefix  uint8  `json:"prefix"`
}

type AttachACLRequest struct {
	AclId     uint32 `json:"acl_id"`
	Direction uint8  `json:"direction"`
}

type DetachACLRequest struct {
	AclId     uint32 `json:"acl_id"`
	Direction uint8  `json:"direction"`
}
