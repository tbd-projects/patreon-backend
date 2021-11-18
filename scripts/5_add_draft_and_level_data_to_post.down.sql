ALTER TABLE posts
    DROP COLUMN is_draft;

ALTER TABLE posts_data
    DROP COLUMN level;

DROP PROCEDURE update_level_post_data;