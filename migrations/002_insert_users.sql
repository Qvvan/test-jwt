-- 002_insert_users.sql

-- Добавление нескольких пользователей
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

insert into users (id, refresh_token, created_at)
values
  (uuid_generate_v4(), 'refresh_token_1_value', current_timestamp),
  (uuid_generate_v4(), 'refresh_token_2_value', current_timestamp),
  (uuid_generate_v4(), 'refresh_token_3_value', current_timestamp),
  (uuid_generate_v4(), 'refresh_token_4_value', current_timestamp),
  (uuid_generate_v4(), 'refresh_token_5_value', current_timestamp);
