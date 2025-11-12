package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	userpb "Acho-mj/2025_Golang_MSA/backend/gen/user"
	userconnect "Acho-mj/2025_Golang_MSA/backend/gen/user/userconnect"
	"Acho-mj/2025_Golang_MSA/backend/internal/storage"
	"Acho-mj/2025_Golang_MSA/backend/services/order/models"

	connect "connectrpc.com/connect"
)

var (
	ErrInvalidInput   = errors.New("잘못된 입력입니다")
	ErrOrderNotFound  = errors.New("주문을 찾을 수 없습니다")
	defaultOrderState = "pending"
)

type OrderService struct {
	storage    *storage.OrderStorage
	userClient userconnect.UserServiceClient
}

func NewOrderService(storage *storage.OrderStorage, userClient userconnect.UserServiceClient) *OrderService {
	return &OrderService{
		storage:    storage,
		userClient: userClient,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, userID string, items []models.OrderItem) (*models.Order, error) {
	if userID == "" {
		return nil, fmt.Errorf("%w: userID는 필수입니다", ErrInvalidInput)
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("%w: 최소 한 개의 상품이 필요합니다", ErrInvalidInput)
	}

	if s.userClient == nil {
		return nil, fmt.Errorf("user 서비스 클라이언트가 초기화되지 않았습니다")
	}

	if err := s.ensureUserExists(ctx, userID); err != nil {
		return nil, err
	}

	recordItems := make([]storage.OrderLine, 0, len(items))
	for _, item := range items {
		if item.ProductID == "" || item.Quantity <= 0 {
			return nil, fmt.Errorf("%w: 상품 ID와 수량은 필수입니다", ErrInvalidInput)
		}
		recordItems = append(recordItems, storage.OrderLine{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	record := &storage.OrderRecord{
		OrderID:   generateOrderID(),
		UserID:    userID,
		Items:     recordItems,
		Status:    defaultOrderState,
		CreatedAt: time.Now().UTC(),
	}

	if err := s.storage.CreateOrder(ctx, record); err != nil {
		return nil, err
	}

	return &models.Order{
		OrderID:   record.OrderID,
		UserID:    record.UserID,
		Items:     items,
		Status:    record.Status,
		CreatedAt: record.CreatedAt,
	}, nil
}

func (s *OrderService) GetOrder(ctx context.Context, orderID string) (*models.Order, error) {
	if orderID == "" {
		return nil, fmt.Errorf("%w: orderID는 필수입니다", ErrInvalidInput)
	}

	record, err := s.storage.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	items := make([]models.OrderItem, 0, len(record.Items))
	for _, item := range record.Items {
		items = append(items, models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	return &models.Order{
		OrderID:   record.OrderID,
		UserID:    record.UserID,
		Items:     items,
		Status:    record.Status,
		CreatedAt: record.CreatedAt,
	}, nil
}

func generateOrderID() string {
	return fmt.Sprintf("order-%d", time.Now().UnixNano())
}

func (s *OrderService) ensureUserExists(ctx context.Context, userID string) error {
	req := connect.NewRequest(&userpb.GetUserRequest{
		UserId: userID,
	})
	_, err := s.userClient.GetUser(ctx, req)
	if err == nil {
		return nil
	}

	var connectErr *connect.Error
	if errors.As(err, &connectErr) {
		switch connectErr.Code() {
		case connect.CodeNotFound:
			return fmt.Errorf("%w: 사용자 %s를 찾을 수 없습니다", ErrInvalidInput, userID)
		case connect.CodeInvalidArgument:
			return fmt.Errorf("%w: userID %s가 올바르지 않습니다", ErrInvalidInput, userID)
		}
	}

	return fmt.Errorf("user 서비스 호출 실패: %w", err)
}
