package kubernetes

import (
	"io/ioutil"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"os"
)

type PodInfo struct {
	Ip       string
	Hostname string
}

type Client interface {
	CheckPodImage() (bool, error)
	GetConfigMap(name, namespace string) (map[string]string, error)
}

type kubeClient struct {
	clientset *kubernetes.Clientset
	namespace string
	logger    *log.Logger
}

func MustNewClient() Client {
	logger := log.New(os.Stderr, "[Kubernetes Client] ", log.Ltime)
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.Fatalf("cannot init kubernetes config: %s", err.Error())
	}
	ns, err := getCurrentNamesapce(logger)
	if err != nil {
		logger.Fatalf("cannot get current namespace: %s", err.Error())
	}
	return &kubeClient{
		clientset: clientset,
		namespace: ns,
		logger:    logger,
	}
}

func getCurrentNamesapce(logger *log.Logger) (string, error) {
	ns, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		return "", err
	}
	logger.Printf("current kubernetes namespace is: %s", string(ns))
	return string(ns), nil
}

func (k *kubeClient) GetCurrentPodInfo() (*PodInfo, error) {
	podName := os.Getenv("HOSTNAME")
	pod, err := k.clientset.CoreV1().Pods(k.namespace).Get(podName, v1.GetOptions{})
	if err != nil {
		k.logger.Printf("cannot get pod info: %s", err.Error())
		return nil, err
	}
	return &PodInfo{
		Ip:       pod.Status.HostIP,
		Hostname: pod.Name,
	}, nil
}

func (k *kubeClient) CheckPodImage() (bool, error) {
	return false, nil
}

func (k *kubeClient) GetConfigMap(name, namespace string) (map[string]string, error) {
	ns, err := getCurrentNamesapce(k.logger)
	if err != nil {
		k.logger.Printf("cannot get current namespace: %s", err.Error())
		return nil, err
	}
	cm, err := k.clientset.CoreV1().ConfigMaps(ns).Get(name, v1.GetOptions{})
	if err != nil {
		k.logger.Printf("cannot get configmap %s, %s", name, err.Error())
		return nil, err
	}
	return cm.Data, nil
}
