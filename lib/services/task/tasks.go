package task

import (
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/golang/glog"
	"github.com/pedrocelso/go-task/lib/http/authcontext"
)

const (
	index           = `Task`
	invalidTaskData = `error: invalid Task data (Name is required)`
	userIndex       = `User`
)

// Task defines task attributes
type Task struct {
	ID							int64  `json:"id"`
	Name						string `json:"name" binding:"required"`
	Description			string `json:"description" binding:"required" datastore:",noindex"`
	CreationTime		int64  `json:"creationTime"`
	UpdateTime			int64  `json:"updateTime"`
	Active			 		bool	 `json:"active"`
	IncidentsCount	int64  `json:"incidentsCount"`
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// Create aa task
func Create(ctx *authcontext.Context, task *Task) (*Task, error) {
	var output *Task
	var keys []*datastore.Key
	if task == nil {
		return nil, fmt.Errorf(invalidTaskData)
	}

	userKey := datastore.NameKey(userIndex, ctx.AuthUser.Email, nil)
	incompleteKey := datastore.IncompleteKey(index, userKey)
	keys = append(keys, incompleteKey)

	completeKeys, err := ctx.DataStoreClient.AllocateIDs(ctx.AppEngineCtx, keys)
	if err != nil {
		glog.Errorf("ERROR ON TASK ID GENERATION: %v", err.Error())
		return nil, err
	}

	task.ID = completeKeys[0].ID
	task.CreationTime = makeTimestamp()
	task.Active = true

	insKey, err := ctx.DataStoreClient.Put(ctx.AppEngineCtx, completeKeys[0], task)
	if err != nil {
		glog.Errorf("ERROR INSERTING TASK: %v", err.Error())
		return nil, err
	}

	output, err = GetByID(ctx, insKey.ID)
	if err != nil {
		glog.Errorf("ERROR GETTING TASK OUTPUT: %v", err.Error())
		return nil, err
	}
	return output, nil
}

// GetByID a task based on its numeric ID
func GetByID(ctx *authcontext.Context, id int64) (*Task, error) {
	userKey := datastore.NameKey(userIndex, ctx.AuthUser.Email, nil)
	key := datastore.IDKey(index, id, userKey)
	var task Task
	err := ctx.DataStoreClient.Get(ctx.AppEngineCtx, key, &task)

	if err != nil {
		if strings.HasPrefix(err.Error(), `datastore: no such entity`) {
			err = fmt.Errorf(`Task '%v' not found`, id)
		}
		return nil, err
	}
	return &task, nil
}

// GetTasks Fetches all tasks for the authenticated user
func GetTasks(ctx *authcontext.Context) ([]Task, error) {
	var output []Task
	q := datastore.NewQuery(index)
	completeQuery := q.Ancestor(datastore.NameKey(userIndex, ctx.AuthUser.Email, nil))
	_, err := ctx.DataStoreClient.GetAll(ctx.AppEngineCtx, completeQuery, &output)

	if err != nil {
		glog.Errorf("error fetching all tasks for %s", ctx.AuthUser.Email)
		return nil, err
	}

	if len(output) <= 0 {
		return nil, fmt.Errorf("no tasks found")
	}
	return output, nil
}

// Update task data
func Update(ctx *authcontext.Context, tsk *Task) (*Task, error) {
	if tsk == nil || (tsk.Name == `` && tsk.Description == ``) {
		return nil, fmt.Errorf(invalidTaskData)
	}

	output, _ := GetByID(ctx, tsk.ID)
	if output != nil {
		userKey := datastore.NameKey(userIndex, ctx.AuthUser.Email, nil)
		key := datastore.IDKey(index, tsk.ID, userKey)
		tsk.UpdateTime = makeTimestamp()
		_, err := ctx.DataStoreClient.Put(ctx.AppEngineCtx, key, tsk)

		if err != nil {
			glog.Errorf("ERROR UPDATING TASK: %v", err.Error())
			return nil, err
		}

		output, err = GetByID(ctx, tsk.ID)
		if err != nil {
			glog.Errorf("ERROR GETTING TASK OUTPUT: %v", err.Error())
			return nil, err
		}
		return output, nil
	}
	return nil, fmt.Errorf(`task '%v' not found`, tsk.ID)
}

// Delete a task based on its id.
func Delete(ctx *authcontext.Context, taskID int64) error {
	var output *Task
	output, _ = GetByID(ctx, taskID)

	if output != nil {
		glog.Infof("Deleting task: %v for %v", taskID, ctx.AuthUser.Email)
		userKey := datastore.NameKey(userIndex, ctx.AuthUser.Email, nil)
		key := datastore.IDKey(index, taskID, userKey)
		err := ctx.DataStoreClient.Delete(ctx.AppEngineCtx, key)

		if err != nil {
			glog.Errorf("ERROR DELETING TASK: %v", err.Error())
			return err
		}
		return nil
	}
	return fmt.Errorf("task '%v' don't exist on the database for %v", taskID, ctx.AuthUser.Email)
}
