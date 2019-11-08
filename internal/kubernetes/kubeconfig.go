package kubernetes

const kubeconfig = `
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: %s
    server: %s
  name: cloud-lab
contexts:
- context:
    cluster: cloud-lab
    user: student
    namespace: %s
  name: cloud-lab
current-context: cloud-lab
kind: Config
preferences: {}
users:
- name: student
  user:
    token: %s
`
