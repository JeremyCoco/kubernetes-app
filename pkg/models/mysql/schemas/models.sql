CREATE TABLE IF NOT EXISTS user
(
    id SMALLINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(80) UNIQUE,
    hashed_password VARCHAR(60) NOT NULL
);

CREATE TABLE IF NOT EXISTS todo
(
    id SMALLINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(255) NOT NULL,
    done BOOlEAN NOT NULL DEFAULT FALSE,
    created_at DATETIME NOT NULL,
    expires DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS user_todo
(
    user_id SMALLINT UNSIGNED,
    todo_id SMALLINT UNSIGNED,
    CONSTRAINT user_id_fk FOREIGN KEY (user_id) REFERENCES user(id),
    CONSTRAINT todo_id_fk FOREIGN KEY (todo_id) REFERENCES todo(id)
);