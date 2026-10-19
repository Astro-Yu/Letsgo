package usecase

// usecase는 @Service에 해당

import (
	"Letsgo/internal/domain"
	"context"
	"errors"
)

// usecase가 필요로 하는 저장소 명세를 직접 인터페이스로 정의
// usecase는 Save라는 메서드, FindById라는 메서드를 갖고 있음
// ctx: Go에는 스레드로컬 개념이 없기 때문에(대신 경량 스레드인 고루틴을 많이 띄움), 레이어간의 컨텍스트를 유지할 객체가 필요하다. 그것이 ctx
// 역할 1. 취소 및 타임아웃 전파: 클라이언트가 HTTP 요청을 끊거나 타임아웃이 발생하면 DB쿼리나 외부 API 호출까지 즉시 중단한다.
// 역할 2. 메타데이터 전달: TraceID, 인증 토큰 등 컨텍스트 데이터를 계층 간 전달한다.
type UserRepository interface {
	Save(ctx context.Context, user *domain.User) error
	FindById(ctx context.Context, id string) (*domain.User, error)
}

// type: 새로운 타입을 선언하는 키워드. (인터페이스, 구조체, 함수타입, 타입별칭 에 사용)
// struct: 데이터 필드들의 묶음. record와 유사
type RegisterUserDTO struct {
	ID    string
	Email string
}

type UserUseCase struct {
	repo UserRepository
}

// 의존성 주입 팩토리 함수(생성자?)
// 그럼 struct를 사용하면 기본 생성자는 알아서 생기는건가. -> 생성자라는 개념 자체가 없다.
func NewUserUseCase(repo UserRepository) *UserUseCase {
	return &UserUseCase{repo: repo}
}

// := 는 어떤 선언방식 -> 단축 변수 선언문. var user *domain.User = ... 을 짧게 줄여서 쓴거임. 원래는 타입이 필요함.
// 새로운 변수를 선언할 때만 사용 가능, 기존 변수에 재할당은 = 를 사용
// (uc *UserUseCase) 는 리시버. Java의 this와 동일. 즉, Register라는 함수는 UserUseCase의 함수임 uc.Register(...) 처럼 사용 가능
func (uc *UserUseCase) Register(ctx context.Context, dto RegisterUserDTO) (*domain.User, error) {
	user, err := domain.NewUser(dto.ID, dto.Email) // NewUser 자체가 2개의 값을 반환하기 때문에 user와 err에 동시에 할당이 된다. 이건 편하네.
	if err != nil {
		return nil, err
	}

	// 조건식이 어떻게 되는거지? ;가 왜 붙은거임 -> if문의 초기화 구문
	// 조건식을 평가하기 직전에 변수를 하나 선언해 실행하고, 그 변수를 바로 조건식에서 검사하는 문법
	// 따라서 err에 repo.Save(...)를 하고 나온 값이 nil이 아니라면 조건문 실행. 해당 값은 조건문 블록에서만 유효하다.
	if err := uc.repo.Save(ctx, user); err != nil {
		return nil, errors.Join(errors.New("failed to persist user"), err)
	}

	return user, nil
}
