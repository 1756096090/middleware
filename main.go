package main

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Forward POST requests
func forwardRequest(c *gin.Context, serviceURL string) {
	// Leer el cuerpo de la solicitud y almacenarlo en un buffer
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
		return
	}
	// Restaurar el cuerpo para permitir múltiples lecturas
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	// Crear la solicitud POST
	req, err := http.NewRequest(http.MethodPost, serviceURL, bytes.NewReader(body))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Copiar los headers originales
	req.Header = c.Request.Header
	req.Header.Set("Content-Type", "application/json")

	// Enviar la solicitud
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward request"})
		return
	}
	defer resp.Body.Close()

	// Leer el cuerpo de la respuesta
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	// Copiar los headers de la respuesta
	for k, v := range resp.Header {
		c.Writer.Header()[k] = v
	}

	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
}

// Forward PUT requests
func forwardPutRequest(c *gin.Context, serviceURL string) {
	// Obtener el ID de la ruta
	id := c.Param("id")
	serviceURL = serviceURL + "/" + id

	// Leer el cuerpo de la solicitud y almacenarlo en un buffer
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
		return
	}
	// Restaurar el cuerpo para permitir múltiples lecturas
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	// Crear la solicitud PUT
	req, err := http.NewRequest(http.MethodPut, serviceURL, bytes.NewReader(body))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Copiar los headers originales
	req.Header = c.Request.Header
	req.Header.Set("Content-Type", "application/json")

	// Enviar la solicitud
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward request"})
		return
	}
	defer resp.Body.Close()

	// Leer el cuerpo de la respuesta
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	// Copiar los headers de la respuesta
	for k, v := range resp.Header {
		c.Writer.Header()[k] = v
	}

	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
}

func main() {
	r := gin.Default()

	// Configurar CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"} // O especifica tus dominios permitidos, ej: []string{"http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true

	// Usar el middleware CORS
	r.Use(cors.New(config))

	// Rutas y handlers
	r.POST("/create", func(c *gin.Context) {
		forwardRequest(c, "http://localhost:8081/create")
	})

	r.PUT("/edit/:id", func(c *gin.Context) {
		forwardPutRequest(c, "http://localhost:8082/edit")
	})

	r.GET("/patient/:id", func(c *gin.Context) {
		patientID := c.Param("id")
		serviceURL := "http://localhost:8083/patient/" + patientID

		// Crear la solicitud GET
		req, err := http.NewRequest(http.MethodGet, serviceURL, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
			return
		}

		// Copiar los headers originales
		req.Header = c.Request.Header

		// Enviar la solicitud
		client := &http.Client{}
		resp, err := client.Do(req)
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

		// Copiar los headers de la respuesta
		for k, v := range resp.Header {
			c.Writer.Header()[k] = v
		}

		c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
	})

	// Ejecutar el servidor en el puerto 9090
	r.Run(":9090")
}