CREATE TABLE IF NOT EXISTS areas (
    id SERIAL PRIMARY KEY,
    area_name VARCHAR(128) NOT NULL,
    top_right_lat FLOAT NOT NULL,
    top_right_lon FLOAT NOT NULL,
    bottom_left_lat FLOAT NOT NULL,
    bottom_left_lon FLOAT NOT NULL,
    deforested_area FLOAT DEFAULT 0.0
);

CREATE TABLE IF NOT EXISTS histories (
    id SERIAL PRIMARY KEY,
    date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    image_path VARCHAR(256) NOT NULL,
    masked_image_path VARCHAR(256) NOT NULL,
    deforested_area FLOAT NOT NULL,
    area_id INTEGER NOT NULL REFERENCES areas(id)
);
