package controllers

import (
	"bytes"
	"encoding/json"
	"my-gin-project/src/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Initialisation de la DB en mémoire pour les tests
func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.Item{}, &models.User{})
	models.DB = db
	return db
}

// Initialisation de Gin pour les tests
func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	ctrl := &Controller{}

	r.GET("/items", ctrl.GetItems)
	r.GET("/items/:id", ctrl.GetItemByID)
	r.POST("/items", ctrl.CreateItem)
	r.PUT("/items/:id", ctrl.UpdateItem)
	r.DELETE("/items/:id", ctrl.DeleteItem)
	r.POST("/register", ctrl.Register)
	r.POST("/login", ctrl.Login)

	return r
}

func TestCreateGetItem(t *testing.T) {
	setupTestDB()
	router := setupRouter()

	t.Log("Création d'un item TestItem")
	item := models.Item{Name: "TestItem", Price: 12.5}
	jsonValue, _ := json.Marshal(item)
	req, _ := http.NewRequest("POST", "/items", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	t.Logf("Code retour création: %d, Body: %s", resp.Code, resp.Body.String())
	if resp.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.Code)
	}

	t.Log("Récupération de l'item créé")
	req, _ = http.NewRequest("GET", "/items/1", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	t.Logf("Code retour récupération: %d, Body: %s", resp.Code, resp.Body.String())
	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.Code)
	}

	var returnedItem models.Item
	json.Unmarshal(resp.Body.Bytes(), &returnedItem)
	t.Logf("Item récupéré: %+v", returnedItem)
	if returnedItem.Name != "TestItem" {
		t.Errorf("Expected item name 'TestItem', got '%s'", returnedItem.Name)
	}
}

func TestUpdateDeleteItem(t *testing.T) {
	setupTestDB()
	router := setupRouter()

	// Créer un item pour update/delete
	models.DB.Create(&models.Item{Name: "OldItem", Price: 5.0})

	// Update
	update := models.Item{Name: "UpdatedItem", Price: 10.0}
	jsonValue, _ := json.Marshal(update)
	req, _ := http.NewRequest("PUT", "/items/1", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.Code)
	}

	// Delete
	req, _ = http.NewRequest("DELETE", "/items/1", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", resp.Code)
	}
}

func TestRegisterLogin(t *testing.T) {
	setupTestDB()
	router := setupRouter()

	// Register
	user := models.User{Username: "testuser", Password: "password"}
	jsonValue, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.Code)
	}

	// Login
	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.Code)
	}

	var body map[string]string
	json.Unmarshal(resp.Body.Bytes(), &body)
	if _, ok := body["token"]; !ok {
		t.Errorf("Expected token in response")
	}
}
