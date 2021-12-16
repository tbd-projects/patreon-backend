DROP INDEX idx_like_search_nick_creators;

CREATE INDEX idx_like_search_nick_creators
    on users (nickname) include (users_id);

ALTER TABLE comments
    DROP COLUMN as_creator,
    DROP COLUMN date;

ALTER TABLE posts
    DROP COLUMN number_comments;