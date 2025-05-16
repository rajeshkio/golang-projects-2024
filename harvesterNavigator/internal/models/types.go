package types

type VMInfo struct {
	Name            string
	ImageId         string
	PodName         string
	StorageClass    string
	ClaimNames      string
	VolumeName      string
	AttachmentInfo  map[string]interface{}
	ReplicaInfo     []ReplicaInfo
	EngineInfo      []EngineInfo
	PodInfo         []PodInfo
	VMStatus        string
	PVCStatus       string
	PrintableStatus string
	VMStatusReason  string
	VMIInfo         []VMIInfo
	MissingResource string
}

type PodInfo struct {
	Name   string
	VMI    string
	NodeID string
	Status string
}

type VMIInfo struct {
	ActivePods  map[string]string
	GuestOSInfo GuestOSInfo
	Interfaces  []Interfaces
	NodeName    string
	Phase       string
	Name        string
}

type Interfaces struct {
	IpAddress string
	Mac       string
}
type GuestOSInfo struct {
	KernelRelease string
	KernelVersion string
	Machine       string
	Name          string
	PrettyName    string
	Version       string
}

type ReplicaInfo struct {
	Name           string
	SpecVolumeName string
	OwnerRefName   string
	NodeID         string
	Active         bool
	EngineName     string
	CurrentState   string
	Started        bool
}

type EngineInfo struct {
	Active       bool
	CurrentState string
	Started      bool
	NodeID       string
	Snapshots    map[string]SnapshotInfo
	Name         string
}
type SnapshotInfo struct {
	Name        string
	Parent      string
	Created     string
	Size        string
	UserCreated bool
	Removed     bool
	Children    map[string]bool
	Labels      map[string]string
}

type VolumeInfo struct {
	Name          string
	VolumeDetails map[string]interface{}
}
