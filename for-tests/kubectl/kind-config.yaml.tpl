kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
networking:
  apiServerAddress: "0.0.0.0"
nodes:
  - role: control-plane
    kubeadmConfigPatches:
    - |
      kind: ClusterConfiguration
      apiServer:
        certSANs:
        - "0.0.0.0"
        - "{{ .Env.LOCAL_IP }}"