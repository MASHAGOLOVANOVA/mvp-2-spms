ALTER TABLE student
ADD CONSTRAINT FK_stud_edprog FOREIGN KEY (educational_programme_id) REFERENCES educational_programme(id) ON DELETE CASCADE ON UPDATE CASCADE;