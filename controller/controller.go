package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"go_skill_test/dao"
	"go_skill_test/models"
	"go_skill_test/utilities"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// Initialize get env vars, get db info from config and get db connection
func Initialize() models.Config {
	//get authorization key and db config
	serviceConfig, err := utilities.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	utilities.DBConfig = serviceConfig.DBConfig

	//connect to database
	utilities.DBSession, err = utilities.GetConnection(utilities.DBConfig)
	if err != nil {
		log.Fatal(err)
	}

	//create uuid generator
	generateUUIDErr := utilities.CreateGenerator()
	if generateUUIDErr != nil {
		log.Fatal(generateUUIDErr)
	}

	return serviceConfig
}

func Router() *mux.Router {

	router := mux.NewRouter()
	router.HandleFunc("/", indexHandler).Methods("GET")
	router.HandleFunc("/upload", uploadImage).Methods("POST")
	return router
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("index.html"))
	err := tpl.ExecuteTemplate(w, "index.html", utilities.AuthKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func uploadImage(w http.ResponseWriter, r *http.Request) {
	var imageInfo models.ImageInfo
	r.Body = http.MaxBytesReader(w, r.Body, 8*1024*1024) // 8 Mb

	err := r.ParseMultipartForm(8 * 1024 * 1024)
	if err != nil {
		fmt.Println("Error in ParseMultipartForm() :", err)
		errorHandler(w, http.StatusBadRequest, err)
		return
	}

	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		fmt.Println("Error in FormFile() :", err)
		errorHandler(w, http.StatusBadRequest, err)
		return
	}

	var buff bytes.Buffer
	fileSize, err := buff.ReadFrom(file)
	imageInfo.ImageSize = fileSize


	defer file.Close()
	fileType := handler.Header["Content-Type"]
	fmt.Println("input file type :", fileType)

	if err := r.ParseForm(); err != nil {
		fmt.Println("Error in ParseForm() :", err)
		errorHandler(w, http.StatusBadRequest, err)
		return
	}

	authKey := r.Form.Get("auth")

	if strings.HasPrefix(fileType[0], "image") && authKey == utilities.AuthKey {
		fmt.Println("Correct file type and successful authentication.")
		imageInfo.FileContentType = fileType[0]
	} else {
		fmt.Println("Incorrect file type or failed authentication.")
		err = errors.New("incorrect file type or failed authentication")
		//return 403 HTTP error code
		errorHandler(w, http.StatusForbidden, err)
		return
	}

	imageInfo.ID = utilities.CreateNewUUID()
	imageInfo.FileName = handler.Filename
	imageInfo.ContentType = r.Header["Content-Type"][0]
	imageInfo.AcceptEncoding = r.Header["Accept-Encoding"][0]
	imageInfo.AcceptLanguage = r.Header["Accept-Language"][0]
	imageInfo.CreatedAt = time.Now()
	imageInfo.CreatedBy = "thiri"
	imageInfoJson, _ := json.Marshal(imageInfo)
	fmt.Println("input image info :", string(imageInfoJson))

	//write the received file data to a temporary file
	tempFile, err := os.Create("./uploaded_images/" + imageInfo.ID + "_" + imageInfo.FileName)
	if err != nil {
		fmt.Println("Error in temp file creation :", err)
		errorHandler(w, http.StatusInternalServerError, err)
		return
	}
	defer tempFile.Close()

	//copy the uploaded file to the temp file
	if _, err := io.Copy(tempFile, file); err != nil {
		fmt.Println("Error in copying to the temp file  :", err)
		errorHandler(w, http.StatusInternalServerError, err)
		return
	}

	//save image meta data into db
	err = dao.SaveImageInfo(imageInfo)
	if err != nil {
		errorHandler(w, http.StatusInternalServerError, err)
		return
	}
	fmt.Fprintf(w, "Saved image successfully!")
}

func errorHandler(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	fmt.Fprintf(w, "Error: %v", err)
}
