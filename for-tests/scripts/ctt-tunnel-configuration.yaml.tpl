configurations:
    kubectl:
    - ssh-tunnel-server: {{ .Env.SSH_TUNNEL_SERVER }}
      origin-server: {{ .Env.K8S_ENDPOINT }}
      name: {{ .Env.K8S_CLUSTER_NAME }}
ssh-config:
  key-path: {{ .Env.SSH_PRIVATE_KEY_PATH }}