package models

import (
	"github.com/tukdesk/mgoutils"
)

var storage *mgoutils.MgoPool

func SetStorage(stg *mgoutils.MgoPool) {
	storage = stg
	return
}

func Insert(collName string, doc ...interface{}) error {
	c := storage.GetCollection(collName)
	defer storage.ReleaseCollection(c)

	return c.Insert(doc...)
}

func FindById(collName string, id, result interface{}) error {
	c := storage.GetCollection(collName)
	defer storage.ReleaseCollection(c)

	return c.FindById(id, result)
}

func FindOne(collName string, query map[string]interface{}, result interface{}) error {
	c := storage.GetCollection(collName)
	defer storage.ReleaseCollection(c)

	return c.FindOne(query, result)
}

func Count(collName string, query map[string]interface{}) (int, error) {
	c := storage.GetCollection(collName)
	defer storage.ReleaseCollection(c)

	return c.Count(query)
}

func List(collName string, query map[string]interface{}, start, limit int, sort []string, result interface{}) error {
	c := storage.GetCollection(collName)
	defer storage.ReleaseCollection(c)

	return c.List(query, start, limit, sort, result)
}

func FindAll(collName string, query map[string]interface{}, sort []string, result interface{}) error {
	c := storage.GetCollection(collName)
	defer storage.ReleaseCollection(c)

	return c.FindAll(query, sort, result)
}

func FindAndModify(collName string, query, change map[string]interface{}, result interface{}) error {
	c := storage.GetCollection(collName)
	defer storage.ReleaseCollection(c)

	return c.FindAndModify(query, change, result)
}

func FindOrInsert(collName string, query map[string]interface{}, doc, result interface{}) (bool, error) {
	c := storage.GetCollection(collName)
	defer storage.ReleaseCollection(c)

	return c.FindOrInsert(query, doc, result)
}

func UpdateById(collName string, id interface{}, change map[string]interface{}) error {
	c := storage.GetCollection(collName)
	defer storage.ReleaseCollection(c)

	return c.UpdateId(id, change)
}

func UpdateOne(collName string, query, update map[string]interface{}) error {
	c := storage.GetCollection(collName)
	defer storage.ReleaseCollection(c)

	return c.Update(query, update)
}

func UpdateAll(collName string, query, update map[string]interface{}) (int, error) {
	c := storage.GetCollection(collName)
	defer storage.ReleaseCollection(c)

	res, err := c.UpdateAll(query, update)
	if res == nil {
		return 0, err
	}
	return res.Updated, err
}

func RemoveById(collName string, id interface{}) error {
	c := storage.GetCollection(collName)
	defer storage.ReleaseCollection(c)

	return c.RemoveId(id)
}

func RemoveOne(collName string, query map[string]interface{}) error {
	c := storage.GetCollection(collName)
	defer storage.ReleaseCollection(c)

	return c.Remove(query)
}

func RemoveAll(collName string, query map[string]interface{}) (int, error) {
	c := storage.GetCollection(collName)
	defer storage.ReleaseCollection(c)

	res, err := c.RemoveAll(query)
	if res == nil {
		return 0, err
	}
	return res.Removed, err
}
