package model

type City struct {
	ID                   int64
	ProvinceId           int64
	Name                 string
	Type                 string
	UnofficialId         *int64
	ProvinceUnofficialId int64
}
