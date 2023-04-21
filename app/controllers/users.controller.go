package controllers

import (
	"api-connect-mongodb-atlas/pkg/middleware"
	"api-connect-mongodb-atlas/pkg/models"
	"api-connect-mongodb-atlas/pkg/utils"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	g4vercel "github.com/tbxark/g4vercel"
)

type IuserController interface {
	RegisterAccount(c *g4vercel.Context)
	GetUserAccount(c *g4vercel.Context)
	GetUsersAccount(c *g4vercel.Context)
}

type PropsUserController struct {
	MongoDB            *mongo.Database
	MainCollectionDB   *mongo.Collection
	HttpResponseWriter http.ResponseWriter
	HttpRequest        *http.Request
	Passport           middleware.CustomPassportInterface
}

func NewUserControllers(DB *mongo.Database, w http.ResponseWriter, r *http.Request) IuserController {
	return &PropsUserController{
		MongoDB:            DB,
		MainCollectionDB:   DB.Collection("users"),
		HttpResponseWriter: w,
		HttpRequest:        r,
	}
}

type response1 struct {
	Page   int
	Fruits []string
}

func (ur *PropsUserController) RegisterAccount(c *g4vercel.Context) {
	requestUser := new(models.ModuleProfile)
	collection := ur.MainCollectionDB

	// parse the request body and bind it to the user instance
	if err := json.NewDecoder(c.Req.Body).Decode(requestUser); err != nil {
		c.JSON(http.StatusBadRequest, models.NewBaseErrorResponse(err.Error(), http.StatusBadRequest))
	}

	hashPassword, _ := utils.HashPassword(requestUser.Password)
	requestUser.Password = hashPassword

	res, err := collection.InsertOne(context.Background(), requestUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewBaseErrorResponse(err.Error(), http.StatusBadRequest))
	}
	id := res.InsertedID

	c.JSON(http.StatusOK, models.NewBaseResponse(id, http.StatusOK))
}

func (ur *PropsUserController) GetUserAccount(c *g4vercel.Context) {
	// ur.Passport.DeserializeUser(c)
	var result models.ModuleProfile
	cookie_user, _ := ur.HttpRequest.Cookie("user_id")
	user_id := cookie_user.Value
	docID, _ := primitive.ObjectIDFromHex(user_id)
	bson := bson.M{"_id": docID}

	collection := ur.MainCollectionDB
	err_db := collection.FindOne(context.Background(), bson).Decode(&result)
	if err_db != nil {
		c.JSON(http.StatusBadRequest, models.NewBaseErrorResponse(err_db.Error(), http.StatusBadRequest))
		return
	}

	c.JSON(http.StatusOK, models.NewBaseResponse(result, http.StatusOK))
}

func (ur *PropsUserController) GetUsersAccount(c *g4vercel.Context) {
	// ur.Passport.DeserializeUser(c)
	collection := ur.MainCollectionDB

	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	var results []models.ModuleProfile
	for cursor.Next(context.Background()) {
		var bson models.ModuleProfile
		err := cursor.Decode(&bson)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, bson)
	}

	c.JSON(http.StatusOK, models.NewBaseResponse(results, http.StatusOK))
}
