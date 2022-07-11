package apple

import (
	"strings"
)

// Check and handle configuration items
func checkAndDealOption() bool {

	// Processing of mode configuration items
	mode := options["mode"]
	if mode != "standalone" && mode != "cluster" {
		panic("The value of the mode configuration item in the apple.ini configuration file is not in the range expected by the system, please correct it: mode=" + mode)
	}

	// If it is a cluster deployment, further processing is required
	if mode == "cluster" {
		dealModeCluster()
		return true
	}
	return false
}

// Handling cluster deployments
func dealModeCluster() {
	// Get the value of a cluster configuration item
	cluster := options["cluster-map"]
	// Data segmentation processing
	clusterArr = strings.Split(cluster, ";")
	for i, v := range clusterArr {
		// Handling possible mistyped space characters
		v = strings.ReplaceAll(v, " ", "")
		clusterArr[i] = strings.TrimSpace(v)
	}
	// Connect to cluster target
	clusterConnection()
}
