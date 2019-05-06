package controllers

import (
	"github.com/golang/glog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pedrocelso/go-task/lib/http/authcontext"
	"github.com/pedrocelso/go-task/lib/services/incident"
	"google.golang.org/appengine/log"
)

// CreateIncident creates a incident
func CreateIncident(c *gin.Context) {
	var newIncident *incident.Incident
	var err error
	var output *incident.Incident
	var taskID int64
	taskID, err = strconv.ParseInt(c.Param("taskID"), 10, 64)
	ctx, _ := authcontext.NewAuthContext(c)

	if err != nil {
		glog.Errorf("ERROR: %v", err.Error())
	}

	if err = c.BindJSON(&newIncident); err == nil {
		if output, err = incident.Create(ctx, taskID, newIncident); err == nil {
			c.JSON(http.StatusOK, ResponseObject{"incident": output})
		}
	}

	if err != nil {
		log.Errorf(ctx.AppEngineCtx, "ERROR: %v", err.Error())
		c.JSON(http.StatusPreconditionFailed, ResponseObject{"error": err.Error()})
	}
}

// GetIncident based on its id
func GetIncident(c *gin.Context) {
	var err error
	var output *incident.Incident
	var incidentID int64
	incidentID, err = strconv.ParseInt(c.Param("incidentId"), 10, 64)
	ctx, _ := authcontext.NewAuthContext(c)

	if output, err = incident.GetByID(ctx, incidentID); err == nil {
		c.JSON(http.StatusOK, output)
	}
	if err != nil {
		c.JSON(http.StatusPreconditionFailed, ResponseObject{"error": err.Error()})
	}
}

// GetIncidents Fetch all incidents
func GetIncidents(c *gin.Context) {
	var err error
	var taskID int64
	taskID, err = strconv.ParseInt(c.Param("taskID"), 10, 64)

	var output []incident.Incident
	ctx, _ := authcontext.NewAuthContext(c)

	if output, err = incident.GetIncidents(ctx, taskID); err == nil {
		c.JSON(http.StatusOK, output)
	}

	if err != nil {
		c.JSON(http.StatusPreconditionFailed, ResponseObject{"error": err.Error()})
	}
}

// UpdateIncident Updates a incident
func UpdateIncident(c *gin.Context) {
	var err error
	var output *incident.Incident
	var updatedIncident incident.Incident

	ctx, _ := authcontext.NewAuthContext(c)

	if err = c.BindJSON(&updatedIncident); err == nil {
		var incidentID int64
		incidentID, err = strconv.ParseInt(c.Param("incidentId"), 10, 64)
		if err != nil {
			log.Errorf(ctx.AppEngineCtx, "ERROR: Failed to parse incidentID.")
			c.JSON(http.StatusPreconditionFailed, ResponseObject{"error": err.Error()})
		}

		updatedIncident.ID = incidentID
		if output, err = incident.Update(ctx, &updatedIncident); err == nil {
			c.JSON(http.StatusOK, ResponseObject{"incident": output})
		}
	}

	if err != nil {
		log.Errorf(ctx.AppEngineCtx, "ERROR: %v", err.Error())
		c.JSON(http.StatusPreconditionFailed, ResponseObject{"error": err.Error()})
	}
}

// DeleteIncident deletes a incident based on its id
func DeleteIncident(c *gin.Context) {
	var err error
	var incidentID int64
	ctx, _ := authcontext.NewAuthContext(c)
	incidentID, err = strconv.ParseInt(c.Param("incidentId"), 10, 64)
	if err != nil {
		log.Errorf(ctx.AppEngineCtx, "ERROR: %v", err.Error())
		c.JSON(http.StatusPreconditionFailed, ResponseObject{"error": err.Error()})
	}

	err = incident.Delete(ctx, incidentID)
	if err != nil {
		log.Errorf(ctx.AppEngineCtx, "ERROR: %v", err.Error())
		c.JSON(http.StatusPreconditionFailed, ResponseObject{"error": err.Error()})
	}
	c.JSON(http.StatusOK, ResponseObject{"result": "ok"})
}
