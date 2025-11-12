package rpchandler

import (
	"context"
	"fmt"

	connect "connectrpc.com/connect"

	orderpb "Acho-mj/2025_Golang_MSA/backend/gen/order"
	orderconnect "Acho-mj/2025_Golang_MSA/backend/gen/order/orderconnect"
	"Acho-mj/2025_Golang_MSA/backend/services/order/models"
	"Acho-mj/2025_Golang_MSA/backend/services/order/store"
)

type OrderHandler struct {
	service *store.OrderService
}

func NewOrderHandler(service *store.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *connect.Request[orderpb.CreateOrderRequest]) (*connect.Response[orderpb.CreateOrderResponse], error) {
	userID := req.Msg.GetUserId()
	items := req.Msg.GetItems()

	if userID == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("user_id는 필수입니다"))
	}
	if len(items) == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("items는 최소 한 개 이상이어야 합니다"))
	}

	modelItems := make([]models.OrderItem, 0, len(items))
	for _, item := range items {
		if item == nil || item.ProductId == "" || item.Quantity <= 0 {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("상품 정보가 올바르지 않습니다"))
		}
		modelItems = append(modelItems, models.OrderItem{
			ProductID: item.ProductId,
			Quantity:  item.Quantity,
		})
	}

	order, err := h.service.CreateOrder(ctx, userID, modelItems)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	resp := connect.NewResponse(&orderpb.CreateOrderResponse{
		Order: order.ToProto(),
	})
	return resp, nil
}

func (h *OrderHandler) GetOrder(ctx context.Context, req *connect.Request[orderpb.GetOrderRequest]) (*connect.Response[orderpb.GetOrderResponse], error) {
	orderID := req.Msg.GetOrderId()
	if orderID == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("order_id는 필수입니다"))
	}

	order, err := h.service.GetOrder(ctx, orderID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	resp := connect.NewResponse(&orderpb.GetOrderResponse{
		Order: order.ToProto(),
	})
	return resp, nil
}

var _ orderconnect.OrderServiceHandler = (*OrderHandler)(nil)
