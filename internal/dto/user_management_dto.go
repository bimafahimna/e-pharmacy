package dto

import (
	"strconv"
	"time"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/model"
)

type UserResponse struct {
	Id         int64     `json:"id"`
	Email      string    `json:"email"`
	Role       string    `json:"role"`
	IsVerified string    `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ListUserParams struct {
	Limit      int    `form:"limit" binding:"omitempty,gte=1"`
	Page       int    `form:"page" binding:"omitempty,gte=1"`
	SortBy     string `form:"sort_by" binding:"omitempty,oneof=created_at"`
	Sort       string `form:"sort" binding:"omitempty,oneof=asc desc"`
	Role       string `form:"role" binding:"omitempty"`
	Email      string `form:"email" binding:"omitempty"`
	IsVerified string `form:"is_verified" binding:"omitempty,boolean"`
}

func ConvertToUserDto(user model.User) UserResponse {
	isVerified := strconv.FormatBool(user.IsVerified)
	return UserResponse{
		Id:         user.ID,
		Role:       user.Role,
		Email:      user.Email,
		IsVerified: isVerified,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}
}

func ConvertToListUsersResponse(usersModel []model.User) []UserResponse {
	users := []UserResponse{}
	for _, userModel := range usersModel {
		user := ConvertToUserDto(userModel)
		users = append(users, user)
	}

	return users
}

func (p *ListUserParams) EnsureDefaults() {
	if p.Limit == 0 {
		p.Limit = 15
	}
	if p.Page == 0 {
		p.Page = 1
	}
	if p.SortBy == "" {
		p.SortBy = "created_at"
	}
	if p.Sort == "" {
		p.Sort = "desc"
	}
}

func (p *ListUserParams) Filters() (filters map[string]interface{}) {
	filters = map[string]interface{}{
		"role":  p.Role,
		"email": p.Email,
	}

	if p.IsVerified != "" {
		isVerified, _ := strconv.ParseBool(p.IsVerified)
		filters["is_verified"] = isVerified
	}
	return
}
