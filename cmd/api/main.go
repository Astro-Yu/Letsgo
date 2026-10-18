package api

import (
	"Letsgo/internal/handler"
	"Letsgo/internal/repository"
	"Letsgo/internal/usecase"
	"log"
	"net/http"
)

func main() {
	// 1. Adaptor 생성
	userRepo := repository.NewInMemoryUserRepository()

	// 2. UseCase 주입
	userUseCase := usecase.NewUserUseCase(userRepo)

	// 3. Inbound Adaptor 주입
	UserHandler := handler.NewUserHandler(userUseCase)

	// 4. 라우팅 등록
	mux := http.NewServeMux()
	mux.HandleFunc("/users", UserHandler.Register)

	log.Println("Server running on: 8080...")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
