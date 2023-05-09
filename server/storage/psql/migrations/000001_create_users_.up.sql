create extension if not exists "uuid-ossp";
CREATE TABLE IF NOT EXISTS users
(
    id uuid primary key default uuid_generate_v4(),
    login varchar (255) unique not null,
    pass_hash varchar (60) not null,
    master_hash varchar (60) not null
    );