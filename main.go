package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Product struct {
	ID              int     `json:"id_producto"`
	Nombre          string  `json:"nombre"`
	PrecioPacaDolar float64 `json:"precio_paca_dolar"`
	CantidadPorPaca int     `json:"cantidad_por_paca"`
	PorcentajeDolar float64 `json:"porcent_dolar"`
	PorcentajeEfect float64 `json:"porcent_efect"`
	PorcentajePunto float64 `json:"porcent_punto"`
	Cantidad        int     `json:"cantidad"`
}

func main() {
	db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/database_name")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		products, err := getProductsFromDB(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	})

	fmt.Println("Server is listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getProductsFromDB(db *sql.DB) ([]Product, error) {
	rows, err := db.Query("SELECT id_producto, nombre, precio_paca_dolar, cantidad_por_paca, porcent_dolar, porcent_efect, porcent_punto, cantidad FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.Nombre, &p.PrecioPacaDolar, &p.CantidadPorPaca, &p.PorcentajeDolar, &p.PorcentajeEfect, &p.PorcentajePunto, &p.Cantidad)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
