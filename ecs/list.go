// ecs/ecs.go
package ecs

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

func List() error {
	// Create ECS Client
	svc, err := Client()
	if err != nil {
		fmt.Println("Error creating ECS client:", err)
		return err
	}

	// List Clusters
	result, err := svc.ListClusters(context.TODO(), &ecs.ListClustersInput{})
	if err != nil {
		return err
	}

	for _, clusterArn := range result.ClusterArns {
		// Extract Cluster Name
		clusterName := extractName(clusterArn)
		fmt.Println("Cluster: ", clusterName)

		// List Services
		serviceResult, err := svc.ListServices(context.TODO(), &ecs.ListServicesInput{
			Cluster: aws.String(clusterName),
		})
		if err != nil {
			return err
		}

		for _, serviceArn := range serviceResult.ServiceArns {
			serviceName := extractName(serviceArn)
			fmt.Println("  Service: ", serviceName)

			// List Tasks
			taskResult, err := svc.ListTasks(context.TODO(), &ecs.ListTasksInput{
				Cluster:     aws.String(clusterName),
				ServiceName: aws.String(serviceName),
			})
			if err != nil {
				return err
			}

			containerNames := make(map[string]struct{})

			for _, taskArn := range taskResult.TaskArns {
				// Describe Task
				taskDetail, err := svc.DescribeTasks(context.TODO(), &ecs.DescribeTasksInput{
					Cluster: aws.String(clusterName),
					Tasks:   []string{taskArn},
				})
				if err != nil {
					return err
				}

				// Collect Container Names
				for _, task := range taskDetail.Tasks {
					for _, container := range task.Containers {
						containerNames[*container.Name] = struct{}{}
					}
				}
			}

			fmt.Println("    Containers:")
			for containerName := range containerNames {
				fmt.Println("      -", containerName)
				fmt.Println("        EXAMPLE: ecs-exec-sh exec --cluster", clusterName, "--service", serviceName, "--container", containerName)
			}
		}
	}

	return nil
}

func extractName(arn string) string {
	elements := strings.Split(arn, "/")
	lastIndex := len(elements) - 1
	return elements[lastIndex]
}
