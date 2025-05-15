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
	VMStatus        string
	PVCStatus       string
	PrintableStatus string
	VMStatusReason  string
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
