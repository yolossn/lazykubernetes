package client

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
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

type PodInfo struct {
	Name      string
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
			Status:    string(pod.Status.Phase),
			Restarts:  restarts,
			Ready:     fmt.Sprintf("%v/%v", ready, totalContianers),
			CreatedAt: pod.ObjectMeta.CreationTimestamp.Time,
		}
		podList = append(podList, p)
	}
	return podList, nil
}

type JobInfo struct {
	Name        string
	Completions string
	Duration    float64 // secs
	Age         string
	CreatedAt   time.Time
}

func (k *K8s) ListJobs(namespace string) ([]JobInfo, error) {
	ctx := context.TODO()
	opts := v1.ListOptions{}
	jobs, _ := k.client.BatchV1().Jobs(namespace).List(ctx, opts)
	jobList := []JobInfo{}
	for _, job := range jobs.Items {
		totalPods := job.Status.Active + job.Status.Succeeded + job.Status.Failed

		var duration float64
		if job.Status.CompletionTime != nil {
			duration = job.Status.CompletionTime.Time.Sub(job.ObjectMeta.CreationTimestamp.Time).Seconds()
		}

		p := JobInfo{
			Name:        job.Name,
			Completions: fmt.Sprintf("%v/%v", job.Status.Succeeded, totalPods),
			Duration:    duration,
			CreatedAt:   job.ObjectMeta.CreationTimestamp.Time,
		}
		jobList = append(jobList, p)
	}
	return jobList, nil
}
