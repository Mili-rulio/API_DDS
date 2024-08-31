package routes

import (
	"encoding/json"	
	"net/http"
	"strconv"

	"github.com/Mili-rulio/API/db"
	"github.com/Mili-rulio/API/models"
	"github.com/gorilla/mux"
)

// GetPuntosHandler obtiene todos los puntos de donación
// @Summary Obtiene todos los puntos de donación recomendados.
// @Description Obtiene todos los puntos de donación recomendados por la ONG.
// @Tags puntos
// @Produce json
// @Success 200 {array} models.PuntoDeDonacion
// @Router /puntosDB [get]
func GetPuntosHandler(w http.ResponseWriter, r *http.Request) {
    var puntos []models.PuntoDeDonacion
    db.DB.Find(&puntos)
    json.NewEncoder(w).Encode(&puntos)
    w.Write([]byte("Bienvenido a la API de Puntos de Donación"))
}

type Coordenadas struct {
    Latitud  float64 `json:"latitud"`
    Longitud float64 `json:"longitud"`
}

type PuntoDeDonacionResult struct {
    Nombre      string      `json:"nombre"`
    Direccion   string      `json:"direccion"`
    Ciudad      string      `json:"ciudad"`
    Coordenadas string      `json:"coordenadas"`
    Distancia   float64     `json:"distancia"`
}

// GetPuntoHandler obtiene puntos de donación dentro de un radio
// @Summary Obtiene puntos de donación dentro de un radio
// @Description Obtiene puntos de donación recomendados dentro de un radio.
// @Tags puntos
// @Produce json
// @Param coordenadas query string true "Coordenadas en formato JSON"
// @Param radio query string true "Radio en metros"
// @Success 200 {array} PuntoDeDonacionResult
// @Failure 400 {string} string "Parámetros inválidos"
// @Failure 404 {string} string "No se encontraron puntos de donación en el radio especificado"
// @Router /puntos [get]
func GetPuntoHandler(w http.ResponseWriter, r *http.Request) {
    coordenadasStr := r.URL.Query().Get("coordenadas")
    radioStr := r.URL.Query().Get("radio")

    if coordenadasStr == "" || radioStr == "" {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Parámetros inválidos"))
        return
    }

    var coordenadas Coordenadas
    if err := json.Unmarshal([]byte(coordenadasStr), &coordenadas); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Parámetros de coordenadas inválidos"))
        return
    }

    latitud := coordenadas.Latitud
    longitud := coordenadas.Longitud

    radio, err := strconv.ParseFloat(radioStr, 64)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Radio inválido"))
        return
    }

    var puntos []PuntoDeDonacionResult

    query := `
        SELECT nombre, direccion, ciudad, coordenadas, 
        ST_Distance(coordenadas, ST_SetSRID(ST_MakePoint(?, ?), 4326)) AS distancia
        FROM puntos_de_donacion
        WHERE ST_Distance(coordenadas, ST_SetSRID(ST_MakePoint(?, ?), 4326)) < ?;
    `

    if err := db.DB.Raw(query, latitud, longitud, latitud, longitud, radio).Scan(&puntos).Error; err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("Error al ejecutar la consulta"))
        return
    }

    if len(puntos) == 0 {
        w.WriteHeader(http.StatusNotFound)
        w.Write([]byte("No se encontraron puntos de donación en el radio especificado"))
        return
    }

    for i, punto := range puntos {
        var coords Coordenadas
        if err := json.Unmarshal([]byte(punto.Coordenadas), &coords); err == nil {
            coordsJSON, err := json.Marshal(coords)
            if err == nil {
                puntos[i].Coordenadas = string(coordsJSON)
            }
        }
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(&puntos)
}

// PostPuntoHandler crea un nuevo punto de donación
// @Summary Crea un nuevo punto de donación
// @Description Crea un nuevo punto de donación
// @Tags puntos
// @Accept json
// @Produce json
// @Param punto body models.PuntoDeDonacion true "Punto de Donación"
// @Success 201 {object} models.PuntoDeDonacion
// @Failure 400 {string} string "Error al crear el punto de donación"
// @Router /puntos [post]
func PostPuntoHandler(w http.ResponseWriter, r *http.Request) {
    var punto models.PuntoDeDonacion
    json.NewDecoder(r.Body).Decode(&punto)
    crearPunto := db.DB.Create(&punto)
    if crearPunto.Error != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Error al crear el punto de donación"))
        return
    }
    json.NewEncoder(w).Encode(&punto)
    w.Write([]byte("POST /punto"))
}

// DeletePuntoHandler elimina un punto de donación
// @Summary Elimina un punto de donación
// @Description Elimina un punto de donación de la base de datos.
// @Tags puntos
// @Param id path int true "ID del Punto de Donación"
// @Success 204 {string} string "Punto de donación eliminado"
// @Failure 404 {string} string "Punto de donación no encontrado"
// @Router /puntos/{id} [delete]
func DeletePuntoHandler(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    var punto models.PuntoDeDonacion
    db.DB.Delete(&punto, params["id"])
    if punto.ID == 0 {
        w.WriteHeader(http.StatusNotFound)
        w.Write([]byte("Punto de donación no encontrado"))
        return
    }
    db.DB.Delete(&punto)
    w.WriteHeader(http.StatusNoContent)
    w.Write([]byte("DELETE /punto"))
}