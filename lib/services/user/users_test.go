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

var usersCollection = map[string]user.Full{
	`pedro@pedrocelso.com.br1`: {
		Name:  `Pedro 1`,
		Email: `pedro@pedrocelso.com.br1`,
	},
	`migeh@pedrocelso.com.br`: {
		Name:  `Mr. Migeh`,
		Email: `migeh@pedrocelso.com.br`,
	},
	`pedro@pedrocelso.com.br0`: {
		Name:  `Pedro 0`,
		Email: `pedro@pedrocelso.com.br0`,
	},
	`full@pedrocelso.com.br`: {
		Name:         `Pedro Full`,
		Email:        `full@pedrocelso.com.br`,
		Password:     `$2a$08$2AH4glNU51oZY0fRMyhc7e/HyCG5.n37mqmuYdJnWiKMBcq1aXNtu`,
		CreationTime: 1556801947331,
		UpdateTime:   0,
	},
}

type MockClient struct {
	T          *testing.T
	collection map[string]user.Full
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

		if reflect.TypeOf(dst).String() == `*user.Full` {
			v.FieldByName("Password").SetString(val.Password)
			v.FieldByName("CreationTime").SetInt(val.CreationTime)
			v.FieldByName("UpdateTime").SetInt(val.UpdateTime)
		}

	} else {
		return fmt.Errorf(`datastore: no such entity '%v'`, email)
	}

	return nil
}

func (mc MockClient) GetAll(ctx context.Context, q *datastore.Query, dst interface{}) (keys []*datastore.Key, err error) {
	assert.Equal(mc.T, `*[]user.Full`, reflect.TypeOf(dst).String())
	v := reflect.ValueOf(dst).Elem()
	var users []user.Full

	for _, v := range mc.collection {
		users = append(users, user.Full{
			Name:  v.Name,
			Email: v.Email,
		})
	}

	v.Set(reflect.ValueOf(users))

	return nil, nil
}

func (mc MockClient) NewTransaction(ctx context.Context, opts ...datastore.TransactionOption) (t *datastore.Transaction, err error) {
	return nil, nil
}

func (mc MockClient) Put(ctx context.Context, key *datastore.Key, src interface{}) (*datastore.Key, error) {
	assert.Equal(mc.T, `*user.Full`, reflect.TypeOf(src).String())

	email := key.Name
	v := reflect.ValueOf(src).Elem()

	mc.collection[email] = user.Full{
		Name:         v.FieldByName("Name").String(),
		Email:        v.FieldByName("Email").String(),
		Password:     v.FieldByName("Password").String(),
		CreationTime: v.FieldByName("CreationTime").Int(),
		UpdateTime:   v.FieldByName("UpdateTime").Int(),
	}

	return nil, nil
}

func (mc MockClient) AllocateIDs(ctx context.Context, keys []*datastore.Key) ([]*datastore.Key, error) {
	return nil, nil
}

func TestMain(m *testing.M) {
	mainCtx.AppEngineCtx = context.Background()
	os.Exit(m.Run())
}

func TestCreateUser(t *testing.T) {
	collection := make(map[string]user.Full)
	for key, value := range usersCollection {
		collection[key] = value
	}
	mainCtx.DataStoreClient = MockClient{
		T:          t,
		collection: collection,
	}

	output, err := user.Create(&mainCtx, &user.Full{
		Name:     `Pedro Costa`,
		Email:    `pedro@pedrocelso.com.br`,
		Password: `test`,
	})

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, "Pedro Costa", output.Name)
	assert.Equal(t, "pedro@pedrocelso.com.br", output.Email)

	output, err = user.Create(&mainCtx, nil)
	assert.NotNil(t, err)
	assert.Equal(t, "error: invalid User data", err.Error())
	assert.Nil(t, output)

	output, err = user.Create(&mainCtx, &user.Full{
		Name:  `Pedro Costa`,
		Email: ``,
	})

	assert.NotNil(t, err)
	assert.Equal(t, "error: invalid User data", err.Error())
	assert.Nil(t, output)

	output, err = user.Create(&mainCtx, &user.Full{
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

func TestGetFullByEmail(t *testing.T) {
	mainCtx.DataStoreClient = MockClient{
		T:          t,
		collection: usersCollection,
	}

	output, err := user.GetFullByEmail(&mainCtx, `full@pedrocelso.com.br`)
	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, "Pedro Full", output.Name)
	assert.Equal(t, "full@pedrocelso.com.br", output.Email)
	assert.Equal(t, "$2a$08$2AH4glNU51oZY0fRMyhc7e/HyCG5.n37mqmuYdJnWiKMBcq1aXNtu", output.Password)
	assert.Equal(t, int64(1556801947331), output.CreationTime)
	assert.Equal(t, int64(0), output.UpdateTime)
}

func TestGetUsers(t *testing.T) {
	mainCtx.DataStoreClient = MockClient{
		T:          t,
		collection: usersCollection,
	}
	output, err := user.GetUsers(&mainCtx)
	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, 4, len(output))

	mainCtx.DataStoreClient = MockClient{
		T:          t,
		collection: map[string]user.Full{},
	}
	output, err = user.GetUsers(&mainCtx)
	assert.NotNil(t, err)
	assert.Nil(t, output)
	assert.Equal(t, `no users found`, err.Error())
}

func TestUpdateUser(t *testing.T) {
	collection := make(map[string]user.Full)
	for key, value := range usersCollection {
		collection[key] = value
	}
	mainCtx.DataStoreClient = MockClient{
		T:          t,
		collection: collection,
	}

	output, err := user.Update(&mainCtx, &user.Full{
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

	output, err = user.Update(&mainCtx, &user.Full{
		Name:  `Mr. Jones`,
		Email: `abrandnew@email.com`,
	})
	assert.NotNil(t, err)
	assert.Equal(t, "user 'abrandnew@email.com' not found", err.Error())
	assert.Nil(t, output)
}

func TestDeleteUser(t *testing.T) {
	collection := make(map[string]user.Full)
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

	err = user.Delete(&mainCtx, `abrandnew@email.com`)
	assert.NotNil(t, err)
	assert.Equal(t, "user 'abrandnew@email.com' don't exist on the database", err.Error())
}
