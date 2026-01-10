package vpp

type VersionResponse struct {
	Version   string `json:"version"`
	BuildDate string `json:"build_date"`
	BuildDir  string `json:"build_dir"`
}
