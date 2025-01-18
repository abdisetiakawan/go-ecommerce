package helper

import (
	"time"

	"github.com/abdisetiakawan/go-ecommerce/internal/model"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

type JwtHelper struct {
    config *viper.Viper
}

type AuthCustomClaims struct {
    ID uint `json:"id"`
    Name  string `json:"name"`
    Username string `json:"username"`
    Role  string `json:"role"`
    Email string `json:"email"`
    jwt.StandardClaims
}

func NewJWTHelper(config *viper.Viper) *JwtHelper {
    return &JwtHelper{
        config: config,
    }
}

func (h *JwtHelper) GenerateTokenUser(user model.AuthResponse) (string, string, error) {
    refreshSecret := h.config.GetString("credentials.refreshsecret")
    accessSecret := h.config.GetString("credentials.accesssecret")

    accessTokenClaims := &AuthCustomClaims{
        ID: user.ID,
        Name:  user.Name,
        Email: user.Email,
        Username: user.Username,
        Role: user.Role,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(time.Minute * 20).Unix(),
            IssuedAt:  time.Now().Unix(),
        },
    }

    refreshTokenClaims := &AuthCustomClaims{
        ID: user.ID,
        Name:  user.Name,
        Email: user.Email,
        Username: user.Username,
        Role: user.Role,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(),
            IssuedAt:  time.Now().Unix(),
        },
    }

    accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
    accessTokenString, err := accessToken.SignedString([]byte(accessSecret))
    if err != nil {
        return "", "", err
    }

    refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
    refreshTokenString, err := refreshToken.SignedString([]byte(refreshSecret))
    if err != nil {
        return "", "", err
    }

    return accessTokenString, refreshTokenString, nil
}
