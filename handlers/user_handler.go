package handlers

import (
	"net/http"

	mgo "gopkg.in/mgo.v2"

	"github.com/labstack/echo"
	"github.com/michelaquino/api_poc_angular_react/context"
)

type userModel struct {
	ID     string `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string `json:"name,omitempty" bson:"name,omitempty"`
	Email  string `json:"email,omitempty" bson:"email,omitempty"`
	Gender string `json:"gender,omitempty" bson:"gender,omitempty"`
}

type UserHandler struct{}

func (UserHandler) GetAllUsers(echoContext echo.Context) error {
	logger := context.GetAPIContext().GetLogger()
	logger.Info("[UserHandler][GetAllUsers] Getting all users")

	user := new(userModel)
	if err := echoContext.Bind(user); err != nil {
		logger.Error("[UserHandler][GetAllUsers] Error on bind request json")
		return echoContext.String(http.StatusBadRequest, "Invalid parameters")
	}

	userList, err := getUserListFromDatabase()
	if err != nil {
		if err == mgo.ErrNotFound {
			logger.Info("[UserHandler][GetAllUsers] Users not found")
			return echoContext.JSON(http.StatusOK, []userModel{})
		}

		logger.Info("[UserHandler][GetAllUsers] Unexpected error: ", err.Error())
		return echoContext.NoContent(http.StatusInternalServerError)
	}

	logger.Info("[UserHandler][CreateUser] userList: ", userList)
	if userList == nil {
		logger.Info("[UserHandler][CreateUser] userList is nil")
		userList = []userModel{}
	}

	logger.Info("[UserHandler][GetAllUsers] Get all users with success")
	return echoContext.JSON(http.StatusOK, userList)
}

func (UserHandler) CreateUser(echoContext echo.Context) error {
	logger := context.GetAPIContext().GetLogger()
	logger.Info("[UserHandler][CreateUser] Getting all users")

	user := new(userModel)
	if err := echoContext.Bind(user); err != nil {
		logger.Error("[UserHandler][CreateUser] Error on bind request json: ", err.Error())
		return echoContext.String(http.StatusBadRequest, "Invalid parameters")
	}

	if user.Name == "" {
		logger.Error("[UserHandler][CreateUser] User name is required")
		return echoContext.String(http.StatusBadRequest, "User name is required")
	}

	if user.Email == "" {
		logger.Error("[UserHandler][CreateUser] User email is required")
		return echoContext.String(http.StatusBadRequest, "User email is required")
	}

	if user.Gender == "" {
		logger.Error("[UserHandler][CreateUser] User gender is required")
		return echoContext.String(http.StatusBadRequest, "User gander is required")
	}

	err := createUserOnDatabase(user)
	if err != nil {
		logger.Info("[UserHandler][CreateUser] Error on created user: ", err.Error())
		return echoContext.NoContent(http.StatusInternalServerError)
	}

	logger.Info("[UserHandler][CreateUser] User created with success")
	return echoContext.NoContent(http.StatusCreated)
}

func getUserListFromDatabase() ([]userModel, error) {
	logger := context.GetAPIContext().GetLogger()
	logger.Info("[getUserListFromDatabase] Getting all users from database")

	session := context.GetAPIContext().GetMongoSession()
	defer session.Close()

	connection := session.DB("api").C("users")
	var userList []userModel

	err := connection.Find(nil).All(&userList)
	if err != nil {
		logger.Error("[getUserListFromDatabase] Error on find users on database: ", err.Error())
		return nil, err
	}

	return userList, nil
}

func createUserOnDatabase(user *userModel) error {
	logger := context.GetAPIContext().GetLogger()
	logger.Info("[createUserOnDatabase] Getting all users from database")

	session := context.GetAPIContext().GetMongoSession()
	defer session.Close()

	connection := session.DB("api").C("users")
	err := connection.Insert(user)
	if err != nil {
		logger.Error("[createUserOnDatabase] Error on insert user on database: ", err.Error())
		return err
	}

	return nil
}
