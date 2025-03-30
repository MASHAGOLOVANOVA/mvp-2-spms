ALTER TABLE drive_integration
ADD type INT NOT NULL DEFAULT 0;
ALTER TABLE planner_integration
ADD type INT NOT NULL DEFAULT 0;
ALTER TABLE repository
ADD repo_hub_type INT NOT NULL DEFAULT 0;