package ecs

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func ExecuteCommand(clusterName, serviceName, containerName, command string) error {
	svc, err := Client()
	if err != nil {
		return fmt.Errorf("failed to create ECS client: %w", err)
	}

	ctx := context.TODO()

	// List tasks
	listTasksOutput, err := svc.ListTasks(ctx, &ecs.ListTasksInput{
		Cluster:     aws.String(clusterName),
		ServiceName: aws.String(serviceName),
	})
	if err != nil {
		return fmt.Errorf("failed to list tasks: %w", err)
	}

	// Get task ID from the list
	taskID := listTasksOutput.TaskArns[0]

	fmt.Printf("Task: %s\n", taskID)

	return runExec(ctx, svc, clusterName, taskID, containerName, command)
}

func runExec(ctx context.Context, svc *ecs.Client, clusterName, taskID, containerName, command string) error {
	out, err := svc.ExecuteCommand(ctx, &ecs.ExecuteCommandInput{
		Cluster:     aws.String(clusterName),
		Interactive: true,
		Task:        aws.String(taskID),
		Command:     aws.String(command),
		Container:   aws.String(containerName),
	})
	if err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}

	tasks, err := svc.DescribeTasks(ctx, &ecs.DescribeTasksInput{
		Cluster: aws.String(clusterName),
		Tasks:   []string{taskID},
	})
	if err != nil {
		return fmt.Errorf("failed to describe task: %w", err)
	}
	task := tasks.Tasks[0]

	target, err := ssmRequestTarget(task, containerName)
	if err != nil {
		return fmt.Errorf("failed to get SSM target: %w", err)
	}
	return runSessionManagerPlugin(ctx, out.Session, target)
}

func ssmRequestTarget(task types.Task, targetContainer string) (string, error) {
	value := strings.Split(*task.TaskArn, "/")
	clusterName := value[1]
	taskID := value[2]
	var runtimeID string
	for _, c := range task.Containers {
		if *c.Name == targetContainer {
			runtimeID = *c.RuntimeId
		}
	}
	return fmt.Sprintf("ecs:%s_%s_%s", clusterName, taskID, runtimeID), nil
}

func runSessionManagerPlugin(ctx context.Context, session *types.Session, target string) error {
	sess, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}
	ssmReq, err := json.Marshal(map[string]string{
		"Target": target,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal ssm request: %w", err)
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	cmd := exec.CommandContext(
		ctx,
		"session-manager-plugin",
		string(sess),
		"ap-northeast-1",
		"StartSession",
		"",
		string(ssmReq),
	)
	// send SIGINT to the process when the context is canceled
	cmd.Cancel = func() error {
		return cmd.Process.Signal(os.Interrupt)
	}
	cmd.WaitDelay = 3 * time.Second
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
