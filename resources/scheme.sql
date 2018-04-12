DROP TABLE IF EXISTS Users CASCADE;

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
  rating FLOAT
);

CREATE TABLE quest_user_link (
  id SERIAL PRIMARY KEY ,
  user_id INT REFERENCES Users(id),
  quest_id INT REFERENCES Quest(id),
  started BOOLEAN,
  completed BOOLEAN,
  marked BOOLEAN,
  mark FLOAT,
  UNIQUE (user_id, quest_id)
);
