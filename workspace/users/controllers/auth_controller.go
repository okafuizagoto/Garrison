package controllers

import (
	"os"
	"time"
	"users/models"

	"github.com/gold-gym/gymkit"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kataras/iris/v12"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login authenticates a user and returns a JWT token.
func (c *AuthController) Login(ctx iris.Context) {
	var input LoginInput
	if err := ctx.ReadJSON(&input); err != nil {
		gymkit.NewResponse(ctx, iris.StatusBadRequest, "Invalid request body")
		return
	}

	var user models.User
	if err := c.DB.Where("email = ? AND status = 1", input.Email).First(&user).Error; err != nil {
		gymkit.NewResponse(ctx, iris.StatusUnauthorized, "Invalid email or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		gymkit.NewResponse(ctx, iris.StatusUnauthorized, "Invalid email or password")
		return
	}

	secret := gymkit.GetEnv("JWT_SECRET", "changeme")
	claims := jwt.MapClaims{
		"sub":   user.Role,
		"uid":   user.ID,
		"email": user.Email,
		"name":  user.Name,
		"roles": user.Role,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		gymkit.NewResponse(ctx, iris.StatusInternalServerError, "Failed to generate token")
		return
	}

	gymkit.NewResponse(ctx, iris.StatusOK, map[string]interface{}{
		"token":      tokenStr,
		"token_type": "Bearer",
		"expires_in": 86400,
		"user": map[string]interface{}{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

// Authenticate validates a Bearer token — used by other services for service-to-service auth.
func (c *AuthController) Authenticate(ctx iris.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" || len(authHeader) < 8 {
		gymkit.NewResponse(ctx, iris.StatusUnauthorized, "Missing authorization header")
		return
	}

	tokenStr := authHeader[7:] // strip "Bearer "
	secret := os.Getenv("JWT_SECRET")
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		gymkit.NewResponse(ctx, iris.StatusUnauthorized, "Invalid or expired token")
		return
	}

	gymkit.NewResponse(ctx, iris.StatusOK, map[string]interface{}{
		"valid":  true,
		"claims": claims,
	})
}
