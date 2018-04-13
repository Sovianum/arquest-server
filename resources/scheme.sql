DROP TABLE IF EXISTS Users CASCADE;
DROP TABLE IF EXISTS quest CASCADE;
DROP TABLE IF EXISTS quest_user_link CASCADE;

DROP TYPE IF EXISTS SEX;

CREATE TYPE SEX AS ENUM ('M', 'F', '');

CREATE TABLE users (
  id       SERIAL PRIMARY KEY,
  login    VARCHAR(50) UNIQUE ,
  password BYTEA,
  sex      SEX NOT NULL DEFAULT '',
  age      INT,
  about    VARCHAR(1000)
);

CREATE TABLE quest (
  id SERIAL PRIMARY KEY,
  name VARCHAR(50),
  description VARCHAR(1000),
  rating FLOAT,
  mark_count INT
);

CREATE TABLE quest_user_link (
  id SERIAL PRIMARY KEY ,
  user_id INT REFERENCES Users(id),
  quest_id INT REFERENCES Quest(id),
  started BOOLEAN DEFAULT TRUE ,
  completed BOOLEAN DEFAULT FALSE ,
  marked BOOLEAN DEFAULT FALSE ,
  mark FLOAT DEFAULT 0,
  CONSTRAINT ux_user_id_quest_id UNIQUE (user_id, quest_id)
);
