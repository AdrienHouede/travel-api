package controllers

import (
	"my-gin-project/src/models"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Controller struct {
	DB *gorm.DB
}

var jwtSecret = []byte("secret") // à mettre dans une variable d'environnement en prod

// GET /items - récupérer tous les items
// @Summary Get all items
// @Description Retrieve list of items (protected route)
// @Tags items
// @Produce json
// @Success 200 {array} models.Item
// @Failure 401 {object} map[string]string
// @Security ApiKeyAuth
// @Router /items [get]
func (c *Controller) GetItems(ctx *gin.Context) {
	var items []models.Item
	if err := models.DB.Find(&items).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items"})
		return
	}
	ctx.JSON(http.StatusOK, items)
}

// GET /items/:id - récupérer un item par ID
// @Summary Get item by ID
// @Description Retrieve a single item
// @Tags items
// @Produce json
// @Param id path int true "Item ID"
// @Success 200 {object} models.Item
// @Failure 404 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Security ApiKeyAuth
// @Router /items/{id} [get]
func (c *Controller) GetItemByID(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var item models.Item
	if err := models.DB.First(&item, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}
	ctx.JSON(http.StatusOK, item)
}

// POST /items - créer un item
// @Summary Create a new item
// @Description Add a new item
// @Tags items
// @Accept json
// @Produce json
// @Param item body models.Item true "Item info"
// @Success 201 {object} models.Item
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Security ApiKeyAuth
// @Router /items [post]
func (c *Controller) CreateItem(ctx *gin.Context) {
	var item models.Item
	if err := ctx.ShouldBindJSON(&item); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if err := models.DB.Create(&item).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item"})
		return
	}
	ctx.JSON(http.StatusCreated, item)
}

// PUT /items/:id - mettre à jour un item
// @Summary Update an item
// @Description Update item details
// @Tags items
// @Accept json
// @Produce json
// @Param id path int true "Item ID"
// @Param item body models.Item true "Item info"
// @Success 200 {object} models.Item
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Security ApiKeyAuth
// @Router /items/{id} [put]
func (c *Controller) UpdateItem(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var item models.Item
	if err := models.DB.First(&item, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	var input models.Item
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	item.Name = input.Name
	item.Price = input.Price

	if err := models.DB.Save(&item).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item"})
		return
	}

	ctx.JSON(http.StatusOK, item)
}

// DELETE /items/:id - supprimer un item
// @Summary Delete an item
// @Description Delete an item by ID
// @Tags items
// @Produce json
// @Param id path int true "Item ID"
// @Success 204 {object} nil
// @Failure 404 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Security ApiKeyAuth
// @Router /items/{id} [delete]
func (c *Controller) DeleteItem(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := models.DB.Delete(&models.Item{}, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}
	ctx.Status(http.StatusNoContent)
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.User true "User info"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /register [post]
func (c *Controller) Register(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}
	user.Password = string(hash)

	if err := models.DB.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving user"})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "User registered"})
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.User true "User info"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /login [post]
func (c *Controller) Login(ctx *gin.Context) {
	var input models.User
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var user models.User
	if err := models.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// Génération du token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(), // expiration 72h
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		ctx.Next()
	}
}
