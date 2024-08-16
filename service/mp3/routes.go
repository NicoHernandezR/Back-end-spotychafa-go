package mp3

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/NicoHernandezR/Back-end-spotychafa-go/types"
	"github.com/NicoHernandezR/Back-end-spotychafa-go/utils"
)

type Handler struct {
	store types.Mp3Store
}

func NewHandler(store types.Mp3Store) *Handler {
	return &Handler{store: store}
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
	// Por ejemplo, podrías generar un nombre único para el archivo y guardarlo
	filePath := fmt.Sprintf("uploads/%d.mp3", payload.ID)
	err = os.WriteFile(filePath, fileBytes, 0644)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error saving file: %v", err))
		return
	}

	// Actualizar el campo Mp3File en la estructura payload con la ruta del archivo guardado
	payload.Mp3File = filePath

	// Insertar en la base de datos
	err = h.store.InsertMp3(payload)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)

}
