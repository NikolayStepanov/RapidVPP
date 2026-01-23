package acl

type ACLResponse struct {
	ID    uint32          `json:"id"`
	Name  string          `json:"name"`
	Rules []RulesResponse `json:"rules"`
}
type ListACLResponse struct {
	ACLs []ACLResponse `json:"acls"`
}
type RulesResponse struct {
	Action        uint8        `json:"action"`
	Proto         uint8        `json:"proto"`
	Src           IPWithPrefix `json:"src"`
	Dst           IPWithPrefix `json:"dst"`
	SrcPortLow    uint16       `json:"src_port_low"`
	SrcPortHigh   uint16       `json:"src_port_high"`
	DstPortLow    uint16       `json:"dst_port_low"`
	DstPortHigh   uint16       `json:"dst_port_high"`
	TCPFlagsMask  uint8        `json:"tcp_flags_mask"`
	TCPFlagsValue uint8        `json:"tcp_flags_value"`
}
