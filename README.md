# ecs-exec-sh

## Install

```bash
go install github.com/snaka/ecs-exec-sh/cmd/ecs-exec-sh
```

## Usage

Execute a command in a container of an ECS service.

```bash
ecs-exec-sh exec --cluster <cluster> --service <service> --container <container> [--command <command>]
```

- `--cluster`: ECS cluster name
- `--service`: ECS service name
- `--container`: ECS container name
- `--command`: Command to execute in the container (default: `/bin/sh`)

List available cluster, service, and container names.

```bash
ecs-exec-sh list
```
