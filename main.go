package main

import (
	"bytes"
	"net/http"
	"github.com/gin-gonic/gin"
	"io"
)

func forwardRequest(c *gin.Context, serviceURL string) {
	var jsonData map[string]interface{}
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	resp, err := http.Post(serviceURL, "application/json", c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward request"})
		return
	}
	defer resp.Body.Close()

	c.JSON(resp.StatusCode, gin.H{"message": "Request forwarded"})
}
func forwardPutRequest(c *gin.Context, serviceURL string) {
	// Obtener el ID dinámico de la URL
	id := c.Param("id")
	serviceURL = serviceURL + "/" + id

	// Leer el cuerpo de la solicitud entrante
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
		return
	}

	// Crear una nueva solicitud PUT para reenviar
	req, err := http.NewRequest(http.MethodPut, serviceURL, bytes.NewReader(body))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Copiar los encabezados originales
	req.Header = c.Request.Header

	// Hacer la solicitud al servicio remoto
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward request"})
		return
	}
	defer resp.Body.Close()

	// Leer la respuesta del servicio remoto
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	// Responder con el código de estado y el cuerpo del servicio remoto
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
}



func main() {
	r := gin.Default()

	r.POST("/create", func(c *gin.Context) {
		forwardRequest(c, "http://localhost:8081/create")
	})

	r.PUT("/edit/:id", func(c *gin.Context) {
		forwardPutRequest(c, "http://localhost:8082/edit")
	})

	r.GET("/patient/:id", func(c *gin.Context) {
		patientID := c.Param("id")
		serviceURL := "http://localhost:8083/patient/" + patientID
	
		resp, err := http.Get(serviceURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve data"})
			return
		}
		defer resp.Body.Close()
	
		// Leer el cuerpo de la respuesta
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
			return
		}
	
		c.Data(resp.StatusCode, "application/json", body)
	})
	

	r.Run(":9090") 
}
