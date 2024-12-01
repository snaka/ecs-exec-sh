package ecs

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/fujiwara/ecsta"
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

    ecstaApp, err := newEcsta(ctx, "ap-northeast-1", clusterName)
    if err != nil {
        return fmt.Errorf("failed to create Ecsta: %w", err)
    }

    // Execute command
    return ecstaApp.RunExec(ctx,&ecsta.ExecOption{
        ID: taskID,
        Service: &serviceName,
        Container: containerName,
        Command: command,
    })
}

func newEcsta(ctx context.Context, region, cluster string) (*ecsta.Ecsta, error) {
    app, err := ecsta.New(ctx, region, cluster)
    if err != nil {
        return nil, fmt.Errorf("failed to create Ecsta, %w", err)
    }
    return app, nil
}
