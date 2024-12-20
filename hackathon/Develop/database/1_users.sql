-- ユーザーテーブルの作成
DROP TABLE IF EXISTS users;
CREATE TABLE users (
    id CHAR(26) NOT NULL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 初期データの挿入
INSERT INTO users (id, name, email) VALUES 
('00000000000000000000000001', 'hanako', 'hanako@example.com'),
('00000000000000000000000002', 'taro', 'taro@example.com'),
('00000000000000000000000003', 'ichiro', 'ichiro@example.com');