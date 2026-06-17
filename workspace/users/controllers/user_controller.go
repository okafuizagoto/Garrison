package controllers

import (
	"time"
	"users/models"

	"github.com/gold-gym/gymkit"
	"github.com/kataras/iris/v12"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

// GetList returns all active users.
func (c *UserController) GetList(ctx iris.Context) {
	var users []models.User
	if err := c.DB.Where("status = 1").Find(&users).Error; err != nil {
		gymkit.NewResponse(ctx, iris.StatusInternalServerError, "Failed to fetch users")
		return
	}
	gymkit.NewResponse(ctx, iris.StatusOK, users)
}

// GetDetail returns a single user by ID.
func (c *UserController) GetDetail(ctx iris.Context) {
	id := ctx.Params().GetUint64Default("id", 0)
	if id == 0 {
		gymkit.NewResponse(ctx, iris.StatusBadRequest, "Invalid user ID")
		return
	}

	var user models.User
	if err := c.DB.Where("id = ? AND status = 1", id).First(&user).Error; err != nil {
		gymkit.NewResponse(ctx, iris.StatusNotFound, "User not found")
		return
	}
	gymkit.NewResponse(ctx, iris.StatusOK, user)
}

type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// Register creates a new user account.
func (c *UserController) Register(ctx iris.Context) {
	var input RegisterInput
	if err := ctx.ReadJSON(&input); err != nil {
		gymkit.NewResponse(ctx, iris.StatusBadRequest, "Invalid request body")
		return
	}

	if gymkit.IsEmptyString(input.Name) || gymkit.IsEmptyString(input.Email) || gymkit.IsEmptyString(input.Password) {
		gymkit.NewResponse(ctx, iris.StatusBadRequest, "name, email, and password are required")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		gymkit.NewResponse(ctx, iris.StatusInternalServerError, "Failed to process password")
		return
	}

	role := input.Role
	if role == "" {
		role = "member"
	}

	user := models.User{
		Name:      input.Name,
		Email:     input.Email,
		Password:  string(hashedPassword),
		Role:      role,
		Status:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := c.DB.Create(&user).Error; err != nil {
		gymkit.NewResponse(ctx, iris.StatusConflict, "Email already registered")
		return
	}

	gymkit.NewResponse(ctx, iris.StatusCreated, user)
}

type UpdateInput struct {
	Name   string `json:"name"`
	Status *int8  `json:"status"`
}

// Update modifies an existing user.
func (c *UserController) Update(ctx iris.Context) {
	id := ctx.Params().GetUint64Default("id", 0)
	if id == 0 {
		gymkit.NewResponse(ctx, iris.StatusBadRequest, "Invalid user ID")
		return
	}

	var input UpdateInput
	if err := ctx.ReadJSON(&input); err != nil {
		gymkit.NewResponse(ctx, iris.StatusBadRequest, "Invalid request body")
		return
	}

	var user models.User
	if err := c.DB.Where("id = ?", id).First(&user).Error; err != nil {
		gymkit.NewResponse(ctx, iris.StatusNotFound, "User not found")
		return
	}

	updates := map[string]interface{}{"updated_at": time.Now()}
	if !gymkit.IsEmptyString(input.Name) {
		updates["name"] = input.Name
	}
	if input.Status != nil {
		updates["status"] = *input.Status
	}

	c.DB.Model(&user).Updates(updates)
	gymkit.NewResponse(ctx, iris.StatusOK, user)
}

// Delete soft-deletes a user by setting status = 0.
func (c *UserController) Delete(ctx iris.Context) {
	id := ctx.Params().GetUint64Default("id", 0)
	if id == 0 {
		gymkit.NewResponse(ctx, iris.StatusBadRequest, "Invalid user ID")
		return
	}

	var user models.User
	if err := c.DB.Where("id = ? AND status = 1", id).First(&user).Error; err != nil {
		gymkit.NewResponse(ctx, iris.StatusNotFound, "User not found")
		return
	}

	c.DB.Model(&user).Updates(map[string]interface{}{"status": 0, "updated_at": time.Now()})
	gymkit.NewResponse(ctx, iris.StatusOK, "User deleted successfully")
}
