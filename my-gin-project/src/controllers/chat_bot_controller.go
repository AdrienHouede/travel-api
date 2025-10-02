package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Message struct {
	User string `json:"user"`
	Text string `json:"text"`
}

type Response struct {
	Bot string `json:"bot"`
}

// Chat
// @Summary      Chat avec le bot
// @Description  Envoie un message au bot et reçoit une réponse
// @Tags         Chatbot
// @Accept       json
// @Produce      json
// @Param        message  body      Message  true  "Message de l'utilisateur"
// @Success      200      {object}  Response
// @Failure      400      {object}  map[string]string
// @Router       /chat [post]
func (ctrl *Controller) Chat(c *gin.Context) {
	var msg Message
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	text := strings.ToLower(msg.Text)
	reply := "Je n'ai pas compris 🤔"

	if strings.Contains(text, "bonjour") {
		reply = "Bonjour " + msg.User + " 👋"
	} else if strings.Contains(text, "ça va") {
		reply = "Oui merci, et toi ?"
	} else if strings.Contains(text, "bye") {
		reply = "Au revoir 👋"
	}

	c.JSON(http.StatusOK, Response{Bot: reply})
}
