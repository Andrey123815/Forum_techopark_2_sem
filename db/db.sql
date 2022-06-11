CREATE EXTENSION IF NOT EXISTS CITEXT;


DROP TABLE IF EXISTS "users" CASCADE;
DROP TABLE IF EXISTS "forums" CASCADE;
DROP TABLE IF EXISTS "threads" CASCADE;
DROP TABLE IF EXISTS "posts" CASCADE;
DROP TABLE IF EXISTS "votes" CASCADE;
DROP TABLE IF EXISTS "forum_users" CASCADE;


DROP FUNCTION IF EXISTS create_post();
DROP TRIGGER IF EXISTS "create_post" ON "posts";
DROP FUNCTION IF EXISTS create_thread();
DROP TRIGGER IF EXISTS "create_thread" ON "threads";


DROP INDEX IF EXISTS idx_users_nickname;
DROP INDEX IF EXISTS idx_users_email;

DROP INDEX IF EXISTS idx_forums_slug;

DROP INDEX IF EXISTS idx_threads_slug;
DROP INDEX IF EXISTS idx_threads_forum;

DROP INDEX IF EXISTS idx_posts_thread;
DROP INDEX IF EXISTS idx_posts_thread_path;

DROP INDEX IF EXISTS idx_votes_user_thread;


CREATE UNLOGGED TABLE IF NOT EXISTS "users"
(
    "id"       BIGSERIAL                  NOT NULL PRIMARY KEY,
    "nickname" CITEXT COLLATE "ucs_basic" NOT NULL UNIQUE,
    "fullname" CITEXT                     NOT NULL,
    "about"    TEXT,
    "email"    CITEXT                     NOT NULL UNIQUE
);

CREATE UNLOGGED TABLE IF NOT EXISTS "forums"
(
    "id"      BIGSERIAL NOT NULL PRIMARY KEY,
    "title"   TEXT      NOT NULL,
    "user"    CITEXT    NOT NULL,
    "slug"    CITEXT    NOT NULL UNIQUE,
    "posts"   BIGINT DEFAULT 0,
    "threads" INT    DEFAULT 0
);

CREATE UNLOGGED TABLE IF NOT EXISTS "threads"
(
    "id"      BIGSERIAL   NOT NULL PRIMARY KEY,
    "title"   TEXT        NOT NULL,
    "author"  CITEXT      NOT NULL,
    "forum"   CITEXT,
    "message" TEXT        NOT NULL,
    "votes"   INT         DEFAULT 0,
    "slug"    CITEXT,
    "created" TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE UNLOGGED TABLE IF NOT EXISTS "posts"
(
    "id"       BIGSERIAL NOT NULL PRIMARY KEY,
    "parent"   BIGINT             DEFAULT 0,
    "author"   CITEXT    NOT NULL,
    "message"  TEXT      NOT NULL,
    "isEdited" BOOL               DEFAULT false,
    "forum"    CITEXT,
    "thread"   INT,
    "created"  TIMESTAMPTZ        DEFAULT now(),
    "path"     BIGINT[]  NOT NULL DEFAULT '{0}'
);

CREATE UNLOGGED TABLE IF NOT EXISTS "votes"
(
    "id"     BIGSERIAL                       NOT NULL PRIMARY KEY,
    "user"   CITEXT REFERENCES "users" (nickname)   NOT NULL,
    "thread" BIGINT REFERENCES "threads" (id) NOT NULL,
    "voice"  INT,
    CONSTRAINT checks UNIQUE ("user", "thread")
);

CREATE UNLOGGED TABLE forum_users
(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    nickname BIGINT REFERENCES "users" (id) NOT NULL,
    forum     BIGINT REFERENCES "forums" (id)NOT NULL
);


CREATE FUNCTION create_post() RETURNS TRIGGER AS
$$
BEGIN
    UPDATE "forums"
    SET "posts" = posts + 1
    WHERE "slug" = new.forum;
    new.path = (SELECT "path" FROM "posts" WHERE "id" = new.parent LIMIT 1) || new.id;
    INSERT INTO "forum_users" ("nickname", "forum")
    VALUES ((SELECT "id" FROM "users" WHERE NEW.author = nickname), (SELECT "id" FROM "forums" WHERE NEW.forum = slug));
    return new;
END
$$ LANGUAGE plpgsql;

CREATE TRIGGER "create_post"
    BEFORE INSERT
    ON "posts"
    FOR EACH ROW
EXECUTE PROCEDURE create_post();


CREATE FUNCTION create_thread() RETURNS TRIGGER AS
$$
BEGIN
    UPDATE "forums"
    SET "threads" = threads + 1
    WHERE "slug" = new.forum;
    INSERT INTO "forum_users" ("nickname", "forum")
    VALUES ((SELECT "id" FROM "users" WHERE nickname = NEW.author), (SELECT "id" FROM "forums" WHERE slug = NEW.forum));
    return new;
END
$$ LANGUAGE plpgsql;

CREATE TRIGGER "create_thread"
    BEFORE INSERT
    ON "threads"
    FOR EACH ROW
EXECUTE PROCEDURE create_thread();


CREATE UNIQUE INDEX IF NOT EXISTS idx_users_nickname ON users ("nickname");
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users ("email");

CREATE UNIQUE INDEX IF NOT EXISTS idx_forums_slug ON forums ("slug");

CREATE INDEX IF NOT EXISTS idx_threads_slug ON threads ("slug");
CREATE INDEX IF NOT EXISTS idx_threads_forum ON threads ("forum");

CREATE INDEX IF NOT EXISTS idx_posts_thread ON posts ("thread");
CREATE UNIQUE INDEX IF NOT EXISTS idx_posts_thread_path ON posts ("thread", "path");

CREATE UNIQUE INDEX IF NOT EXISTS idx_votes_user_thread ON votes ("user", "thread");
