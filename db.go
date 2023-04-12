package goja_onchain_vm

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type GojaDB struct {
	Client     *mongo.Database
	Collection string
	Project    Project
	KVMap      map[string]string
}

type Project struct {
	Id   string
	Name string
}

type storageList struct {
	ID          primitive.ObjectID `json:"_id"          bson:"_id"`
	ProjectId   string             `json:"project_id"   bson:"project_id"`
	ProjectName string             `json:"project_name" bson:"project_name"`
	Key         string             `json:"key"          bson:"key"`
	Value       string             `json:"value"        bson:"value"`
	CreatedAt   string             `json:"created_at"   bson:"created_at"`
	UpdatedAt   string             `json:"updated_at"   bson:"updated_at"`
	DeletedAt   string             `json:"-"            bson:"deleted_at"`
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

	kvData := storageList{
		ID:          primitive.NewObjectID(),
		ProjectId:   gdb.Project.Id,
		ProjectName: gdb.Project.Name,
		Key:         key,
		Value:       value,
		CreatedAt:   time.Now().Format(timeLayout),
		UpdatedAt:   time.Now().Format(timeLayout),
		DeletedAt:   time.Now().Format(timeLayout),
	}
	_, err := gdb.Client.Collection(gdb.Collection).InsertOne(context.TODO(), kvData)
	if err != nil {
		return err
	}

	gdb.KVMap[key] = value
	return nil
}
