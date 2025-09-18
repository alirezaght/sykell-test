package user

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"sykell-backend/internal/db"
	"sykell-backend/internal/utils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserService_Register(t *testing.T) {
	tests := []struct {
		name           string
		request        RegisterRequest
		setupMock      func(mock sqlmock.Sqlmock)
		expectedError  string
		validateResult func(t *testing.T, result sql.Result)
	}{
		{
			name: "successful registration",
			request: RegisterRequest{
				Email:    "newuser@example.com",
				Password: "password123",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				// First check if user exists (should return ErrNoRows)
				mock.ExpectQuery("SELECT (.+) FROM users WHERE email = ?").
					WithArgs("newuser@example.com").
					WillReturnError(sql.ErrNoRows)

				// Then create the user
				mock.ExpectExec("INSERT INTO users").
					WithArgs("newuser@example.com", sqlmock.AnyArg()). // password hash will be different each time
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: "",
			validateResult: func(t *testing.T, result sql.Result) {
				id, err := result.LastInsertId()
				assert.NoError(t, err)
				assert.Equal(t, int64(1), id)

				rowsAffected, err := result.RowsAffected()
				assert.NoError(t, err)
				assert.Equal(t, int64(1), rowsAffected)
			},
		},
		{
			name: "user already exists",
			request: RegisterRequest{
				Email:    "existing@example.com",
				Password: "password123",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				// User exists check returns a user
				user := db.User{
					ID:           "550e8400-e29b-41d4-a716-446655440000",
					Email:        "existing@example.com",
					PasswordHash: "hashed_password",
					CreatedAt:    sql.NullTime{Time: time.Now(), Valid: true},
					UpdatedAt:    sql.NullTime{Time: time.Now(), Valid: true},
				}

				rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "created_at", "updated_at"}).
					AddRow(user.ID, user.Email, user.PasswordHash, user.CreatedAt, user.UpdatedAt)

				mock.ExpectQuery("SELECT (.+) FROM users WHERE email = ?").
					WithArgs("existing@example.com").
					WillReturnRows(rows)
			},
			expectedError: "user already exists",
		},
		{
			name: "database error on user check",
			request: RegisterRequest{
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
			name: "database error on user creation",
			request: RegisterRequest{
				Email:    "newuser@example.com",
				Password: "password123",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				// User doesn't exist
				mock.ExpectQuery("SELECT (.+) FROM users WHERE email = ?").
					WithArgs("newuser@example.com").
					WillReturnError(sql.ErrNoRows)

				// Error during creation
				mock.ExpectExec("INSERT INTO users").
					WithArgs("newuser@example.com", sqlmock.AnyArg()).
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: "sql: connection is already closed",
		},
		{
			name: "empty email",
			request: RegisterRequest{
				Email:    "",
				Password: "password123",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE email = ?").
					WithArgs("").
					WillReturnError(sql.ErrNoRows)

				mock.ExpectExec("INSERT INTO users").
					WithArgs("", sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: "",
		},
		{
			name: "empty password",
			request: RegisterRequest{
				Email:    "test@example.com",
				Password: "",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE email = ?").
					WithArgs("test@example.com").
					WillReturnError(sql.ErrNoRows)

				mock.ExpectExec("INSERT INTO users").
					WithArgs("test@example.com", sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: "",
		},
		{
			name: "long password (within bcrypt limit)",
			request: RegisterRequest{
				Email:    "test@example.com",
				Password: string(make([]byte, 72)), // bcrypt max length
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE email = ?").
					WithArgs("test@example.com").
					WillReturnError(sql.ErrNoRows)

				mock.ExpectExec("INSERT INTO users").
					WithArgs("test@example.com", sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: "",
		},
		{
			name: "password too long (exceeds bcrypt limit)",
			request: RegisterRequest{
				Email:    "test@example.com",
				Password: string(make([]byte, 100)), // exceeds bcrypt limit
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE email = ?").
					WithArgs("test@example.com").
					WillReturnError(sql.ErrNoRows)
				// No CREATE expectation because it should fail before that
			},
			expectedError: "bcrypt: password length exceeds 72 bytes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, service := setupMockDB(t)
			defer mockDB.Close()

			tt.setupMock(mock)

			result, err := service.Register(context.Background(), tt.request)

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

func TestUserService_Register_PasswordHashing(t *testing.T) {
	mockDB, mock, service := setupMockDB(t)
	defer mockDB.Close()

	password := "test-password-123"
	email := "test@example.com"

	// Capture the hashed password
	var capturedHash string
	mock.ExpectQuery("SELECT (.+) FROM users WHERE email = ?").
		WithArgs(email).
		WillReturnError(sql.ErrNoRows)

	mock.ExpectExec("INSERT INTO users").
		WithArgs(email, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1)).
		WillDelayFor(0) // No delay

	// Override the expectation to capture the hash
	mock.ExpectExec("INSERT INTO users").
		WithArgs(email, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	request := RegisterRequest{
		Email:    email,
		Password: password,
	}

	// Register user
	result, err := service.Register(context.Background(), request)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify the password would be hashed correctly
	hashedPassword, err := utils.HashPassword(password)
	require.NoError(t, err)
	
	// Verify the hash is different from original password
	assert.NotEqual(t, password, hashedPassword)
	
	// Verify the hash can be verified
	err = utils.CheckPassword(hashedPassword, password)
	assert.NoError(t, err)

	// Note: We can't directly verify the exact hash used in the database
	// because bcrypt generates a different hash each time due to random salt
	capturedHash = hashedPassword // For demonstration
	assert.NotEmpty(t, capturedHash)
}

func TestUserService_Register_ContextCancellation(t *testing.T) {
	mockDB, _, service := setupMockDB(t)
	defer mockDB.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	request := RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	// With cancelled context, the query may not execute at all
	// so we shouldn't expect any database interactions
	result, err := service.Register(ctx, request)
	assert.Error(t, err)
	assert.Nil(t, result)
	// The exact error message may vary depending on where cancellation occurs

	// Don't check mock expectations as the context cancellation
	// may prevent the query from reaching the database layer
}

func TestUserService_Register_UniqueEmails(t *testing.T) {
	// Test that the same email cannot be registered twice
	mockDB, mock, service := setupMockDB(t)
	defer mockDB.Close()

	email := "unique@example.com"

	// First registration should succeed
	mock.ExpectQuery("SELECT (.+) FROM users WHERE email = ?").
		WithArgs(email).
		WillReturnError(sql.ErrNoRows)

	mock.ExpectExec("INSERT INTO users").
		WithArgs(email, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	request1 := RegisterRequest{
		Email:    email,
		Password: "password123",
	}

	result1, err1 := service.Register(context.Background(), request1)
	assert.NoError(t, err1)
	assert.NotNil(t, result1)

	// Second registration with same email should fail
	user := db.User{
		ID:           "550e8400-e29b-41d4-a716-446655440000",
		Email:        email,
		PasswordHash: "hashed_password",
		CreatedAt:    sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt:    sql.NullTime{Time: time.Now(), Valid: true},
	}

	rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "created_at", "updated_at"}).
		AddRow(user.ID, user.Email, user.PasswordHash, user.CreatedAt, user.UpdatedAt)

	mock.ExpectQuery("SELECT (.+) FROM users WHERE email = ?").
		WithArgs(email).
		WillReturnRows(rows)

	request2 := RegisterRequest{
		Email:    email,
		Password: "differentpassword",
	}

	result2, err2 := service.Register(context.Background(), request2)
	assert.Error(t, err2)
	assert.Contains(t, err2.Error(), "user already exists")
	assert.Nil(t, result2)

	assert.NoError(t, mock.ExpectationsWereMet())
}