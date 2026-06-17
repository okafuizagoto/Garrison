package controllers

import (
	"members/helpers"
	"members/models"
	"time"

	"github.com/gold-gym/gymkit"
	"github.com/kataras/iris/v12"
	"gorm.io/gorm"
)

type MemberController struct {
	DB *gorm.DB
}

// GetList returns all active members. Admin can see all; regular users see only their own.
func (c *MemberController) GetList(ctx iris.Context) {
	appOrigin := ctx.Values().GetStringDefault("app_origin", "")

	var members []models.Member
	q := c.DB.Where("status = 1")

	if helpers.IsUserAppID(appOrigin) {
		uid := ctx.Values().GetUint64Default("user_id", 0)
		q = q.Where("user_id = ?", uid)
	}

	if err := q.Find(&members).Error; err != nil {
		gymkit.NewResponse(ctx, iris.StatusInternalServerError, "Failed to fetch members")
		return
	}
	gymkit.NewResponse(ctx, iris.StatusOK, members)
}

// GetDetail returns a single member record.
func (c *MemberController) GetDetail(ctx iris.Context) {
	id := ctx.Params().GetUint64Default("id", 0)
	if id == 0 {
		gymkit.NewResponse(ctx, iris.StatusBadRequest, "Invalid member ID")
		return
	}

	var member models.Member
	if err := c.DB.Where("id = ? AND status = 1", id).First(&member).Error; err != nil {
		gymkit.NewResponse(ctx, iris.StatusNotFound, "Member not found")
		return
	}
	gymkit.NewResponse(ctx, iris.StatusOK, member)
}

type CreateMemberInput struct {
	UserID         uint64    `json:"user_id"`
	MembershipType string    `json:"membership_type"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	Notes          string    `json:"notes"`
}

// Create registers a new gym member.
func (c *MemberController) Create(ctx iris.Context) {
	var input CreateMemberInput
	if err := ctx.ReadJSON(&input); err != nil {
		gymkit.NewResponse(ctx, iris.StatusBadRequest, "Invalid request body")
		return
	}

	if input.UserID == 0 {
		gymkit.NewResponse(ctx, iris.StatusBadRequest, "user_id is required")
		return
	}

	membershipType := input.MembershipType
	if membershipType == "" {
		membershipType = "basic"
	}

	createdBy := ctx.Values().GetUint64Default("user_id", 0)

	member := models.Member{
		UserID:         input.UserID,
		MembershipType: membershipType,
		StartDate:      input.StartDate,
		EndDate:        input.EndDate,
		Notes:          input.Notes,
		Status:         1,
		CreatedBy:      createdBy,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := c.DB.Create(&member).Error; err != nil {
		gymkit.NewResponse(ctx, iris.StatusInternalServerError, "Failed to create member")
		return
	}
	gymkit.NewResponse(ctx, iris.StatusCreated, member)
}

type UpdateMemberInput struct {
	MembershipType string     `json:"membership_type"`
	EndDate        *time.Time `json:"end_date"`
	Status         *int8      `json:"status"`
	Notes          string     `json:"notes"`
}

// Update modifies an existing member record.
func (c *MemberController) Update(ctx iris.Context) {
	id := ctx.Params().GetUint64Default("id", 0)
	if id == 0 {
		gymkit.NewResponse(ctx, iris.StatusBadRequest, "Invalid member ID")
		return
	}

	var input UpdateMemberInput
	if err := ctx.ReadJSON(&input); err != nil {
		gymkit.NewResponse(ctx, iris.StatusBadRequest, "Invalid request body")
		return
	}

	var member models.Member
	if err := c.DB.Where("id = ?", id).First(&member).Error; err != nil {
		gymkit.NewResponse(ctx, iris.StatusNotFound, "Member not found")
		return
	}

	updates := map[string]interface{}{"updated_at": time.Now()}
	if input.MembershipType != "" {
		updates["membership_type"] = input.MembershipType
	}
	if input.EndDate != nil {
		updates["end_date"] = *input.EndDate
	}
	if input.Status != nil {
		updates["status"] = *input.Status
	}
	if input.Notes != "" {
		updates["notes"] = input.Notes
	}

	c.DB.Model(&member).Updates(updates)
	gymkit.NewResponse(ctx, iris.StatusOK, member)
}

// Delete soft-deletes a member by setting status = 0.
func (c *MemberController) Delete(ctx iris.Context) {
	id := ctx.Params().GetUint64Default("id", 0)
	if id == 0 {
		gymkit.NewResponse(ctx, iris.StatusBadRequest, "Invalid member ID")
		return
	}

	var member models.Member
	if err := c.DB.Where("id = ? AND status = 1", id).First(&member).Error; err != nil {
		gymkit.NewResponse(ctx, iris.StatusNotFound, "Member not found")
		return
	}

	c.DB.Model(&member).Updates(map[string]interface{}{"status": 0, "updated_at": time.Now()})
	gymkit.NewResponse(ctx, iris.StatusOK, "Member deleted successfully")
}
