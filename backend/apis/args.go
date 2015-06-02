package apis

import (
	"github.com/tukdesk/tukdesk/backend/models"
)

// base
type SignupArgs struct {
	Password string `json:"password"`
}

type BrandInitArgs struct {
	BrandName string `json:"brandName"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// brand
type BrandUpdateArgs struct {
	Base models.BrandBaseInfo `json:"base"`
}

// profile
type ProfileUpdateArgs struct {
	Base models.UserBaseInfo `json:"base"`
}
