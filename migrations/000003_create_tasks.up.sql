create table tasks(
    task_id UUID primary key,
    title Text not null,
    description Text not null,
    status varchar(50) not null default 'TODO',
    assignee_id UUID  null,
    workspace_id UUID not null,
    create_at Timestamp not null,
    constraint task_status_contraint 
        check (status in ('TODO', 'IN_PROGRESS', 'DONE', 'BLOCKED')),

    constraint fk_tasks_workspaces 
        foreign key (workspace_id) references workspaces(workspace_id)
        on delete cascade,

    constraint fk_tasks_users 
        foreign key (assignee_id) references users(user_id)
        on delete cascade
);



