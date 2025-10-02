package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AIMessage struct {
	User string `json:"user" example:"Alice"`
	Text string `json:"text" example:"Bonjour, raconte-moi une blague."`
}

type AIResponse struct {
	Bot string `json:"bot" example:"Voici une blague..."`
}

// ChatAI : envoie le message à Ollama en local via Docker
// @Summary      Chat avec modèle IA local
// @Description  Envoie un message au modèle IA exécuté dans Docker (Ollama)
// @Tags         Chatbot
// @Accept       json
// @Produce      json
// @Param        message  body      AIMessage  true  "Message de l'utilisateur"
// @Success      200      {object}  AIResponse
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /chat-ai [post]
func (ctrl *Controller) ChatAI(c *gin.Context) {
	var msg AIMessage
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	iaHost := "ia"
	iaPort := 11434
	iaURL := fmt.Sprintf("http://%s:%s/api/generate", iaHost, iaPort)

	payload := map[string]interface{}{
		"model":       "mistral",
		"prompt":      msg.Text,
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
		return
	}
	defer resp.Body.Close()

	// Lire la réponse brute
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de lire la réponse de l'IA"})
		return
	}

	// Afficher dans les logs pour debug
	fmt.Println("Réponse brute de l'IA :", string(bodyBytes))

	// Renvoyer la réponse brute à l'utilisateur
	c.Data(http.StatusOK, "application/json", bodyBytes)
}
