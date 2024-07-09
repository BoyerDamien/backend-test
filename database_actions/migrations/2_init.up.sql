CREATE TABLE IF NOT EXISTS core.breeds (
    id BIGINT NOT NULL AUTO_INCREMENT,
    species enum('dog', 'cat'),
    pet_size enum('small', 'medium', 'tall'),
    name VARCHAR(255) NOT NULL,
    average_male_adult_weight INT CHECK (average_male_adult_weight > 0),
    average_female_adult_weight INT CHECK (average_female_adult_weight > 0),
    
    PRIMARY KEY (id)
);