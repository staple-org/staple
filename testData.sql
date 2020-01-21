create table users (email varchar(255), password text, confirm_link text);
create table staples (name varchar(255), id serial, content text, created_at timestamp, archived bool, user_email varchar(255));
create user staple with password 'password123';
create database staples;
GRANT ALL PRIVILEGES ON DATABASE staples TO staple;
ALTER USER staple WITH SUPERUSER;