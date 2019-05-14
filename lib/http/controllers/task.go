package controllers

import (
	"net/http"
	"strconv"

	"github.com/golang/glog"

	"github.com/gin-gonic/gin"
	"github.com/pedrocelso/go-task/lib/http/authcontext"
	"github.com/pedrocelso/go-task/lib/services/task"
)

// CreateTask creates a task
func CreateTask(c *gin.Context) {
	var newTask *task.Task
	var err error
	var output *task.Task
	ctx, err := authcontext.NewAuthContext(c)

	if err != nil {
		glog.Errorf("ERROR: %v", err.Error())
	}

	if err = c.BindJSON(&newTask); err == nil {
		if output, err = task.Create(ctx, newTask); err == nil {
			c.JSON(http.StatusOK, ResponseObject{"task": output})
		}
	}

	if err != nil {
		c.JSON(http.StatusPreconditionFailed, ResponseObject{"error": err.Error()})
		glog.Errorf("ERROR: %v", err.Error())
	}
}

// GetTask based on its id
func GetTask(c *gin.Context) {
	var output *task.Task
	var taskID int64
	ctx, err := authcontext.NewAuthContext(c)

	if err != nil {
		glog.Errorf("ERROR: %v", err.Error())
	}

	taskID, err = strconv.ParseInt(c.Param("taskId"), 10, 64)

	if output, err = task.GetByID(ctx, taskID); err == nil {
		c.JSON(http.StatusOK, output)
	}
	if err != nil {
		c.JSON(http.StatusPreconditionFailed, ResponseObject{"error": err.Error()})
	}
}

// GetTasks Fetch all tasks
func GetTasks(c *gin.Context) {
	var output []task.Task
	ctx, err := authcontext.NewAuthContext(c)

	if err != nil {
		glog.Errorf("ERROR: %v", err.Error())
	}

	if output, err = task.GetTasks(ctx); err == nil {
		c.JSON(http.StatusOK, output)
	}

	if err != nil {
		c.JSON(http.StatusPreconditionFailed, ResponseObject{"error": err.Error()})
	}
}

// UpdateTask Updates a task
func UpdateTask(c *gin.Context) {
	var output *task.Task
	var tsk task.Task

	ctx, err := authcontext.NewAuthContext(c)

	if err != nil {
		glog.Errorf("ERROR: %v", err.Error())
	}

	if err = c.BindJSON(&tsk); err == nil {
		var taskID int64
		taskID, err = strconv.ParseInt(c.Param("taskId"), 10, 64)
		if err != nil {
			c.JSON(http.StatusPreconditionFailed, ResponseObject{"error": err.Error()})
			glog.Errorf("ERROR: Failed to parse taskID.")
		}

		tsk.ID = taskID
		if output, err = task.Update(ctx, &tsk); err == nil {
			c.JSON(http.StatusOK, ResponseObject{"task": output})
		}
	}

	if err != nil {
		c.JSON(http.StatusPreconditionFailed, ResponseObject{"error": err.Error()})
		glog.Errorf("ERROR: %v", err.Error())
	}
}

// DeleteTask deletes a task based on its id
func DeleteTask(c *gin.Context) {
	var taskID int64
	ctx, err := authcontext.NewAuthContext(c)

	if err != nil {
		glog.Errorf("ERROR: %v", err.Error())
	}

	taskID, err = strconv.ParseInt(c.Param("taskId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusPreconditionFailed, ResponseObject{"error": err.Error()})
		glog.Errorf("ERROR: %v", err.Error())
	}

	err = task.Delete(ctx, taskID)
	if err != nil {
		c.JSON(http.StatusPreconditionFailed, ResponseObject{"error": err.Error()})
		glog.Errorf("ERROR: %v", err.Error())
	}
	c.JSON(http.StatusOK, ResponseObject{"result": "ok"})
}
