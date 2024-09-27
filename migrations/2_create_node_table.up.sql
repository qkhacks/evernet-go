CREATE TABLE nodes
(
    identifier          TEXT PRIMARY KEY,
    display_name        TEXT NOT NULL,
    signing_private_key TEXT NOT NULL,
    signing_public_key  TEXT NOT NULL,
    creator             TEXT NOT NULL,
    created_at          INT  NOT NULL,
    updated_at          INT  NOT NULL
);
