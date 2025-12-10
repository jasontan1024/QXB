package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"gorm.io/gorm"
)

// User 用户模型（兼容旧代码）
type User struct {
	ID            int64
	Email         string
	Address       string
	EncPrivKeyB64 string
	EncSaltB64    string
	PassSaltB64   string
	PasswordHash  string
	CreatedAt     time.Time
}

// Service 负责用户注册/登录以及密钥管理
type Service struct {
	db *gorm.DB
}

// NewService 创建服务并初始化表结构
func NewService(db *gorm.DB) (*Service, error) {
	s := &Service{db: db}
	if err := s.initSchema(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Service) initSchema() error {
	// 使用 GORM AutoMigrate 自动创建表
	if err := s.db.AutoMigrate(&UserModel{}, &ClaimLockModel{}); err != nil {
		return fmt.Errorf("自动迁移表结构失败: %w", err)
	}
	return nil
}

// Register 注册并返回用户与地址
func (s *Service) Register(email, password string) (*User, error) {
	// 生成密钥对
	key, err := crypto.GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("生成密钥失败: %w", err)
	}
	privBytes := crypto.FromECDSA(key)
	address := crypto.PubkeyToAddress(key.PublicKey).Hex()

	passHash, passSalt, err := hashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("密码哈希失败: %w", err)
	}
	encPriv, encSalt, err := encryptPrivateKey(password, privBytes)
	if err != nil {
		return nil, fmt.Errorf("加密私钥失败: %w", err)
	}

	userModel := &UserModel{
		Email:         email,
		Address:       address,
		EncPrivKeyB64: encPriv,
		EncSaltB64:    encSalt,
		PassSaltB64:   passSalt,
		PasswordHash:  passHash,
		CreatedAt:     time.Now(),
	}

	if err := s.db.Create(userModel).Error; err != nil {
		return nil, err
	}

	return &User{
		ID:            userModel.ID,
		Email:         userModel.Email,
		Address:       userModel.Address,
		EncPrivKeyB64: userModel.EncPrivKeyB64,
		EncSaltB64:    userModel.EncSaltB64,
		PassSaltB64:   userModel.PassSaltB64,
		PasswordHash:  userModel.PasswordHash,
		CreatedAt:     userModel.CreatedAt,
	}, nil
}

// Authenticate 验证用户并返回用户信息
func (s *Service) Authenticate(email, password string) (*User, error) {
	var userModel UserModel
	if err := s.db.Where("email = ?", email).First(&userModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在或密码错误")
		}
		return nil, err
	}

	if !verifyPassword(password, userModel.PasswordHash, userModel.PassSaltB64) {
		return nil, errors.New("用户不存在或密码错误")
	}

	return &User{
		ID:            userModel.ID,
		Email:         userModel.Email,
		Address:       userModel.Address,
		EncPrivKeyB64: userModel.EncPrivKeyB64,
		EncSaltB64:    userModel.EncSaltB64,
		PassSaltB64:   userModel.PassSaltB64,
		PasswordHash:  userModel.PasswordHash,
		CreatedAt:     userModel.CreatedAt,
	}, nil
}

// GetByID 获取用户
func (s *Service) GetByID(id int64) (*User, error) {
	var userModel UserModel
	if err := s.db.First(&userModel, id).Error; err != nil {
		return nil, err
	}

	return &User{
		ID:            userModel.ID,
		Email:         userModel.Email,
		Address:       userModel.Address,
		EncPrivKeyB64: userModel.EncPrivKeyB64,
		EncSaltB64:    userModel.EncSaltB64,
		PassSaltB64:   userModel.PassSaltB64,
		PasswordHash:  userModel.PasswordHash,
		CreatedAt:     userModel.CreatedAt,
	}, nil
}

// DecryptPrivateKey 使用用户密码解密私钥
func (s *Service) DecryptPrivateKey(u *User, password string) ([]byte, error) {
	// 密码错误会导致解密失败，直接返回错误
	return decryptPrivateKey(password, u.EncPrivKeyB64, u.EncSaltB64)
}

// IsClaimLocked 检查用户在指定日期是否已提交领取
func (s *Service) IsClaimLocked(userID, claimDay int64) (bool, error) {
	var lock ClaimLockModel
	err := s.db.Where("user_id = ? AND claim_day = ?", userID, claimDay).First(&lock).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// AddClaimLock 为用户在指定日期添加领取锁
func (s *Service) AddClaimLock(userID, claimDay int64) error {
	lock := &ClaimLockModel{
		UserID:    userID,
		ClaimDay:  claimDay,
		CreatedAt: time.Now(),
	}
	// 使用 FirstOrCreate 实现 INSERT OR IGNORE 的效果
	return s.db.Where("user_id = ? AND claim_day = ?", userID, claimDay).FirstOrCreate(lock).Error
}

// RemoveClaimLock 删除用户在指定日期的领取锁（用于失败回滚）
func (s *Service) RemoveClaimLock(userID, claimDay int64) error {
	return s.db.Where("user_id = ? AND claim_day = ?", userID, claimDay).Delete(&ClaimLockModel{}).Error
}
