create table workspaces (
    workspace_id UUID primary key,
    name Text not null,
    description Text not null,
    owner_id UUID not null,
    create_at Timestamp not null,

    constraint fk_workspaces_users
        foreign key (owner_id) 
        references users(user_id) 
        on delete cascade
);

