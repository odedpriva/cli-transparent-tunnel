kubectl:
  path: {{ .Env.KUBECTL_LOCATION }}
  supported-subcommands:
    - create
    - cp
    - expose
    - run
    - set
    - explain
    - get
    - delete
    - rollout
    - scale
    - certificate
    - cluster-info
    - cordon
    - uncordon
    - drain
    - taint
    - describe
    - logs
    - auth
    - apply
    - patch
    - replace
  tunnel-configurations:
    - ssh-tunnel-server: {{ .Env.SSH_TUNNEL_SERVER }}
      origin-server: {{ .Env.K8S_ENDPOINT }}
      name: {{ .Env.K8S_CLUSTER_NAME }}
ssh-config:
  key-path: {{ .Env.SSH_PRIVATE_KEY_PATH }}