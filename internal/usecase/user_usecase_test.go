package usecase

import (
	"Letsgo/internal/domain"
	"context"
	"errors"
	"testing"
)

// 가짜 저장소 작성
// UserRepository 인터페이스를 만족하는 Mock 구조체
type mockUserRepository struct {
	saveFunc     func(ctx context.Context, user *domain.User) error
	findByIdFunc func(ctx context.Context, id string) (*domain.User, error)
}

// 인터페이스의 Save 함수 구현
func (m *mockUserRepository) Save(ctx context.Context, user *domain.User) error {
	if m.saveFunc != nil {
		return m.saveFunc(ctx, user)
	}

	return nil
}

// 인터페이스의 FindById 구현
func (m *mockUserRepository) FindById(ctx context.Context, id string) (*domain.User, error) {
	if m.findByIdFunc != nil {
		return m.findByIdFunc(ctx, id)
	}

	return nil, nil
}

// 테이블 주도 테스트 작성
func TestUserUseCase_Register(t *testing.T) {
	errDBFailed := errors.New("database connection failed")

	// 1. 테스트 케이스 정의 (테이블 형태)
	tests := []struct {
		name        string
		input       RegisterUserDTO
		mockRepo    UserRepository
		wantErr     bool
		expectedErr error
	}{
		// 테스트 시나리오에 맞춰 값 등록
		{
			name: "성공: 정상적인 유저 등록",
			input: RegisterUserDTO{
				ID:    "user-1",
				Email: "test@example.com",
			},
			mockRepo: &mockUserRepository{
				saveFunc: func(ctx context.Context, user *domain.User) error {
					// 정상 저장이므로 nil 반환
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "실패: 유효하지 않은 이메일 포맷",
			input: RegisterUserDTO{
				ID:    "user-2",
				Email: "invalid-email-format",
			},
			mockRepo: &mockUserRepository{
				saveFunc: func(ctx context.Context, user *domain.User) error {
					t.Fatal("이메일 검증 실패 시 Save가 호출되어서는 안 됩니다")
					return nil
				},
			},
			wantErr:     true,
			expectedErr: domain.ErrInvalidEmail,
		},
		{
			name: "실패: 저장소(DB) 저장 실패",
			input: RegisterUserDTO{
				ID:    "user-3",
				Email: "db-fail@example.com",
			},
			mockRepo: &mockUserRepository{
				saveFunc: func(ctx context.Context, user *domain.User) error {
					return errDBFailed // DB 에러 시뮬레이션
				},
			},
			wantErr: true,
		},
	}

	// 테스트 케이스를 순회하면서 테스트 시작
	for _, tt := range tests {
		// 각 테스트를 독립된 서브테스크로 실행
		t.Run(tt.name, func(t *testing.T) {
			uc := NewUserUseCase(tt.mockRepo)

			user, err := uc.Register(context.Background(), tt.input)

			// 에러 검증
			if (err != nil) != tt.wantErr {
				t.Fatalf("Register() error = %v, wantErr %v", err, tt.wantErr)
			}

			// 특정 에러가 와야 하는 경우 에러 일치 여부 확인
			if tt.expectedErr != nil && !errors.Is(err, tt.expectedErr) {
				t.Errorf("Register() error = %v, expectedErr %v", err, tt.expectedErr)
			}

			// 성공 케이스인 경우 유저 객체 검증
			if !tt.wantErr {
				if user == nil {
					t.Fatal("user가 nil이면 안 됩니다")
				}
				if user.Email != tt.input.Email {
					t.Errorf("user.Email = %s, want %s", user.Email, tt.input.Email)
				}
			}
		})
	}
}
