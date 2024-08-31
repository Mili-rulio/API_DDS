package models
//gorm orm para GO


// PuntoDeDonacion representa un punto de donación en la base de datos
type PuntoDeDonacion struct {
    ID          uint       `gorm:"primaryKey"`
    Nombre      string     `gorm:"not null; unique_index"`
    Direccion   string     `gorm:"not null; unique_index"`
    Ciudad      string     `gorm:"not null"`
    Coordenadas string     `json:"coordenadas" gorm:"type:geography(POINT,4326);not null"` // Coordenadas en formato WKT
  //  Provincia   string     `gorm:"not null"`
}


func (PuntoDeDonacion) TableName() string {
    return "puntos_de_donacion"
}