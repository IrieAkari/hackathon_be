-- 投稿テーブルの作成
DROP TABLE IF EXISTS posts;
CREATE TABLE posts (
    id CHAR(26) NOT NULL PRIMARY KEY,
    parent_id CHAR(26),
    user_id CHAR(26) NOT NULL,
    content TEXT NOT NULL,
    likes_count INT DEFAULT 0,
    replys_count INT DEFAULT 0, -- リプライ数を追加
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (parent_id) REFERENCES posts(id)
);

-- 初期データの挿入
INSERT INTO posts (id, parent_id, user_id, content, likes_count, replys_count) VALUES 
('00000000000000000000000001', NULL, '00000000000000000000000001', 'a   This is Hanako first post', 1, 1),
('00000000000000000000000002', '00000000000000000000000001', '00000000000000000000000002', 'a-a   This is Taro reply to Hanako', 1, 0),
('00000000000000000000000003', NULL, '00000000000000000000000003', 'b   Three days left until the hackathon', 0, 0),
('00000000000000000000000004', NULL, '00000000000000000000000003', 'c   This is a test post', 0, 0),
('00000000000000000000000005', NULL, '00000000000000000000000001', 'd   This is Hanako', 0, 0),
('00000000000000000000000006', NULL, '00000000000000000000000001', 'e   This is a longer test post. The character limit is 200 characters. What is the actual Twitter character limit? 140 characters?', 0, 0);