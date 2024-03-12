BEGIN;

CREATE TYPE "status" AS ENUM ('pending', 'accepted', 'rejected');
CREATE TYPE "chat_type" AS ENUM ('private', 'group');

CREATE TABLE IF NOT EXISTS "users" (
    "id" SERIAL PRIMARY KEY,

    "email" TEXT NOT NULL,
    "username" TEXT NOT NULL,
    "password" TEXT NOT NULL,
    "active" BOOLEAN NOT NULL DEFAULT false,

    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "friendships" (
    "id" SERIAL PRIMARY KEY,

    "status" status NOT NULL DEFAULT 'pending',
    "seen" BOOLEAN NOT NULL DEFAULT false,
    "inviter_id" INTEGER NOT NULL,
    "friend_id" INTEGER NOT NULL,

    "requested_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "accepted_at" TIMESTAMP(3),

    FOREIGN KEY ("inviter_id") REFERENCES "users" ("id"),
    FOREIGN KEY ("friend_id") REFERENCES "users" ("id")
);

CREATE TABLE IF NOT EXISTS "chats" (
    "id" SERIAL PRIMARY KEY,

    "name" TEXT NOT NULL,
    "type" chat_type NOT NULL,

    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "messages" (
    "id" SERIAL PRIMARY KEY,

    "text" TEXT,
    "image" TEXT,
    "sender_id" INTEGER NOT NULL,
    "chat_id" INTEGER NOT NULL,

    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY ("sender_id") REFERENCES "users"("id") ON DELETE CASCADE,
    FOREIGN KEY ("chat_id") REFERENCES "chats"("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "chat_to_user" (
    "chat_id" INTEGER NOT NULL,
    "user_id" INTEGER NOT NULL,

    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY ("chat_id") REFERENCES "chats" ("id"),
    FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);

CREATE UNIQUE INDEX "user_email_index" ON "users"("email");

CREATE UNIQUE INDEX "inviter_id_to_friend_id" ON "friendships"("inviter_id", "friend_id");

CREATE UNIQUE INDEX "chat_to_user_connection" ON "chat_to_user"("chat_id", "user_id");

COMMIT;