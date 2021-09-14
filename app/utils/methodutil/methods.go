package methodutil

import (
	"boilerplate/app/utils/hash"
	"boilerplate/app/utils/msgutil"
	"boilerplate/infra/errors"
	"boilerplate/infra/logger"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func IsInvalid(value string) bool {
	return value == ""
}

func InArray(needle interface{}, haystack interface{}) bool {
	switch reflect.TypeOf(haystack).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(haystack)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(needle, s.Index(i).Interface()) {
				return true
			}
		}
	}

	return false
}

func IsEmpty(x interface{}) bool {
	return x == nil || reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

func MapToStruct(input map[string]interface{}, output interface{}) error {
	if b, err := json.Marshal(input); err == nil {
		return json.Unmarshal(b, &output)
	} else {
		return err
	}
}

func StructToStruct(input interface{}, output interface{}) error {
	if b, err := json.Marshal(input); err == nil {
		return json.Unmarshal(b, &output)
	} else {
		return err
	}
}

func ParseJwtToken(token, secret string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.ErrInvalidJwtSigningMethod
		}
		return []byte(secret), nil
	})
}

func StringToIntArray(stringArray []string) []int {
	var res []int

	for _, v := range stringArray {
		if i, err := strconv.Atoi(v); err == nil {
			res = append(res, i)
		}
	}

	return res
}

func GenerateHashedImageName(fileName string) string {
	ext := filepath.Ext(fileName)
	name := strings.TrimSuffix(fileName, ext)

	return hash.GetSha1Hash(name) + ext
}

func GenerateFilePathForLocalStorage(image *multipart.FileHeader, uploadPath string) (string, *errors.RestErr) {
	// Source
	src, err := image.Open()
	if err != nil {
		logger.Error(msgutil.EntityGenericFailedMsg("opening file"), err)
		return "", errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}
	defer func() {
		_ = src.Close()
	}()

	hashedFileName := GenerateHashedImageName(image.Filename)

	basePath := "uploads/" + uploadPath

	logger.Info("creating file base path if not exists")
	err = os.MkdirAll(basePath, 0755)
	if err != nil {
		logger.Error(msgutil.EntityGenericFailedMsg("creating file base path if not exists"), err)
		return "", errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	// Destination
	dst, err := os.Create(basePath + hashedFileName)
	if err != nil {
		logger.Error(msgutil.EntityGenericFailedMsg("creating file to local folder"), err)
		return "", errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}
	defer func() {
		_ = dst.Close()
	}()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		logger.Error(msgutil.EntityGenericFailedMsg("copying file to destination"), err)
		return "", errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return uploadPath + hashedFileName, nil
}

func DeleteFileFromLocal(filePath string) *errors.RestErr {
	err := os.Remove(filePath)

	if err != nil {
		logger.Error(msgutil.EntityGenericFailedMsg("delete local stored file"), err)
		return errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}
	logger.Info(msgutil.EntityDeleteSuccessMsg(fmt.Sprintf("%s path file deleted", filePath)))
	return nil
}

func ValidateImageFileType(fle *multipart.FileHeader, Type string) *errors.RestErr {
	file, err := fle.Open()
	if err != nil {
		logger.Error("can't open temporary file", err)
	}
	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		logger.Error("can't read temporary file", err)
	}

	fileType := http.DetectContentType(buff)

	switch fileType {
	case "image/jpeg", "image/jpg", "image/gif", "image/png":
		logger.Info(fileType)
		break
	case "application/octet-stream": // for csv or xlsx files
		logger.Info(fileType)
		if Type == "img" {
			return errors.NewBadRequestError("Please provide a image type file to upload")
		}
		break
	case "application/pdf": // not image, but application !
		logger.Info(fileType)
		return errors.NewBadRequestError("pdf file type not accepted")
	default:
		logger.Info("unknown file type uploaded")
		return errors.NewBadRequestError("unknown file type uploaded")
	}

	return nil
}

func GenerateRandomStringOfLength(length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

	if length == 0 {
		length = 8
	}

	var b strings.Builder

	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}

	return b.String()
}
