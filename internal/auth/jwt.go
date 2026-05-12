package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecretBytes []byte
	jwtSecretOnce  sync.Once
)

func jwtSecret() []byte {
	jwtSecretOnce.Do(func() {
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			b := make([]byte, 32)
			if _, err := rand.Read(b); err != nil {
				log.Printf("FATAL: cannot generate JWT secret: %v", err)
				jwtSecretBytes = []byte("e8c9e3a7f2a6d47bf8c9a02b3e4d5f67890abcdef1234567890abcdef")
				return
			}
			secret = hex.EncodeToString(b)
			log.Printf("WARNING: JWT_SECRET not set, generated random key")
		}
		jwtSecretBytes = []byte(secret)
	})
	return jwtSecretBytes
}

type Claims struct {
	AccountID uint   `json:"account_id"`
	Username  string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateToken(accountID uint, username string) (string, error) {
	now := time.Now()

	claims := Claims{
		AccountID: accountID,
		Username:  username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret())
}

func GenerateRefreshToken(accountID uint) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			if token.Method == nil || token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, errors.New("unexpected signing method")
			}
			return jwtSecret(), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims) //转换成我的Claims类型
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "authorization header required"})
			return
		}

		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid authorization format"})
			return
		}

		tokenString := authHeader[7:]
		claims, err := ParseToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token: " + err.Error()})
			return
		}

		c.Set("accountID", claims.AccountID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
