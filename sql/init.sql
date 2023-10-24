CREATE TABLE study_sets (
    `id` INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
    `author_id` TEXT NOT NULL,
    `name` VARCHAR(256) NOT NULL,
    `description` VARCHAR(512) NOT NULL,
    `phrase_language`  ENUM('PL', 'EN') NOT NULL,
    `definition_language` ENUM ('PL', 'EN') NOT NULL,
    `definitions` JSON DEFAULT(JSON_ARRAY()),
    KEY `author_id` (`author_id`)
);

CREATE TABLE user_study_sets (
    `user_id` TEXT NOT NULL,
    `study_set_id` INT NOT NULL
);