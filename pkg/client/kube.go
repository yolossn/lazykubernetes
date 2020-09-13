package client

import (
	"context"
	"time"

	"github.com/yolossn/lazykubernetes/pkg/utils"
	v1Core "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type K8s struct {
	client *kubernetes.Clientset
}

func Newk8s() (*K8s, error) {

	kubeConfigPath, err := utils.FindKubeConfig()
	if err != nil {
		return nil, err
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &K8s{clientset}, nil
}

type ServerInfo struct {
	Context    string
	Server     string
	ServerInfo *version.Info
}

func (k *K8s) GetServerInfo() (*ServerInfo, error) {
	kubeConfigPath, err := utils.FindKubeConfig()
	if err != nil {
		return nil, err
	}

	kubeConfig, err := clientcmd.LoadFromFile(kubeConfigPath)
	if err != nil {
		return nil, err
	}

	version, err := k.client.DiscoveryClient.ServerVersion()
	if err != nil {
		return nil, err
	}

	currentCluster := *kubeConfig.Clusters[kubeConfig.CurrentContext]
	s := ServerInfo{
		Context:    kubeConfig.CurrentContext,
		Server:     currentCluster.Server,
		ServerInfo: version,
	}

	return &s, nil
}

type NamespaceInfo struct {
	Name      string
	Status    string
	CreatedAt time.Time
}

func (k *K8s) ListNamespace() ([]NamespaceInfo, error) {
	ctx := context.TODO()
	opts := v1.ListOptions{}
	list, _ := k.client.CoreV1().Namespaces().List(ctx, opts)
	ns := []NamespaceInfo{}
	for _, item := range list.Items {
		n := NamespaceInfo{
			Name:      item.ObjectMeta.Name,
			Status:    string(item.Status.Phase),
			CreatedAt: item.ObjectMeta.CreationTimestamp.Time,
		}
		ns = append(ns, n)
	}
	return ns, nil
}

// TODO: Verify timeout and handle it
func (k *K8s) WatchNamespace() (watch.Interface, error) {
	ctx := context.TODO()
	opts := v1.ListOptions{}
	wi, err := k.client.CoreV1().Namespaces().Watch(ctx, opts)
	if err != nil {
		return nil, err
	}
	return wi, nil
}

func (k *K8s) GetNamespace(ns string) (*v1Core.Namespace, error) {
	ctx := context.TODO()
	opts := v1.GetOptions{}
	return k.client.CoreV1().Namespaces().Get(ctx, ns, opts)
}

func (k *K8s) WatchPods(namespace string) (watch.Interface, error) {
	ctx := context.TODO()
	opts := v1.ListOptions{}

	wi, err := k.client.CoreV1().Pods(namespace).Watch(ctx, opts)
	if err != nil {
		return nil, err
	}
	return wi, nil
}

type PodInfo struct {
	Name            string
	Namespace       string
	Status          string
	ReadyContainers int32
	TotalContainers int32
	Restarts        int32
	CreatedAt       time.Time
}

func (k *K8s) ListPods(namespace string) ([]PodInfo, error) {

	ctx := context.TODO()
	opts := v1.ListOptions{}
	pods, _ := k.client.CoreV1().Pods(namespace).List(ctx, opts)
	podList := []PodInfo{}
	for _, pod := range pods.Items {
		restarts := int32(0)
		ready := int32(0)
		totalContianers := int32(len(pod.Status.ContainerStatuses))
		for _, container := range pod.Status.ContainerStatuses {
			if container.RestartCount > restarts {
				restarts = container.RestartCount
			}
			if container.State.Running != nil {
				ready++
			}
		}

		p := PodInfo{
			Name:            pod.Name,
			Namespace:       pod.Namespace,
			Status:          string(pod.Status.Phase),
			Restarts:        restarts,
			ReadyContainers: ready,
			TotalContainers: totalContianers,
			CreatedAt:       pod.ObjectMeta.CreationTimestamp.Time,
		}
		podList = append(podList, p)
	}
	return podList, nil
}

func (k *K8s) DescribePod(ns string, podname string) (*v1Core.Pod, error) {
	ctx := context.TODO()
	opts := v1.GetOptions{}

	out, err := k.client.CoreV1().Pods(ns).Get(ctx, podname, opts)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (k *K8s) StreamPodLogs(ns string, podname string) *restclient.Request {
	opts := &v1Core.PodLogOptions{Follow: true}
	request := k.client.CoreV1().Pods(ns).GetLogs(podname, opts)
	return request
}

type JobInfo struct {
	Name        string
	Namespace   string
	Active      int32
	Succeeded   int32
	Failed      int32
	Age         string
	CompletedAt time.Time
	CreatedAt   time.Time
}

func (k *K8s) ListJobs(namespace string) ([]JobInfo, error) {
	ctx := context.TODO()
	opts := v1.ListOptions{}
	jobs, _ := k.client.BatchV1().Jobs(namespace).List(ctx, opts)
	jobList := []JobInfo{}
	for _, job := range jobs.Items {

		p := JobInfo{
			Name:        job.Name,
			Namespace:   job.Namespace,
			Active:      job.Status.Active,
			Succeeded:   job.Status.Succeeded,
			Failed:      job.Status.Failed,
			CompletedAt: job.Status.CompletionTime.Time,
			CreatedAt:   job.ObjectMeta.CreationTimestamp.Time,
		}
		jobList = append(jobList, p)
	}
	return jobList, nil
}

type DeploymentInfo struct {
	Name            string
	Namespace       string
	Available       int32
	ReadyReplicas   int32
	Replicas        int32
	UpdatedReplicas int32
	CreatedAt       time.Time
}

func (k *K8s) ListDeployments(namespace string) ([]DeploymentInfo, error) {
	ctx := context.TODO()
	opts := v1.ListOptions{}
	deployments, _ := k.client.AppsV1().Deployments(namespace).List(ctx, opts)
	deploymentList := []DeploymentInfo{}
	for _, deployment := range deployments.Items {
		d := DeploymentInfo{
			Name:            deployment.Name,
			Namespace:       deployment.Namespace,
			Available:       deployment.Status.AvailableReplicas,
			ReadyReplicas:   deployment.Status.ReadyReplicas,
			Replicas:        deployment.Status.Replicas,
			UpdatedReplicas: deployment.Status.UpdatedReplicas,
			CreatedAt:       deployment.ObjectMeta.CreationTimestamp.Time,
		}
		deploymentList = append(deploymentList, d)
	}
	return deploymentList, nil
}

type StatefulsetInfo struct {
	Name            string
	Namespace       string
	CurrentReplicas int32
	ReadyReplicas   int32
	Replicas        int32
	UpdatedReplicas int32
	CreatedAt       time.Time
}

func (k *K8s) ListStatefulsets(namespace string) ([]StatefulsetInfo, error) {
	ctx := context.TODO()
	opts := v1.ListOptions{}
	statefulsets, _ := k.client.AppsV1().StatefulSets(namespace).List(ctx, opts)
	statefulsetList := []StatefulsetInfo{}
	// https: //kubernetes.io/docs/reference/generated/kubernetes-api/v1.19/#statefulsetstatus-v1-apps
	for _, statefulset := range statefulsets.Items {
		s := StatefulsetInfo{
			Name:            statefulset.Name,
			Namespace:       statefulset.Namespace,
			CurrentReplicas: statefulset.Status.CurrentReplicas,
			ReadyReplicas:   statefulset.Status.ReadyReplicas,
			Replicas:        statefulset.Status.Replicas,
			UpdatedReplicas: statefulset.Status.UpdatedReplicas,
			CreatedAt:       statefulset.ObjectMeta.CreationTimestamp.Time,
		}
		statefulsetList = append(statefulsetList, s)
	}
	return statefulsetList, nil
}

type SecretInfo struct {
	Name      string
	Namespace string
	Type      string
	Data      int32
	CreatedAt time.Time
}

func (k *K8s) ListSecrets(namespace string) ([]SecretInfo, error) {
	ctx := context.TODO()
	opts := v1.ListOptions{}
	secrets, _ := k.client.CoreV1().Secrets(namespace).List(ctx, opts)
	secretList := []SecretInfo{}
	// https: //kubernetes.io/docs/reference/generated/kubernetes-api/v1.19/#secret-v1-core
	for _, secret := range secrets.Items {
		s := SecretInfo{
			Name:      secret.Name,
			Namespace: secret.Namespace,
			Type:      string(secret.Type),
			Data:      int32(len(secret.Data)),
			CreatedAt: secret.ObjectMeta.CreationTimestamp.Time,
		}
		secretList = append(secretList, s)
	}
	return secretList, nil
}

type ConfigMapInfo struct {
	Name      string
	Namespace string
	Data      int32
	CreatedAt time.Time
}

func (k *K8s) ListConfigMap(namespace string) ([]ConfigMapInfo, error) {
	ctx := context.TODO()
	opts := v1.ListOptions{}
	configmaps, _ := k.client.CoreV1().ConfigMaps(namespace).List(ctx, opts)
	configmapList := []ConfigMapInfo{}
	// https: //kubernetes.io/docs/reference/generated/kubernetes-api/v1.19/#configmap-v1-core
	for _, configmap := range configmaps.Items {
		s := ConfigMapInfo{
			Name:      configmap.Name,
			Namespace: configmap.Namespace,
			Data:      int32(len(configmap.Data)),
			CreatedAt: configmap.ObjectMeta.CreationTimestamp.Time,
		}
		configmapList = append(configmapList, s)
	}
	return configmapList, nil
}

type NodeInfo struct {
	Name      string
	Status    string
	Version   string
	CreatedAt time.Time
}

func (k *K8s) ListNode() ([]NodeInfo, error) {
	ctx := context.TODO()
	opts := v1.ListOptions{}
	nodes, _ := k.client.CoreV1().Nodes().List(ctx, opts)
	nodeList := []NodeInfo{}
	// https: //kubernetes.io/docs/reference/generated/kubernetes-api/v1.19/#node-v1-core
	for _, node := range nodes.Items {
		var status string
		if len(node.Status.Conditions) > 0 {
			status = string(node.Status.Conditions[len(node.Status.Conditions)-1].Type)
		}

		n := NodeInfo{
			Name:      node.Name,
			Status:    status,
			Version:   node.Status.NodeInfo.KubeletVersion,
			CreatedAt: node.ObjectMeta.CreationTimestamp.Time,
		}
		nodeList = append(nodeList, n)
	}
	return nodeList, nil
}
