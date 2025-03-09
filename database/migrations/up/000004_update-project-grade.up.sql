ALTER TABLE project 
ADD final_grade FLOAT,
RENAME COLUMN grade to defence_grade;

ALTER TABLE review_criteria
MODIFY COLUMN grade FLOAT;
