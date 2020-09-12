package client

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	v1Core "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type K8s struct {
	client *kubernetes.Clientset
}

func Newk8s() (*K8s, error) {
	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}
	// } else {
	// 	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	// }

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &K8s{clientset}, nil
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
	Name      string
	Namespace string
	Status    string
	Ready     string
	Restarts  int32
	CreatedAt time.Time
}

func (k *K8s) ListPods(namespace string) ([]PodInfo, error) {

	ctx := context.TODO()
	opts := v1.ListOptions{}
	pods, _ := k.client.CoreV1().Pods(namespace).List(ctx, opts)
	podList := []PodInfo{}
	for _, pod := range pods.Items {
		restarts := int32(0)
		ready := 0
		totalContianers := len(pod.Status.ContainerStatuses)
		for _, container := range pod.Status.ContainerStatuses {
			if container.RestartCount > restarts {
				restarts = container.RestartCount
			}
			if container.State.Running != nil {
				ready++
			}
		}

		p := PodInfo{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			Status:    string(pod.Status.Phase),
			Restarts:  restarts,
			Ready:     fmt.Sprintf("%v/%v", ready, totalContianers),
			CreatedAt: pod.ObjectMeta.CreationTimestamp.Time,
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
