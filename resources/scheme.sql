CREATE EXTENSION IF NOT EXISTS Postgis;

DROP TABLE IF EXISTS Users;
DROP TYPE IF EXISTS SEX;

DROP INDEX IF EXISTS position_user_idx;

CREATE TYPE SEX AS ENUM ('M', 'F', '');

CREATE TABLE Users (
  id       SERIAL PRIMARY KEY,
  login    VARCHAR(50) UNIQUE ,
  password VARCHAR(50),
  sex      SEX NOT NULL DEFAULT '',
  age      INT,
  about    VARCHAR(1000)
);

CREATE TABLE Position (
  id     SERIAL PRIMARY KEY,
  userId INTEGER REFERENCES Users (id),
  point  GEOMETRY,
  time   TIMESTAMP
);

CREATE INDEX position_user_idx
  ON Position (userId);
