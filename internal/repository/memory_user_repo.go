package repository

import (
	"Letsgo/internal/domain"
	"context"
	"sync"
)

// InMemoryUserRepository 구조체를 만듦
// 필드로 mu와 users를 갖고 있음. users는 Map<String, User>
// mu는 무엇? -> 상호 배제의 약어(Mutual Exclusion). Go에선 map은 thread safe 하지 않기 때문에 필요하다.
type InMemoryUserRepository struct {
	mu    sync.RWMutex // 읽기 쓰기 분리 자물쇠. 동시 접근 충돌을 방지해줌
	users map[string]*domain.User
}

// 반환형이 InMemoryUserRepository의 주소(포인터)
// 생성자?
// 생성자로 만들 때, 빈 Map을 만듦. New Map<String, User> 같은거

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*domain.User), // 근데 뒤에 쉼표는 왜 붙음? -> Go 문법 규칙, 구조체/슬라이스/맵 등을 여러 줄로 작성할 때는 마지막 줄에도 반드시 , 가 들어가야함.
	}
}

// InMemoryUserRepository가 사용하는 Save라는 함수.
// 매개변수로 ctx와 user를 받음.
// error는 뭐지? -> 반환형
func (r *InMemoryUserRepository) Save(ctx context.Context, user *domain.User) error {
	// 이게 하는 일이 뭐지? -> 쓰기 잠금 획득
	r.mu.Lock()

	// defer 예약어. 당장 실행하지 않고, 함수가 끝난 뒤에 역순으로 실행됨. 그니까 뭔가를 unlock하는건데 뭔지 잘 모르겠음 -> 쓰기 잠금을 unlock
	defer r.mu.Unlock()

	// repository의 Map에 유저 할당
	r.users[user.ID] = user
	return nil
}

// 반환형이 2개 user객체와 error
func (r *InMemoryUserRepository) FindById(ctx context.Context, id string) (*domain.User, error) {
	// 뭔가 읽기 락?
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Go에선 Map에서 id로 조회할때 두 가지 값을 받게 선언하면 첫번째가 value, 두번째가 실제로 키가 존재했는지 아닌지 boolean
	user, ok := r.users[id]

	// 따라서 ok가 아니라면, 유저는 null, 에러는 유저 없음
	if !ok {
		return nil, domain.ErrUserNotFound
	}

	// 유저가 있다면 유저 리턴, 에러는 없음
	return user, nil
}

// memory_user_repo는 user_usercase에서 정의한 userRepository의 구현체이고, 명시적 implement 없이 내부 메서드 종류가 같다면 구현체 취급한다.
