package vm

import (
	"context"
	"encoding/json"
	"fmt"

	types "github.com/rk280392/harvesterNavigator/internal/models"
	"k8s.io/client-go/kubernetes"
)

func fetchVMData(client *kubernetes.Clientset, vmname, namespace string) (map[string]interface{}, error) {

	vm, err := client.RESTClient().Get().
		AbsPath("apis/kubevirt.io/v1").
		Namespace(namespace).
		Resource("virtualmachines").
		Name(vmname).
		Do(context.TODO()).Raw()
	if err != nil {
		return nil, err
	}
	var vmData map[string]interface{}
	if err := json.Unmarshal(vm, &vmData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal VM Data: %w", err)
	}
	return vmData, nil
}

func parseVMMetadata(vmData map[string]interface{}, vmInfo *types.VMInfo) (string, error) {
	metadata := vmData["metadata"].(map[string]interface{})
	//fmt.Println("VM Name:", name)
	//nameTest := metadata["name"].(string)
	//fmt.Println("VM name: ", nameTest)
	//fmt.Println("")

	annotations := metadata["annotations"].(map[string]interface{})
	volumeClaimTemplateInterface := annotations["harvesterhci.io/volumeClaimTemplates"].(string)
	return volumeClaimTemplateInterface, nil
}
func GetVMInfo(client *kubernetes.Clientset, vmname, namespace string) (*types.VMInfo, error) {
	vmData, err := fetchVMData(client, vmname, namespace)
	if err != nil {
		return nil, err
	}

	vmInfo := &types.VMInfo{
		Name: vmname,
	}
	parseVMMetadata(vmData, vmInfo)
}
