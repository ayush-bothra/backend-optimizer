package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/ayush-bothra/backend-optimizer/internal/ai"
	"github.com/ayush-bothra/backend-optimizer/internal/cache"
	"github.com/gin-gonic/gin"
)

type AIrequest struct {
	Messages []ai.Message `json:"messages"`
}

func (h *Handler) AIquery(c *gin.Context)  {
	var message ai.Message
	// get the message from the user
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		return 
	}

	val, exists := c.Get("jwt")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error":"user claims not found"})
		return 
	}

	claims, ok := val.(*validator.ValidatedClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error":ok})
		return 
	}

	sub := claims.RegisteredClaims.Subject
	fmt.Printf("DEBUG - User subject: %s\n", sub)
	b, mark, err := cache.GetFromRedis(h.rdb, c, sub)
	fmt.Printf("DEBUG - Redis check: mark=%v, err=%v, data_length=%d\n", mark, err, len(b))

	var aireq AIrequest
	if mark == false {
			if err == cache.ErrCacheMiss {
				aireq.Messages = []ai.Message{message}
				updatedBytes, err := json.Marshal(&aireq)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
					return 
				}

				if err := cache.SetToRedis(h.rdb, c, sub, updatedBytes, 1800 * time.Second); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
					return 
				}
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
				return 
			}
	} else {
		if err := json.Unmarshal(b, &aireq); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
			return 
		}
		aireq.Messages = append(aireq.Messages, message)
		updatedBytes, err := json.Marshal(&aireq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
			return 
		}
		if err := cache.SetToRedis(h.rdb, c, sub, updatedBytes, 1800 * time.Second); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
			return 
		}
	}
	fmt.Printf("DEBUG - Number of messages: %d\n", len(aireq.Messages))
	// prepare to send msg to groq/ollama
	req := ai.NewRequest(aireq.Messages)
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
	defer res.Body.Close()
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
	// After getting the AI response
	aiResponse := result.Choices[len(result.Choices)-1].Message
	aireq.Messages = append(aireq.Messages, aiResponse)
	updatedBytes, _ := json.Marshal(&aireq)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"reply": aiResponse.Content,
			"warning": "response saved but conversation history may not persist",
		})
		return
	}
	if err := cache.SetToRedis(h.rdb, c, sub, updatedBytes, 1800 * time.Second); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"reply": aiResponse.Content,
			"warning": "response saved but conversation history may not persist",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"reply":aiResponse.Content})
}