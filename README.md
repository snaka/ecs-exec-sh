# ecs-exec-sh

## Install

```bash
go install github.com/snaka/ecs-exec-sh/cmd/ecs-exec-sh
```

## Usage

Execute a command in a container of an ECS service.

```bash
ecs-exec-sh exec -c <cluster> -s <service> -C <container> [-x <command>]
```

- `--cluster`, `-c`: ECS cluster name
- `--service`, `-s`: ECS service name
- `--container`, `-C`: ECS container name
- `--command`, `-x`: Command to execute in the container (default: `/bin/sh`)

List available cluster, service, and container names.

```bash
ecs-exec-sh list
```
