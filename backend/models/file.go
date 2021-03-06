package models

import (
	"gopkg.in/mgo.v2/bson"
)

const (
	FileCollectionName = "file"
)

type File struct {
	Id         bson.ObjectId `json:"id" bson:"_id"`
	IsInternal bool          `json:"isInternal" bson:"isInternal"`
	MimeType   string        `json:"mimeType" bson:"mimeType"`
	FileSize   int64         `json:"fileSize" bson:"fileSize"`
	FileName   string        `json:"fileName" bson:"fileName"`
	FileKey    string        `json:"fileKey" bson:"fileKey"`
	Created    int64         `json:"created" bson:"created"`
}

func NewFile() *File {
	return &File{
		Id:      NewId(),
		Created: Now().Unix(),
	}
}

func (this *File) Insert() error {
	return Insert(FileCollectionName, this)
}
