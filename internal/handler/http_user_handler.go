package handler

import (
	"Letsgo/internal/usecase"
	"encoding/json"
	"net/http"
)

// Java로 치면 request dto
type RegisterRequest struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

// Controller, UserUsecase를 필드로 갖고 있음.
type UserHandler struct {
	uc *usecase.UserUseCase
}

// 생성자
// uc: uc는 무엇인지?
func NewUserHandler(uc *usecase.UserUseCase) *UserHandler {
	// &를 통해서 생성된 struct의 메모리 주소를 가져옴. 주소 추출
	// Go에는 생성자 메서드 문법이 없어서, 중괄호를 써서 직접 구조체 필드에 값을 넣는 문법을 쓴다.
	return &UserHandler{uc: uc} // 이 과정이 DI
}

// UserHandler가 사용할 수 있는 Register 함수.
// 반환형이 없다? 에러도 없다? void 반환형은 아무것도 안 써도 되나? -> 안써도 됨
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Post 요청이 아닐 경우 예외 처리
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	// 클라이언트가 보낸 http 요청 바디를 Go의 구조체에 맞게 역직렬화 하는 과정
	// 이 과정에서 json 문법이 맞지 않거나, 형식이 맞지 않으면 err가 발생
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	user, err := h.uc.Register(r.Context(), usecase.RegisterUserDTO{
		ID:    req.ID,
		Email: req.Email,
	})

	// 비지니스 로직 수행 중 예외가 나면 발생
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 응답 생성
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(user) // 마지막 이 과정은 왜 필요한지?
}
