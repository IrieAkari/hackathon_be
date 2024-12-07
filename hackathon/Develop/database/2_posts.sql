-- Create posts table
DROP TABLE IF EXISTS posts;
CREATE TABLE posts (
    id CHAR(26) NOT NULL PRIMARY KEY,
    parent_id CHAR(26),
    user_id CHAR(26) NOT NULL,
    content TEXT NOT NULL,
    likes_count INT DEFAULT 0,
    replys_count INT DEFAULT 0, -- Add reply count
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    trust_score INT DEFAULT 100, -- Column to represent trust score
    trust_description TEXT, -- Column to represent trust description
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (parent_id) REFERENCES posts(id)
);

-- Insert initial data
INSERT INTO posts (id, parent_id, user_id, content, likes_count, replys_count, trust_score, trust_description) VALUES 
('00000000000000000000000001', NULL, '00000000000000000000000001', 'a   This is Hanako first post', 1, 1, 90, 'Post from a reliable source'),
('00000000000000000000000002', '00000000000000000000000001', '00000000000000000000000002', 'a-a   This is Taro reply to Hanako', 1, 0, 85, 'Reply from a reliable source'),
('00000000000000000000000003', NULL, '00000000000000000000000003', 'b   Three days left until the hackathon', 0, 0, 95, 'Reliable information about an event'),
('00000000000000000000000004', NULL, '00000000000000000000000003', 'c   This is a test post', 0, 0, 50, 'Test post, accuracy is uncertain'),
('00000000000000000000000005', NULL, '00000000000000000000000001', 'd   This is Hanako', 0, 0, 100, 'Self-introduction'),
('00000000000000000000000006', NULL, '00000000000000000000000001', 'e   This is a longer test post. The character limit is 200 characters. What is the actual Twitter character limit? 140 characters?', 0, 0, 30, 'Test post, accuracy of information is unknown'),
('00000000000000000000000007', NULL, '00000000000000000000000001', 'f   Hanako second post', 2, 0, 80, 'Post from a reliable source'),
('00000000000000000000000008', NULL, '00000000000000000000000002', 'g   Taro new post', 1, 0, 75, 'Post from a reliable source'),
('00000000000000000000000009', NULL, '00000000000000000000000003', 'h   Another test post', 0, 0, 40, 'Test post, accuracy is uncertain'),
('00000000000000000000000010', NULL, '00000000000000000000000001', 'i   Hanako third post', 1, 0, 85, 'Post from a reliable source'),
('00000000000000000000000011', NULL, '00000000000000000000000003', 'j   Final test post', 0, 0, 20, 'Test post, accuracy is very low');
