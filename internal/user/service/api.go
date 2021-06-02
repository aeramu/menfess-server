package service

type (
	CreateReq struct {
		Type string
	}
	GetReq struct {
		ID string
	}
	GetMenfessReq    struct{}
	UpdateProfileReq struct {
		ID     string
		Name   *string
		Avatar *string
		Bio    *string
	}
)
