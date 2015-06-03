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

// ticket
type TicketAddArgs struct {
	Email    string            `json:"email"`
	Channel  string            `json:"channel`
	Subject  string            `json:"subject"`
	Content  string            `json:"content"`
	Status   string            `json:"status"`
	IsPublic bool              `json:"isPublic,omitempty"`
	Extend   map[string]string `json:"extend"`
}
