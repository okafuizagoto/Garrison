// Package auth provides JWT authentication middleware for Iris.
package auth

import (
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"

	"github.com/gold-gym/gymkit"
)

type Config struct {
	// SecretKey is the JWT HMAC signing secret.
	SecretKey string
	// AllowedRoles limits which roles may access the route (empty = all roles allowed).
	AllowedRoles []string
}

type authMiddleware struct {
	config Config
}

// New returns a new Iris handler that validates Bearer JWT tokens.
func New(c Config) context.Handler {
	if c.SecretKey == "" {
		c.SecretKey = gymkit.GetEnv("JWT_SECRET", "changeme")
	}
	am := &authMiddleware{config: c}
	return am.Serve
}

func (am *authMiddleware) Serve(ctx *context.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		gymkit.NewResponse(ctx, iris.StatusUnauthorized, "Authorization header is required")
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		gymkit.NewResponse(ctx, iris.StatusUnauthorized, "Authorization header must be Bearer <token>")
		return
	}

	tokenStr := parts[1]
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, iris.ErrNotFound
		}
		return []byte(am.config.SecretKey), nil
	})

	if err != nil || !token.Valid {
		gymkit.NewResponse(ctx, iris.StatusUnauthorized, "Invalid or expired token")
		return
	}

	// Inject user info into request values for downstream handlers
	if sub, ok := claims["sub"].(string); ok {
		ctx.Values().Set("app_origin", sub)
		ctx.Request().Header.Set("X-App-Origin", sub)
	}
	if uid, ok := claims["uid"]; ok {
		ctx.Values().Set("user_id", uid)
	}
	if roles, ok := claims["roles"].(string); ok {
		ctx.Values().Set("user_roles", roles)
		if len(am.config.AllowedRoles) > 0 && !checkRoles(roles, am.config.AllowedRoles) {
			gymkit.NewResponse(ctx, iris.StatusForbidden, "Insufficient permissions")
			return
		}
	}

	ctx.Next()
}

func checkRoles(userRoles string, allowedRoles []string) bool {
	roles := strings.Split(userRoles, ",")
	for _, role := range roles {
		for _, allowed := range allowedRoles {
			if strings.TrimSpace(role) == strings.TrimSpace(allowed) {
				return true
			}
		}
	}
	return false
}
