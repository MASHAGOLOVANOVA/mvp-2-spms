ALTER TABLE review_criteria
RENAME COLUMN description to criteria_id;

ALTER TABLE review_criteria
MODIFY COLUMN criteria_id INT NOT NULL,
ADD CONSTRAINT FK_ReviewcCriteria FOREIGN KEY (criteria_id) REFERENCES criteria(id) ON DELETE CASCADE ON UPDATE CASCADE;
