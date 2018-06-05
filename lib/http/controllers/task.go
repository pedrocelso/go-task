package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pedrocelso/go-rest-service/lib/http/authcontext"
	"github.com/pedrocelso/go-rest-service/lib/services/task"
	"google.golang.org/appengine/log"
)

// CreateTask creates a task
func CreateTask(c *gin.Context) {
	var mewTask *task.Task
	var err error
	var output *task.Task
	ctx := authcontext.NewAuthContext(c)

	if err = c.BindJSON(&mewTask); err == nil {
		if output, err = task.Create(ctx, mewTask); err == nil {
			c.JSON(http.StatusOK, ResponseObject{"task": output})
		}
	}

	if err != nil {
		log.Errorf(ctx.AppEngineCtx, "ERROR: %v", err.Error())
		c.JSON(http.StatusPreconditionFailed, ResponseObject{"error": err.Error()})
	}
}

// GetTask based on its id
func GetTask(c *gin.Context) {
	var err error
	var output *task.Task
	var taskID int64
	taskID, err = strconv.ParseInt(c.Param("taskId"), 10, 64)
	ctx := authcontext.NewAuthContext(c)

	if output, err = task.GetByID(ctx, taskID); err == nil {
		c.JSON(http.StatusOK, output)
	}
	if err != nil {
		c.JSON(http.StatusPreconditionFailed, ResponseObject{"error": err.Error()})
	}
}

// GetTasks Fetch all tasks
func GetTasks(c *gin.Context) {
	var err error

	var output []task.Task

	if output, err = task.GetTasks(authcontext.NewAuthContext(c)); err == nil {
		c.JSON(http.StatusOK, output)
	}

	if err != nil {
		c.JSON(http.StatusPreconditionFailed, ResponseObject{"error": err.Error()})
	}
}

// UpdateTask Updates a task
func UpdateTask(c *gin.Context) {
	var err error
	var output *task.Task
	var tsk task.Task

	ctx := authcontext.NewAuthContext(c)

	if err = c.BindJSON(&tsk); err == nil {
		var taskID int64
		taskID, err = strconv.ParseInt(c.Param("taskId"), 10, 64)
		if err != nil {
			log.Errorf(ctx.AppEngineCtx, "ERROR: Failed to parse taskID.")
			c.JSON(http.StatusPreconditionFailed, ResponseObject{"error": err.Error()})
		}

		tsk.ID = taskID
		if output, err = task.Update(ctx, &tsk); err == nil {
			c.JSON(http.StatusOK, ResponseObject{"task": output})
		}
	}

	if err != nil {
		log.Errorf(ctx.AppEngineCtx, "ERROR: %v", err.Error())
		c.JSON(http.StatusPreconditionFailed, ResponseObject{"error": err.Error()})
	}
}

// DeleteTask deletes a task based on its id
func DeleteTask(c *gin.Context) {
	var err error
	var taskID int64
	ctx := authcontext.NewAuthContext(c)
	taskID, err = strconv.ParseInt(c.Param("taskId"), 10, 64)
	if err != nil {
		log.Errorf(ctx.AppEngineCtx, "ERROR: %v", err.Error())
		c.JSON(http.StatusPreconditionFailed, ResponseObject{"error": err.Error()})
	}

	err = task.Delete(ctx, taskID)
	if err != nil {
		log.Errorf(ctx.AppEngineCtx, "ERROR: %v", err.Error())
		c.JSON(http.StatusPreconditionFailed, ResponseObject{"error": err.Error()})
	}
	c.JSON(http.StatusOK, ResponseObject{"result": "ok"})
}
