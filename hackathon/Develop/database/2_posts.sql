-- 投稿テーブルの作成
DROP TABLE IF EXISTS posts;
CREATE TABLE posts (
    id CHAR(26) NOT NULL PRIMARY KEY,
    parent_id CHAR(26),
    user_id CHAR(26) NOT NULL,
    content TEXT NOT NULL,
    likes_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (parent_id) REFERENCES posts(id)
);

-- 初期データの挿入
INSERT INTO posts (id, parent_id, user_id, content, likes_count) VALUES 
('00000000000000000000000001', NULL, '00000000000000000000000001', 'This is Hanako\'s first post', 10),
('00000000000000000000000002', '00000000000000000000000001', '00000000000000000000000002', 'This is Taro\'s reply to Hanako', 5);