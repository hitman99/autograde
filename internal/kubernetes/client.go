package kubernetes

import (
	stdErr "errors"
	"fmt"
	cfg "github.com/hitman99/autograde/internal/config"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	v1rbac "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1meta "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	CreateNamespace(namespace string) error
	DeleteNamespace(namespace string) error
	GetKubeconfig(namespace string) (string, error)
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
	pod, err := k.clientset.CoreV1().Pods(k.namespace).Get(podName, v1meta.GetOptions{})
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
	depList, err := k.clientset.AppsV1().Deployments(namespace).List(v1meta.ListOptions{
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
	var cfgNamespace string
	if namespace == "" {
		ns, err := getCurrentNamesapce(k.logger)
		if err != nil {
			k.logger.Printf("cannot get current namespace: %s", err.Error())
			return nil, err
		}
		cfgNamespace = ns
	} else {
		cfgNamespace = namespace
	}
	cm, err := k.clientset.CoreV1().ConfigMaps(cfgNamespace).Get(name, v1meta.GetOptions{})
	if err != nil {
		k.logger.Printf("cannot get configmap %s, %s", name, err.Error())
		return nil, err
	}
	return cm.Data, nil
}

func (k *kubeClient) CreateNamespace(namespace string) error {
	_, err := k.clientset.CoreV1().Namespaces().Get(namespace, v1meta.GetOptions{})
	if err == nil ||  errors.IsNotFound(err){
		_, err := k.clientset.CoreV1().Namespaces().Create(&v1.Namespace{
			ObjectMeta: v1meta.ObjectMeta{
				Name: namespace,
			},
		})

		if err != nil {
			k.logger.Printf("cannot create namespace %s, %s", namespace, err.Error())
			return err
		}
		_, err = k.clientset.RbacV1().RoleBindings(namespace).Create(&v1rbac.RoleBinding{
			ObjectMeta: v1meta.ObjectMeta{
				Name:                       "admin",
			},
			Subjects:   []v1rbac.Subject{{
				Kind:      "ServiceAccount",
				APIGroup:  "",
				Name:      "default",
				Namespace: namespace,
			}},
			RoleRef:    v1rbac.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "ClusterRole",
				Name:     "admin",
			},
		})
		if err != nil {
			k.logger.Printf("cannot create rolebinding %s, %s", namespace, err.Error())
			return err
		}
	}
	return nil
}

func (k *kubeClient) GetKubeconfig(namespace string) (string, error) {
	sa, err := k.clientset.CoreV1().ServiceAccounts(namespace).Get("default", v1meta.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("kubernetes error", err)
	}
	if len(sa.Secrets) == 1 {
		secret, err := k.clientset.CoreV1().Secrets(namespace).Get(sa.Secrets[0].Name, v1meta.GetOptions{})
		if err != nil {
			return "", fmt.Errorf("kubernetes error", err)
		}
		return fmt.Sprintf(kubeconfig, cfg.GetConfig().KubeApiServerCA, cfg.GetConfig().KubeApiServer, namespace, secret.Data["token"]), nil
	} else {
		return "", stdErr.New("no secrets found")
	}
}

func (k *kubeClient) DeleteNamespace(namespace string) error {
	err := k.clientset.CoreV1().Namespaces().Delete(namespace, &v1meta.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		k.logger.Printf("cannot delete namespace %s, %s", namespace, err.Error())
		return err
	}
	return nil
}
