ALTER TABLE project 
DROP COLUMN final_grade,
RENAME COLUMN defence_grade to grade;

ALTER TABLE review_criteria
MODIFY COLUMN grade FLOAT NOT NULL;
