CREATE TABLE study_set
(
    `id`                  INT  AUTO_INCREMENT NOT NULL,
    `author_id`           TEXT                           NOT NULL,
    `name`                VARCHAR(256)                   NOT NULL,
    `description`         VARCHAR(512)                   NOT NULL,
    `phrase_language`     ENUM ('pl-PL', 'en-US')        NOT NULL,
    `definition_language` ENUM ('pl-PL', 'en-US')        NOT NULL,

    PRIMARY KEY (`id`)
);

CREATE TABLE definition (
    `id`           INT AUTO_INCREMENT NOT NULL,
    `study_set_id` INT                            NOT NULL,
    `phrase`       VARCHAR(256)                   NOT NULL,
    `meaning`      VARCHAR(256)                   NOT NULL,
    `sentences`    JSON NOT NULL,

    PRIMARY KEY (`id`),
    FOREIGN KEY (`study_set_id`) REFERENCES study_set(`id`) ON DELETE CASCADE
);


CREATE TABLE study_set_user
(
    `user_id`      TEXT NOT NULL,
    `study_set_id` INT  NOT NULL
);