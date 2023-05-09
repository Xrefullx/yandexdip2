CREATE TABLE IF NOT EXISTS secrets
(
    id uuid primary key default uuid_generate_v4(),
    ver int not null,
    user_id      uuid NOT NULL,
    "data"     TEXT,
    is_deleted boolean not null default false,

    FOREIGN KEY (user_id) REFERENCES users (id)
    );