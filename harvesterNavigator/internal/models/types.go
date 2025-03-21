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
}

type ReplicaInfo struct {
	Name             string
	SpecVolumeName   string
	OwnerName        string
	HasMismatch      string
	LonghornNode     string
	LonghornDiskUUID string
}
