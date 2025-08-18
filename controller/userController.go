package controller

import (
	"context"
	"fmt"
	"jwtauth/Model"
	"jwtauth/database"
	"jwtauth/helper"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var TableOpen *mongo.Collection = database.CreateCollection(database.ClientInstance, "UserTable")
var validateVal = validator.New()

func GetUserDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")

		errval := helper.MatchUserId(c, userId)
		if errval != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": errval.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var userModel = Model.UserDb{}
		err := TableOpen.FindOne(ctx, bson.M{"user_id": userId}).Decode(&userModel)
		if err != nil {
			c.JSON(http.StatusGatewayTimeout, gin.H{
				"error": "User name not found",
			})
			return
		}
		c.JSON(http.StatusAccepted, userModel)

	}
}

func HashPassword(password string) string {
	passwordEncrypt, _ := bcrypt.GenerateFromPassword([]byte(password), 2)
	return string(passwordEncrypt)
}
func SingupUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userStruct Model.UserDb
		err := c.BindJSON(&userStruct)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		errorValidator := validateVal.Struct(userStruct)
		if errorValidator != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": errorValidator.Error(),
			})
			return
		}
		ctxdec, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		countUser, err := TableOpen.CountDocuments(ctxdec, bson.M{
			"mongoEmail": userStruct.Email,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		if countUser > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "This user email already exist"})
			return
		}
		crypticpwd := HashPassword(*userStruct.Password)
		userStruct.Password = &crypticpwd
		userStruct.UpdatedAt = time.Now()
		userStruct.CreatedAt = time.Now()
		userStruct.Id = primitive.NewObjectID()
		userStruct.User_id = userStruct.Id.Hex()
		access_token, refresh_token, _ := helper.TokenGeneration(*userStruct.Name, *userStruct.Email)
		userStruct.Token_gen = &access_token
		userStruct.Refresh_token = &refresh_token
		insertVal, err := TableOpen.InsertOne(ctxdec, userStruct)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Error": "Error while inserting the entry",
			})
		}
		c.JSON(http.StatusAccepted, insertVal)
	}

}
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var useractual Model.UserDb
		var userExpected Model.UserDb

		context, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		err := TableOpen.FindOne(context, bson.M{
			"mongoEmail": useractual.Email,
		}).Decode(&userExpected)
		if err != nil {
			fmt.Println("**********************Email not found login not possible**********************")
			return
		}
		status := CheckPassword(*useractual.Password, *userExpected.Password)

		if !status {
			c.JSON(http.StatusForbidden, gin.H{"Error": "Invalid Password"})
			return
		}
		c.JSON(http.StatusAccepted, gin.H{"Message": "Login Successfully"})

		if userExpected.Email == nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "emal not found"})
			log.Panic("Error in missing email")
		}
		access_token, refresh_token, _ := helper.TokenGeneration(*userExpected.Name, *userExpected.Email)
		UpdateTokenAfterLogin(access_token, refresh_token, userExpected.User_id)

	}
}
func UpdateTokenAfterLogin(access_token string, refresh_token string, userid string) {
	contextval, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	var ValueEnter bson.D

	ValueEnter = append(ValueEnter, bson.E{Key: "Token_gen", Value: access_token})
	ValueEnter = append(ValueEnter, bson.E{Key: "refresh_token", Value: refresh_token})
	UpdatedAt := time.Now()
	ValueEnter = append(ValueEnter, bson.E{Key: "UpdateAt", Value: UpdatedAt})
	filterEntry := bson.E{Key: "User_id", Value: userid}
	upserStatus := true
	noDocPresentThenUpdatea := options.UpdateOptions{
		Upsert: &upserStatus,
	}
	_, err := TableOpen.UpdateOne(
		contextval,
		filterEntry,
		bson.D{
			{Key: "$set", Value: ValueEnter},
		},
		&noDocPresentThenUpdatea,
	)
	if err != nil {
		log.Fatal("Error while updating the data")
	}
	defer cancel()

}
func CheckPassword(actualPwd string, expectedPwd string) bool {
	status := true
	erroval := bcrypt.CompareHashAndPassword([]byte(actualPwd), []byte(expectedPwd))
	if erroval != nil {
		status = false
		log.Fatal("*********************************Password is not matched***************************")
	}
	return status
}

func AggreGator() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), time.Millisecond*10)
		defer cancel()
		err := helper.CheckUserPermission(c, "ADMIN")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		recordEntry, err := strconv.Atoi(c.Query("recordEntry"))
		if err != nil || recordEntry < 1 {
			recordEntry = 10
		}
		recordPage, err := strconv.Atoi(c.Query("Page"))
		if err != nil || recordPage < 1 {
			recordPage = 1
		}
		matchstage := bson.D{
			{Key: "$match", Value: bson.D{{}}},
		}
		grpStage := bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: bson.D{{Key: "_id", Value: "null"}}},
				{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
				{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
			}},
		}
		third := bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "total_count", Value: 1},
				{Key: "user_items", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", recordEntry, recordPage}}}},
			}},
		}
		result, err := TableOpen.Aggregate(ctx, mongo.Pipeline{
			matchstage, grpStage, third,
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		var users []bson.D
		errval := result.All(ctx, &users)
		if errval != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error:": "Error in result.All function"})
		}
		if len(users) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "User not found"})
			return
		}
		c.JSON(http.StatusOK, users[0])
	}

}
