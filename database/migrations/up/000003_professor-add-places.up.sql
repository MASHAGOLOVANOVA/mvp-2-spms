CREATE TABLE IF NOT EXISTS professor_available_places (
    professor_id INT NOT NULL,
    ed_prog_id INT NOT NULL,
    cource INT NOT NULL,
    available_places INT NOT NULL,
    PRIMARY KEY(professor_id, ed_prog_id, cource),
    FOREIGN KEY (professor_id) REFERENCES professor(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (ed_prog_id) REFERENCES educational_programme(id) ON DELETE CASCADE ON UPDATE CASCADE
);