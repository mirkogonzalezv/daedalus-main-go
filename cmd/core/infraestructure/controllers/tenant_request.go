package controllers

type CreateTenantRequest struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug" biding:"required"`
	Plan string `json:"plan" biding:"required"`
}
