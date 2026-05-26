create table users (
    user_id UUID primary key,
    email varchar(255) unique not null,
    password_hash varchar(255) not null,
    full_name Text not null,
    create_at Timestamp not null
);

create index idx_users_email
on users(email);