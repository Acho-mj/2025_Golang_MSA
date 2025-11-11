package models

import (
	"time"

	orderpb "Acho-mj/2025_Golang_MSA/backend/gen/services/order/api"
)

type OrderItem struct {
	ProductID string `dynamodbav:"product_id"`
	Quantity  int32  `dynamodbav:"quantity"`
}

type Order struct {
	OrderID   string      `dynamodbav:"order_id"`
	UserID    string      `dynamodbav:"user_id"`
	Items     []OrderItem `dynamodbav:"items"`
	Status    string      `dynamodbav:"status"`
	CreatedAt time.Time   `dynamodbav:"created_at"`
}

func (o *Order) ToProto() *orderpb.Order {
	if o == nil {
		return nil
	}

	items := make([]*orderpb.OrderItem, 0, len(o.Items))
	for _, item := range o.Items {
		items = append(items, &orderpb.OrderItem{
			ProductId: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	return &orderpb.Order{
		OrderId:   o.OrderID,
		UserId:    o.UserID,
		Items:     items,
		Status:    o.Status,
		CreatedAt: o.CreatedAt.UTC().Format(time.RFC3339),
	}
}

func OrderFromProto(p *orderpb.Order) *Order {
	if p == nil {
		return nil
	}

	items := make([]OrderItem, 0, len(p.Items))
	for _, item := range p.Items {
		if item == nil {
			continue
		}
		items = append(items, OrderItem{
			ProductID: item.ProductId,
			Quantity:  item.Quantity,
		})
	}

	createdAt, err := time.Parse(time.RFC3339, p.CreatedAt)
	if err != nil {
		createdAt = time.Time{}
	}

	return &Order{
		OrderID:   p.OrderId,
		UserID:    p.UserId,
		Items:     items,
		Status:    p.Status,
		CreatedAt: createdAt,
	}
}

