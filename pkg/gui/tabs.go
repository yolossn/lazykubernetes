package gui

func getResourceTabs() []string {
	return []string{"pod", "job", "deploy", "service", "secret", "configMap"}
}

func getClusterInfoTabs() []string {
	return []string{"info-dump"}
}

func getNamespaceInfoTabs() []string {
	return []string{"description"}
}

func getPodInfoTabs() []string {
	return []string{"description", "logs"}
}

func getJobInfoTabs() []string {
	return []string{"logs", "description", "cron"}
}

func getDeployInfoTabs() []string {
	return []string{"description"}
}

func getServiceInfoTabs() []string {
	return []string{"description"}
}

func getSecretInfoTabs() []string {
	return []string{"description"}
}

func getConfigMapInfoTabs() []string {
	return []string{"description"}
}
