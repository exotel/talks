-- !selectOne UserByID
-- $1: id int64

SELECT id, name FROM users WHERE id = $1 LIMIT 1;
