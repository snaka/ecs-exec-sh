package ecs

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

func Client() (*ecs.Client, error) {
	// Load AWS Config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Error loading AWS config:", err)
		return nil, err
	}

	// Create ECS Client
	svc := ecs.NewFromConfig(cfg)

	return svc, nil
}
