package user

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"sykell-backend/internal/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUserService_GetProfile(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		setupMock      func(mock sqlmock.Sqlmock)
		expectedError  string
		validateResult func(t *testing.T, result *UserResponse)
	}{
		{
			name:   "successful profile retrieval",
			userID: "550e8400-e29b-41d4-a716-446655440000",
			setupMock: func(mock sqlmock.Sqlmock) {
				user := db.User{
					ID:           "550e8400-e29b-41d4-a716-446655440000",
					Email:        "test@example.com",
					PasswordHash: "hashed_password",
					CreatedAt:    sql.NullTime{Time: time.Now(), Valid: true},
					UpdatedAt:    sql.NullTime{Time: time.Now(), Valid: true},
				}

				rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "created_at", "updated_at"}).
					AddRow(user.ID, user.Email, user.PasswordHash, user.CreatedAt, user.UpdatedAt)

				mock.ExpectQuery("SELECT (.+) FROM users WHERE id = ?").
					WithArgs("550e8400-e29b-41d4-a716-446655440000").
					WillReturnRows(rows)
			},
			expectedError: "",
			validateResult: func(t *testing.T, result *UserResponse) {
				assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", result.ID)
				assert.Equal(t, "test@example.com", result.Email)
				// UserResponse doesn't have PasswordHash field, which is correct for security
				assert.NotNil(t, result.CreatedAt)
				assert.NotNil(t, result.UpdatedAt)
			},
		},
		{
			name:   "user not found",
			userID: "nonexistent-id",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE id = ?").
					WithArgs("nonexistent-id").
					WillReturnError(sql.ErrNoRows)
			},
			expectedError: "user not found",
		},
		{
			name:   "database error",
			userID: "550e8400-e29b-41d4-a716-446655440000",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE id = ?").
					WithArgs("550e8400-e29b-41d4-a716-446655440000").
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: "sql: connection is already closed",
		},
		{
			name:   "empty user ID",
			userID: "",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE id = ?").
					WithArgs("").
					WillReturnError(sql.ErrNoRows)
			},
			expectedError: "user not found",
		},
		{
			name:   "invalid UUID format",
			userID: "invalid-uuid-format",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE id = ?").
					WithArgs("invalid-uuid-format").
					WillReturnError(sql.ErrNoRows)
			},
			expectedError: "user not found",
		},
		{
			name:   "user with null timestamps",
			userID: "550e8400-e29b-41d4-a716-446655440000",
			setupMock: func(mock sqlmock.Sqlmock) {
				user := db.User{
					ID:           "550e8400-e29b-41d4-a716-446655440000",
					Email:        "test@example.com",
					PasswordHash: "hashed_password",
					CreatedAt:    sql.NullTime{Valid: false}, // NULL timestamp
					UpdatedAt:    sql.NullTime{Valid: false}, // NULL timestamp
				}

				rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "created_at", "updated_at"}).
					AddRow(user.ID, user.Email, user.PasswordHash, user.CreatedAt, user.UpdatedAt)

				mock.ExpectQuery("SELECT (.+) FROM users WHERE id = ?").
					WithArgs("550e8400-e29b-41d4-a716-446655440000").
					WillReturnRows(rows)
			},
			expectedError: "",
			validateResult: func(t *testing.T, result *UserResponse) {
				assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", result.ID)
				assert.Equal(t, "test@example.com", result.Email)
				// UserResponse doesn't have PasswordHash field, which is correct for security
				assert.Nil(t, result.CreatedAt, "created_at should be null")
				assert.Nil(t, result.UpdatedAt, "updated_at should be null")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, service := setupMockDB(t)
			defer mockDB.Close()

			tt.setupMock(mock)

			result, err := service.GetProfile(context.Background(), tt.userID)

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


func TestUserService_GetProfile_ContextCancellation(t *testing.T) {
	mockDB, _, service := setupMockDB(t)
	defer mockDB.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	userID := "550e8400-e29b-41d4-a716-446655440000"

	// With cancelled context, the query may not execute at all
	// so we shouldn't expect any database interactions
	result, err := service.GetProfile(ctx, userID)
	assert.Error(t, err)
	assert.Nil(t, result)
	// The exact error message may vary depending on where cancellation occurs

	// Don't check mock expectations as the context cancellation
	// may prevent the query from reaching the database layer
}

func TestUserService_GetProfile_DatabaseRowScan(t *testing.T) {
	// Test handling of database row scanning errors
	mockDB, mock, service := setupMockDB(t)
	defer mockDB.Close()

	userID := "550e8400-e29b-41d4-a716-446655440000"

	// Create a row with wrong column types to cause scan error
	rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "created_at", "updated_at"}).
		AddRow(123, "test@example.com", "hash", "invalid-time", "invalid-time") // id as int instead of string

	mock.ExpectQuery("SELECT (.+) FROM users WHERE id = ?").
		WithArgs(userID).
		WillReturnRows(rows)

	result, err := service.GetProfile(context.Background(), userID)
	assert.Error(t, err)
	assert.Nil(t, result)
	// The exact error message may vary depending on the SQL driver

	assert.NoError(t, mock.ExpectationsWereMet())
}
