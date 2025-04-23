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

create table slot (
    id INT NOT NULL auto_increment,
    event_id varchar(2000) NOT NULL,
    is_online BOOLEAN NOT NULL,
    description varchar(2000) DEFAULT NULL,
    professor_id INT NOT NULL,
    planner_id VARCHAR(200),
    PRIMARY KEY (id),
    FOREIGN KEY (professor_id) REFERENCES professor(id) ON DELETE CASCADE ON UPDATE CASCADE
);

create table student_meeting (
    id INT NOT NULL auto_increment,
    student_id INT NOT NULL,
    slot_id INT NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (student_id) REFERENCES student(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (slot_id) REFERENCES slot(id) ON DELETE CASCADE ON UPDATE CASCADE
);