package task_test

import (
	"os"
	"testing"

	"github.com/pedrocelso/go-task/lib/http/authcontext"
	"github.com/pedrocelso/go-task/lib/services/task"
	"github.com/stretchr/testify/assert"

	"fmt"

	"cloud.google.com/go/datastore"
)

var mainCtx authcontext.Context

func TestMain(m *testing.M) {
	ctx, done, _ := aetest.NewContext()
	mainCtx.AppEngineCtx = ctx
	mainCtx.AuthUser = authcontext.AuthUser{
		Name:  `Pedro`,
		Email: `1@gmail.com`,
	}
	os.Exit(m.Run())
	done()
}

func TestCreateTask(t *testing.T) {
	output, err := task.Create(mainCtx, &task.Task{
		Name:        `Test`,
		Description: `Hey, Michael, what you gonna do?`,
	})

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, int64(1), output.ID)
	assert.Equal(t, "Test", output.Name)
	assert.Equal(t, "Hey, Michael, what you gonna do?", output.Description)

	output, err = task.Create(mainCtx, nil)
	assert.NotNil(t, err)
	assert.Equal(t, `error: invalid Task data (Name is required)`, err.Error())
	assert.Nil(t, output)
}

func TestGetById(t *testing.T) {
	err := createTasks(mainCtx)
	if err != nil {
		t.Fatal(err)
	}

	output, err := task.GetByID(mainCtx, int64(4))
	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, "Task 4", output.Name)
	assert.Equal(t, "description 4", output.Description)

	output, err = task.GetByID(mainCtx, int64(99))
	assert.NotNil(t, err)
	assert.Equal(t, "Task '99' not found", err.Error())
	assert.Nil(t, output)
}

// // This test run on a different context ot ensure that only
// // the created users will be saved on the datastore
// func TestGetUsers(t *testing.T) {
// 	mainCtx.AuthUser.Email = `2@gmail.com`
// 	err := createTasks(mainCtx)
// 	if err != nil {
// 		mainCtx.AuthUser.Email = `1@gmail.com`
// 		t.Fatal(err)
// 	}
// 	// This sleep is needed because it take some milliseconds for the objects
// 	// created on `createUsers` to be indexed and returned on query
// 	time.Sleep(time.Millisecond * 5e2)
// 	output, err := task.GetTasks(mainCtx)
// 	mainCtx.AuthUser.Email = `1@gmail.com`
// 	assert.Nil(t, err)
// 	assert.NotNil(t, output)
// 	assert.Equal(t, 5, len(output))
// }

func TestUpdateUser(t *testing.T) {
	err := createTasks(mainCtx)

	output, err := task.Update(mainCtx, &task.Task{
		ID:          int64(4),
		Name:        `Migeh`,
		Description: `Description 1`,
	})
	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, int64(4), output.ID)
	assert.Equal(t, "Migeh", output.Name)
	assert.Equal(t, "Description 1", output.Description)

	tsk, err := task.GetByID(mainCtx, int64(4))
	assert.Nil(t, err)
	assert.NotNil(t, tsk)
	assert.Equal(t, "Migeh", tsk.Name)
	assert.Equal(t, "Description 1", tsk.Description)

	output, err = task.Update(mainCtx, nil)
	assert.NotNil(t, err)
	assert.Equal(t, "error: invalid Task data (Name is required)", err.Error())
	assert.Nil(t, output)
}

func TestDeleteTask(t *testing.T) {
	err := createTasks(mainCtx)

	tsk, err := task.GetByID(mainCtx, int64(4))
	assert.Nil(t, err)
	assert.NotNil(t, tsk)
	assert.Equal(t, "Task 4", tsk.Name)
	assert.Equal(t, "description 4", tsk.Description)

	err = task.Delete(mainCtx, int64(4))
	assert.Nil(t, err)

	tsk, err = task.GetByID(mainCtx, int64(4))
	assert.NotNil(t, err)
	assert.Equal(t, "Task '4' not found", err.Error())
	assert.Nil(t, tsk)
}

func createTasks(ctx authcontext.Context) error {
	var userKey *datastore.Key
	userKey = datastore.NewKey(ctx.AppEngineCtx, `User`, ctx.AuthUser.Email, 0, nil)
	for i := 0; i < 5; i++ {
		name := fmt.Sprintf(`Task %v`, i)
		description := fmt.Sprintf(`description %v`, i)

		key := datastore.NewKey(ctx.AppEngineCtx, `Task`, "", int64(i), userKey)
		if _, err := datastore.Put(ctx.AppEngineCtx, key, &task.Task{
			ID:          int64(i),
			Name:        name,
			Description: description,
		}); err != nil {
			return err
		}
	}
	return nil
}
