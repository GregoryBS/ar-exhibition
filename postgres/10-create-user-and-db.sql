-- file: 10-create-user-and-db.sql
CREATE DATABASE art;
CREATE ROLE program WITH PASSWORD 'test';
GRANT ALL PRIVILEGES ON DATABASE art TO program;
ALTER ROLE program WITH LOGIN;