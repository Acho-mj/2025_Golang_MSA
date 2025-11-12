package main

import (
	"context"
	"log"
	"net/http"

	userconnect "Acho-mj/2025_Golang_MSA/backend/gen/user/userconnect"

	"Acho-mj/2025_Golang_MSA/backend/internal/config"
	"Acho-mj/2025_Golang_MSA/backend/internal/storage"
	"Acho-mj/2025_Golang_MSA/backend/services/user/rpchandler"
	"Acho-mj/2025_Golang_MSA/backend/services/user/store"
)

func main() {
	// 최상위 Context 민들기
	ctx := context.Background()

	// 환경 변수 설정
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("config load 실패: %v", err)
	}

	// DynamoDB 연결
	dynamoClient, err := storage.NewDynamoClient(ctx, cfg)
	if err != nil {
		log.Fatalf("dynamodb 초기화 실패: %v", err)
	}

	userStorage, err := storage.NewUserStorage(dynamoClient, cfg.DynamoUserTable)
	if err != nil {
		log.Fatalf("user storage 초기화 실패: %v", err)
	}

	// 핸들러
	userService := store.NewUserService(userStorage)
	userHandler := rpchandler.NewUserHandler(userService)

	mux := http.NewServeMux()
	path, handler := userconnect.NewUserServiceHandler(userHandler)
	mux.Handle(path, handler)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	addr := ":" + cfg.Port
	log.Printf("user service listening on %s", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("서버 종료: %v", err)
	}
}
