BEGIN;

CREATE TYPE "notification_type" AS ENUM ('friend_request', 'new_message');

CREATE TABLE IF NOT EXISTS notifications (
    "id" SERIAL PRIMARY KEY,

    "data" JSONB NOT NULL,
    "type" notification_type NOT NULL,
    "seen" BOOLEAN DEFAULT FALSE,

    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    "user_id" INTEGER NOT NULL,

    FOREIGN KEY ("user_id") REFERENCES users ("id")
);

COMMIT;