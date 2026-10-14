package domain // 현재 패키지

import ( // 사용할 라이브러리 설정
	"errors"
	"strings"
	"time"
)

var ( // 예외 정의
	ErrInvalidEmail = errors.New("invalid email address")
	ErrUserNotFound = errors.New("user not found")
)

type User struct { // User의 필드 설정
	ID        string
	Email     string
	CreatedAt time.Time
}

// 생성자. 생성할때 id와 email을 string 타입으로 받음.
// 반환형은 User(포인터)와 에러. 따라서 return 할때, User 혹은 nil(== null) 과 에러를 반환함. 순서는 반드시 지켜야함.
// 대문자로 시작하면 public, 소문자로 시작하면 private
// *User
// &User = 주소 연산자. User 인스턴스를 만들고, 그 인스턴스의 메모리 주소(=포인터)를 반환함. 반환타입 *User.
// Go는 C 처럼 값 복사가 기본이기 때문에(Java의 경우 객체 변수가 알아서 힙 메모리를 가리킴) *User 처럼 해줘야 주소를 가리키게 된다.
func NewUser(id, email string) (*User, error) {
	trimmedEmail := strings.TrimSpace(email)
	if !strings.Contains(trimmedEmail, "@") {
		return nil, ErrInvalidEmail
	}

	return &User{
		ID:        id,
		Email:     trimmedEmail,
		CreatedAt: time.Now(),
	}, nil
}
