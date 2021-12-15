CREATE TABLE user_settings (
    user_id bigint not null references users (users_id) on delete cascade,
    get_post bool not null,
    get_comment bool not null,
    get_sub bool not null
);

CREATE PROCEDURE insert_settings_for_user()
AS $$
DECLARE
    Id bigint;
BEGIN
    FOR Id IN SELECT DISTINCT users_id FROM users LOOP
        INSERT INTO user_settings (user_id, get_sub, get_comment, get_post)  VALUES  (Id, true, true, true);
    END LOOP;
END;
$$
    LANGUAGE plpgsql;
CALL insert_settings_for_user()