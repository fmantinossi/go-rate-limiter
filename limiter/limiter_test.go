package limiter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-rate-limiter/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Increment(ctx context.Context, key string) (int64, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockStorage) Get(ctx context.Context, key string) (int64, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockStorage) Set(ctx context.Context, key string, value int64, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockStorage) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockStorage) Exists(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

func TestRateLimiter_Allow(t *testing.T) {
	cfg := &config.Config{
		RateLimitIPRequests:      5,
		RateLimitIPWindow:        time.Second,
		RateLimitIPBlockDuration: time.Minute * 5,
	}

	tests := []struct {
		name       string
		setupMock  func(*MockStorage)
		identifier string
		isToken    bool
		want       bool
		wantErr    bool
	}{
		{
			name: "Allow request within limit",
			setupMock: func(m *MockStorage) {
				m.On("Exists", mock.Anything, "blocked:127.0.0.1").Return(false, nil)
				m.On("Increment", mock.Anything, "counter:127.0.0.1").Return(int64(1), nil)
				m.On("Set", mock.Anything, "counter:127.0.0.1", int64(1), time.Second).Return(nil)
			},
			identifier: "127.0.0.1",
			isToken:    false,
			want:       true,
			wantErr:    false,
		},
		{
			name: "Block request when limit exceeded",
			setupMock: func(m *MockStorage) {
				m.On("Exists", mock.Anything, "blocked:127.0.0.1").Return(false, nil)
				m.On("Increment", mock.Anything, "counter:127.0.0.1").Return(int64(6), nil)
				m.On("Set", mock.Anything, "blocked:127.0.0.1", int64(1), time.Minute*5).Return(nil)
				m.On("Delete", mock.Anything, "counter:127.0.0.1").Return(nil)
			},
			identifier: "127.0.0.1",
			isToken:    false,
			want:       false,
			wantErr:    false,
		},
		{
			name: "Block request when already blocked",
			setupMock: func(m *MockStorage) {
				m.On("Exists", mock.Anything, "blocked:127.0.0.1").Return(true, nil)
			},
			identifier: "127.0.0.1",
			isToken:    false,
			want:       false,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := new(MockStorage)
			tt.setupMock(mockStorage)

			rl := NewRateLimiter(mockStorage, cfg)
			got, err := rl.Allow(context.Background(), tt.identifier, tt.isToken)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
			mockStorage.AssertExpectations(t)
		})
	}
}

func TestRateLimiter_Middleware(t *testing.T) {
	cfg := &config.Config{
		RateLimitIPRequests:      5,
		RateLimitIPWindow:        time.Second,
		RateLimitIPBlockDuration: time.Minute * 5,
	}

	tests := []struct {
		name           string
		setupMock      func(*MockStorage)
		requestHeaders map[string]string
		wantStatus     int
	}{
		{
			name: "Allow request within IP limit",
			setupMock: func(m *MockStorage) {
				m.On("Exists", mock.Anything, mock.Anything).Return(false, nil)
				m.On("Increment", mock.Anything, mock.Anything).Return(int64(1), nil)
				m.On("Set", mock.Anything, mock.Anything, int64(1), mock.Anything).Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Block request when IP limit exceeded",
			setupMock: func(m *MockStorage) {
				m.On("Exists", mock.Anything, mock.Anything).Return(true, nil)
			},
			wantStatus: http.StatusTooManyRequests,
		},
		{
			name: "Allow request with valid token",
			setupMock: func(m *MockStorage) {
				m.On("Exists", mock.Anything, mock.Anything).Return(false, nil)
				m.On("Increment", mock.Anything, mock.Anything).Return(int64(1), nil)
				m.On("Set", mock.Anything, mock.Anything, int64(1), mock.Anything).Return(nil)
			},
			requestHeaders: map[string]string{
				"API_KEY": "test-token",
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := new(MockStorage)
			tt.setupMock(mockStorage)

			rl := NewRateLimiter(mockStorage, cfg)
			handler := rl.Middleware()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/", nil)
			for k, v := range tt.requestHeaders {
				c.Request.Header.Set(k, v)
			}

			handler(c)

			assert.Equal(t, tt.wantStatus, w.Code)
			mockStorage.AssertExpectations(t)
		})
	}
}
