package vmi

import (
	"context"
	"encoding/json"
	"fmt"

	types "github.com/rk280392/harvesterNavigator/internal/models"
	"k8s.io/client-go/kubernetes"
)

func FetchVMIDetails(client *kubernetes.Clientset, name, absPath, namespace, resource string) (map[string]interface{}, error) {
	vm, err := client.RESTClient().Get().
		AbsPath(absPath).
		Namespace(namespace).
		Name(name).
		Resource(resource).
		Do(context.Background()).Raw()

	if err != nil {
		return nil, fmt.Errorf("failed to get VMIDetails: %s", err)
	}

	var vmiData map[string]interface{}
	err = json.Unmarshal(vm, &vmiData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall VMI details: %s", err)
	}
	return vmiData, nil
}

func ParseVMIData(vmiData map[string]interface{}) ([]types.VMIInfo, error) {
	var vmiInfos []types.VMIInfo
	if vmiData == nil {
		return vmiInfos, fmt.Errorf("no VMI data available")
	}

	vmiMetadata, ok := vmiData["metadata"].(map[string]interface{})
	vmiName := vmiMetadata["name"].(string)
	vmiStatus, ok := vmiData["status"].(map[string]interface{})
	if !ok {
		return vmiInfos, fmt.Errorf("status not found or invalid format")
	}

	phase, ok := vmiStatus["phase"].(string)
	if !ok {
		phase = "Unknown" // Default value if not found
	}

	nodeName, _ := vmiStatus["nodeName"].(string)
	vmiInfo := types.VMIInfo{
		Phase:       phase,
		NodeName:    nodeName,
		ActivePods:  make(map[string]string),
		GuestOSInfo: types.GuestOSInfo{},
		Interfaces:  []types.Interfaces{},
		Name:        vmiName,
	}

	if activePodsRaw, ok := vmiStatus["activePods"].(map[string]interface{}); ok {
		for podUID, nodeNameVal := range activePodsRaw {
			if nodeNameStr, ok := nodeNameVal.(string); ok {
				vmiInfo.ActivePods[podUID] = nodeNameStr
			}
		}
	}
	if guestOSInfoRaw, ok := vmiStatus["guestOSInfo"].(map[string]interface{}); ok {
		if name, ok := guestOSInfoRaw["name"].(string); ok {
			vmiInfo.GuestOSInfo.Name = name
		}
		if version, ok := guestOSInfoRaw["version"].(string); ok {
			vmiInfo.GuestOSInfo.Version = version
		}
		if prettyName, ok := guestOSInfoRaw["prettyName"].(string); ok {
			vmiInfo.GuestOSInfo.PrettyName = prettyName
		}
	}

	fmt.Println("vmiStatus", vmiStatus)

	if interfacesRaw, ok := vmiStatus["interfaces"].(map[string]interface{}); ok {
		fmt.Println("interfacesRaw", interfacesRaw)
		for _, ifaceRaw := range interfacesRaw {
			ifaceMap, ok := ifaceRaw.(map[string]interface{})
			if !ok {
				continue
			}
			iface := types.Interfaces{}
			if ipAddress, ok := ifaceMap["ipAddress"].(string); ok {
				iface.IpAddress = ipAddress
			}
			if mac, ok := ifaceMap["mac"].(string); ok {
				iface.Mac = mac
			}
			vmiInfo.Interfaces = append(vmiInfo.Interfaces, iface)
		}
	}

	vmiInfos = append(vmiInfos, vmiInfo)
	return vmiInfos, nil
}
