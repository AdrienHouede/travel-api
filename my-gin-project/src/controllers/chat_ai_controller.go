package controllers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// Modèles GORM
type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Password string
}

type ConversationHistory struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	Sender    string // "user" ou "bot"
	Message   string
	CreatedAt time.Time
}

func (ConversationHistory) TableName() string {
	return "conversation_history"
}

type AIMessage struct {
	User string `json:"user" example:"Alice"`
	Text string `json:"text" example:"Bonjour, raconte-moi une blague."`
}

type AIResponse struct {
	Bot string `json:"bot" example:"Voici une blague..."`
}

type OllamaStreamResp struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
}

// ChatAI : envoie le message à Ollama en local via Docker avec contexte
func (ctrl *Controller) ChatAI(c *gin.Context) {
	var msg AIMessage
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println("[ERROR] JSON binding failed:", err)
		return
	}

	db := ctrl.DB
	fmt.Println("[INFO] Nouveau message reçu de:", msg.User, "Texte:", msg.Text)
	if db == nil {
		fmt.Println("[ERROR] ctrl.DB est nil !")
		return
	}

	// 1️⃣ Récupérer ou créer l'utilisateur
	var user User
	err := db.Where("username = ?", msg.User).First(&user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Println("[INFO] Utilisateur inconnu, création de l'utilisateur:", msg.User)
			user = User{
				Username: msg.User,
				Password: "",
			}
			if err := db.Create(&user).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de créer l'utilisateur"})
				fmt.Println("[ERROR] Impossible de créer l'utilisateur:", err)
				return
			}
			fmt.Println("[INFO] Utilisateur créé avec succès, UserID:", user.ID)
		}
	} else {
		fmt.Println("[INFO] Utilisateur existant trouvé:", user.Username)
	}

	// 2️⃣ Récupérer l'historique
	var history []ConversationHistory
	db.Where("user_id = ?", user.ID).Order("created_at asc").Find(&history)
	fmt.Println("[INFO] Nombre de messages historiques récupérés:", len(history))

	// 3️⃣ Construire le prompt
	var fullPrompt strings.Builder
	for _, h := range history {
		fullPrompt.WriteString(fmt.Sprintf("%s: %s\n", h.Sender, h.Message))
	}
	fullPrompt.WriteString(fmt.Sprintf("%s: %s\n", msg.User, msg.Text))
	fullPrompt.WriteString("Bot:")

	fmt.Println("[INFO] Prompt construit, envoi au modèle IA...")

	// 4️⃣ Appel au modèle IA
	iaHost := "ia"
	iaPort := 11434
	iaURL := fmt.Sprintf("http://%s:%d/api/generate", iaHost, iaPort)

	payload := map[string]interface{}{
		"model":       "mistral",
		"prompt":      fullPrompt.String(),
		"temperature": 0.7,
		"max_tokens":  300,
	}

	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", iaURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println("[ERROR] Erreur lors de l'appel IA:", err)
		return
	}
	defer resp.Body.Close()

	// 5️⃣ Lecture du flux IA
	scanner := bufio.NewScanner(resp.Body)
	var finalResp strings.Builder
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		var chunk OllamaStreamResp
		if err := json.Unmarshal([]byte(line), &chunk); err != nil {
			fmt.Println("[WARN] Impossible de parser la ligne:", line)
			continue
		}
		finalResp.WriteString(chunk.Response)
	}

	if err := scanner.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la lecture du flux IA"})
		fmt.Println("[ERROR] Scanner IA:", err)
		return
	}

	botResponse := finalResp.String()
	fmt.Println("[INFO] Réponse IA générée:", botResponse)

	// 6️⃣ Sauvegarder les messages
	if err := db.Create(&ConversationHistory{
		UserID:  user.ID,
		Sender:  msg.User,
		Message: msg.Text,
	}).Error; err != nil {
		fmt.Println("[ERROR] Impossible de sauvegarder message utilisateur:", err)
	}

	if err := db.Create(&ConversationHistory{
		UserID:  user.ID,
		Sender:  "bot",
		Message: botResponse,
	}).Error; err != nil {
		fmt.Println("[ERROR] Impossible de sauvegarder message bot:", err)
	}

	c.JSON(http.StatusOK, AIResponse{
		Bot: botResponse,
	})
	fmt.Println("[INFO] Conversation sauvegardée avec succès pour l'utilisateur:", user.Username)
}
