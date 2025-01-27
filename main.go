package main

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Reenviar la solicitud PUT con parámetros dinámicos
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

	// Ruta para reenviar solicitudes PUT con parámetros dinámicos
	r.PUT("/edit/:id", func(c *gin.Context) {
		forwardPutRequest(c, "http://localhost:8082/edit")
	})

	r.Run(":9090") // Ejecutar el servidor en el puerto 9090
}
