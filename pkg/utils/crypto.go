package utils

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"strconv"
	"time"
)

// HashType 哈希类型
type HashType string

const (
	HashMD5    HashType = "md5"
	HashSHA1   HashType = "sha1"
	HashSHA256 HashType = "sha256"
	HashSHA512 HashType = "sha512"
)

// Hash 计算哈希值
func Hash(data string, hashType HashType) string {
	var h hash.Hash

	switch hashType {
	case HashMD5:
		h = md5.New()
	case HashSHA1:
		h = sha1.New()
	case HashSHA256:
		h = sha256.New()
	case HashSHA512:
		h = sha512.New()
	default:
		h = sha256.New()
	}

	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// MD5 计算MD5哈希
func MD5(data string) string {
	return Hash(data, HashMD5)
}

// SHA1 计算SHA1哈希
func SHA1(data string) string {
	return Hash(data, HashSHA1)
}

// SHA256 计算SHA256哈希
func SHA256(data string) string {
	return Hash(data, HashSHA256)
}

// SHA512 计算SHA512哈希
func SHA512(data string) string {
	return Hash(data, HashSHA512)
}

// GenerateRandomString 生成随机字符串
func GenerateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	for i, b := range bytes {
		bytes[i] = charset[b%byte(len(charset))]
	}

	return string(bytes), nil
}

// GenerateRandomBytes 生成随机字节
func GenerateRandomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}
	return bytes, nil
}

// GenerateUUID 生成UUID（简单版本）
func GenerateUUID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return ""
	}

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// GenerateToken 生成令牌
func GenerateToken(length int) (string, error) {
	return GenerateRandomString(length)
}

// GenerateAPIKey 生成API密钥
func GenerateAPIKey() (string, error) {
	return GenerateRandomString(32)
}

// GenerateSecret 生成密钥
func GenerateSecret() (string, error) {
	return GenerateRandomString(64)
}

// HashPassword 哈希密码
func HashPassword(password string) string {
	// 添加盐值
	salt := time.Now().Format("20060102150405")
	return SHA256(password + salt)
}

// VerifyPassword 验证密码
func VerifyPassword(password, hashedPassword string) bool {
	// 这里需要根据实际的哈希算法来实现
	// 简化版本，实际应该使用bcrypt等专业库
	return HashPassword(password) == hashedPassword
}

// GenerateOTP 生成OTP验证码
func GenerateOTP(length int) string {
	// 生成数字OTP
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}

	otp := ""
	for _, b := range bytes {
		otp += strconv.Itoa(int(b) % 10)
	}

	return otp
}

// GenerateNumericCode 生成数字验证码
func GenerateNumericCode(length int) string {
	return GenerateOTP(length)
}

// GenerateAlphanumericCode 生成字母数字验证码
func GenerateAlphanumericCode(length int) (string, error) {
	return GenerateRandomString(length)
}

// HashFile 计算文件哈希
func HashFile(reader io.Reader, hashType HashType) (string, error) {
	var h hash.Hash

	switch hashType {
	case HashMD5:
		h = md5.New()
	case HashSHA1:
		h = sha1.New()
	case HashSHA256:
		h = sha256.New()
	case HashSHA512:
		h = sha512.New()
	default:
		h = sha256.New()
	}

	if _, err := io.Copy(h, reader); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// CompareHash 比较哈希值
func CompareHash(data, hash string, hashType HashType) bool {
	return Hash(data, hashType) == hash
}

// GenerateSalt 生成盐值
func GenerateSalt(length int) (string, error) {
	return GenerateRandomString(length)
}

// HashWithSalt 使用盐值哈希
func HashWithSalt(data, salt string, hashType HashType) string {
	return Hash(data+salt, hashType)
}

// VerifyHashWithSalt 验证带盐值的哈希
func VerifyHashWithSalt(data, salt, hash string, hashType HashType) bool {
	return HashWithSalt(data, salt, hashType) == hash
}
