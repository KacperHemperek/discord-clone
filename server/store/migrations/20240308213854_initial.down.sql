BEGIN;

DROP INDEX IF EXISTS "chat_to_user_connection";
DROP INDEX IF EXISTS "inviter_id_to_friend_id";
DROP INDEX IF EXISTS "user_email_index";

DROP TABLE IF EXISTS "chat_to_user";
DROP TABLE IF EXISTS "messages";
DROP TABLE IF EXISTS "chats";
DROP TABLE IF EXISTS "friendships";
DROP TABLE IF EXISTS "users";

DROP TYPE IF EXISTS status;
DROP TYPE IF EXISTS chat_type;

COMMIT;