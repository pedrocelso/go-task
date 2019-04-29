package user_test

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/pedrocelso/go-task/lib/http/authcontext"
	"github.com/pedrocelso/go-task/lib/services/user"
	"github.com/stretchr/testify/assert"

	"fmt"

	"cloud.google.com/go/datastore"
)

const email = `pedro@pedrocelso.com.br`

var mainCtx authcontext.Context
var c context.Context

var usersCollection = map[string]user.User{
	`pedro@pedrocelso.com.br1`: user.User{
		Name:  `Pedro 1`,
		Email: `pedro@pedrocelso.com.br1`,
	},
	`migeh@pedrocelso.com.br`: user.User{
		Name:  `Mr. Migeh`,
		Email: `migeh@pedrocelso.com.br`,
	},
	`pedro@pedrocelso.com.br0`: user.User{
		Name:  `Pedro 0`,
		Email: `pedro@pedrocelso.com.br0`,
	},
}

type MockClient struct {
	T          *testing.T
	collection map[string]user.User
}

func (mc MockClient) Delete(ctx context.Context, key *datastore.Key) error {
	email := key.Name

	if _, ok := mc.collection[email]; ok {
		delete(mc.collection, email)
	} else {
		return fmt.Errorf(`datastore: no such entity '%v'`, email)
	}
	return nil
}

func (mc MockClient) Get(ctx context.Context, key *datastore.Key, dst interface{}) (err error) {
	v := reflect.ValueOf(dst).Elem()
	email := key.Name

	if val, ok := mc.collection[email]; ok {
		v.FieldByName("Name").SetString(val.Name)
		v.FieldByName("Email").SetString(val.Email)
	} else {
		return fmt.Errorf(`datastore: no such entity '%v'`, email)
	}

	return nil
}

func (mc MockClient) GetAll(ctx context.Context, q *datastore.Query, dst interface{}) (keys []*datastore.Key, err error) {
	v := reflect.ValueOf(dst).Elem()
	var users []user.User

	for _, v := range mc.collection {
		users = append(users, v)
	}

	v.Set(reflect.ValueOf(users))

	return nil, nil
}

func (mc MockClient) Put(ctx context.Context, key *datastore.Key, src interface{}) (*datastore.Key, error) {
	assert.Equal(mc.T, `*user.User`, reflect.TypeOf(src).String())

	email := key.Name
	v := reflect.ValueOf(src).Elem()

	mc.collection[email] = user.User{
		Name:  v.FieldByName("Name").String(),
		Email: v.FieldByName("Email").String(),
	}

	return nil, nil
}

func TestMain(m *testing.M) {
	c := context.Background()
	mainCtx.AppEngineCtx = c
	os.Exit(m.Run())
}

func TestCreateUser(t *testing.T) {
	collection := make(map[string]user.User)
	for key, value := range usersCollection {
		collection[key] = value
	}
	mainCtx.DataStoreClient = MockClient{
		T:          t,
		collection: collection,
	}

	output, err := user.Create(&mainCtx, &user.User{
		Name:  `Pedro Costa`,
		Email: `pedro@pedrocelso.com.br`,
	})

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, "Pedro Costa", output.Name)
	assert.Equal(t, "pedro@pedrocelso.com.br", output.Email)

	output, err = user.Create(&mainCtx, nil)
	assert.NotNil(t, err)
	assert.Equal(t, "error: invalid User data", err.Error())
	assert.Nil(t, output)

	output, err = user.Create(&mainCtx, &user.User{
		Name:  `Pedro Costa`,
		Email: ``,
	})

	assert.NotNil(t, err)
	assert.Equal(t, "error: invalid User data", err.Error())
	assert.Nil(t, output)

	output, err = user.Create(&mainCtx, &user.User{
		Name:  `Pedro 1`,
		Email: `pedro@pedrocelso.com.br1`,
	})

	assert.NotNil(t, err)
	assert.Equal(t, "User 'pedro@pedrocelso.com.br1' already exists", err.Error())
	assert.NotNil(t, output)
	assert.Equal(t, "Pedro 1", output.Name)
	assert.Equal(t, "pedro@pedrocelso.com.br1", output.Email)
}

func TestGetByEmail(t *testing.T) {
	mainCtx.DataStoreClient = MockClient{
		T:          t,
		collection: usersCollection,
	}

	output, err := user.GetByEmail(&mainCtx, `pedro@pedrocelso.com.br1`)
	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, "Pedro 1", output.Name)
	assert.Equal(t, "pedro@pedrocelso.com.br1", output.Email)

	output, err = user.GetByEmail(&mainCtx, `bad_email@gmail.com`)
	assert.NotNil(t, err)
	assert.Equal(t, "user 'bad_email@gmail.com' not found", err.Error())
	assert.Nil(t, output)

	output, err = user.GetByEmail(&mainCtx, ``)
	assert.NotNil(t, err)
	assert.Equal(t, `error: invalid User data`, err.Error())
	assert.Nil(t, output)
}

func TestGetUsers(t *testing.T) {
	mainCtx.DataStoreClient = MockClient{
		T:          t,
		collection: usersCollection,
	}
	output, err := user.GetUsers(&mainCtx)
	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, 3, len(output))
}

func TestUpdateUser(t *testing.T) {
	collection := make(map[string]user.User)
	for key, value := range usersCollection {
		collection[key] = value
	}
	mainCtx.DataStoreClient = MockClient{
		T:          t,
		collection: collection,
	}

	output, err := user.Update(&mainCtx, &user.User{
		Name:  `Migeh`,
		Email: `migeh@pedrocelso.com.br`,
	})
	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, "Migeh", output.Name)
	assert.Equal(t, "migeh@pedrocelso.com.br", output.Email)

	output, err = user.Update(&mainCtx, nil)
	assert.NotNil(t, err)
	assert.Equal(t, "error: invalid User data", err.Error())
	assert.Nil(t, output)
}

func TestDeleteUser(t *testing.T) {
	collection := make(map[string]user.User)
	for key, value := range usersCollection {
		collection[key] = value
	}
	mainCtx.DataStoreClient = MockClient{
		T:          t,
		collection: collection,
	}

	usr, err := user.GetByEmail(&mainCtx, `pedro@pedrocelso.com.br0`)
	assert.Nil(t, err)
	assert.NotNil(t, usr)
	assert.Equal(t, "Pedro 0", usr.Name)
	assert.Equal(t, "pedro@pedrocelso.com.br0", usr.Email)

	err = user.Delete(&mainCtx, `pedro@pedrocelso.com.br0`)
	assert.Nil(t, err)

	usr, err = user.GetByEmail(&mainCtx, `pedro@pedrocelso.com.br0`)
	assert.NotNil(t, err)
	assert.Equal(t, "user 'pedro@pedrocelso.com.br0' not found", err.Error())
	assert.Nil(t, usr)
}
