ALTER TABLE review_criteria
RENAME COLUMN criteria_id to description;

ALTER TABLE review_criteria
MODIFY COLUMN description VARCHAR(500) NOT NULL,
DROP FOREIGN KEY FK_ReviewcCriteria;

DROP TABLE criteria;