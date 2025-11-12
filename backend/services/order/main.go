package main

import (
	"context"
	"log"
	"net/http"

	orderconnect "Acho-mj/2025_Golang_MSA/backend/gen/order/orderconnect"
	userconnect "Acho-mj/2025_Golang_MSA/backend/gen/user/userconnect"

	"Acho-mj/2025_Golang_MSA/backend/internal/config"
	"Acho-mj/2025_Golang_MSA/backend/internal/storage"
	"Acho-mj/2025_Golang_MSA/backend/services/order/rpchandler"
	"Acho-mj/2025_Golang_MSA/backend/services/order/store"
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("config load 실패: %v", err)
	}

	dynamoClient, err := storage.NewDynamoClient(ctx, cfg)
	if err != nil {
		log.Fatalf("dynamodb 초기화 실패: %v", err)
	}

	orderStorage, err := storage.NewOrderStorage(dynamoClient, cfg.DynamoOrderTable)
	if err != nil {
		log.Fatalf("order storage 초기화 실패: %v", err)
	}

	userClient := userconnect.NewUserServiceClient(
		http.DefaultClient,
		cfg.UserServiceURL,
	)

	orderService := store.NewOrderService(orderStorage, userClient)
	orderHandler := rpchandler.NewOrderHandler(orderService)

	mux := http.NewServeMux()
	path, handler := orderconnect.NewOrderServiceHandler(orderHandler)
	mux.Handle(path, handler)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	addr := ":" + cfg.Port
	log.Printf("order service listening on %s", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("서버 종료: %v", err)
	}
}
