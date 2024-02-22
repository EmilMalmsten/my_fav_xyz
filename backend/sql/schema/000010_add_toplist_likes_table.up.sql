CREATE TABLE toplist_likes (
    user_id INT,
    toplist_id INT,
    liked_at TIMESTAMP,
    PRIMARY KEY (user_id, toplist_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (toplist_id) REFERENCES toplists(id)
);