package models

import (
	"gopkg.in/mgo.v2"
)

var (
	MgoIndexes = map[string][]mgo.Index{
		BrandCollectionName: []mgo.Index{
			mgo.Index{
				Key:    []string{"on"},
				Unique: true,
			},
		},
		UserCollectionName: []mgo.Index{
			mgo.Index{
				Key:    []string{"channel.name", "channel.id"},
				Unique: true,
			},
		},
		TicketCollectionName: []mgo.Index{
			mgo.Index{
				Key:    []string{"channel.name", "channel.id"},
				Unique: true,
			},
		},
	}
)

func EnsureIndexes() error {
	for name, indexes := range MgoIndexes {
		c := storage.GetCollection(name)
		defer storage.ReleaseCollection(c)

		for _, index := range indexes {
			if err := c.EnsureIndex(index); err != nil {
				return err
			}
		}
	}

	return nil
}
