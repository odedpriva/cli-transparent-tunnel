# cli-transparent-tunnel 

* CLI wrapper that:
  * create a ssl tunnel -> run the command through the tunnel -> close the tunnel.
  * basically it can act as an inline replacer fot the CLI 
    * the wrapper check if there is a tunnel configured for the required endpoint, if it has:
      * it adds the ssh tunnel flags to the CLI and invoke the command with all args.
      * if not, it just invokes the CLI 

* THIS IS A WIP PROJECT, I believe it can solve accessing k8s clusters through ssh tunnel in certain scenarios such as CI / CD etc .
* Let me know what you think.

## supported CLIs:
  - kubectl


[![demo](https://asciinema.org/a/3EG9Jckd4Oy3uKkQIHbdYgs5q.svg)](https://asciinema.org/a/3EG9Jckd4Oy3uKkQIHbdYgs5q?autoplay=1)

## install

```shell
>> brew tap odedpriva/ctt
>> brew install ctt
```
## Setup

```shell 
>> ctt ctt-init # this will create a config file under ${HOME}/.ctt/config.yaml
```
* update config file with tunnel configuration. 
* `ctt` creates tunnel only when the k8s cluster name equals the name of the tunnel configuration
* `ctt` creates tunnel for supported subcommand only 

* environment Variables

| Variable      | Description | Default Value |
|---------------|-------------|---------------|
| CTT_LOG_LEVEL |             | panic         |
| CTT_CONFIG    |             |               |

## Using

* to see the tool in action, I suggest setting up CTT_LOG_LEVEL to debug.

```shell
>> export CTT_LOG_LEVEL=debug
>> alias kubectl=ctt
>> kubectl get pods
kubectl get pods
DEBU[0000] using context docker-desktop
DEBU[0000] using server docker-desktop
DEBU[0000] using ssh key /Users/odedpriva/.ssh/secure-access-cloud-key.pem
DEBU[0000] /usr/local/bin/kubectl --server https://127.0.0.1:65152 --tls-server-name kubernetes.docker.internal get pods
DEBU[0000] accepted connection
DEBU[0000] forwording connection on ssh tunnel kubernetes-docker-internal.tcp.symchatbotdemo.luminatesite.com:22
DEBU[0000] using ssh configuration &{Config:{Rand:<nil> RekeyThreshold:0 KeyExchanges:[] Ciphers:[] MACs:[]} User:tcptunnel@kubernetes-docker-internal Auth:[0x2499ba0] HostKeyCallback:0x2498980 BannerCallback:<nil> ClientVersion: HostKeyAlgorithms:[] Timeout:0s}
DEBU[0000] connected to kubernetes-docker-internal.tcp.symchatbotdemo.luminatesite.com:22 (1 of 2)
DEBU[0001] connected to kubernetes.docker.internal:6443 (2 of 2)
NAME                     READY   STATUS    RESTARTS   AGE
nginx-6799fc88d8-p4l67   1/1     Running   0          42s
```

## Caveats

* Supports subset of kubectl subcommands listed in ${HOME}/.ctt/config.yaml
* when aliasing, the shell completion is not working

