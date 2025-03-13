package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	metricsv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

// PodInfo holds pod information with its resource usage
type PodInfo struct {
	Name        string
	Namespace   string
	Node        string
	CPUUsage    int64 // in millicores
	MemoryUsage int64 // in MB
}

// truncateString truncates a string if it's longer than the specified length
func truncateString(str string, length int) string {
	if len(str) <= length {
		return str
	}
	return str[:length-3] + "..."
}

func main() {
	fmt.Println("Starting Kubernetes Pod Monitor...")

	// Handle the kubeconfig flag
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "path to your kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "path to your kubeconfig file")
	}

	// Add filter and sort options
	namespace := flag.String("namespace", "", "Filter by namespace (empty for all namespaces)")
	nodeFilter := flag.String("node", "", "Filter by node name (empty for all nodes)")
	sortBy := flag.String("sort", "cpu", "Sort by: cpu or memory")
	limit := flag.Int("limit", 0, "Limit the number of results (0 for no limit)")

	flag.Parse()

	// Get the current context from the raw config
	rawConfig, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: *kubeconfig},
		&clientcmd.ConfigOverrides{}).RawConfig()
	if err == nil {
		fmt.Printf("Using context: %s\n", rawConfig.CurrentContext)
	}

	// Build the config from the kubeconfig file
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		fmt.Printf("Error building kubeconfig: %v\n", err)
		os.Exit(1)
	}

	// Create a clientset (client)
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error creating kubernetes client: %v\n", err)
		os.Exit(1)
	}

	// Create the metrics client using the same config
	metricsClient, err := metricsv.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error creating metrics client: %v\n", err)
		os.Exit(1)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Check if metrics-server is available
	_, err = metricsClient.MetricsV1beta1().PodMetricses("").List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error connecting to metrics-server: %v\n", err)
		fmt.Println("Is metrics-server installed in your cluster?")
		os.Exit(1)
	}

	// Determine which namespaces to process
	var namespacesToProcess []string
	if *namespace != "" {
		namespacesToProcess = []string{*namespace}
	} else {
		// Get all namespaces
		namespaces, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
		if err != nil {
			fmt.Printf("Error listing namespaces: %v\n", err)
			os.Exit(1)
		}
		for _, ns := range namespaces.Items {
			namespacesToProcess = append(namespacesToProcess, ns.Name)
		}
	}

	// Create a slice to hold all pod info
	var allPodInfo []PodInfo

	// Process each namespace in our list
	for _, ns := range namespacesToProcess {
		// Get pods in this namespace
		pods, err := clientset.CoreV1().Pods(ns).List(ctx, metav1.ListOptions{})
		if err != nil {
			fmt.Printf("Error listing pods in namespace %s: %v\n", ns, err)
			continue
		}

		// Get pod metrics in this namespace
		podMetrics, err := metricsClient.MetricsV1beta1().PodMetricses(ns).List(ctx, metav1.ListOptions{})
		if err != nil {
			fmt.Printf("Error getting pod metrics in namespace %s: %v\n", ns, err)
			continue
		}

		// Create a map for quick lookup of pod metrics by pod name
		podMetricsMap := make(map[string]*metricsv1beta1.PodMetrics)
		for i := range podMetrics.Items {
			podMetricsMap[podMetrics.Items[i].Name] = &podMetrics.Items[i]
		}

		// Process pods in this namespace
		for _, pod := range pods.Items {
			// Skip pods that aren't running
			if pod.Status.Phase != "Running" {
				continue
			}

			// Skip if it doesn't match the node filter
			if *nodeFilter != "" && pod.Spec.NodeName != *nodeFilter {
				continue
			}

			// Lookup metrics for this pod
			metrics, exists := podMetricsMap[pod.Name]
			if !exists {
				// Pod doesn't have metrics yet
				continue
			}

			// Calculate total CPU and memory for the pod
			var cpuUsage int64
			var memoryUsage int64
			for _, container := range metrics.Containers {
				cpuUsage += container.Usage.Cpu().MilliValue()
				memoryUsage += container.Usage.Memory().Value() / (1024 * 1024) // Convert to MB
			}

			// Add to our collection
			allPodInfo = append(allPodInfo, PodInfo{
				Name:        pod.Name,
				Namespace:   pod.Namespace,
				Node:        pod.Spec.NodeName,
				CPUUsage:    cpuUsage,
				MemoryUsage: memoryUsage,
			})
		}
	}

	// Sort based on the sort option
	if *sortBy == "memory" {
		sort.Slice(allPodInfo, func(i, j int) bool {
			return allPodInfo[i].MemoryUsage > allPodInfo[j].MemoryUsage
		})
	} else {
		// Default sort by CPU
		sort.Slice(allPodInfo, func(i, j int) bool {
			return allPodInfo[i].CPUUsage > allPodInfo[j].CPUUsage
		})
	}

	// Apply the limit if specified
	if *limit > 0 && *limit < len(allPodInfo) {
		allPodInfo = allPodInfo[:*limit]
	}

	// Print the results
	fmt.Printf("\nFound %d pods with metrics\n", len(allPodInfo))
	fmt.Printf("%-40s %-30s %-10s %-10s %-s\n", "POD", "NAMESPACE", "CPU(m)", "MEM(MB)", "NODE")
	fmt.Println(strings.Repeat("-", 100))

	for _, info := range allPodInfo {
		fmt.Printf("%-40s %-30s %-10d %-10d %-s\n",
			truncateString(info.Name, 39),
			truncateString(info.Namespace, 29),
			info.CPUUsage,
			info.MemoryUsage,
			info.Node)
	}
}
