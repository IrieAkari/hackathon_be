-- Create posts table
DROP TABLE IF EXISTS posts;
CREATE TABLE posts (
    id CHAR(26) NOT NULL PRIMARY KEY,
    parent_id CHAR(26),
    is_parent_deleted BOOLEAN DEFAULT FALSE, -- Column to indicate if parent post is deleted
    user_id CHAR(26) NOT NULL,
    content TEXT NOT NULL,
    likes_count INT DEFAULT 0,
    replys_count INT DEFAULT 0, -- Add reply count
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    trust_score INT DEFAULT -1, -- Column to represent trust score
    trust_description TEXT DEFAULT 'Unassessed', -- Column to represent trust description
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (parent_id) REFERENCES posts(id)
);

-- Insert initial data
INSERT INTO posts (id, parent_id, is_parent_deleted, user_id, content, likes_count, replys_count, trust_score, trust_description) VALUES 
('00000000000000000000000001', NULL, FALSE, '00000000000000000000000001', 'a   This is Hanako first post', 1, 1, 90, 'Post from a reliable source'),
('00000000000000000000000002', '00000000000000000000000001', FALSE, '00000000000000000000000002', 'a-a   This is Taro reply to Hanako', 1, 0, 85, 'Reply from a reliable source'),
('00000000000000000000000003', NULL, FALSE, '00000000000000000000000003', 'b   Three days left until the hackathon', 0, 0, 95, 'Reliable information about an event'),
('00000000000000000000000004', NULL, FALSE, '00000000000000000000000003', 'c   This is a test post', 0, 0, 50, 'Test post, accuracy is uncertain'),
('00000000000000000000000005', NULL, FALSE, '00000000000000000000000001', 'd   This is Hanako', 0, 0, 100, 'Self-introduction'),
('00000000000000000000000006', NULL, FALSE, '00000000000000000000000001', 'e   This is a longer test post. The character limit is 200 characters. What is the actual Twitter character limit? 140 characters?', 0, 0, 30, 'Test post, accuracy of information is unknown'),
('00000000000000000000000007', NULL, FALSE, '00000000000000000000000001', 'f   Hanako second post', 2, 0, 80, 'Post from a reliable source'),
('00000000000000000000000008', '00000000000000000000000002', FALSE, '00000000000000000000000002', 'g   Taro new post', 1, 0, 75, 'Post from a reliable source'),
('00000000000000000000000009', NULL, FALSE, '00000000000000000000000003', 'h   Another test post', 0, 0, 40, 'Test post, accuracy is uncertain'),
('00000000000000000000000010', NULL, FALSE, '00000000000000000000000001', 'i   Hanako third post', 1, 0, 85, 'Post from a reliable source'),
('00000000000000000000000011', NULL, FALSE, '00000000000000000000000003', 'j   Final test post', 0, 0, 20, 'Test post, accuracy is very low'),
('00000000000000000000000012', '00000000000000000000000001', FALSE, '00000000000000000000000002', 'a-b   Another reply to Hanako', 0, 0, 70, 'Another reply from a reliable source'),
('00000000000000000000000013', '00000000000000000000000012', FALSE, '00000000000000000000000003', 'a-b-a   Reply to another reply', 0, 0, 60, 'Reply to another reply, accuracy is uncertain'),
('00000000000000000000000014', '00000000000000000000000003', FALSE, '00000000000000000000000001', 'b-a   Reply to hackathon post', 0, 0, 80, 'Reply to hackathon post, reliable source'),
('00000000000000000000000015', '00000000000000000000000014', FALSE, '00000000000000000000000002', 'b-a-a   Further reply to hackathon post', 0, 0, 75, 'Further reply to hackathon post, reliable source'),
('00000000000000000000000016', '00000000000000000000000007', FALSE, '00000000000000000000000003', 'f-a   Reply to Hanako second post', 0, 0, 65, 'Reply to Hanako second post, accuracy is uncertain'),
('00000000000000000000000017', '00000000000000000000000016', FALSE, '00000000000000000000000001', 'f-a-a   Further reply to Hanako second post', 0, 0, 60, 'Further reply to Hanako second post, accuracy is uncertain');