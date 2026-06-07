CREATE TABLE comments (
    comment_id UUID PRIMARY KEY,
    task_id UUID NOT NULL,
    user_id UUID NOT NULL,
    content TEXT NOT NULL,
    create_at TIMESTAMP NOT NULL DEFAULT NOW(),
    update_at TIMESTAMP NULL,

    CONSTRAINT fk_comments_tasks
        FOREIGN KEY (task_id) REFERENCES tasks(task_id)
        ON DELETE CASCADE,

    CONSTRAINT fk_comments_users
        FOREIGN KEY (user_id) REFERENCES users(user_id)
        ON DELETE CASCADE
);

CREATE INDEX idx_comments_task_id ON comments(task_id);
