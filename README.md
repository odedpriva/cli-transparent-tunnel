# cli-transparent-tunnel 

## Introduction

SSH tunneling is a great method of transporting arbitrary networking data over an encrypted SSH connection.
It can be used to add 
- encryption to legacy applications. 
- access network servers behind an ssh bastion server 
- and more .. 

To do so, we need to create an SSH tunnel and point the local client to the tunnel's local port. 

for examples, running psql command through a ssl tunnel

```shell
# open a tunnel in a one terminal 
# make sure 3307 is not in use
>> ssh -N -L 5432:my-bastion-server:5432 -p 22 <>USER>@<IP> 

## run psql command in another terminal 
>> psql -U username -h 127.0.0.1 -P 5432 -p password -f commands.sql 
```

I found it a bit tedious and a bit complicated for scenarios such as CI etc .  

`ctt` allows you to prefix any supported command, the tool will preform all the heavy lifting of creating the tunnel and 
adjusting the cli with the proper host and port.
```shell
>> ctt --tunnel-config psql-us psql -U username -p password -f commands.sql  

psql -h 127.0.0.1 -P 65152 -U username -p password -f commands.sql 
```

## How

`ctt` has 2 config files, one for tunnel configurations, and other for command config 

e.g:

cli-config
```yaml
commands-configuration:
  redis-cli:
    path: /usr/local/bin/redis-cli
    flags:
      host:
        - -h
      port:
        - -p
      sni:
        - --sni
  psql:
    path: /usr/local/bin/psql
    flags:
      host:
        - --host
        - -h
      port:
        - -p
        - --port
  kubectl:
    path: /usr/local/bin/kubectl
    flags:
      address:
        - --server
        - -s
      sni:
        - --tls-server-name
```

In theory, `ctt` should support any cli that allows passing endpoint using a flag.   

tunnel-config
```yaml
configurations:
  redis-cli:
    - ssh-tunnel-server: my-user@eu-bastion:22
      name: redis-eu
      origin-server: redis:6379
  psql:
    - ssh-tunnel-server: my-user@us-bastion:22
      origin-server: postgres:5432
      name: psql-us
  kubectl:
    - ssh-tunnel-server: my-user@us-bastion:22
      origin-server: k8s:443
      name: k8s-us
ssh-config:
  key-path: ~/ssh/id_rsa
```

## Installation

```shell
>> brew tap odedpriva/ctt
>> brew install ctt
```

## Usage

```shell
NAME:
   ctt - make a command run through an ssh tunnel

USAGE:
   ctt [global options] command [command options] [arguments...]

COMMANDS:
   setup
   tunnel
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug     debug mode (default: false)
   --help, -h  show help (default: false)
```

### tunnel
```shell
NAME:
   ctt tunnel

USAGE:
   ctt tunnel [command options] command-to-tunnel [command-to-tunnel-options]

OPTIONS:
   --tunnel-config value  tunnel config name

```

## Setup

TODO .. 






