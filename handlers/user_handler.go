package handlers

import (
	"net/http"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/labstack/echo"
	"github.com/michelaquino/api_poc_angular_react/context"
)

type userModel struct {
	ID     bson.ObjectId `bson:"_id,omitempty"`
	Name   string        `bson:"name,omitempty"`
	Email  string        `bson:"email,omitempty"`
	Gender string        `bson:"gender,omitempty"`
}

type userRepr struct {
	ID     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Email  string `json:"email,omitempty"`
	Gender string `json:"gender,omitempty"`
}

type UserHandler struct{}

func (UserHandler) GetAllUsers(echoContext echo.Context) error {
	logger := context.GetAPIContext().GetLogger()
	logger.Info("[UserHandler][GetAllUsers] Getting all users")

	userList, err := getUserListFromDatabase()
	if err != nil {
		if err == mgo.ErrNotFound {
			logger.Info("[UserHandler][GetAllUsers] Users not found")
			return echoContext.JSON(http.StatusOK, []userModel{})
		}

		logger.Info("[UserHandler][GetAllUsers] Unexpected error: ", err.Error())
		return echoContext.NoContent(http.StatusInternalServerError)
	}

	if userList == nil {
		logger.Info("[UserHandler][CreateUser] userList is nil")
		userList = []userModel{}
	}

	userReprList := []userRepr{}
	for _, user := range userList {
		userRepresentation := userRepr{
			ID:     user.ID.Hex(),
			Email:  user.Email,
			Name:   user.Name,
			Gender: user.Gender,
		}

		userReprList = append(userReprList, userRepresentation)
	}

	logger.Info("[UserHandler][GetAllUsers] Get all users with success")
	return echoContext.JSON(http.StatusOK, userReprList)
}

func (UserHandler) CreateUser(echoContext echo.Context) error {
	logger := context.GetAPIContext().GetLogger()
	logger.Info("[UserHandler][CreateUser] Getting all users")

	user := new(userRepr)
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

	newUserModel := &userModel{
		Email:  user.Email,
		Name:   user.Name,
		Gender: user.Gender,
	}
	err := createUserOnDatabase(newUserModel)
	if err != nil {
		logger.Info("[UserHandler][CreateUser] Error on created user: ", err.Error())
		return echoContext.NoContent(http.StatusInternalServerError)
	}

	userRepresentation := userRepr{
		ID:     newUserModel.ID.Hex(),
		Email:  newUserModel.Email,
		Name:   newUserModel.Name,
		Gender: newUserModel.Gender,
	}

	logger.Info("[UserHandler][CreateUser] User created with success")
	return echoContext.JSON(http.StatusCreated, userRepresentation)
}

func (UserHandler) DeleteUser(echoContext echo.Context) error {
	userID := echoContext.Param("id")

	if err := deleteUser(userID); err != nil {
		return echoContext.NoContent(http.StatusInternalServerError)
	}

	return echoContext.NoContent(http.StatusOK)
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

	id := bson.NewObjectId()
	user.ID = id

	connection := session.DB("api").C("users")
	err := connection.Insert(user)
	if err != nil {
		logger.Error("[createUserOnDatabase] Error on insert user on database: ", err.Error())
		return err
	}

	return nil
}

func deleteUser(userID string) error {
	logger := context.GetAPIContext().GetLogger()
	logger.Infof("[deleteUser] Delete user with id %s from database", userID)

	session := context.GetAPIContext().GetMongoSession()
	defer session.Close()

	objID := bson.ObjectIdHex(userID)

	connection := session.DB("api").C("users")
	err := connection.RemoveId(objID)
	if err != nil {
		logger.Error("[deleteUser] Error on delete user on database: ", err.Error())
		return err
	}

	logger.Infof("[deleteUser] User with id %s delete with success from database", userID)
	return nil
}
