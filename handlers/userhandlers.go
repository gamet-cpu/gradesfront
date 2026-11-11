package handler

import (
	"context"
	"net/http"
	"time"

	"grades/model"
	"grades/repository"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	repo *repository.UserRepository
}

func NewUserHandler(r *repository.UserRepository) *UserHandler {
	return &UserHandler{repo: r}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "hash error"})
		return
	}

	err = h.repo.Create(context.Background(),
		req.Email,
		string(hash),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "registered"})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, hash, err := h.repo.GetByEmail(context.Background(), req.Email)
	if err != nil ||
		bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	t, err := token.SignedString([]byte(viper.GetString("jwt.secret")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": t})
}

func (h *UserHandler) GetMe(c *gin.Context) {
	id := c.GetInt("user_id")

	user, err := h.repo.GetByID(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) Update(c *gin.Context) {
	id := c.GetInt("user_id")

	var req struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Grade string `json:"grade"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.repo.Update(
		context.Background(),
		id,
		req.Name,
		req.Age,
		req.Grade,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func (h *UserHandler) Delete(c *gin.Context) {
	id := c.GetInt("user_id")

	err := h.repo.Delete(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *UserHandler) CreateSubj(c *gin.Context) {

	var subject model.Subjects

	if err := c.ShouldBindJSON(&subject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат данных",
		})
		return
	}

	userIDValue, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Пользователь не авторизован",
		})
		return
	}

	userID, ok := userIDValue.(int)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Некорректный user_id",
		})
		return
	}

	subject.UserId = userID

	err := h.repo.CreateSubj(
		c.Request.Context(),
		subject,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Оценки успешно добавлены",
	})
}

func (h *UserHandler) GetSubj(c *gin.Context) {

	userIDValue, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Пользователь не авторизован",
		})
		return
	}

	userID, ok := userIDValue.(int)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Некорректный user_id",
		})
		return
	}

	subject, err := h.repo.GetSubj(
		c.Request.Context(),
		userID,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Оценки не найдены",
		})
		return
	}

	c.JSON(http.StatusOK, subject)
}

func (h *UserHandler) UpdateSubj(c *gin.Context) {

	var subject model.Subjects

	if err := c.ShouldBindJSON(&subject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат данных",
		})
		return
	}

	userIDValue, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Пользователь не авторизован",
		})
		return
	}

	userID, ok := userIDValue.(int)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Некорректный user_id",
		})
		return
	}

	subject.UserId = userID

	err := h.repo.UpdateSubj(
		c.Request.Context(),
		subject,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Оценки успешно изменены",
	})
}

func (h *UserHandler) DeleteSubj(c *gin.Context) {

	userIDValue, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Пользователь не авторизован",
		})
		return
	}

	userID, ok := userIDValue.(int)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Некорректный user_id",
		})
		return
	}

	err := h.repo.DeleteSubj(
		c.Request.Context(),
		userID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Оценки удалены",
	})
}
