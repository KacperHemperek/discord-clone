BEGIN;

ALTER TABLE chat_to_user
    DROP CONSTRAINT chat_to_user_chat_id_fkey;

ALTER TABLE chat_to_user
    ADD CONSTRAINT chat_to_user_chat_id_fkey
        FOREIGN KEY ("chat_id") REFERENCES "chats" ("id");

ALTER TABLE chat_to_user
    DROP CONSTRAINT chat_to_user_user_id_fkey;

ALTER TABLE chat_to_user
    ADD CONSTRAINT chat_to_user_user_id_fkey
        FOREIGN KEY ("user_id") REFERENCES "users" ("id");

COMMIT;