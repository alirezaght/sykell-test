package user

import (
	"context"
	"database/sql"
	"sykell-backend/internal/config"
	"sykell-backend/internal/db"
)

// Repo defines the interface for user repository operations
type Repo interface {
	// GetByEmail retrieves a user by their email
	GetByEmail(ctx context.Context, email string) (UserResponse, error)
	// Create adds a new user to the database
	Create(ctx context.Context, email string, passwordHash string) error
	// GetByID retrieves a user by their ID
	GetByID(ctx context.Context, id string) (UserResponse, error)
}

// userRepo is the concrete implementation of the Repo interface
type userRepo struct {
	sqlDB *sql.DB
}

// NewRepo creates a new instance of userRepo
func NewRepo(db *sql.DB) Repo {
	return &userRepo{
		sqlDB: db,
	}
}

// GetByEmail retrieves a user by their email
func (r *userRepo) GetByEmail(ctx context.Context, email string) (UserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DefaultTimeout)
	defer cancel()
	queries := db.New(r.sqlDB)
	user, err := queries.GetUserByEmail(context.Background(), email)
	if err != nil {
		return UserResponse{}, err
	}
	return toUserResponse(user), nil
}

// Create adds a new user to the database
func (r *userRepo) Create(ctx context.Context, email string, passwordHash string) error {
	ctx, cancel := context.WithTimeout(ctx, config.DefaultTimeout)
	defer cancel()
	queries := db.New(r.sqlDB)
	_, err := queries.CreateUser(ctx, db.CreateUserParams{
		Email:        email,
		PasswordHash: passwordHash,
	})
	
	return err
}

// GetByID retrieves a user by their ID
func (r *userRepo) GetByID(ctx context.Context, id string) (UserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DefaultTimeout)
	defer cancel()
	queries := db.New(r.sqlDB)
	user, err := queries.GetUser(ctx, id)
	if err != nil {
		return UserResponse{}, err
	}
	return toUserResponse(user), nil
}



// toUserResponse converts a db.User to UserResponse, handling sql.NullTime properly
func toUserResponse(user db.User) UserResponse {
	resp := UserResponse{
		ID:    user.ID,
		Email: user.Email,
		PasswordHash: user.PasswordHash,
	}
	
	// Handle created_at
	if user.CreatedAt.Valid {
		resp.CreatedAt = &user.CreatedAt.Time
	}
	
	// Handle updated_at
	if user.UpdatedAt.Valid {
		resp.UpdatedAt = &user.UpdatedAt.Time
	}
	
	return resp
}