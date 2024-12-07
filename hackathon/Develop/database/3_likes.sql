-- いいねテーブルの作成
DROP TABLE IF EXISTS likes;
CREATE TABLE likes (
    id CHAR(26) NOT NULL PRIMARY KEY,
    user_id CHAR(26) NOT NULL,
    post_id CHAR(26) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (post_id) REFERENCES posts(id)
);

-- 初期データの挿入
INSERT INTO likes (id, user_id, post_id) VALUES 
('00000000000000000000000001', '00000000000000000000000001', '00000000000000000000000002'),
('00000000000000000000000002', '00000000000000000000000002', '00000000000000000000000001'),
('00000000000000000000000003', '00000000000000000000000001', '00000000000000000000000007'),
('00000000000000000000000004', '00000000000000000000000002', '00000000000000000000000007'),
('00000000000000000000000005', '00000000000000000000000002', '00000000000000000000000008'),
('00000000000000000000000006', '00000000000000000000000001', '00000000000000000000000010');