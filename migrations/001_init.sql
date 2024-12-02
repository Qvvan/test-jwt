-- 001_create_table_example.sql
create table users  (
  id uuid primary key,
  refresh_token varchar(100),
  created_at timestamp DEFAULT current_timestamp
);

---- create above / drop below ----

drop table users;
