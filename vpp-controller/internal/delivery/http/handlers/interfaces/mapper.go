package interfaces

import "github.com/NikolayStepanov/RapidVPP/internal/domain"

func ACLInterfaceListToDTO(aclInterfaceList domain.ACLInterfaceList) ACLInterfaceListResponses {
	return ACLInterfaceListResponses{
		InterfaceID: aclInterfaceList.InterfaceID,
		Count:       aclInterfaceList.Count,
		InputACLs:   aclInterfaceList.InputACLs,
		OutputACLs:  aclInterfaceList.OutputACLs,
	}
}
