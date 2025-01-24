package internal

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"microservice/initializers"
	kafkabroker "microservice/kafka-broker"
	"microservice/model"
	"net/http"
	"os"
	"time"
)

func SignUp(c *gin.Context) {
	var user model.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to generate password",
		})
	}
	user.Password = string(hash)

	result := initializers.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})

}

func LogIn(c *gin.Context) {

	var body struct {
		Link     string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
	}

	var user model.User
	initializers.DB.First(&user, "link = ?", body.Link)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid link/password",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid password",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	fmt.Println(os.Getenv("SECRET"))
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{})
}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")

	fmt.Println(user)

	c.JSON(http.StatusOK, gin.H{
		"message": "User validation success",
	})
}

func Recommendation(c *gin.Context) {
	var users []model.User
	result := initializers.DB.Select("username", "age", "about").Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get users",
		})
	}

	userInBytes, err := json.Marshal(users)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Failed to convert into bytes user": err.Error()})
		return
	}

	err = kafkabroker.PushUserToQueue("user_recommendation", userInBytes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Failed to push user to queue": err.Error()})
		return
	}

	response := map[string]interface{}{
		"success": true,
		//"msg":     strconv.FormatInt(user.ID, 10) + user.Username + "was successfully registered",
	}

	c.JSON(http.StatusOK, response)
	return
}
