package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/pedrocelso/go-task/lib/http/authcontext"
	"github.com/pedrocelso/go-task/lib/services/user"
)

type claims struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	jwt.StandardClaims
}

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
	spew.Dump(c)
	var err error
	var output *user.Basic
	usrEmail := c.Param("userEmail")
	ctx, _ := authcontext.NewAuthContext(c)
	spew.Dump(ctx)

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

// AuthenticateUser authenticates an user
func AuthenticateUser(c *gin.Context) {
	var usr *user.Full
	var err error
	var tokenString string
	var output *user.Full
	var ctx *authcontext.Context

	ctx, err = authcontext.NewAuthContext(c)

	if err != nil {
		glog.Errorf("ERROR: %v", err.Error())
	}

	if err = c.BindJSON(&usr); err == nil {
		if output, err = user.GetFullByEmail(ctx, usr.Email); err == nil {
			if err = bcrypt.CompareHashAndPassword([]byte(output.Password), []byte(usr.Password)); err == nil {
				expirationTime := time.Now().Add(5 * time.Minute)
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
					Name:  usr.Name,
					Email: usr.Email,
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: expirationTime.Unix(),
					},
				})
				if tokenString, err = token.SignedString(os.Getenv("JWT_SECRET")); err == nil {
					c.JSON(http.StatusOK, ResponseObject{"token": tokenString})
				}
			}
		}
	}

	if err != nil {
		glog.Errorf("ERROR: %v", err.Error())
		c.JSON(http.StatusUnauthorized, ResponseObject{"error": err.Error()})
	}
}
