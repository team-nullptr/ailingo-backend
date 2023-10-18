-- table representing a single study set
CREATE TABLE study_sets (
    id INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
    author_id INT NOT NULL,
    name VARCHAR(256) NOT NULL,
    description VARCHAR(512) NOT NULL,
    phrase_language  ENUM('PL', 'EN-GB') NOT NULL,
    meaning_language ENUM ('PL', 'EN-GB') NOT NULL,
    definitions JSON DEFAULT JSON_ARRAY()
);

-- user_study_sets represents a many-to-many relationship between users and study sets.
CREATE TABLE user_study_sets (
    user_id INT NOT NULL,
    study_set_id INT NOT NULL
    -- TODO: Some information about learning progress?
);