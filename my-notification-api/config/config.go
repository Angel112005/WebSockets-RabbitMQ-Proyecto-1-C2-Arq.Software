package config

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectDatabase() {
	dsn := "root:01toto01@tcp(127.0.0.1:3306)/notification_db_arq_soft?parseTime=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(fmt.Sprintf("❌ Error conectando a la base de datos: %v", err))
	}

	err = db.Ping()
	if err != nil {
		panic(fmt.Sprintf("❌ No se pudo conectar a MySQL: %v", err))
	}

	fmt.Println("✅ Conectado a MySQL correctamente.")
	DB = db
}
