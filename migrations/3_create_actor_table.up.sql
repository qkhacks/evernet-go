CREATE TABLE actors
(
    identifier      TEXT PRIMARY KEY,
    display_name    TEXT NOT NULL,
    type            TEXT NOT NULL,
    password        TEXT NOT NULL,
    node_identifier TEXT NOT NULL,
    creator         TEXT NOT NULL,
    created_at      INT  NOT NULL,
    updated_at      INT  NOT NULL
);
