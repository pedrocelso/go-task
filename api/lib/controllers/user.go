package controllers

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"github.com/gin-gonic/gin"
	"github.com/pedrocelso/go-rest-service/api/lib/services/user"
)

// CreateUser creates an User
func CreateUser(c *gin.Context) {
	var usr *user.User
	var err error
	var output *user.User
	ctx := appengine.NewContext(c.Request)

	if err = c.BindJSON(&usr); err == nil {
		if output, err = user.Create(ctx, usr); err == nil {
			c.JSON(http.StatusOK, ResponseObject{"user": output})
		}
	}

	if err != nil {
		log.Errorf(ctx, "ERROR: %v", err.Error())
		c.JSON(http.StatusPreconditionFailed, ResponseObject{"error": err.Error()})
	}
}

// GetUser based on its email
func GetUser(c *gin.Context) {
	var err error
	var output *user.User
	usrEmail := c.Param("userEmail")
	ctx := appengine.NewContext(c.Request)

	if output, err = user.GetByEmail(ctx, usrEmail); err == nil {
		c.JSON(http.StatusOK, output)
	}
	if err != nil {
		c.JSON(http.StatusPreconditionFailed, ResponseObject{"error": err.Error()})
	}
}

// GetUsers Fectch all users
func GetUsers(c *gin.Context) {
	var err error
	ctx := appengine.NewContext(c.Request)
	var output *[]user.User

	if output, err = user.GetUsers(ctx); err == nil {
		c.JSON(http.StatusOK, output)
	}

	if err != nil {
		c.JSON(http.StatusPreconditionFailed, ResponseObject{"error": err.Error()})
	}
}

// UpdateUser Updates an user
func UpdateUser(c *gin.Context) {
	var usr *user.User
	var err error
	var output *user.User
	ctx := appengine.NewContext(c.Request)

	if err = c.BindJSON(&usr); err == nil {
		if output, err = user.Update(ctx, usr); err == nil {
			c.JSON(http.StatusOK, ResponseObject{"user": output})
		}
	}

	if err != nil {
		log.Errorf(ctx, "ERROR: %v", err.Error())
		c.JSON(http.StatusPreconditionFailed, ResponseObject{"error": err.Error()})
	}
}

// DeleteUser deletes an user based on its email
func DeleteUser(c *gin.Context) {
	usrEmail := c.Param("userEmail")
	ctx := appengine.NewContext(c.Request)

	err := user.Delete(ctx, usrEmail)

	if err != nil {
		log.Errorf(ctx, "ERROR: %v", err.Error())
		c.JSON(http.StatusPreconditionFailed, ResponseObject{"error": err.Error()})
	}
	c.JSON(http.StatusOK, ResponseObject{"result": "ok"})
}
