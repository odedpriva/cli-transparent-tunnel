# cli-transparent-tunnel 

* cli wrapper that:
  * create a ssl tunnel -> run the command through the tunnel -> close the tunnel.
  * basically it can act as an inline replacer fot the CLI 
    * the wrapper check if there is a tunnel configured for the required endpoint, if it has:
      * it adds the ssh tunnel flags to the CLI and invoke the command with all args.
      * if not, it just invokes the CLI 

## supported CLIs:
  - kubectl

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
* 