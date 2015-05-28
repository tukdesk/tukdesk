package models

import (
	"gopkg.in/mgo.v2/bson"
)

const (
	AttachmentCollectionName = "attachment"
)

type Attachment struct {
	Id        bson.ObjectId `json:"id" bson:"_id"`
	CreatorId bson.ObjectId `json:"creatorId" bson:"creatorId"`
	MimeType  string        `json:"mimeType" bson:"mimeType"`
	FileSize  int64         `json:"fileSize" bson:"fileSize"`
	FileKey   string        `json:"fileKey" bson:"fileKey"`
	Created   int64         `json:"created" bson:"created"`
}
