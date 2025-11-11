package storage

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"Acho-mj/2025_Golang_MSA/backend/internal/config"
)

// AWS DynamoDb에 접근하기 위한 클라이언트 객체 초기화 (연결)
func NewDynamoClient(ctx context.Context, cfg *config.Config) (*dynamodb.Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config가 nil입니다")
	}

	// AWS SDK의 설정 로딩 옵션들을 모아둔 슬라이스
	loadOpts := []func(*awscfg.LoadOptions) error{
		awscfg.WithRegion(cfg.AWSRegion),
	}

	awsCfg, err := awscfg.LoadDefaultConfig(ctx, loadOpts...)
	if err != nil {
		return nil, fmt.Errorf("AWS 설정 로드 실패: %w", err)
	}

	if cfg.AWSEndpoint != "" {
		return dynamodb.NewFromConfig(awsCfg, func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String(cfg.AWSEndpoint)
		}), nil
	}

	return dynamodb.NewFromConfig(awsCfg), nil
}
