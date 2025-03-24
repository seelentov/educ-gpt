package impl

import (
	"educ-gpt/models"
	"educ-gpt/services"
	"educ-gpt/utils/securityutils"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type UserServiceImpl struct {
	db          *gorm.DB
	logger      *zap.Logger
	defaultRole string
}

func (u UserServiceImpl) ChangePassword(userId uint, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		u.logger.Error("Error hash password", zap.Error(err))
		return fmt.Errorf("%w: %w", services.ErrHashingPassword, err)
	}

	if err := u.db.Model(&models.User{}).Where("id = ?", userId).Update("password", hashedPassword).Error; err != nil {
		u.logger.Error("Update password failed", zap.Error(err))
		return fmt.Errorf("%w:%w", services.ErrUpdateUser, err)
	}

	return nil
}

func (u UserServiceImpl) Activate(key string) error {
	var user *models.User
	if err := u.db.Model(&models.User{}).Where("activation_key = ?", key).First(&user).Error; err != nil {
		u.logger.Error("activate failed", zap.Error(err))
		return fmt.Errorf("%w:%w", services.ErrActivate, err)
	}

	user.ActivationKey = ""
	timeNow := time.Now()
	user.ActivateAt = &timeNow

	if err := u.db.Save(&user).Error; err != nil {
		u.logger.Error("activate failed", zap.Error(err))
		return fmt.Errorf("%w:%w", services.ErrActivate, err)
	}

	return nil
}

func (u UserServiceImpl) Create(user *models.User) (string, error) {
	err := u.checkUnique(user)
	if err != nil {
		u.logger.Warn("User duplicate", zap.Error(err))
		return "", fmt.Errorf("%w: %w", services.ErrDuplicate, err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		u.logger.Error("Error hash password", zap.Error(err))
		return "", fmt.Errorf("%w: %w", services.ErrHashingPassword, err)
	}

	key := securityutils.GenerateKey(200)

	user.ActivationKey = key

	user.Password = string(hashedPassword)

	tx := u.db.Begin()

	result := tx.Create(user)
	if result.Error != nil {
		tx.Rollback()
		u.logger.Error("Error creating user", zap.Error(result.Error))
		return "", fmt.Errorf("%w: %w", services.ErrCreatingUser, result.Error)
	}

	if err := tx.First(&user, user.ID).Error; err != nil {
		tx.Rollback()
		u.logger.Error("Error reloading user", zap.Error(err))
		return "", fmt.Errorf("failed to reload user: %w", err)
	}

	role := &models.Role{}
	result = u.db.Where("name = ?", u.defaultRole).First(&role)

	if result.Error != nil {
		tx.Rollback()
		u.logger.Error("Error retrieving role", zap.Error(result.Error))
		return "", fmt.Errorf("%w: %w", services.ErrRetrievingRole, result.Error)
	}

	err = tx.Model(user).Association("Roles").Append(role)
	if err != nil {
		tx.Rollback()
		u.logger.Error("Error assigning role to user", zap.Error(err))
		return "", fmt.Errorf("failed to assign role: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return "", fmt.Errorf("%w:%w", services.ErrCreatingUser, err)
	}

	return key, nil
}

func (u UserServiceImpl) Update(id uint, data map[string]interface{}) error {
	result := u.db.Model(&models.User{}).Where("id = ?", id).Updates(data)
	if result.Error != nil {
		u.logger.Error("Error update user", zap.Error(result.Error))
		return fmt.Errorf("%w: %w", services.ErrUpdateUser, result.Error)
	}

	return nil
}

func (u UserServiceImpl) GetById(id uint) (*models.User, error) {
	return u.getBy("id", id)
}

func (u UserServiceImpl) GetByName(s string) (*models.User, error) {
	return u.getBy("name", s)
}

func (u UserServiceImpl) GetByEmail(s string) (*models.User, error) {
	return u.getBy("email", s)
}

func (u UserServiceImpl) GetByCredential(s string) (*models.User, error) {
	query := "email = ? OR name = ?"

	user, err := u.getByWhere(query, s, s)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", services.ErrRetrievingUser, err)
	}
	return user, nil
}

func (u UserServiceImpl) VerifyPassword(input string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(input))
	if err != nil {
		u.logger.Error("Error comparing password hashes", zap.Error(err))
		return fmt.Errorf("%w: %w", services.ErrPasswordComparison, err)
	}

	return nil
}

func (u UserServiceImpl) checkUnique(user *models.User) error {
	var existingUser models.User
	result := u.db.Where("name = ?", user.Name).First(&existingUser)
	if result.Error == nil {
		u.logger.Warn("Duplicate name", zap.String("name", user.Name))
		return services.ErrDuplicateName
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) { // Some other database error occurred
		u.logger.Error("Error checking for duplicate name", zap.Error(result.Error))
		return fmt.Errorf("%w: %w", services.ErrRetrievingUser, result.Error)
	}

	result = u.db.Where("email = ?", user.Email).First(&existingUser)
	if result.Error == nil {
		u.logger.Warn("Duplicate email", zap.String("email", user.Email))
		return services.ErrDuplicateEmail
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		u.logger.Error("Error checking for duplicate email", zap.Error(result.Error))
		return fmt.Errorf("error checking email uniqueness: %w", result.Error)
	}

	return nil
}

func (u UserServiceImpl) getBy(key string, value interface{}) (*models.User, error) {
	return u.getByWhere(fmt.Sprintf("%s = ?", key), value)
}

func (u UserServiceImpl) getByWhere(query interface{}, args ...interface{}) (*models.User, error) {
	user := &models.User{}
	result := u.db.Where(query, args...).Preload("Roles").First(user)
	if result.Error != nil {
		u.logger.Error("Error retrieving user", zap.Error(result.Error))
		return nil, fmt.Errorf("%w: %w", services.ErrRetrievingUser, result.Error)
	}
	return user, nil
}

func (u UserServiceImpl) ClearNonActivatedUsers() error {
	threshold := time.Now().Add(-2 * time.Hour)

	result := u.db.Where("created_at < ? AND activate_at IS NULL", threshold).Delete(&models.User{})
	if result.Error != nil {
		u.logger.Error("Cant remove non active users", zap.Error(result.Error))
		return fmt.Errorf("%w:%w", services.ErrDeleteUsers, result.Error)
	}
	return nil
}

func NewUserServiceImpl(db *gorm.DB, logger *zap.Logger, defaultRole string) *UserServiceImpl {
	return &UserServiceImpl{db, logger, defaultRole}
}
