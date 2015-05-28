package models

import (
	"fmt"
)

type TypePriority string

const (
	PriorityNon    TypePriority = ""
	PriorityLow    TypePriority = "1"
	PriorityNormal TypePriority = "2"
	PriorityHign   TypePriority = "3"
	PriorityUrgent TypePriority = "4"

	PriorityNonStr    = ""
	PriorityLowStr    = "LOW"
	PriorityNormalStr = "NORMAL"
	PriorityHignStr   = "HIGH"
	PriorityUrgentStr = "URGENT"
)

func NewTypePriority(s string) TypePriority {
	t := TypePriority("")
	t.FromOutput(s)
	return t
}

func (this TypePriority) ToOutput() string {
	switch this {
	case PriorityLow:
		return PriorityLowStr
	case PriorityNormal:
		return PriorityNormalStr
	case PriorityHign:
		return PriorityHignStr
	case PriorityUrgent:
		return PriorityUrgentStr
	default:
		return PriorityNonStr
	}
}

func (this *TypePriority) FromOutput(s string) {
	switch s {
	case PriorityLowStr:
		*this = PriorityLow
	case PriorityNormalStr:
		*this = PriorityNormal
	case PriorityHignStr:
		*this = PriorityHign
	case PriorityUrgentStr:
		*this = PriorityUrgent
	default:
		*this = PriorityNon
	}
}

func (this TypePriority) MarshalJSON() ([]byte, error) {
	return []byte(`"` + this.ToOutput() + `"`), nil
}

func (this *TypePriority) UnmarshalJSON(b []byte) error {
	if len(b) < 2 {
		return fmt.Errorf("invalid TypePriority")
	}
	if b[0] != '"' || b[len(b)-1] != '"' {
		return fmt.Errorf("TypePriority should be a string")
	}
	b = b[1 : len(b)-1]
	this.FromOutput(string(b))
	return nil
}
