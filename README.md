# ecs-exec-sh

Easily execute shell commands in your Amazon ECS containers.

`ecs-exec-sh` simplifies running commands in containers in your ECS containers. Interactively select the cluster, service, and container -- no need to memorize names or IDs.

## Install

```bash
go install github.com/snaka/ecs-exec-sh/cmd/ecs-exec-sh@latest
```

## Usage

Execute a command in a container of an ECS service.

```bash
ecs-exec-sh
```

**Options:**

- `--cluster`, `-c`: ECS cluster name
- `--service`, `-s`: ECS service name
- `--container`, `-C`: ECS container name
- `--command`, `-x`: Command to execute in the container (default: `/bin/sh`)
