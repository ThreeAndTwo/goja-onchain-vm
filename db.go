package goja_onchain_vm

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type GojaDB struct {
	Client     *mongo.Database
	Collection string
	TeamSlug   string
	Project    Project
	KVMap      map[string]string
}

type Project struct {
	Id   string
	Name string
}

const timeLayout = "2006-01-02 15:04:05"

var (
	errDBClientIsNull     = errors.New("db client is null")
	errDBCollectionIsNull = errors.New("db collection is null")
	errKVMapIsNull        = errors.New("kv map is null")
	errKeyIsNull          = errors.New("key is null")
	errValueIsNull        = errors.New("value is null")
)

func (gdb *GojaDB) check() error {
	if gdb.Client == nil {
		return errDBClientIsNull
	}

	if gdb.Collection == "" {
		return errDBCollectionIsNull
	}

	if gdb.KVMap == nil {
		return errKVMapIsNull
	}
	return nil
}

func (gdb *GojaDB) Get(key string) (string, error) {
	if err := gdb.check(); err != nil {
		return "", err
	}

	if key == "" {
		return "", errKeyIsNull
	}

	if _, ok := gdb.KVMap[key]; ok {
		return gdb.KVMap[key], nil
	}

	filter := bson.D{{"key", key}}
	res := bson.M{}
	err := gdb.Client.Collection(gdb.Collection).FindOne(context.TODO(), filter).Decode(&res)
	if err != nil {
		return "", err
	}

	if _, ok := res["value"]; !ok {
		return "", fmt.Errorf("failed get value by %s key, unknown value field", key)
	}
	return res["value"].(string), nil
}

func (gdb *GojaDB) Set(key, value string) error {
	if err := gdb.check(); err != nil {
		return err
	}

	if key == "" {
		return errKeyIsNull
	}

	if value == "" {
		return errValueIsNull
	}

	filter := bson.M{"team_slug": gdb.TeamSlug, "project_id": gdb.Project.Id, "key": key}
	update := bson.M{
		"$set": bson.M{
			"project_id":   gdb.Project.Id,
			"project_name": gdb.Project.Name,
			"value":        value,
			"team_slug":    gdb.TeamSlug,
			"created_at":   time.Now().Format(timeLayout),
			"updated_at":   time.Now().Format(timeLayout),
		},
	}

	_, err := gdb.Client.Collection(gdb.Collection).UpdateOne(
		context.TODO(),
		filter,
		update,
		options.Update().SetUpsert(true),
	)

	if err != nil {
		return err
	}

	gdb.KVMap[key] = value
	return nil
}
