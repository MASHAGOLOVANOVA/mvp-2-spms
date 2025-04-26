create table if not exists slot (
    id INT NOT NULL auto_increment,
    event_id varchar(2000) NOT NULL,
    is_online BOOLEAN NOT NULL,
    description varchar(2000) DEFAULT NULL,
    professor_id INT NOT NULL,
    planner_id VARCHAR(200),
    status INT NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (professor_id) REFERENCES professor(id) ON DELETE CASCADE ON UPDATE CASCADE
);

create table if not exists student_meeting(
    id INT NOT NULL auto_increment,
    student_id INT NOT NULL,
    slot_id INT NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (student_id) REFERENCES student(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (slot_id) REFERENCES slot(id) ON DELETE CASCADE ON UPDATE CASCADE
);

create table if not exists project_meeting (
    id INT NOT NULL auto_increment,
    project_id INT NOT NULL,
    stud_meeting_id INT NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (stud_meeting_id) REFERENCES student_meeting(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (project_id) REFERENCES project(id) ON DELETE CASCADE ON UPDATE CASCADE
);