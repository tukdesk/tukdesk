package models

import (
	"fmt"
)

type TypeUserImportance string

const (
	UserImportanceNon       TypeUserImportance = "0"
	UserImportanceNormal    TypeUserImportance = "1"
	UserImportanceImportant TypeUserImportance = "2"
	UserImportanceVIP       TypeUserImportance = "3"

	UserImportanceNonStr       = "NON"
	UserImportanceNormalStr    = "NORMAL"
	UserImportanceImportantStr = "IMPORTANT"
	UserImportanceVIPStr       = "VIP"
)

func NewTypeUserImportance(s string) TypeUserImportance {
	importance := TypeUserImportance("")
	importance.FromOutput(s)
	return importance
}

func (this TypeUserImportance) ToOutput() string {
	switch this {
	case UserImportanceNormal:
		return UserImportanceNormalStr
	case UserImportanceImportant:
		return UserImportanceImportantStr
	case UserImportanceVIP:
		return UserImportanceVIPStr
	default:
		return UserImportanceNonStr
	}
}

func (this *TypeUserImportance) FromOutput(s string) {
	switch s {
	case UserImportanceNormalStr:
		*this = UserImportanceNormal
	case UserImportanceImportantStr:
		*this = UserImportanceImportant
	case UserImportanceVIPStr:
		*this = UserImportanceVIP
	default:
		*this = UserImportanceNon
	}
}

func (this TypeUserImportance) MarshalJSON() ([]byte, error) {
	return []byte(`"` + this.ToOutput() + `"`), nil
}

func (this *TypeUserImportance) UnmarshalJSON(b []byte) error {
	if len(b) < 2 {
		return fmt.Errorf("invalid TypeUserImportance")
	}
	if b[0] != '"' || b[len(b)-1] != '"' {
		return fmt.Errorf("TypeUserImportance should be a string")
	}
	b = b[1 : len(b)-1]
	this.FromOutput(string(b))
	return nil
}
