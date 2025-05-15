package display

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	types "github.com/rk280392/harvesterNavigator/internal/models"
)

func DisplayVMInfo(info *types.VMInfo) {
	// Create a tabwriter for consistent alignment of all information
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// Header
	fmt.Fprintln(w, strings.Repeat("=", 80))
	fmt.Fprintf(w, "VIRTUAL MACHINE DETAILS: %s\n", info.Name)
	fmt.Fprintln(w, strings.Repeat("=", 80))

	// VM details section
	fmt.Fprintln(w, "\nVIRTUAL MACHINE INFO:")
	fmt.Fprintln(w, "------------------------")
	fmt.Fprintf(w, "Name:\t%s\n", info.Name)
	fmt.Fprintf(w, "Image ID:\t%s\n", info.ImageId)
	fmt.Fprintf(w, "Storage Class:\t%s\n", info.StorageClass)
	fmt.Fprintf(w, "Status:\t%s\n", formatVMStatus(info.VMStatus))
	fmt.Fprintf(w, "Status Reason:\t%s\n", formatVMStatusReason(info.VMStatusReason))
	fmt.Fprintf(w, "Printable Status:\t%s\n", formatPrintableStatus(info.PrintableStatus))
	fmt.Fprintf(w, "Pod Name:\t%s\n", info.PodName)

	// Storage section
	fmt.Fprintln(w, "\nSTORAGE INFO:")
	fmt.Fprintln(w, "-------------")
	fmt.Fprintf(w, "PVC Claim Names:\t%s\n", info.ClaimNames)
	fmt.Fprintf(w, "Volume Name:\t%s\n", info.VolumeName)
	fmt.Fprintf(w, "PVC Status:\t%s\n", formatPVCStatus(info.PVCStatus))

	w.Flush()

	// Helper function to pad a string to a specific visual width
	padToVisualWidth := func(s string, width int) string {
		// Calculate the number of invisible characters (ANSI escape codes)
		invisibleChars := 0
		if strings.Contains(s, "\033[") {
			// Each color code sequence is typically something like "\033[32m" and "\033[0m"
			// We'll approximate by counting occurrences of "\033["
			invisibleChars = strings.Count(s, "\033[") * 4
			// Add 1 for each "m" character
			invisibleChars += strings.Count(s, "m")
		}

		totalWidth := width + invisibleChars
		actualLen := len(s)
		padding := ""

		if actualLen < totalWidth {
			padding = strings.Repeat(" ", totalWidth-actualLen)
		}

		return s + padding
	}

	if len(info.EngineInfo) > 0 {
		fmt.Println("\nENGINE INFORMATION:")
		fmt.Println("-----------------")

		for i, engine := range info.EngineInfo {
			if i > 0 {
				fmt.Println("\n--- Engine", i+1, "---")
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintf(w, "Name:\t%s\n", engine.Name)
			fmt.Fprintf(w, "Node ID:\t%s\n", engine.NodeID)
			fmt.Fprintf(w, "Current State:\t%s\n", formatState(engine.CurrentState))
			fmt.Fprintf(w, "Active:\t%s\n", formatBool(engine.Active))
			fmt.Fprintf(w, "Started:\t%s\n", formatBool(engine.Started))
			w.Flush()

			if len(engine.Snapshots) > 0 {
				// Define column widths
				nameWidth := 38
				createdWidth := 22
				userCreatedWidth := 16
				removedWidth := 20

				// Print headers
				fmt.Printf("%s %s %s %s\n",
					padToVisualWidth("NAME", nameWidth),
					padToVisualWidth("CREATED", createdWidth),
					padToVisualWidth("USER CREATED", userCreatedWidth),
					padToVisualWidth("REMOVED", removedWidth))

				fmt.Printf("%s %s %s %s\n",
					padToVisualWidth("----", nameWidth),
					padToVisualWidth("-------", createdWidth),
					padToVisualWidth("------------", userCreatedWidth),
					padToVisualWidth("-------", removedWidth))

				for _, snapshot := range engine.Snapshots {
					fmt.Printf("%s %s %s %s\n",
						padToVisualWidth(snapshot.Name, nameWidth),
						padToVisualWidth(snapshot.Created, createdWidth),
						padToVisualWidth(formatBool(snapshot.UserCreated), userCreatedWidth),
						padToVisualWidth(formatBool(snapshot.Removed), removedWidth))
				}
				fmt.Println("\nSNAPSHOT TREE:")
				fmt.Println("-------------")
				displaySnapshotTree(engine.Snapshots)
			}
		}
	}

	// Create a new tabwriter specifically for the replica table with better spacing
	if len(info.ReplicaInfo) > 0 {
		fmt.Println("\nREPLICAS:")
		fmt.Println("---------")

		// Define column widths
		nameWidth := 15
		stateWidth := 10
		nodeWidth := 12
		ownerWidth := 45
		startedWidth := 12
		engineWidth := 45
		activeWidth := 12

		// Print headers
		fmt.Printf("%s %s %s %s %s %s %s\n",
			padToVisualWidth("NAME", nameWidth),
			padToVisualWidth("STATE", stateWidth),
			padToVisualWidth("NODE", nodeWidth),
			padToVisualWidth("OWNER", ownerWidth),
			padToVisualWidth("STARTED", startedWidth),
			padToVisualWidth("ENGINE", engineWidth),
			padToVisualWidth("ACTIVE", activeWidth))

		fmt.Printf("%s %s %s %s %s %s %s\n",
			padToVisualWidth("----", nameWidth),
			padToVisualWidth("-----", stateWidth),
			padToVisualWidth("----", nodeWidth),
			padToVisualWidth("-----", ownerWidth),
			padToVisualWidth("-------", startedWidth),
			padToVisualWidth("------", engineWidth),
			padToVisualWidth("------", activeWidth))

		// Print each replica row
		for _, replica := range info.ReplicaInfo {
			shortName := shortenName(replica.Name)
			stateFormatted := formatState(replica.CurrentState)
			startedFormatted := formatBool(replica.Started)
			activeFormatted := formatBool(replica.Active)

			fmt.Printf("%s %s %s %s %s %s %s\n",
				padToVisualWidth(shortName, nameWidth),
				padToVisualWidth(stateFormatted, stateWidth),
				padToVisualWidth(replica.NodeID, nodeWidth),
				padToVisualWidth(replica.OwnerRefName, ownerWidth),
				padToVisualWidth(startedFormatted, startedWidth),
				padToVisualWidth(replica.EngineName, engineWidth),
				padToVisualWidth(activeFormatted, activeWidth))
		}
	} else {
		fmt.Println("\nNo replicas found for this volume")
	}

	// Footer
	fmt.Println("\n" + strings.Repeat("=", 80))
}

// Helper function to shorten the replica names for better display
func shortenName(name string) string {
	// For longhorn replicas, extract just the unique part at the end
	if strings.Contains(name, "-r-") {
		parts := strings.Split(name, "-r-")
		if len(parts) == 2 {
			return "r-" + parts[1]
		}
	}

	return name
}

func formatVMStatus(status string) string {
	switch strings.ToLower(status) {
	case "true":
		return "\033[32mTrue\033[0m" // Green
	case "false":
		return "\033[31mFalse\033[0m" // Red
	default:
		return status
	}
}

func formatVMStatusReason(reason string) string {
	switch reason {
	case "GuestNotRunning":
		return "\033[31mGuestNotRunning\033[0m" // Red
	case "Running":
		return "\033[32mRunning\033[0m" // Green
	case "Starting":
		return "\033[33mStarting\033[0m" // Yellow
	case "Stopping":
		return "\033[33mStopping\033[0m" // Yellow
	case "Error":
		return "\033[31;1mError\033[0m" // Bold red
	default:
		return reason
	}
}

func formatPrintableStatus(status string) string {
	lower := strings.ToLower(status)
	if strings.Contains(lower, "starting") {
		return "\033[33mStarting\033[0m" // Yellow
	} else if strings.Contains(lower, "running") {
		return "\033[32mRunning\033[0m" // Green
	} else if strings.Contains(lower, "stopped") || strings.Contains(lower, "stopping") {
		return "\033[31mStopped\033[0m" // Red
	} else if strings.Contains(lower, "error") || strings.Contains(lower, "fail") {
		return "\033[31;1m" + status + "\033[0m" // Bold red
	}
	return status
}

// Format state with colors
func formatState(state string) string {
	switch state {
	case "running":
		return "\033[32mRUNNING\033[0m" // Green
	case "stopped":
		return "\033[31mSTOPPED\033[0m" // Red
	case "error":
		return "\033[31;1mERROR\033[0m" // Bold red
	default:
		return state
	}
}

func formatPVCStatus(status string) string {
	switch strings.ToLower(status) {
	case "bound":
		return "\033[32mBound\033[0m" // Green
	case "pending":
		return "\033[33mPending\033[0m" // Yellow
	case "lost":
		return "\033[31mLost\033[0m" // Red
	default:
		return status
	}
}

// Format boolean with colors and symbols
func formatBool(b bool) string {
	if b {
		return "\033[32mYES ✓\033[0m"
	}
	return "\033[31mNO ✗\033[0m"
}

func displaySnapshotTree(snapshots map[string]types.SnapshotInfo) {
	// Find the root node (one without a parent)
	var rootID string
	for id, snapshot := range snapshots {
		if snapshot.Parent == "" {
			rootID = id
			break
		}
	}

	if rootID == "" {
		fmt.Println("Could not determine snapshot tree root")
		return
	}

	// Recursively display the tree
	displaySnapshotNode(snapshots, rootID, "")
}

// Recursive function to display a node and its children
func displaySnapshotNode(snapshots map[string]types.SnapshotInfo, nodeID string, indent string) {
	node, exists := snapshots[nodeID]
	if !exists {
		return
	}

	// Display the current node
	label := nodeID
	if node.UserCreated {
		label += " (user)"
	}
	if node.Removed {
		label += " (removed)"
	}

	fmt.Println(indent + "└── " + label)

	// Display children with increased indentation
	childIndent := indent + "    "
	for childID := range node.Children {
		displaySnapshotNode(snapshots, childID, childIndent)
	}
}
