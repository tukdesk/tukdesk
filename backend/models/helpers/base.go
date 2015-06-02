package helpers

import (
	"github.com/tukdesk/mgoutils"
	"gopkg.in/mgo.v2/bson"
)

const (
	LimitedDataFieldMaxLength = 200
)

var (
	EmptyId     = mgoutils.EmptyObjectId
	ErrNotFound = mgoutils.ErrNotFound
)

func IdFromString(s string) (bson.ObjectId, bool) {
	return mgoutils.IdFromString(s)
}

func IsNotFound(err error) bool {
	return mgoutils.IsNotFound(err)
}

func IsDup(err error) bool {
	return mgoutils.IsDup(err)
}

func IsEmptyId(id bson.ObjectId) bool {
	return mgoutils.IsEmptyObjectId(id)
}

type M bson.M
