CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    text_message TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_messages_created_at ON messages(created_at);