## Ethbaas
Ethererum blockchain as a service.

## Prerequisites
1. this project build on k8s
2. run this cli in the k8s master node
3. you need to install sqlite3

## Using
1. create project dir and db file
```
mkdir .projects
touch .projects/dbstore
```

2. run main program
```
# go run .
A generator of baas using ethereum.

Usage:
  ethbaas [flags]
  ethbaas [command]

Available Commands:
  chain       Chain Operations.
  completion  Generate the autocompletion script for the specified shell
  contract    Contract Operation.
  help        Help about any command
  proj        Project operations.
  server      Server Operations.
  store       Store contract operations.

Flags:
  -h, --help   help for ethbaas

Use "ethbaas [command] --help" for more information about a command.

```

