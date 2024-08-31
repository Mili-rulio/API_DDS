package main

//descargar go mux: framework para crear servidores web
//golang air (para que se actualice sin que bajemos la terminal)
//contenedor de docker para conectarnos con la base de datos
//go gorm es una biblioteca para conectarnos con la base de datos

// @title API de Puntos de Donación
// @version 1.0
// @description Esta es una API para gestionar puntos de donación.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
import (
	"log"
	"net/http"

	"github.com/Mili-rulio/API/db"
	_ "github.com/Mili-rulio/API/docs" // Importa los documentos generados por swag
	"github.com/Mili-rulio/API/models"
	"github.com/Mili-rulio/API/routes"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Puntos de Donación API
// @version 2.0
// @description Esta es una API para gestionar puntos de donación y obtener recomendaciones en base a los mismos.
// @host localhost:8080
// @BasePath /api

func main() {
	// Inicializar la base de datos
	db.DBConnection()
    // Migrar la base de datos
    err := db.DB.AutoMigrate(models.PuntoDeDonacion{})
    if err != nil {
        log.Fatalf("Error en la migración de la base de datos: %v", err)
    } else {
        log.Println("Migración de la base de datos exitosa")
    }

    // Crear un router
    r := mux.NewRouter()
    r.HandleFunc("/", routes.HomeHandler)

    r.HandleFunc("/api/puntosDB", routes.GetPuntosHandler).Methods("GET")
    r.HandleFunc("/api/puntos", routes.GetPuntoHandler).Methods("GET")
    r.HandleFunc("/api/puntos", routes.PostPuntoHandler).Methods("POST")
    r.HandleFunc("/api/puntos/{id}", routes.DeletePuntoHandler).Methods("DELETE")

    // Añadir la ruta para Swagger
    r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

    // Iniciar el servidor HTTP
    log.Println("Servidor iniciado en el puerto 8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}