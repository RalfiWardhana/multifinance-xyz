package repository

import (
	"context"
	"fmt"
	"pt-xyz-multifinance/internal/domain/entity"
	"pt-xyz-multifinance/internal/domain/repository"

	"gorm.io/gorm"
)

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) Create(ctx context.Context, user *entity.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *userRepositoryImpl) GetByID(ctx context.Context, id uint64) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return &user, nil
}

func (r *userRepositoryImpl) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	return &user, nil
}

func (r *userRepositoryImpl) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, nil
}

func (r *userRepositoryImpl) Update(ctx context.Context, user *entity.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *userRepositoryImpl) Delete(ctx context.Context, id uint64) error {
	if err := r.db.WithContext(ctx).Delete(&entity.User{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (r *userRepositoryImpl) GetAll(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	var users []*entity.User
	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	return users, nil
}

func (r *userRepositoryImpl) GetByRole(ctx context.Context, role entity.UserRole, limit, offset int) ([]*entity.User, error) {
	var users []*entity.User
	if err := r.db.WithContext(ctx).Where("role = ?", role).Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to get users by role: %w", err)
	}
	return users, nil
}
