package interfaces

type ACLInterfaceListResponses struct {
	InterfaceID uint32   `json:"interface_id"`
	Count       uint8    `json:"count"`
	InputACLs   []uint32 `json:"input_acls"`
	OutputACLs  []uint32 `json:"output_acls"`
}
