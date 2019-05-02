package controllers

import (
	"net/http"

	"github.com/golang/glog"

	"github.com/gin-gonic/gin"
	"github.com/pedrocelso/go-task/lib/http/authcontext"
	"github.com/pedrocelso/go-task/lib/services/user"
)

// CreateUser creates an User
func CreateUser(c *gin.Context) {
	var usr *user.Full
	var err error
	var output *user.Basic
	var ctx *authcontext.Context
	ctx, err = authcontext.NewAuthContext(c)

	if err != nil {
		glog.Errorf("ERROR: %v", err.Error())
	}

	if err = c.BindJSON(&usr); err == nil {
		if output, err = user.Create(ctx, usr); err == nil {
			c.JSON(http.StatusOK, ResponseObject{"user": output})
		}
	}

	if err != nil {
		glog.Errorf("ERROR: %v", err.Error())
		c.JSON(http.StatusPreconditionFailed, ResponseObject{"error": err.Error()})
	}
}

// GetUser based on its email
func GetUser(c *gin.Context) {
	var err error
	var output *user.Basic
	usrEmail := c.Param("userEmail")
	ctx, _ := authcontext.NewAuthContext(c)

	if output, err = user.GetByEmail(ctx, usrEmail); err == nil {
		c.JSON(http.StatusOK, output)
	}
	if err != nil {
		c.JSON(http.StatusPreconditionFailed, ResponseObject{"error": err.Error()})
	}
}

// GetUsers Fetch all users
func GetUsers(c *gin.Context) {
	var err error
	var ctx *authcontext.Context

	var output []user.Basic
	ctx, err = authcontext.NewAuthContext(c)

	if err != nil {
		glog.Error(err.Error())
	}

	if output, err = user.GetUsers(ctx); err == nil {
		c.JSON(http.StatusOK, output)
	}

	if err != nil {
		c.JSON(http.StatusPreconditionFailed, ResponseObject{"error": err.Error()})
	}
}

// UpdateUser Updates an user
func UpdateUser(c *gin.Context) {
	var usr *user.Full
	var err error
	var output *user.Basic
	ctx, _ := authcontext.NewAuthContext(c)

	if err = c.BindJSON(&usr); err == nil {
		if output, err = user.Update(ctx, usr); err == nil {
			c.JSON(http.StatusOK, ResponseObject{"user": output})
		}
	}

	if err != nil {
		glog.Errorf("ERROR: %v", err.Error())
		c.JSON(http.StatusPreconditionFailed, ResponseObject{"error": err.Error()})
	}
}

// DeleteUser deletes an user based on its email
func DeleteUser(c *gin.Context) {
	usrEmail := c.Param("userEmail")
	ctx, _ := authcontext.NewAuthContext(c)

	err := user.Delete(ctx, usrEmail)

	if err != nil {
		glog.Errorf("ERROR: %v", err.Error())
		c.JSON(http.StatusPreconditionFailed, ResponseObject{"error": err.Error()})
	}
	c.JSON(http.StatusOK, ResponseObject{"result": "ok"})
}
