package incident_test

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"

	"cloud.google.com/go/datastore"

	"github.com/pedrocelso/go-task/lib/http/authcontext"
	"github.com/pedrocelso/go-task/lib/services/incident"
	"github.com/stretchr/testify/assert"
)

var mainCtx authcontext.Context

var incidentCollection = map[int64]map[int64]incident.Incident{
	33: {
		1: incident.Incident{
			ID:          0,
			Name:        `Old Incident`,
			Description: `Plain Old Incident`,
		},
		4: incident.Incident{
			ID:          4,
			Name:        `Incident 4`,
			Description: `Description 4`,
		},
	},
}

func getMockCollection() map[int64]map[int64]incident.Incident {
	incidents := make(map[int64]incident.Incident)
	collection := make(map[int64]map[int64]incident.Incident)

	for key, value := range incidentCollection[33] {
		incidents[key] = value
	}

	collection[33] = incidents

	return collection
}

// type MockTransaction struct {
// 	T          *testing.T
// }

// func (mt MockTransaction) Commit() (c *datastore.Commit, err error) {
// 	return c, nil
// }

// func (mt MockTransaction) Delete(key *datastore.Key) error {
// 	return nil
// }

// func (mt MockTransaction) Get(key *datastore.Key, dst interface{}) (err error) {
// 	return nil
// }

// func (mt MockTransaction) Put(key *datastore.Key, src interface{}) (*datastore.PendingKey, error) {
// 	return nil, nil
// }

type MockClient struct {
	T          *testing.T
	collection map[int64]map[int64]incident.Incident
}

func (mc MockClient) Delete(ctx context.Context, key *datastore.Key) error {
	if _, ok := mc.collection[key.ID]; ok {
		delete(mc.collection, key.ID)
	} else {
		return fmt.Errorf(`datastore: no such entity '%v'`, key.ID)
	}
	return nil
}

func (mc MockClient) Get(ctx context.Context, key *datastore.Key, dst interface{}) (err error) {
	id := key.ID
	taskID := key.Parent.ID
	v := reflect.ValueOf(dst).Elem()

	if val, ok := mc.collection[taskID][key.ID]; ok {
		v.FieldByName("ID").SetInt(val.ID)
		v.FieldByName("Name").SetString(val.Name)
		v.FieldByName("Description").SetString(val.Description)
	} else {
		return fmt.Errorf(`datastore: no such entity '%v'`, id)
	}

	return nil
}

func (mc MockClient) GetAll(ctx context.Context, q *datastore.Query, dst interface{}) (keys []*datastore.Key, err error) {
	v := reflect.ValueOf(dst).Elem()
	var incidents []incident.Incident

	for _, v := range mc.collection[33] {
		incidents = append(incidents, v)
	}

	v.Set(reflect.ValueOf(incidents))

	return nil, nil
}

func (mc MockClient) NewTransaction(ctx context.Context, opts ...datastore.TransactionOption) (t *datastore.Transaction, err error) {
	return nil, nil
}

func (mc MockClient) Put(ctx context.Context, key *datastore.Key, src interface{}) (*datastore.Key, error) {
	assert.Equal(mc.T, `*incident.Incident`, reflect.TypeOf(src).String())

	v := reflect.ValueOf(src).Elem()

	mc.collection[33][key.ID] = incident.Incident{
		ID:          key.ID,
		Name:        v.FieldByName("Name").String(),
		Description: v.FieldByName("Description").String(),
	}

	return key, nil
}

func (mc MockClient) AllocateIDs(ctx context.Context, keys []*datastore.Key) ([]*datastore.Key, error) {
	for key := range keys {
		keys[key].ID = int64(key + 1)
	}

	return keys, nil
}

func TestMain(m *testing.M) {
	mainCtx.AppEngineCtx = context.Background()
	mainCtx.AuthUser = authcontext.AuthUser{
		Name:  `Pedro`,
		Email: `1@gmail.com`,
	}
	os.Exit(m.Run())
}

// func TestCreateIncident(t *testing.T) {
// 	mainCtx.DataStoreClient = MockClient{
// 		T:          t,
// 		collection: getMockCollection(),
// 	}

// 	output, err := incident.Create(&mainCtx, int64(1), &incident.Incident{
// 		Name:        `Test`,
// 		Description: `Hey, Michael, what you gonna do?`,
// 	})

// 	assert.Nil(t, err)
// 	assert.NotNil(t, output)
// 	assert.Equal(t, int64(1), output.ID)
// 	assert.Equal(t, "Test", output.Name)
// 	assert.Equal(t, "Hey, Michael, what you gonna do?", output.Description)

// 	output, err = incident.Create(&mainCtx, int64(1), nil)
// 	assert.NotNil(t, err)
// 	assert.Equal(t, `error: invalid Incident data`, err.Error())
// 	assert.Nil(t, output)
// }

func TestGetById(t *testing.T) {
	mainCtx.DataStoreClient = MockClient{
		T:          t,
		collection: getMockCollection(),
	}

	output, err := incident.GetByID(&mainCtx, 33, int64(4))
	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, "Incident 4", output.Name)
	assert.Equal(t, "Description 4", output.Description)

	output, err = incident.GetByID(&mainCtx, 33, int64(99))
	assert.NotNil(t, err)
	assert.Equal(t, "Incident '99' not found", err.Error())
	assert.Nil(t, output)
}

func TestGetIncidents(t *testing.T) {
	mainCtx.DataStoreClient = MockClient{
		T:          t,
		collection: getMockCollection(),
	}

	output, err := incident.GetIncidents(&mainCtx, 1)
	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, 2, len(output))
}

func TestUpdateIncident(t *testing.T) {
	mainCtx.DataStoreClient = MockClient{
		T:          t,
		collection: getMockCollection(),
	}

	output, err := incident.Update(&mainCtx, 33, &incident.Incident{
		ID:          int64(4),
		Name:        `Migeh`,
		Description: `Description 1`,
	})
	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, int64(4), output.ID)
	assert.Equal(t, "Migeh", output.Name)
	assert.Equal(t, "Description 1", output.Description)

	tsk, err := incident.GetByID(&mainCtx, 33, int64(4))
	assert.Nil(t, err)
	assert.NotNil(t, tsk)
	assert.Equal(t, "Migeh", tsk.Name)
	assert.Equal(t, "Description 1", tsk.Description)

	output, err = incident.Update(&mainCtx, 33, nil)
	assert.NotNil(t, err)
	assert.Equal(t, "error: invalid Incident data", err.Error())
	assert.Nil(t, output)
}

// func TestDeleteIncident(t *testing.T) {
// 	mainCtx.DataStoreClient = MockClient{
// 		T:          t,
// 		collection: getMockCollection(),
// 	}

// 	tsk, err := incident.GetByID(&mainCtx, int64(4))
// 	assert.Nil(t, err)
// 	assert.NotNil(t, tsk)
// 	assert.Equal(t, "Incident 4", tsk.Name)
// 	assert.Equal(t, "Description 4", tsk.Description)

// 	err = incident.Delete(&mainCtx, int64(4))
// 	assert.Nil(t, err)

// 	tsk, err = incident.GetByID(&mainCtx, int64(4))
// 	assert.NotNil(t, err)
// 	assert.Equal(t, "Incident '4' not found", err.Error())
// 	assert.Nil(t, tsk)

// 	err = incident.Delete(&mainCtx, int64(10))
// 	assert.NotNil(t, err)
// 	assert.Equal(t, "incident '10' don't exist on the database for 1@gmail.com", err.Error())
// }
