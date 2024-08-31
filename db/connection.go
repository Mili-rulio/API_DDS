package db

// "\l": me muestra las bases de datos
// "\c": me conecta a una base de datos
//password 12345
//docker exec -it some-postgres bash
//psql -U postgres mrulio --password
// some-postgres-new
//createuser --superuser postgres
//psql -U postgres
//1 docker exec -it some-postgres-new bash
//2 psql -U murlio --password
//\c gorm
//docker run --name some-postgres-new -e POSTGRES_USER=mrulio -e POSTGRES_PASSWORD=12345 -d postgres
//psql -U mrulio -d gorm
import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)
var DSN = "host=localhost user=mrulio password=12345 dbname=gorm port=5432"
var DB *gorm.DB

func DBConnection(){
	var error error
	DB, error = gorm.Open(postgres.Open(DSN), &gorm.Config{}) //conectar a la base de datos
	if error != nil {
		log.Fatal("Error en la conexión a la base de datos")
	} else {
		log.Println("Conexión exitosa a la base de datos")
	}
}
