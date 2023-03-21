package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/lib/pq"

	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	host     = "db"
	port     = 5432
	user     = "divrhinotrivia"
	password = "divrhinotrivia"
	dbname   = "divrhinotrivia"
)

type Product struct {
	SKU         string  `json:"sku"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	DisplayCase string  `json:"display_case"`
	ImageID     int     `json:"image_id"`
	Weight      float64 `json:"weight"`
	Price       float64 `json:"price"`
	ReviewID    int     `json:"review_id"`
}

var router *chi.Mux
var db *sql.DB

func routers() *chi.Mux {
	router.Get("/products", getProducts)
	router.Get("/products/{id}", getProductByID)
	router.Post("/products", createProduct)
	router.Put("/products/{id}", updateProduct)
	router.Delete("/products/{id}", deleteProduct)

	return router
}

func initialize() {
	router = chi.NewRouter()
	router.Use(middleware.Recoverer)

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)
	var err error
	db, err = sql.Open("postgres", psqlInfo)

	catch(err)
}

func main() {
	initialize()
	routers()
	http.ListenAndServe(":3000", Logger())
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM product")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	products := make([]Product, 0)
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.SKU, &product.Title, &product.Description, &product.Category, &product.DisplayCase, &product.ImageID, &product.Weight, &product.Price, &product.ReviewID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert the products slice to JSON
	data, err := json.Marshal(products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func getProductByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	query := fmt.Sprintf("SELECT * FROM product WHERE sku = '%s'", id)

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(query)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	products := make([]Product, 0)
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.SKU, &product.Title, &product.Description, &product.Category, &product.DisplayCase, &product.ImageID, &product.Weight, &product.Price, &product.ReviewID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert the products slice to JSON
	data, err := json.Marshal(products[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	var p Product
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query, err := db.Prepare(
		`
		Insert INTO product (SKU, Title, Description, Category, Display_case, Image_ID, Weight, Price, Review_ID)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`,
	)
	catch(err)

	_, er := query.Exec(p.SKU, p.Title, p.Description, p.Category, p.DisplayCase, p.ImageID, p.Weight, p.Price, p.ReviewID)
	catch(er)
	defer query.Close()

	respondwithJSON(w, http.StatusCreated, map[string]string{"message": "successfully created"})
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	var p Product
	id := chi.URLParam(r, "id")
	json.NewDecoder(r.Body).Decode(&p)

	query, err := db.Prepare(
		`Update product set title=$1, Description=$2 , Category=$3 , Display_Case=$4 , Image_ID=$5 , Weight=$6 , Price=$7 , Review_ID=$8
		where sku=$9`,
	)
	catch(err)
	_, er := query.Exec(p.Title, p.Description, p.Category, p.DisplayCase, p.ImageID, p.Weight, p.Price, p.ReviewID, id)
	catch(er)

	defer query.Close()

	respondwithJSON(w, http.StatusOK, map[string]string{"message": "updated successfully"})
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	query, err := db.Prepare("delete from product where sku=$1")
	catch(err)
	_, er := query.Exec(id)
	catch(er)
	query.Close()

	respondwithJSON(w, http.StatusOK, map[string]string{"message": "deleted successfully"})
}
