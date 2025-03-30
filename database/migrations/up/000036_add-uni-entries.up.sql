INSERT INTO city (`name`)
VALUES ("Perm");
SET @last_id_in_city = LAST_INSERT_ID();

INSERT INTO university (`name`, `city_id`)
VALUES ("HSE", @last_id_in_city);

SET @last_id_in_uni = LAST_INSERT_ID();

INSERT INTO department (`name`, `uni_id`)
VALUES ("unknown", @last_id_in_uni);

SET @last_id_in_dept = LAST_INSERT_ID();

INSERT INTO faculty (`name`, `dept_id`)
VALUES ("unknown", @last_id_in_dept);

SET @last_id_in_faculty = LAST_INSERT_ID();

INSERT INTO educational_programme (`name`, `ed_level`, `faculty_id`)
VALUES ("Software Engineering", 1, @last_id_in_faculty), ("Business Informatics", 1, @last_id_in_faculty);