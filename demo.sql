-- Подключаемся к базе данных
\c cities_db;

-- Создаем таблицу регионов
CREATE TABLE IF NOT EXISTS regions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

-- Создаем таблицу городов
CREATE TABLE IF NOT EXISTS cities (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    region_id INTEGER NOT NULL REFERENCES regions(id) ON DELETE CASCADE,
    distance_to_moscow INTEGER NOT NULL CHECK (distance_to_moscow >= 0),
    population_million DECIMAL(4,2) NOT NULL CHECK (population_million >= 0)
);

-- Заполняем таблицу регионов
INSERT INTO regions (name) VALUES
('Центральный федеральный округ'),
('Северо-Западный федеральный округ'),
('Южный федеральный округ'),
('Приволжский федеральный округ'),
('Уральский федеральный округ'),
('Сибирский федеральный округ'),
('Дальневосточный федеральный округ');

-- Заполняем таблицу городов (15 строк)
INSERT INTO cities (name, region_id, distance_to_moscow, population_million) VALUES
('Москва', 1, 0, 12.6),
('Санкт-Петербург', 2, 635, 5.4),
('Новосибирск', 6, 2700, 1.6),
('Екатеринбург', 5, 1660, 1.5),
('Казань', 4, 790, 1.2),
('Нижний Новгород', 4, 400, 1.2),
('Челябинск', 5, 1800, 1.2),
('Омск', 6, 2500, 1.1),
('Самара', 4, 800, 1.1),
('Ростов-на-Дону', 3, 1000, 1.1),
('Уфа', 4, 1300, 1.1),
('Красноярск', 6, 3400, 1.0),
('Пермь', 5, 1300, 1.0),
('Воронеж', 3, 500, 1.0),
('Волгоград', 3, 900, 1.0);