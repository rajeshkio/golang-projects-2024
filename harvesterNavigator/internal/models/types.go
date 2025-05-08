package types

type VMInfo struct {
	Name           string
	ImageId        string
	PodName        string
	StorageClass   string
	ClaimNames     string
	VolumeName     string
	AttachmentInfo map[string]interface{}
	ReplicaInfo    []ReplicaInfo
	VMStatus       string
	PVCStatus      string
}

type ReplicaInfo struct {
	Name             string
	SpecVolumeName   string
	OwnerName        string
	HasMismatch      string
	LonghornNode     string
	LonghornDiskUUID string
}

type VolumeInfo struct {
	Name          string
	VolumeDetails map[string]interface{}
}
