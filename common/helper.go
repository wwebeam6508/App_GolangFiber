package common

import (
	"PBD_backend_go/configuration"
	"PBD_backend_go/exception"
	"context"
	"encoding/base64"
	"math"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/google/uuid"
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

func PageArray(totalSize int32, pageSize int, page int, maxLength int) []interface{} {
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

func IsEmpty(x interface{}) bool {
	//check if x is array
	if reflect.TypeOf(x).Kind() == reflect.Slice {
		return reflect.ValueOf(x).Len() == 0
	}
	return x == reflect.Zero(reflect.TypeOf(x)).Interface()
}

// EncryptPassword encrypts the password using bcrypt
func EncryptPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func UploadImageToStorage(path string, filename string, image string) (string, error) {
	//check is image is base64
	if !strings.Contains(image, "data:image/") {
		return "", exception.ValidationError{Message: "Invalid image format"}
	}
	data, err := base64.StdEncoding.DecodeString(image[strings.IndexByte(image, ',')+1:])
	if err != nil {
		return "", err
	}
	storage, err := configuration.ConnectToStorage()
	if err != nil {
		return "", err
	}
	storageBucket := os.Getenv("STORAGEBUCKET")
	bucket := storage.Bucket(storageBucket)
	obj := bucket.Object(`` + path + `/` + filename)
	wc := obj.NewWriter(context.Background())
	id := uuid.New()
	wc.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadTokens": id.String()}
	defer wc.Close()

	// Determine the image type
	contentType := http.DetectContentType(data)
	wc.ContentType = contentType

	//get URL after upload
	if _, err := wc.Write(data); err != nil {
		return "", err
	}

	return getPathStorageFromUrl(path, filename, id.String()), nil

}
func DeleteImageFromStorage(path string, filename string) error {
	storage, err := configuration.ConnectToStorage()
	if err != nil {
		return err
	}
	storageBucket := os.Getenv("STORAGEBUCKET")
	bucket := storage.Bucket(storageBucket)
	obj := bucket.Object(`` + path + `/` + filename)
	if err := obj.Delete(context.Background()); err != nil {
		return err
	}
	return nil
}

func getPathStorageFromUrl(path string, filename string, uidd string) string {

	storageBucket := os.Getenv("STORAGEBUCKET")
	baseURL := `https://firebasestorage.googleapis.com/v0/b/` + storageBucket + `/o/` + path + `%2F` + filename + `?alt=media` + `&token=` + uidd

	return baseURL
}

func Contains(arr []int, x int) bool {
	for _, n := range arr {
		if x == n {
			return true
		}
	}
	return false
}

func SortIntDesc(input *[]int) {
	for i := 0; i < len(*input); i++ {
		for j := i + 1; j < len(*input); j++ {
			if (*input)[i] < (*input)[j] {
				(*input)[i], (*input)[j] = (*input)[j], (*input)[i]
			}
		}
	}
}

func FindIndex[T any](slice []T, predicate func(T) bool) int {
	for i, item := range slice {
		if predicate(item) {
			return i
		}
	}
	return -1
}
