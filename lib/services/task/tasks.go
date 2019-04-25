package task

// import (
// 	"fmt"
// 	"strings"

// 	"cloud.google.com/go/datastore"
// 	"github.com/pedrocelso/go-task/lib/http/authcontext"
// 	"google.golang.org/appengine/log"
// )

// const (
// 	index           = `Task`
// 	invalidTaskData = `error: invalid Task data (Name is required)`
// 	userIndex       = `User`
// )

// // Task defines task attributes
// type Task struct {
// 	ID          int64  `json:"id"`
// 	Name        string `json:"name" binding:"required"`
// 	Description string `json:"description" binding:"required"`
// }

// // Create aa task
// func Create(c authcontext.Context, task *Task) (*Task, error) {
// 	var output *Task
// 	var keys []*datastore.Key
// 	if task == nil {
// 		return nil, fmt.Errorf(invalidTaskData)
// 	}

// 	userKey := datastore.NameKey(userIndex, c.AuthUser.Email, nil)
// 	incompleteKey := datastore.IncompleteKey(index, userKey)
// 	keys = append(keys, incompleteKey)

// 	completeKeys, err := c.DataStoreClient.AllocateIDs(c.AppEngineCtx, keys)
// 	if err != nil {
// 		log.Errorf(c.AppEngineCtx, "ERROR ON TASK ID GENERATION", err.Error())
// 		return nil, err
// 	}

// 	log.Infof(c.AppEngineCtx, "LOW: %v / HIGH: %v", low, high)
// 	task.ID = low

// 	key := datastore.NewKey(c.AppEngineCtx, index, "", low, userKey)
// 	insKey, err := datastore.Put(c.AppEngineCtx, key, task)
// 	if err != nil {
// 		log.Errorf(c.AppEngineCtx, "ERROR INSERTING TASK: %v", err.Error())
// 		return nil, err
// 	}

// 	output, err = GetByID(c, insKey.IntID())
// 	if err != nil {
// 		log.Errorf(c.AppEngineCtx, "ERROR GETTING TASK OUTPUT: %v", err.Error())
// 		return nil, err
// 	}
// 	return output, nil
// }

// // GetByID a task based on its numeric ID
// func GetByID(c authcontext.Context, id int64) (*Task, error) {
// 	userKey := datastore.NewKey(c.AppEngineCtx, userIndex, c.AuthUser.Email, 0, nil)
// 	key := datastore.NewKey(c.AppEngineCtx, index, "", id, userKey)
// 	var task Task
// 	err := datastore.Get(c.AppEngineCtx, key, &task)

// 	if err != nil {
// 		if strings.HasPrefix(err.Error(), `datastore: no such entity`) {
// 			err = fmt.Errorf(`Task '%v' not found`, id)
// 		}
// 		return nil, err
// 	}
// 	return &task, nil
// }

// // GetTasks Fetches all tasks for the authenticated user
// func GetTasks(c authcontext.Context) ([]Task, error) {
// 	log.Debugf(c.AppEngineCtx, "GETTING ALL Tasks FOR %s", c.AuthUser.Email)
// 	var output []Task
// 	q := datastore.NewQuery(index)
// 	completeQuery := q.Ancestor(datastore.NewKey(c.AppEngineCtx, userIndex, c.AuthUser.Email, 0, nil))
// 	_, err := completeQuery.GetAll(c.AppEngineCtx, &output)

// 	if err != nil {
// 		log.Errorf(c.AppEngineCtx, "error fetching all tasks for %s", c.AuthUser.Email)
// 		return nil, err
// 	}

// 	if len(output) <= 0 {
// 		return nil, fmt.Errorf("no tasks found")
// 	}
// 	return output, nil
// }

// // Update task data
// func Update(c authcontext.Context, tsk *Task) (*Task, error) {
// 	if tsk == nil || (tsk.Name == `` && tsk.Description == ``) {
// 		return nil, fmt.Errorf(invalidTaskData)
// 	}

// 	output, _ := GetByID(c, tsk.ID)
// 	if output != nil {
// 		userKey := datastore.NewKey(c.AppEngineCtx, userIndex, c.AuthUser.Email, 0, nil)
// 		key := datastore.NewKey(c.AppEngineCtx, index, "", tsk.ID, userKey)
// 		_, err := datastore.Put(c.AppEngineCtx, key, tsk)

// 		if err != nil {
// 			log.Errorf(c.AppEngineCtx, "ERROR UPDATING TASK: %v", err.Error())
// 			return nil, err
// 		}

// 		output, err = GetByID(c, tsk.ID)
// 		if err != nil {
// 			log.Errorf(c.AppEngineCtx, "ERROR GETTING TASK OUTPUT: %v", err.Error())
// 			return nil, err
// 		}
// 		return output, nil
// 	}
// 	return nil, fmt.Errorf(`task '%v' not found`, tsk.ID)
// }

// // Delete a task based on its id.
// func Delete(c authcontext.Context, taskID int64) error {
// 	var output *Task
// 	output, _ = GetByID(c, taskID)

// 	if output != nil {
// 		log.Infof(c.AppEngineCtx, "Deleting task: %v for %v", taskID, c.AuthUser.Email)
// 		userKey := datastore.NewKey(c.AppEngineCtx, userIndex, c.AuthUser.Email, 0, nil)
// 		key := datastore.NewKey(c.AppEngineCtx, index, "", taskID, userKey)
// 		err := datastore.Delete(c.AppEngineCtx, key)

// 		if err != nil {
// 			log.Errorf(c.AppEngineCtx, "ERROR DELETING TASK: %v", err.Error())
// 			return err
// 		}
// 		return nil
// 	}
// 	return fmt.Errorf("task '%v' don't exist on the database for %v", taskID, c.AuthUser.Email)
// }
