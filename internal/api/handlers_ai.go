package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ayush-bothra/backend-optimizer/internal/ai"
	"github.com/gin-gonic/gin"
)

type AIrequest struct {
	Messages []ai.Message `json:"messages"`
}

func (h *Handler) AIquery(c *gin.Context)  {
	var messages AIrequest
	// get the message from the user
	if err := c.ShouldBindJSON(&messages); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		return 
	}

	// prepare to send msg to groq/ollama
	req := ai.NewRequest(messages.Messages)
	jsonBody, err := json.Marshal(req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error":err.Error()})
		return 
	}

	// send the query to groq via HTTP post
	postreq, err := ai.NewQuery(jsonBody)
	if err != nil {
		c.JSON((http.StatusBadRequest), gin.H{"error":err.Error()})
		return 
	}
	res, err := http.DefaultClient.Do(postreq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		return
	}

	// unmarshal to get the required data from the reply
	/* res.Body is a io.ReadCloser, meaning it will close once it reaches error or EOF
	this read amount is then sent to the ReadAll to return the []byte format for it
	the correct run of the function will return err == nil
	*/
	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		c.JSON(res.StatusCode, gin.H{"error":string(body)})
		return 
	}

	defer res.Body.Close()
	resbytes, err := io.ReadAll(res.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		return 
	}

	// display the data as a JSON answer
	var result ai.RespBody
	if err := json.Unmarshal(resbytes, &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		return 
	} 

	if len(result.Choices) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error":"no choices returned"})
		return 
	}
	c.JSON(http.StatusOK, gin.H{"reply":result.Choices[len(result.Choices)-1].Message.Content})
}