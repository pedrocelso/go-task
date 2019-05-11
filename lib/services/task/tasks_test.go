package task_test

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"

	"cloud.google.com/go/datastore"

	"github.com/pedrocelso/go-task/lib/http/authcontext"
	"github.com/pedrocelso/go-task/lib/services/task"
	"github.com/stretchr/testify/assert"
)

var mainCtx authcontext.Context

var taskCollection = map[string]map[int64]task.Task{
	`1@gmail.com`: {
		1: task.Task{
			ID:             0,
			Name:           `Old Task`,
			Description:    `Plain Old Task`,
			IncidentsCount: 1,
		},
		4: task.Task{
			ID:          4,
			Name:        `Task 4`,
			Description: `Description 4`,
		},
	},
}

func getMockCollection() map[string]map[int64]task.Task {
	tasks := make(map[int64]task.Task)
	collection := make(map[string]map[int64]task.Task)
	for key, value := range taskCollection[`1@gmail.com`] {
		tasks[key] = value
	}

	collection[`1@gmail.com`] = tasks
	return collection
}

type MockClient struct {
	T          *testing.T
	collection map[string]map[int64]task.Task
}

func (mc MockClient) Count(ctx context.Context, q *datastore.Query) (int, error) {
	return 1, nil
}

func (mc MockClient) Delete(ctx context.Context, key *datastore.Key) error {
	email := key.Parent.Name

	if _, ok := mc.collection[email][key.ID]; ok {
		delete(mc.collection[email], key.ID)
	} else {
		return fmt.Errorf(`datastore: no such entity '%v'`, email)
	}
	return nil
}

func (mc MockClient) Get(ctx context.Context, key *datastore.Key, dst interface{}) (err error) {
	v := reflect.ValueOf(dst).Elem()
	email := key.Parent.Name

	if val, ok := mc.collection[email][key.ID]; ok {
		v.FieldByName("ID").SetInt(val.ID)
		v.FieldByName("Name").SetString(val.Name)
		v.FieldByName("Description").SetString(val.Description)
	} else {
		return fmt.Errorf(`datastore: no such entity '%v'`, email)
	}

	return nil
}

func (mc MockClient) GetAll(ctx context.Context, q *datastore.Query, dst interface{}) (keys []*datastore.Key, err error) {
	v := reflect.ValueOf(dst).Elem()
	var tasks []task.Task

	for _, v := range mc.collection[`1@gmail.com`] {
		tasks = append(tasks, v)
	}

	v.Set(reflect.ValueOf(tasks))

	return nil, nil
}

func (mc MockClient) NewTransaction(ctx context.Context, opts ...datastore.TransactionOption) (t *datastore.Transaction, err error) {
	return nil, nil
}

func (mc MockClient) Put(ctx context.Context, key *datastore.Key, src interface{}) (*datastore.Key, error) {
	assert.Equal(mc.T, `*task.Task`, reflect.TypeOf(src).String())

	email := key.Parent.Name
	v := reflect.ValueOf(src).Elem()

	mc.collection[email][key.ID] = task.Task{
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
func TestCreateTask(t *testing.T) {
	mainCtx.DataStoreClient = MockClient{
		T:          t,
		collection: getMockCollection(),
	}

	output, err := task.Create(&mainCtx, &task.Task{
		Name:        `Test`,
		Description: `Hey, Michael, what you gonna do?`,
	})

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, int64(1), output.ID)
	assert.Equal(t, "Test", output.Name)
	assert.Equal(t, "Hey, Michael, what you gonna do?", output.Description)

	output, err = task.Create(&mainCtx, nil)
	assert.NotNil(t, err)
	assert.Equal(t, `error: invalid Task data (Name is required)`, err.Error())
	assert.Nil(t, output)
}

func TestGetById(t *testing.T) {
	mainCtx.DataStoreClient = MockClient{
		T:          t,
		collection: getMockCollection(),
	}

	output, err := task.GetByID(&mainCtx, int64(4))
	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, "Task 4", output.Name)
	assert.Equal(t, "Description 4", output.Description)

	output, err = task.GetByID(&mainCtx, int64(99))
	assert.NotNil(t, err)
	assert.Equal(t, "Task '99' not found", err.Error())
	assert.Nil(t, output)
}

func TestGetTasks(t *testing.T) {
	mainCtx.DataStoreClient = MockClient{
		T:          t,
		collection: getMockCollection(),
	}

	output, err := task.GetTasks(&mainCtx)
	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, 2, len(output))
	assert.Equal(t, 1, output[0].PendingIncidentsCount)
}

func TestUpdateTask(t *testing.T) {
	mainCtx.DataStoreClient = MockClient{
		T:          t,
		collection: getMockCollection(),
	}

	output, err := task.Update(&mainCtx, &task.Task{
		ID:          int64(4),
		Name:        `Migeh`,
		Description: `Description 1`,
	})
	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, int64(4), output.ID)
	assert.Equal(t, "Migeh", output.Name)
	assert.Equal(t, "Description 1", output.Description)

	tsk, err := task.GetByID(&mainCtx, int64(4))
	assert.Nil(t, err)
	assert.NotNil(t, tsk)
	assert.Equal(t, "Migeh", tsk.Name)
	assert.Equal(t, "Description 1", tsk.Description)

	output, err = task.Update(&mainCtx, nil)
	assert.NotNil(t, err)
	assert.Equal(t, "error: invalid Task data (Name is required)", err.Error())
	assert.Nil(t, output)
}

func TestDeleteTask(t *testing.T) {
	mainCtx.DataStoreClient = MockClient{
		T:          t,
		collection: getMockCollection(),
	}

	tsk, err := task.GetByID(&mainCtx, int64(4))
	assert.Nil(t, err)
	assert.NotNil(t, tsk)
	assert.Equal(t, "Task 4", tsk.Name)
	assert.Equal(t, "Description 4", tsk.Description)

	err = task.Delete(&mainCtx, int64(4))
	assert.Nil(t, err)

	tsk, err = task.GetByID(&mainCtx, int64(4))
	assert.NotNil(t, err)
	assert.Equal(t, "Task '4' not found", err.Error())
	assert.Nil(t, tsk)

	err = task.Delete(&mainCtx, int64(10))
	assert.NotNil(t, err)
	assert.Equal(t, "task '10' don't exist on the database for 1@gmail.com", err.Error())
}
