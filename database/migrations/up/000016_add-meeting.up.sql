CREATE TABLE meeting (
    id INT NOT NULL auto_increment,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(100) NOT NULL,
    meeting_time DATETIME NOT NULL, 
    student_id INT NOT NULL,
    is_online BOOLEAN NOT NULL,
    professor_id INT NOT NULL,
    status INT NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (student_id) REFERENCES student(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (professor_id) REFERENCES professor(id) ON DELETE CASCADE ON UPDATE CASCADE
);