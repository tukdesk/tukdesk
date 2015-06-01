package models

import (
	"time"

	"github.com/tukdesk/mgoutils"
)

var (
	Now = time.Now

	NewId = mgoutils.NewId
)
