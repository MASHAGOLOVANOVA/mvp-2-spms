CREATE TABLE IF NOT EXISTS department (
    id INT NOT NULL auto_increment,
    name VARCHAR(150) NOT NULL,
    uni_id INT NOT NULL DEFAULT 1,
    PRIMARY KEY(id),
    FOREIGN KEY (uni_id) REFERENCES university(id) ON DELETE CASCADE ON UPDATE CASCADE
);
