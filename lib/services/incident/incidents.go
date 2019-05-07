package incident

import (
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/golang/glog"
	"github.com/pedrocelso/go-task/lib/http/authcontext"
	"github.com/pedrocelso/go-task/lib/services/task"
)

const (
	index           		= `Incident`
	invalidIncidentData = `error: invalid Incident data`
	taskIndex       		= `Task`
	userIndex       		= `User`
)

// Incident defines incident attributes
type Incident struct {
	ID           int64  `json:"id"`
	TaskID			 int64	`json:"taskId"`
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description" binding:"required" datastore:",noindex"`
	Active			 bool 	`json:"active"`
	CreationTime int64  `json:"creationTime"`
	UpdateTime   int64  `json:"updateTime"`
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// Create an incident
func Create(ctx *authcontext.Context, taskID int64, incident *Incident) (*Incident, error) {
	var output *Incident
	var keys []*datastore.Key
	var insKey *datastore.PendingKey
	var tx *datastore.Transaction
	var commit *datastore.Commit
	var parentTask task.Task
	var err error

	if incident == nil {
		return nil, fmt.Errorf(invalidIncidentData)
	}

	userKey := datastore.NameKey(userIndex, ctx.AuthUser.Email, nil)
	taskKey := datastore.IDKey(taskIndex, taskID, userKey)
	incompleteKey := datastore.IncompleteKey(index, taskKey)
	keys = append(keys, incompleteKey)

	completeKeys, err := ctx.DataStoreClient.AllocateIDs(ctx.AppEngineCtx, keys)
	if err != nil {
		return nil, fmt.Errorf("ERROR ON INCIDENT ID GENERATION: %v", err.Error())
	}

	incident.ID = completeKeys[0].ID
	incident.TaskID = taskID
	incident.CreationTime = makeTimestamp()
	incident.Active = true

	tx, err = ctx.DataStoreClient.NewTransaction(ctx.AppEngineCtx)

	if err != nil {
		return nil, err
	}

	if err = tx.Get(taskKey, &parentTask); err != nil {
		return nil, fmt.Errorf("ERROR GETTING INCIDENT PARENT TASK: %v", err.Error())
	}

	parentTask.IncidentsCount = parentTask.IncidentsCount + 1

	if _, err = tx.Put(taskKey, &parentTask); err != nil {
		return nil, fmt.Errorf("ERROR INCREASING PARENT TASK INCIDENT COUNT: %v", err.Error())
	}

	insKey, err = tx.Put(completeKeys[0], incident)
	if err != nil {
		return nil, fmt.Errorf("ERROR INSERTING INCIDENT: %v", err.Error())
	}

	if commit, err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("ERROR INSERTING INCIDENT: %v", err.Error())
	}

	output, err = GetByID(ctx, taskID, commit.Key(insKey).ID)
	if err != nil {
		return nil, fmt.Errorf("ERROR GETTING INCIDENT OUTPUT: %v", err.Error())
	}
	return output, nil
}

// GetByID a incident based on its numeric ID
func GetByID(ctx *authcontext.Context, taskID int64, id int64) (*Incident, error) {
	userKey := datastore.NameKey(userIndex, ctx.AuthUser.Email, nil)
	taskKey := datastore.IDKey(taskIndex, taskID, userKey)
	key := datastore.IDKey(index, id, taskKey)
	var incident Incident
	err := ctx.DataStoreClient.Get(ctx.AppEngineCtx, key, &incident)

	if err != nil {
		if strings.HasPrefix(err.Error(), `datastore: no such entity`) {
			err = fmt.Errorf(`Incident '%v' not found`, id)
		}
		return nil, err
	}
	return &incident, nil
}

// GetIncidents Fetches all incidents for the authenticated user
func GetIncidents(ctx *authcontext.Context, taskID int64) ([]Incident, error) {
	var output []Incident
	taskKey := datastore.NameKey(userIndex, ctx.AuthUser.Email, nil)
	q := datastore.NewQuery(index)
	completeQuery := q.Ancestor(datastore.IDKey(taskIndex, taskID, taskKey))
	_, err := ctx.DataStoreClient.GetAll(ctx.AppEngineCtx, completeQuery, &output)

	if err != nil {
		return nil, fmt.Errorf("error fetching all incidents for %s: %v", ctx.AuthUser.Email, err.Error())
	}

	if len(output) <= 0 {
		return nil, fmt.Errorf("no incidents found")
	}
	return output, nil
}

// Update incident data
func Update(ctx *authcontext.Context, taskID int64, incident *Incident) (*Incident, error) {
	if incident == nil || (incident.Name == `` && incident.Description == ``) {
		return nil, fmt.Errorf(invalidIncidentData)
	}

	output, _ := GetByID(ctx, taskID, incident.ID)
	if output != nil {
		key := datastore.IDKey(index, incident.ID, nil)
		incident.UpdateTime = makeTimestamp()
		_, err := ctx.DataStoreClient.Put(ctx.AppEngineCtx, key, incident)

		if err != nil {
			return nil, fmt.Errorf("ERROR UPDATING INCIDENT: %v", err.Error())
		}

		output, err = GetByID(ctx, taskID, incident.ID)
		if err != nil {
			return nil, fmt.Errorf("ERROR GETTING INCIDENT OUTPUT: %v", err.Error())
		}
		return output, nil
	}
	return nil, fmt.Errorf(`incident '%v' not found`, incident.ID)
}

// Delete a incident based on its id.
func Delete(ctx *authcontext.Context, taskID int64, incidentID int64) error {
	var output *Incident
	var tx *datastore.Transaction
	var parentTask *task.Task
	var err error

	output, err = GetByID(ctx, taskID, incidentID)

	if err != nil {
		return err
	}

	if output != nil {
		glog.Infof("Deleting incident: %v for %v", incidentID, ctx.AuthUser.Email)

		tx, err = ctx.DataStoreClient.NewTransaction(ctx.AppEngineCtx)
		if err != nil {
			return err
		}

		if err = tx.Get(datastore.IDKey(index, incidentID, nil), output); err != nil {
			return fmt.Errorf("ERROR GETTING INCIDENT : %v", err.Error())
		}

		taskKey := datastore.IDKey(taskIndex, output.TaskID, nil)
		if err = tx.Get(taskKey, parentTask); err != nil {
			return fmt.Errorf("ERROR GETTING INCIDENT PARENT TASK: %v", err.Error())
		}

		parentTask.IncidentsCount = parentTask.IncidentsCount - 1
		if _, err = tx.Put(taskKey, parentTask); err != nil {
			return fmt.Errorf("ERROR DECREASING PARENT TASK INCIDENT COUNT: %v", err.Error())
		}

		key := datastore.IDKey(index, incidentID, nil)
		err = tx.Delete(key)

		if err != nil {
			return fmt.Errorf("ERROR DELETING INCIDENT: %v", err.Error())
		}

		if _, err = tx.Commit(); err != nil {
			return fmt.Errorf("ERROR INSERTING INCIDENT: %v", err.Error())
		}

		return nil
	}
	return fmt.Errorf("incident '%v' don't exist on the database for %v", incidentID, ctx.AuthUser.Email)
}
