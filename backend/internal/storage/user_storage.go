package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type UserStorage struct {
	// client 객체가 있어야 DB쿼리를 AWS에 보낼 수 있음
	client *dynamodb.Client
	// 어떤 테이블에서 작업을 할지 명시함
	tableName string
}

// 실제 테이블 구조와 1:1 대응
type UserItem struct {
	UserID    string    `dynamodbav:"user_id"`
	Email     string    `dynamodbav:"email"`
	Name      string    `dynamodbav:"name"`
	CreatedAt time.Time `dynamodbav:"created_at"`
}

// UserStorage 객체를 생성하고 초기화
func NewUserStorage(client *dynamodb.Client, tableName string) (*UserStorage, error) {
	if client == nil {
		return nil, errors.New("dynamodb client가 nil입니다")
	}
	if tableName == "" {
		return nil, errors.New("tableName이 비어 있습니다")
	}

	return &UserStorage{
		client:    client,
		tableName: tableName,
	}, nil
}

// 실제 테이블에 CRUD 로직을 수행
func (s *UserStorage) GetUserByID(ctx context.Context, userID string) (*UserItem, error) {
	if s == nil || s.client == nil {
		return nil, errors.New("UserStorage가 초기화되지 않았습니다")
	}
	if userID == "" {
		return nil, errors.New("userID가 비어 있습니다")
	}

	out, err := s.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(s.tableName),
		Key: map[string]types.AttributeValue{
			"user_id": &types.AttributeValueMemberS{Value: userID},
		},
		ConsistentRead: aws.Bool(true),
	})
	if err != nil {
		return nil, fmt.Errorf("GetItem 실패: %w", err)
	}
	if out.Item == nil {
		return nil, fmt.Errorf("사용자 %s를 찾을 수 없습니다", userID)
	}

	var user UserItem
	if err := attributevalue.UnmarshalMap(out.Item, &user); err != nil {
		return nil, fmt.Errorf("사용자 언마샬 실패: %w", err)
	}

	return &user, nil
}

func (s *UserStorage) CreateUser(ctx context.Context, item *UserItem) error {
	if s == nil || s.client == nil {
		return errors.New("UserStorage가 초기화되지 않았습니다")
	}
	if item == nil {
		return errors.New("UserItem이 nil입니다")
	}
	if item.UserID == "" {
		return errors.New("UserItem.UserID가 비어 있습니다")
	}
	if item.CreatedAt.IsZero() {
		item.CreatedAt = time.Now().UTC()
	}

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("사용자 marshal 실패: %w", err)
	}

	_, err = s.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(s.tableName),
		Item:                av,
		ConditionExpression: aws.String("attribute_not_exists(user_id)"),
	})
	if err != nil {
		var ccfe *types.ConditionalCheckFailedException
		if errors.As(err, &ccfe) {
			return fmt.Errorf("이미 존재하는 사용자: %s", item.UserID)
		}
		return fmt.Errorf("PutItem 실패: %w", err)
	}

	return nil
}

func (s *UserStorage) UpdateUser(ctx context.Context, userID string, email, name *string) (*UserItem, error) {
	if s == nil || s.client == nil {
		return nil, errors.New("UserStorage가 초기화되지 않았습니다")
	}
	if userID == "" {
		return nil, errors.New("userID가 비어 있습니다")
	}
	if email == nil && name == nil {
		return nil, errors.New("업데이트할 필드가 없습니다")
	}

	updateBuilder := expression.UpdateBuilder{}
	if email != nil {
		updateBuilder = updateBuilder.Set(expression.Name("email"), expression.Value(*email))
	}
	if name != nil {
		updateBuilder = updateBuilder.Set(expression.Name("name"), expression.Value(*name))
	}

	expr, err := expression.NewBuilder().
		WithUpdate(updateBuilder).
		WithCondition(expression.AttributeExists(expression.Name("user_id"))).
		Build()
	if err != nil {
		return nil, fmt.Errorf("expression 빌드 실패: %w", err)
	}

	out, err := s.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                 aws.String(s.tableName),
		Key:                       map[string]types.AttributeValue{"user_id": &types.AttributeValueMemberS{Value: userID}},
		UpdateExpression:          expr.Update(),
		ConditionExpression:       expr.Condition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ReturnValues:              types.ReturnValueAllNew,
	})
	if err != nil {
		var ccfe *types.ConditionalCheckFailedException
		if errors.As(err, &ccfe) {
			return nil, fmt.Errorf("사용자 %s가 존재하지 않습니다", userID)
		}
		return nil, fmt.Errorf("UpdateItem 실패: %w", err)
	}

	var updated UserItem
	if err := attributevalue.UnmarshalMap(out.Attributes, &updated); err != nil {
		return nil, fmt.Errorf("업데이트 결과 언마샬 실패: %w", err)
	}

	return &updated, nil
}

func (s *UserStorage) DeleteUser(ctx context.Context, id string) error {
	if s == nil || s.client == nil {
		return errors.New("UserStorage가 초기화되지 않았습니다")
	}
	if id == "" {
		return errors.New("id가 비어 있습니다")
	}

	_, err := s.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(s.tableName),
		Key: map[string]types.AttributeValue{
			"user_id": &types.AttributeValueMemberS{Value: id},
		},
		ConditionExpression: aws.String("attribute_exists(user_id)"),
	})
	if err != nil {
		var ccfe *types.ConditionalCheckFailedException
		if errors.As(err, &ccfe) {
			return fmt.Errorf("삭제할 사용자 %s가 존재하지 않습니다", id)
		}
		return fmt.Errorf("DeleteItem 실패: %w", err)
	}

	return nil
}
