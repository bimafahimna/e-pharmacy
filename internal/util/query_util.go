package util

func ToOffset(page, limit int) int {
	return (page - 1) * limit
}
