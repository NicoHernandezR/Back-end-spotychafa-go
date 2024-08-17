package mp3

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/NicoHernandezR/Back-end-spotychafa-go/service/s3"
	"github.com/NicoHernandezR/Back-end-spotychafa-go/types"
	"github.com/NicoHernandezR/Back-end-spotychafa-go/utils"
)

type Handler struct {
	store types.Mp3Store
	awsS3 *s3.S3Client
}

func NewHandler(store types.Mp3Store, awsS3 *s3.S3Client) *Handler {
	return &Handler{store: store, awsS3: awsS3}
}

func (h *Handler) RegisterRouter(router *http.ServeMux) {
	router.HandleFunc("GET /mp3/{mp3ID}", h.handlerGetMp3)
	router.HandleFunc("POST /mp3", h.handlerInsertMp3)
}

func (h *Handler) handlerGetMp3(w http.ResponseWriter, r *http.Request) {
	// Convierte el valor del ID obtenido del path a int64 primero
	mp3ID64, err := strconv.ParseInt(r.PathValue("mp3ID"), 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("MP3 not found"))
		return
	}

	// Convierte int64 a int
	mp3ID := int(mp3ID64)

	mp3, err := h.store.GetMp3ByID(mp3ID)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("MP3 not found"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, mp3)

}
func (h *Handler) handlerInsertMp3(w http.ResponseWriter, r *http.Request) {
	// Obtener el archivo y el encabezado asociado
	file, _, err := r.FormFile("file")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error retrieving the file: %v", err))
		return
	}
	defer file.Close()

	// Leer el JSON con la información de la canción
	songData := r.FormValue("song")
	var payload types.Mp3
	if err := json.Unmarshal([]byte(songData), &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error parsing json: %v", err))
		return
	}

	dir := "uploads"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error creating directory: %v", err))
			return
		}
	}

	// Leer los datos del archivo en memoria
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error reading file: %v", err))
		return
	}

	// Aquí puedes procesar y guardar el archivo en algún directorio o en la base de datos

	fileName := fmt.Sprintf("%s.mp3", payload.ID)

	h.awsS3.Upload(fileBytes, "spotychafa", fileName)

	// Actualizar el campo Mp3File en la estructura payload con la ruta del archivo guardado
	payload.Mp3File = fileName

	// Insertar en la base de datos
	err = h.store.InsertMp3(payload)
	if err != nil {
		//TODO Cuando ocurra error al insertar el mp3 en la base de datos, eliminar el mp3 del servidor
		// Para que no se guarde el mp3, siendo que no se guardo en la BD
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)

}
