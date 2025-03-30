ALTER TABLE user_account 
DROP FOREIGN KEY user_account_ibfk_1;
ALTER TABLE user_account 
DROP PRIMARY KEY,
ADD PRIMARY KEY (login),
ADD CONSTRAINT user_account_ibfk_1 FOREIGN KEY (professor_id) REFERENCES professor (id) ON DELETE CASCADE ON UPDATE CASCADE;