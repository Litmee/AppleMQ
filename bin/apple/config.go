package apple

func checkAndDealOption() {

	// Processing of mode configuration items
	mode := options["mode"]
	if mode != "standalone" && mode != "cluster" {
		panic("The value of the mode configuration item in the apple.ini configuration file is not in the range expected by the system, please correct it: mode=" + mode)
	}

	// 如果是集群部署,则需要处理一下
	if mode == "cluster" {
		dealModeCluster()
	}
}

func dealModeCluster() {
}
