package utils

import (
	"crypto/sha256"
	"fmt"
)

// HashSHA256 genera un hash SHA-256 de la cadena de entrada
// y lo retorna como una cadena hexadecimal de 64 caracteres.
// Se usa para almacenar y verificar contraseñas de administrador sin guardarlas en texto plano.
func HashSHA256(input string) string {
	hash := sha256.Sum256([]byte(input))
	return fmt.Sprintf("%x", hash)
}
