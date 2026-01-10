package vpp

import "github.com/NikolayStepanov/RapidVPP/internal/domain"

func ToVersionResponse(v domain.Version) VersionResponse {
	return VersionResponse{
		Version:   v.Version,
		BuildDate: v.BuildDate,
		BuildDir:  v.BuildDir,
	}
}
