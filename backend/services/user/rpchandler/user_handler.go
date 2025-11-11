package rpchandler

import (
	"context"
	"fmt"

	connect "connectrpc.com/connect"

	userpb "Acho-mj/2025_Golang_MSA/backend/gen/services/user/api"
	userconnect "Acho-mj/2025_Golang_MSA/backend/gen/services/user/connect"
	"Acho-mj/2025_Golang_MSA/backend/services/user/store"
)

type UserHandler struct {
	service *store.UserService
}

func NewUserHandler(service *store.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) CreateUser(ctx context.Context, req *connect.Request[userpb.CreateUserRequest]) (*connect.Response[userpb.CreateUserResponse], error) {
	email := req.Msg.GetEmail()
	name := req.Msg.GetName()

	if email == "" || name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("email과 name은 필수입니다"))
	}

	user, err := h.service.CreateUser(ctx, email, name)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	resp := connect.NewResponse(&userpb.CreateUserResponse{
		User: user.ToProto(),
	})

	return resp, nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *connect.Request[userpb.GetUserRequest]) (*connect.Response[userpb.GetUserResponse], error) {
	userID := req.Msg.GetUserId()
	if userID == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("user_id는 필수입니다"))
	}

	user, err := h.service.GetUser(ctx, userID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	resp := connect.NewResponse(&userpb.GetUserResponse{
		User: user.ToProto(),
	})

	return resp, nil
}

var _ userconnect.UserServiceHandler = (*UserHandler)(nil)
