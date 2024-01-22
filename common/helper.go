package common

import (
	"math"
	"math/rand"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func GetOffset(currentPage int, listPerPage int) int {
	return (currentPage - 1) * listPerPage
}

func EmptyOrRows(rows []interface{}) []interface{} {
	if rows == nil {
		return []interface{}{}
	}
	return rows
}

func PageArray(totalSize int, pageSize int, page int, maxLength int) []interface{} {
	currentPage := page
	currentPosition := maxLength / 2
	totalPage := int(math.Ceil(float64(totalSize) / float64(pageSize)))

	startPoint := 1
	if currentPage-currentPosition >= 1 {
		startPoint = currentPage - currentPosition
	}
	endPoint := totalPage
	if currentPage+currentPosition <= totalPage {
		endPoint = currentPage + currentPosition
	}

	var pages []interface{}
	if startPoint != 1 {
		pages = append(pages, "...")
	}
	for i := startPoint; i <= endPoint; i++ {
		pages = append(pages, i)
	}
	if endPoint != totalPage {
		pages = append(pages, "...")
	}
	return pages
}

func RandomString(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func IsEmpty(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}

// EncryptPassword encrypts the password using bcrypt
func EncryptPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// Implement the remaining functions...
