ALTER TABLE professor 
ADD university_id INT,
ADD CONSTRAINT FK_profuid FOREIGN KEY (university_id) REFERENCES university(id) ON DELETE CASCADE ON UPDATE CASCADE;