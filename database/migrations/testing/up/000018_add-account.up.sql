CREATE TABLE IF NOT EXISTS user_account(
    professor_id INT NOT NULL,
    login VARCHAR(100) NOT NULL,
    hash VARCHAR(200) NOT NULL,
    PRIMARY KEY(professor_id),
    FOREIGN KEY (professor_id) REFERENCES professor(id)ON DELETE CASCADE ON UPDATE CASCADE
);