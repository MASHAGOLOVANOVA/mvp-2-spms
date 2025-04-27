CREATE TABLE IF NOT EXISTS git_repository_integration(
    id INT NOT NULL auto_increment,
    account_id INT NOT NULL,
    api_key VARCHAR(200) NOT NULL,
    type INT NOT NULL DEFAULT 0,
    PRIMARY KEY(id),
    FOREIGN KEY (account_id) REFERENCES professor(id)ON DELETE CASCADE ON UPDATE CASCADE
);