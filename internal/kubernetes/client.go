package kubernetes

import (
    cfg "github.com/hitman99/autograde/internal/config"
    "io/ioutil"
    v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/clientcmd"
    "log"
    "os"
    "strings"
)

type PodInfo struct {
    Ip       string
    Hostname string
}

type Client interface {
    CheckContainerImage(namespace, labelSelector, expectedImage string) (bool, error)
    GetConfigMap(name, namespace string) (map[string]string, error)
}

type kubeClient struct {
    clientset *kubernetes.Clientset
    namespace string
    logger    *log.Logger
}

func MustNewClient() Client {
    logger := log.New(os.Stderr, "[Kubernetes Client] ", log.Ltime)
    var (
        err    error
        config *rest.Config
    )
    if cfg.GetConfig().DevMode {
        config, err = clientcmd.BuildConfigFromFlags("", cfg.GetConfig().KubeconfigPath)
    } else {
        config, err = rest.InClusterConfig()
    }
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
    var ns []byte
    if cfg.GetConfig().DevMode {
        ns = []byte("autograde")
    } else {
        ns, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
        if err != nil {
            return "", err
        }
        logger.Printf("current kubernetes namespace is: %s", string(ns))
    }
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

func (k *kubeClient) CheckContainerImage(namespace, labelSelector, expectedImage string) (bool, error) {
    depList, err := k.clientset.AppsV1().Deployments(namespace).List(v1.ListOptions{
        LabelSelector: labelSelector,
    })
    if err != nil {
        return false, nil
    }
    for _, dp := range depList.Items {
		for _, container := range dp.Spec.Template.Spec.Containers {
			if strings.Contains(container.Image, expectedImage) {
				return true, nil
			}
		}
	}
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
