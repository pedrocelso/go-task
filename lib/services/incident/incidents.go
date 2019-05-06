package incident

import (
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/golang/glog"
	"github.com/pedrocelso/go-task/lib/http/authcontext"
)

const (
	index           		= `Incident`
	invalidIncidentData = `error: invalid Incident data`
	taskIndex       		= `Task`
)

// Incident defines incident attributes
type Incident struct {
	ID           int64  `json:"id"`
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description" binding:"required"`
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
	if incident == nil {
		return nil, fmt.Errorf(invalidIncidentData)
	}

	taskKey := datastore.IDKey(taskIndex, taskID, nil)
	incompleteKey := datastore.IncompleteKey(index, taskKey)
	keys = append(keys, incompleteKey)

	completeKeys, err := ctx.DataStoreClient.AllocateIDs(ctx.AppEngineCtx, keys)
	if err != nil {
		glog.Errorf("ERROR ON INCIDENT ID GENERATION: %v", err.Error())
		return nil, err
	}

	incident.ID = completeKeys[0].ID
	incident.CreationTime = makeTimestamp()

	insKey, err := ctx.DataStoreClient.Put(ctx.AppEngineCtx, completeKeys[0], incident)
	if err != nil {
		glog.Errorf("ERROR INSERTING INCIDENT: %v", err.Error())
		return nil, err
	}

	output, err = GetByID(ctx, insKey.ID)
	if err != nil {
		glog.Errorf("ERROR GETTING INCIDENT OUTPUT: %v", err.Error())
		return nil, err
	}
	return output, nil
}

// GetByID a incident based on its numeric ID
func GetByID(ctx *authcontext.Context, id int64) (*Incident, error) {
	key := datastore.IDKey(index, id, nil)
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
	q := datastore.NewQuery(index)
	completeQuery := q.Ancestor(datastore.IDKey(taskIndex, taskID, nil))
	_, err := ctx.DataStoreClient.GetAll(ctx.AppEngineCtx, completeQuery, &output)

	if err != nil {
		glog.Errorf("error fetching all incidents for %s", ctx.AuthUser.Email)
		return nil, err
	}

	if len(output) <= 0 {
		return nil, fmt.Errorf("no incidents found")
	}
	return output, nil
}

// Update incident data
func Update(ctx *authcontext.Context, incident *Incident) (*Incident, error) {
	if incident == nil || (incident.Name == `` && incident.Description == ``) {
		return nil, fmt.Errorf(invalidIncidentData)
	}

	output, _ := GetByID(ctx, incident.ID)
	if output != nil {
		key := datastore.IDKey(index, incident.ID, nil)
		incident.UpdateTime = makeTimestamp()
		_, err := ctx.DataStoreClient.Put(ctx.AppEngineCtx, key, incident)

		if err != nil {
			glog.Errorf("ERROR UPDATING INCIDENT: %v", err.Error())
			return nil, err
		}

		output, err = GetByID(ctx, incident.ID)
		if err != nil {
			glog.Errorf("ERROR GETTING INCIDENT OUTPUT: %v", err.Error())
			return nil, err
		}
		return output, nil
	}
	return nil, fmt.Errorf(`incident '%v' not found`, incident.ID)
}

// Delete a incident based on its id.
func Delete(ctx *authcontext.Context, incidentID int64) error {
	var output *Incident
	output, _ = GetByID(ctx, incidentID)

	if output != nil {
		glog.Infof("Deleting incident: %v for %v", incidentID, ctx.AuthUser.Email)
		key := datastore.IDKey(index, incidentID, nil)
		err := ctx.DataStoreClient.Delete(ctx.AppEngineCtx, key)

		if err != nil {
			glog.Errorf("ERROR DELETING INCIDENT: %v", err.Error())
			return err
		}
		return nil
	}
	return fmt.Errorf("incident '%v' don't exist on the database for %v", incidentID, ctx.AuthUser.Email)
}
