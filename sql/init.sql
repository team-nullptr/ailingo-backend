CREATE TABLE study_set
(
    `id`                  INT AUTO_INCREMENT      NOT NULL,
    `author_id`           VARCHAR(32)             NOT NULL,
    `name`                VARCHAR(256)            NOT NULL,
    `description`         VARCHAR(512)            NOT NULL,
    `phrase_language`     ENUM ('pl-PL', 'en-US') NOT NULL,
    `definition_language` ENUM ('pl-PL', 'en-US') NOT NULL,

    INDEX (`author_id`(20)),
    PRIMARY KEY (`id`)
);

CREATE TABLE definition
(
    `id`           INT AUTO_INCREMENT NOT NULL,
    `study_set_id` INT                NOT NULL,
    `phrase`       VARCHAR(256)       NOT NULL,
    `meaning`      VARCHAR(256)       NOT NULL,
    `sentences`    JSON               NOT NULL,

    PRIMARY KEY (`id`),
    FOREIGN KEY (`study_set_id`) REFERENCES study_set (`id`) ON DELETE CASCADE
);

CREATE TABLE star
(
    `user_id`      VARCHAR(32) NOT NULL,
    `study_set_id` INT         NOT NULL,

    INDEX (`user_id`(20)),
    UNIQUE (`user_id`, `study_set_id`),
    FOREIGN KEY (`study_set_id`) REFERENCES study_set (`id`) ON DELETE CASCADE
);

CREATE TABLE study_session
(
    `user_id`         VARCHAR(32) NOT NULL,
    `study_set_id`    INT         NOT NULL,
    `last_session_at` DATETIME DEFAULT (NULL),

    INDEX (`user_id`(20)),
    FOREIGN KEY (`study_set_id`) REFERENCES study_set (`id`) ON DELETE CASCADE
);

CREATE TABLE user
(
    `id`        VARCHAR(32) UNIQUE NOT NULL,
    `username`  TEXT               NOT NULL,
    `image_url` TEXT               NOT NULL
);