# cli-transparent-tunnel 

- cli wrapper that:
  - create a ssl tunnel -> run the command through the tunnel -> close the tunnel.
- basically it can act as an inline replacer fot the CLI 
  - the wrapper check if there is a tunnel configured for the required endpoint, if it has:
    - it adds the ssh tunnel flags to the CLI and invoke the command with all args.
    - if not, it just invokes the CLI 

## supported CLIs:
  - kubectl

## install

```shell
>> brew tap odedpriva/ctt
>> brew install ctt
```

## Environment Variables

| Variable      | Description | Default Value |
|---------------|-------------|---------------|
| CTT_LOG_LEVEL | panic       |               |
| CTT_CONFIG    |             |               | 