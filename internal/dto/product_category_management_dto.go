package dto

type ProductCategory struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type AddProductCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateProductCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}
