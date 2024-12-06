CREATE TABLE url 
(
    id serial not null unique,
    alias TEXT NOT NULL UNIQUE,
    url TEXT NOT NULL
);
CREATE INDEX idx_alias ON url(alias);