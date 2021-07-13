package app

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"image"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	// Нужны для метода DecodeConfig, пакета image.
	// Чтобы доставать размеры картинок (height * width).
	_ "image/jpeg"
	_ "image/png"

	"github.com/SardorMS/project/cmd/app/middleware"
	"github.com/SardorMS/project/pkg/types"
	"github.com/gorilla/mux"
)

// handleUserRegistration - метод регистрации пользователя.
func (s *Server) handleUserRegistration(writer http.ResponseWriter, request *http.Request) {

	country, ok := mux.Vars(request)["region"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	country = strings.ToUpper(country)

	var item *types.Registration

	if err := json.NewDecoder(request.Body).Decode(&item); err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if item.DisplayName == "" || item.Email == "" || item.Password == "" ||
		item.Product == "" || item.Birthdate == "" {
		http.Error(writer, "Body Params is Empty", http.StatusBadRequest)
		return
	}

	result, err := s.userSvc.Register(request.Context(), item, country)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	respondJSON(writer, result)
}

// handleUserLogin - метод для идентификации конкретного пользователя.
func (s *Server) handleUserLogin(writer http.ResponseWriter, request *http.Request) {

	item := &types.Registration{}

	if err := json.NewDecoder(request.Body).Decode(&item); err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if item.AccessToken == "" || item.Email == "" || item.Password == "" {
		http.Error(writer, "Body Params is Empty", http.StatusBadRequest)
		return
	}

	login, err := s.userSvc.Login(request.Context(), item.AccessToken, item.Email, item.Password)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	respondJSON(writer, login)
}

// handleGetCurrentUser - метод для вывода информации для конкретного пользователя.
func (s *Server) handleGetCurrentUser(writer http.ResponseWriter, request *http.Request) {

	id, err := middleware.Authentication(request.Context())
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	result, err := s.userSvc.GetCurrentUser(request.Context(), id)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	respondJSON(writer, result)
}

// handleGetUserByID - метод для вывода информации для конкретного пользователя по его user_id.
func (s *Server) handleGetUserByID(writer http.ResponseWriter, request *http.Request) {

	userID, ok := mux.Vars(request)["user_id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	result, err := s.userSvc.GetUserByID(request.Context(), userID)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	respondJSON(writer, result)
}

// handleUserImageUpload - метод для загрузки изображения для профиля, под определённого пользователя.
func (s *Server) handleUserImageUpload(writer http.ResponseWriter, request *http.Request) {

	id, err := middleware.Authentication(request.Context())
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Максимально допустимый размер фотографии.
	maxUploadSize := int64(2 * 1024 * 1024) // 2 MB
	uploadPath := "../tmp"

	if err := request.ParseMultipartForm(maxUploadSize); err != nil {
		http.Error(writer, "Can not Parse the Form!", http.StatusInternalServerError)
		return
	}

	file, fileHeader, err := request.FormFile("images")
	if err != nil {
		http.Error(writer, "Invalid File!", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Размер файла
	fileSize := fileHeader.Size
	log.Printf("File size (bytes): %v\n", fileSize)

	if fileSize > maxUploadSize {
		http.Error(writer, "File is too big!", http.StatusBadRequest)
		return
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(writer, "Invalid file!", http.StatusBadRequest)
		return
	}

	// Тип файла
	fileType := http.DetectContentType(fileBytes)
	switch fileType {
	case "image/jpeg", "image/jpg":
	case "image/png":
		break
	default:
		http.Error(writer, "Invalid file type!", http.StatusBadRequest)
		return
	}

	buffer := make([]byte, 8)
	n, err := rand.Read(buffer)
	if n != len(buffer) || err != nil {
		http.Error(writer, "Can not generate rand!", http.StatusBadRequest)
		return
	}
	randomFileName := hex.EncodeToString(buffer)
	extantions := filepath.Ext(fileHeader.Filename)
	fileFullName := randomFileName + extantions

	newPath := filepath.Join(uploadPath, fileFullName)
	log.Printf("FileType: %s, File: %s\n", fileType, newPath)

	// Достаем разрешения картинки width * height
	reader := bytes.NewReader(fileBytes)
	img, _, err := image.DecodeConfig(reader)
	if err != nil {
		return
	}
	width := img.Width
	height := img.Height

	url := "https://accounts.sprotify.com/v1/me/image/" + fileFullName

	item := &types.ImageUpdate{
		Images: []types.Image{
			{
				URL:    url,
				Width:  width,
				Height: height,
			},
		},
	}

	result, message, err := s.userSvc.ImageUpload(request.Context(), item, id)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if message == "1" {
		data := &types.Message{
			Status:  "400 Bad Request.",
			Message: "You Already Have An Image. This Operation Can Not Be Done! Please Remove Your Image First. Thank You!",
		}

		respondJSON(writer, data)
	} else {

		newFile, err := os.Create(newPath)
		if err != nil {
			http.Error(writer, "Can not write a file!", http.StatusInternalServerError)
			return
		}
		defer newFile.Close()
		if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
			http.Error(writer, "Can not write a file!", http.StatusInternalServerError)
			return
		}

		respondJSON(writer, result)
	}
}

// handleUserImageRemove - метод для удаления изображения из профиля для конкретного пользователя.
func (s *Server) handleUserImageRemove(writer http.ResponseWriter, request *http.Request) {

	id, err := middleware.Authentication(request.Context())
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	result, err := s.userSvc.ImageRemove(request.Context(), id)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	respondJSON(writer, result)
}
