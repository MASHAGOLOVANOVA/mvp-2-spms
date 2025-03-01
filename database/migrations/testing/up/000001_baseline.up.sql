CREATE TABLE IF NOT EXISTS
    professor (
        id INT NOT NULL auto_increment,
        name VARCHAR(50) NOT NULL,
        surname VARCHAR(50) NOT NULL,
        middlename VARCHAR(50) NOT NULL,
        science_degree VARCHAR(100) NOT NULL,
        PRIMARY KEY(id)
    );

CREATE TABLE IF NOT EXISTS
    student (
        id INT NOT NULL auto_increment,
        name VARCHAR(50) NOT NULL,
        surname VARCHAR(50) NOT NULL,
        middlename VARCHAR(50) NOT NULL,
        enrollment_year INT UNSIGNED NOT NULL,
        PRIMARY KEY(id)
    );

CREATE TABLE IF NOT EXISTS
    project_status (
        id INT NOT NULL,
        name VARCHAR(50) NOT NULL,
        PRIMARY KEY(id)
    );

CREATE TABLE IF NOT EXISTS
    repository (
        id INT NOT NULL auto_increment,
        name VARCHAR(100) NOT NULL,
        is_public BOOLEAN NOT NULL,
        PRIMARY KEY(id)
    );

CREATE TABLE IF NOT EXISTS
    project_stage (
        id INT NOT NULL,
        name VARCHAR(50) NOT NULL,
        PRIMARY KEY(id)
    );
CREATE TABLE IF NOT EXISTS
    supervisor_review (
        id INT NOT NULL auto_increment,
        creation_date DATETIME NOT NULL,
        PRIMARY KEY(id)
    );
CREATE TABLE IF NOT EXISTS
    review_criteria (
        id INT NOT NULL auto_increment,
        description VARCHAR(500) NOT NULL,
        grade FLOAT NOT NULL,
        grade_weight FLOAT NOT NULL,
        supervisor_review_id INT NOT NULL,
        PRIMARY KEY(id),
        FOREIGN KEY (supervisor_review_id) REFERENCES supervisor_review(id) ON DELETE CASCADE ON UPDATE CASCADE
    );

CREATE TABLE IF NOT EXISTS
    project (
        id INT NOT NULL auto_increment,
        theme VARCHAR(100) NOT NULL,
        year INT NOT NULL,
        supervisor_id INT NOT NULL,
        student_id INT NOT NULL,
        status_id INT NOT NULL,
        stage_id INT NOT NULL,
        repo_id INT,
        grade FLOAT,
        supervisor_review_id INT,
        PRIMARY KEY(id),
        FOREIGN KEY (supervisor_id) REFERENCES professor(id)ON DELETE CASCADE ON UPDATE CASCADE,
        FOREIGN KEY (student_id) REFERENCES student(id)ON DELETE CASCADE ON UPDATE CASCADE,
        FOREIGN KEY (status_id) REFERENCES project_status(id)ON DELETE CASCADE ON UPDATE CASCADE,
        FOREIGN KEY (stage_id) REFERENCES project_stage(id)ON DELETE CASCADE ON UPDATE CASCADE,
        FOREIGN KEY (supervisor_review_id) REFERENCES supervisor_review(id)ON DELETE CASCADE ON UPDATE CASCADE,
        FOREIGN KEY (repo_id) REFERENCES repository(id)ON DELETE CASCADE ON UPDATE CASCADE
    );

CREATE TABLE IF NOT EXISTS
    task (
        id INT NOT NULL auto_increment,
        name VARCHAR(50) NOT NULL,
        description VARCHAR(300) NOT NULL,
        deadline DATETIME NOT NULL,
        project_id INT NOT NULL,
        PRIMARY KEY(id),
        FOREIGN KEY (project_id) REFERENCES project(id) ON DELETE CASCADE ON UPDATE CASCADE
    );
