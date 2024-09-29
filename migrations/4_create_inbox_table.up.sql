CREATE TABLE inboxes
(
    identifier      TEXT PRIMARY KEY,
    display_name    TEXT NOT NULL,
    node_identifier TEXT NOT NULL,
    actor_address   TEXT NOT NULL,
    created_at      INT  NOT NULL,
    updated_at      INT  NOT NULL
);
