DROP TRIGGER update_creator_search on creator_profile;
DROP TRIGGER insert_search on creator_profile;
DROP TRIGGER update_user_search on users;
DROP TRIGGER delete_search on creator_profile;

DROP FUNCTION upd_creators_search_creators();
DROP FUNCTION upd_user_search_creators();
DROP FUNCTION insert_search_creators();
DROP FUNCTION delete_search_creators();
DROP PROCEDURE update_search_creators();

DROP INDEX idx_rum_search_creators;

DROP TABLE search_creators;

DROP EXTENSION rum;