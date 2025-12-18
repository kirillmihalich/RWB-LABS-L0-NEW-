package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq" // Импортируем драйвер PostgreSQL
)

// Структура для данных города
type City struct {
	ID               int     `json:"id"`
	Name             string  `json:"name"`
	RegionName       string  `json:"region_name"`
	DistanceToMoscow int     `json:"distance_to_moscow"`
	PopulationMillion float64 `json:"population_million"`
}

var db *sql.DB

func main() {
	var err error
	// Подключение к базе данных
	db, err = sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=cities_db sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Проверка подключения
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Подключение к базе данных успешно.")

	// Создаем роутер
	r := mux.NewRouter()

	// Главная страница с формой поиска
	r.HandleFunc("/", homeHandler).Methods("GET")

	// API эндпоинт для поиска (возвращает JSON)
	r.HandleFunc("/api/search", searchAPIHandler).Methods("GET")

	// Запуск сервера
	fmt.Println("Сервер запущен на :8080")
	http.ListenAndServe(":8080", r)
}

// Обработчик главной страницы
func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Ошибка шаблона", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// Обработчик API для поиска
func searchAPIHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q") // Получаем параметр поиска

	// Если запрос пустой, возвращаем все города
	if query == "" {
		cities, err := getAllCities()
		if err != nil {
			http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
			return
		}
		sendJSONResponse(w, cities)
		return
	}

	// Иначе ищем по части названия
	cities, err := searchCitiesByName(query)
	if err != nil {
		http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, cities)
}

// Функция для получения всех городов
func getAllCities() ([]City, error) {
	rows, err := db.Query(`
        SELECT c.id, c.name, r.name AS region_name, c.distance_to_moscow, c.population_million
        FROM cities c
        JOIN regions r ON c.region_id = r.id
        ORDER BY c.name ASC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cities []City
	for rows.Next() {
		var city City
		err := rows.Scan(&city.ID, &city.Name, &city.RegionName, &city.DistanceToMoscow, &city.PopulationMillion)
		if err != nil {
			return nil, err
		}
		cities = append(cities, city)
	}
	return cities, nil
}

// Функция для поиска городов по части названия
func searchCitiesByName(namePart string) ([]City, error) {
	rows, err := db.Query(`
        SELECT c.id, c.name, r.name AS region_name, c.distance_to_moscow, c.population_million
        FROM cities c
        JOIN regions r ON c.region_id = r.id
        WHERE c.name ILIKE $1
        ORDER BY c.name ASC
    `, "%"+namePart+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cities []City
	for rows.Next() {
		var city City
		err := rows.Scan(&city.ID, &city.Name, &city.RegionName, &city.DistanceToMoscow, &city.PopulationMillion)
		if err != nil {
			return nil, err
		}
		cities = append(cities, city)
	}
	return cities, nil
}

// Вспомогательная функция для отправки JSON
func sendJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}