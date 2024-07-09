CREATE TABLE IF NOT EXISTS core.breeds (
    id BIGINT NOT NULL AUTO_INCREMENT,
    species enum('dog', 'cat') NOT NULL,
    pet_size enum('small', 'medium', 'tall') NOT NULL,
    name VARCHAR(255) UNIQUE NOT NULL,
    average_male_adult_weight INT DEFAULT 0,
    average_female_adult_weight INT DEFAULT 0,
    
    PRIMARY KEY (id)
);