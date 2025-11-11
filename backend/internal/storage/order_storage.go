package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type OrderStorage struct {
	client    *dynamodb.Client
	tableName string
}

type OrderRecord struct {
	OrderID   string      `dynamodbav:"order_id"`
	UserID    string      `dynamodbav:"user_id"`
	Items     []OrderLine `dynamodbav:"items"`
	Status    string      `dynamodbav:"status"`
	CreatedAt time.Time   `dynamodbav:"created_at"`
}

type OrderLine struct {
	ProductID string `dynamodbav:"product_id"`
	Quantity  int32  `dynamodbav:"quantity"`
}

func NewOrderStorage(client *dynamodb.Client, tableName string) (*OrderStorage, error) {
	if client == nil {
		return nil, errors.New("dynamodb client가 nil입니다")
	}
	if tableName == "" {
		return nil, errors.New("tableName이 비어 있습니다")
	}

	return &OrderStorage{
		client:    client,
		tableName: tableName,
	}, nil
}

func (s *OrderStorage) GetOrderByID(ctx context.Context, orderID string) (*OrderRecord, error) {
	if s == nil || s.client == nil {
		return nil, errors.New("OrderStorage가 초기화되지 않았습니다")
	}
	if orderID == "" {
		return nil, errors.New("orderID가 비어 있습니다")
	}

	out, err := s.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(s.tableName),
		Key: map[string]types.AttributeValue{
			"order_id": &types.AttributeValueMemberS{Value: orderID},
		},
		ConsistentRead: aws.Bool(true),
	})
	if err != nil {
		return nil, fmt.Errorf("GetItem 실패: %w", err)
	}
	if out.Item == nil {
		return nil, fmt.Errorf("주문 %s를 찾을 수 없습니다", orderID)
	}

	var record OrderRecord
	if err := attributevalue.UnmarshalMap(out.Item, &record); err != nil {
		return nil, fmt.Errorf("주문 언마샬 실패: %w", err)
	}

	return &record, nil
}

func (s *OrderStorage) CreateOrder(ctx context.Context, record *OrderRecord) error {
	if s == nil || s.client == nil {
		return errors.New("OrderStorage가 초기화되지 않았습니다")
	}
	if record == nil {
		return errors.New("OrderRecord가 nil입니다")
	}
	if record.OrderID == "" {
		return errors.New("OrderRecord.OrderID가 비어 있습니다")
	}
	if record.CreatedAt.IsZero() {
		record.CreatedAt = time.Now().UTC()
	}

	av, err := attributevalue.MarshalMap(record)
	if err != nil {
		return fmt.Errorf("주문 marshal 실패: %w", err)
	}

	_, err = s.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(s.tableName),
		Item:                av,
		ConditionExpression: aws.String("attribute_not_exists(order_id)"),
	})
	if err != nil {
		var ccfe *types.ConditionalCheckFailedException
		if errors.As(err, &ccfe) {
			return fmt.Errorf("이미 존재하는 주문: %s", record.OrderID)
		}
		return fmt.Errorf("PutItem 실패: %w", err)
	}

	return nil
}
