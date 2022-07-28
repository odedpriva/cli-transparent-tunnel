commands-configuration:
  kubectl:
    path: {{ .Env.KUBECTL_LOCATION }}
    flags:
      address:
        - --server
        - -s
      sni:
        - --tls-server-name
