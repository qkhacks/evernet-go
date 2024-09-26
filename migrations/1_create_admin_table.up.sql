CREATE TABLE admins
(
    username   TEXT PRIMARY KEY,
    password   TEXT NOT NULL,
    creator    TEXT NOT NULL,
    created_at INT  NOT NULL,
    updated_at INT  NOT NULL
)