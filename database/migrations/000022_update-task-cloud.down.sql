ALTER TABLE task
MODIFY COLUMN description VARCHAR(300),
MODIFY COLUMN folder_id INT,
MODIFY COLUMN task_file_id INT;