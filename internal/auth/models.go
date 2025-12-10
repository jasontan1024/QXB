package auth

import (
	"time"
)

// UserModel GORM 用户模型
type UserModel struct {
	ID            int64     `gorm:"primaryKey;autoIncrement"`
	Email         string    `gorm:"uniqueIndex;not null;column:email"`
	Address       string    `gorm:"not null;column:address"`
	EncPrivKeyB64 string    `gorm:"not null;column:enc_priv_key"`
	EncSaltB64    string    `gorm:"not null;column:enc_salt"`
	PassSaltB64   string    `gorm:"not null;column:pass_salt"`
	PasswordHash  string    `gorm:"not null;column:password_hash"`
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP;column:created_at"`
}

// TableName 指定表名
func (UserModel) TableName() string {
	return "users"
}

// ClaimLockModel GORM 领取锁模型
type ClaimLockModel struct {
	UserID    int64     `gorm:"primaryKey;column:user_id"`
	ClaimDay  int64     `gorm:"primaryKey;column:claim_day"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP;column:created_at"`
}

// TableName 指定表名
func (ClaimLockModel) TableName() string {
	return "claim_locks"
}
