package user

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"sykell-backend/internal/config"
	"sykell-backend/internal/db"
	"sykell-backend/internal/utils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *UserService) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	cfg := &config.Config{
		JWTSecret: "test-jwt-secret",
	}

	service := &UserService{
		db:     mockDB,
		config: cfg,
	}

	return mockDB, mock, service
}

func TestUserService_Login(t *testing.T) {
	tests := []struct {
		name           string
		request        LoginRequest
		setupMock      func(mock sqlmock.Sqlmock)
		expectedError  string
		validateResult func(t *testing.T, result *LoginResponse)
	}{
		{
			name: "successful login",
			request: LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				hashedPassword, _ := utils.HashPassword("password123")
				user := db.User{
					ID:           "550e8400-e29b-41d4-a716-446655440000",
					Email:        "test@example.com",
					PasswordHash: hashedPassword,
					CreatedAt:    sql.NullTime{Time: time.Now(), Valid: true},
					UpdatedAt:    sql.NullTime{Time: time.Now(), Valid: true},
				}

				rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "created_at", "updated_at"}).
					AddRow(user.ID, user.Email, user.PasswordHash, user.CreatedAt, user.UpdatedAt)

				mock.ExpectQuery("SELECT (.+) FROM users WHERE email = ?").
					WithArgs("test@example.com").
					WillReturnRows(rows)
			},
			expectedError: "",
			validateResult: func(t *testing.T, result *LoginResponse) {
				assert.NotEmpty(t, result.Token)
				assert.Equal(t, "test@example.com", result.User.Email)
				assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", result.User.ID)
				// UserResponse doesn't have PasswordHash field, which is correct for security
				assert.Greater(t, result.ExpiresAt, time.Now().Unix())
			},
		},
		{
			name: "user not found",
			request: LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE email = ?").
					WithArgs("nonexistent@example.com").
					WillReturnError(sql.ErrNoRows)
			},
			expectedError: "invalid email or password",
		},
		{
			name: "wrong password",
			request: LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				hashedPassword, _ := utils.HashPassword("correctpassword")
				user := db.User{
					ID:           "550e8400-e29b-41d4-a716-446655440000",
					Email:        "test@example.com",
					PasswordHash: hashedPassword,
					CreatedAt:    sql.NullTime{Time: time.Now(), Valid: true},
					UpdatedAt:    sql.NullTime{Time: time.Now(), Valid: true},
				}

				rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "created_at", "updated_at"}).
					AddRow(user.ID, user.Email, user.PasswordHash, user.CreatedAt, user.UpdatedAt)

				mock.ExpectQuery("SELECT (.+) FROM users WHERE email = ?").
					WithArgs("test@example.com").
					WillReturnRows(rows)
			},
			expectedError: "invalid email or password",
		},
		{
			name: "database error",
			request: LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE email = ?").
					WithArgs("test@example.com").
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: "sql: connection is already closed",
		},
		{
			name: "empty email",
			request: LoginRequest{
				Email:    "",
				Password: "password123",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE email = ?").
					WithArgs("").
					WillReturnError(sql.ErrNoRows)
			},
			expectedError: "invalid email or password",
		},
		{
			name: "empty password",
			request: LoginRequest{
				Email:    "test@example.com",
				Password: "",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				hashedPassword, _ := utils.HashPassword("realpassword")
				user := db.User{
					ID:           "550e8400-e29b-41d4-a716-446655440000",
					Email:        "test@example.com",
					PasswordHash: hashedPassword,
					CreatedAt:    sql.NullTime{Time: time.Now(), Valid: true},
					UpdatedAt:    sql.NullTime{Time: time.Now(), Valid: true},
				}

				rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "created_at", "updated_at"}).
					AddRow(user.ID, user.Email, user.PasswordHash, user.CreatedAt, user.UpdatedAt)

				mock.ExpectQuery("SELECT (.+) FROM users WHERE email = ?").
					WithArgs("test@example.com").
					WillReturnRows(rows)
			},
			expectedError: "invalid email or password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, service := setupMockDB(t)
			defer mockDB.Close()

			tt.setupMock(mock)

			result, err := service.Login(context.Background(), tt.request)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tt.validateResult != nil {
					tt.validateResult(t, result)
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserService_Login_JWTGeneration(t *testing.T) {
	mockDB, mock, service := setupMockDB(t)
	defer mockDB.Close()

	hashedPassword, _ := utils.HashPassword("password123")
	user := db.User{
		ID:           "550e8400-e29b-41d4-a716-446655440000",
		Email:        "test@example.com",
		PasswordHash: hashedPassword,
		CreatedAt:    sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt:    sql.NullTime{Time: time.Now(), Valid: true},
	}

	rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "created_at", "updated_at"}).
		AddRow(user.ID, user.Email, user.PasswordHash, user.CreatedAt, user.UpdatedAt)

	mock.ExpectQuery("SELECT (.+) FROM users WHERE email = ?").
		WithArgs("test@example.com").
		WillReturnRows(rows)

	request := LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	result, err := service.Login(context.Background(), request)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Validate JWT token
	claims, err := utils.ValidateJWT(result.Token, []byte(service.config.JWTSecret))
	assert.NoError(t, err)
	assert.Equal(t, user.ID, claims.UserID)
	assert.Equal(t, user.Email, claims.Email)
	assert.Equal(t, "sykell-backend", claims.Issuer)
	assert.Equal(t, user.Email, claims.Subject)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserService_Login_ContextCancellation(t *testing.T) {
	mockDB, _, service := setupMockDB(t)
	defer mockDB.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	req := LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	// With cancelled context, the query may not execute at all
	// so we shouldn't expect any database interactions
	result, err := service.Login(ctx, req)
	assert.Error(t, err)
	assert.Nil(t, result)
	// The exact error message may vary depending on where cancellation occurs

	// Don't check mock expectations as the context cancellation
	// may prevent the query from reaching the database layer
}