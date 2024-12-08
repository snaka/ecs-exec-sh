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
	"github.com/erikgeiser/promptkit/selection"
)

func ExecuteCommand(clusterName, serviceName, containerName, command string) error {
	svc, err := Client()
	if err != nil {
		return fmt.Errorf("failed to create ECS client: %w", err)
	}

	ctx := context.TODO()

	// cluster selection
	if clusterName == "" {
		choice, err := selectCluster(ctx, svc)
		if err != nil {
			return fmt.Errorf("failed to select cluster: %w", err)
		}
		clusterName = choice
	}

	// service selection
	if serviceName == "" {
		choice, err := selectService(ctx, svc, clusterName)
		if err != nil {
			return fmt.Errorf("failed to select service: %w", err)
		}
		serviceName = choice
	}

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

	// container selection
	if containerName == "" {
		choice, err := selectContainer(ctx, svc, clusterName, taskID)
		if err != nil {
			return fmt.Errorf("failed to select container: %w", err)
		}
		containerName = choice
	}

	return runExec(ctx, svc, clusterName, taskID, containerName, command)
}

func selectCluster(ctx context.Context, svc *ecs.Client) (string, error) {
	clusters, err := svc.ListClusters(ctx, &ecs.ListClustersInput{})
	if err != nil {
		return "", fmt.Errorf("failed to list clusters: %w", err)
	}

	clusterNames := make([]string, 0, len(clusters.ClusterArns))
	for _, c := range clusters.ClusterArns {
		clusterNames = append(clusterNames, extractName(c))
	}

	// selection
	sp := selection.New("Select cluster:", clusterNames)
	choice, err := sp.RunPrompt()
	if err != nil {
		return "", fmt.Errorf("failed to select cluster: %w", err)
	}
	return choice, nil
}

func selectService(ctx context.Context, svc *ecs.Client, clusterName string) (string, error) {
	services, err := svc.ListServices(ctx, &ecs.ListServicesInput{
		Cluster: aws.String(clusterName),
	})
	if err != nil {
		return "", fmt.Errorf("failed to list services: %w", err)
	}

	serviceNames := make([]string, 0, len(services.ServiceArns))
	for _, s := range services.ServiceArns {
		serviceNames = append(serviceNames, extractName(s))
	}

	// selection
	sp := selection.New("Select service:", serviceNames)
	choice, err := sp.RunPrompt()
	if err != nil {
		return "", fmt.Errorf("failed to select service: %w", err)
	}
	return choice, nil
}

func selectContainer(ctx context.Context, svc *ecs.Client, clusterName, taskID string) (string, error) {
	tasks, err := svc.DescribeTasks(ctx, &ecs.DescribeTasksInput{
		Cluster: aws.String(clusterName),
		Tasks:   []string{taskID},
	})
	if err != nil {
		return "", fmt.Errorf("failed to describe task: %w", err)
	}
	task := tasks.Tasks[0]

	containerNames := make([]string, 0, len(task.Containers))
	for _, c := range task.Containers {
		containerNames = append(containerNames, *c.Name)
	}

	// selection
	sp := selection.New("Select container:", containerNames)
	choice, err := sp.RunPrompt()
	if err != nil {
		return "", fmt.Errorf("failed to select container: %w", err)
	}
	return choice, nil
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
