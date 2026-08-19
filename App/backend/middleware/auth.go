package middleware

import (
	"net/http"

	"realico-comidas-backend/repository"

	"github.com/gin-gonic/gin"
)

// BasicAuthMiddleware valida las credenciales del administrador contra la tabla 'admins' en la base de datos.
// Retorna HTTP 401 si no se envían credenciales o si el usuario/contraseña no existen en la base.
func BasicAuthMiddleware(adminRepo repository.AdminRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, password, hasAuth := c.Request.BasicAuth()

		if !hasAuth {
			c.Header("WWW-Authenticate", `Basic realm="Restricted Admin Area"`)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Se requieren credenciales de administrador"})
			c.Abort()
			return
		}

		// Validar que el usuario exista en la base de datos con esa contraseña
		admin, err := adminRepo.FindByCredentials(user, password)
		if err != nil || admin == nil {
			c.Header("WWW-Authenticate", `Basic realm="Restricted Admin Area"`)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales de administrador inválidas"})
			c.Abort()
			return
		}

		c.Next()
	}
}
