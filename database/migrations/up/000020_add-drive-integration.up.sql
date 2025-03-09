CREATE TABLE IF NOT EXISTS drive_integration(
    account_id INT NOT NULL,
    base_folder_id VARCHAR(100) NOT NULL,
    api_key VARCHAR(200) NOT NULL,
    FOREIGN KEY (account_id) REFERENCES professor(id)ON DELETE CASCADE ON UPDATE CASCADE
);