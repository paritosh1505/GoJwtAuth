package controller

import (
	"context"
	"fmt"
	"jwtauth/Model"
	"jwtauth/database"
	"jwtauth/helper"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
func GetUsesrDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		userval := c.Param("userId")
	}
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
		defer cancel()
	}

}
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var useractual Model.UserDb
		var userExpected Model.UserDb

		context, cancel := context.WithTimeout(context.Background(), 100*time.Second)
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
		defer cancel()

	}
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
