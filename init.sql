BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users(
  "id" uuid DEFAULT uuid_generate_v4(),
  "name" varchar NOT NULL,
  "screen_name" varchar NOT NULL,
  "email" varchar NOT NULL,
  "bio" varchar NOT NULL,
  "location" varchar NOT NULL,
  "website" varchar NOT NULL,
  "birth_date" date NOT NULL,
  "profile_image_url" text NOT NULL,
  "profile_banner_url" text NOT NULL,
  "followers_count" int NOT NULL DEFAULT 0,
  "followings_count" int NOT NULL DEFAULT 0,
  "created_at" timestamptz DEFAULT NOW(),
  "updated_at" timestamptz DEFAULT NOW(),
  PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX IF NOT EXISTS index_users_email ON users USING btree("email");

CREATE INDEX IF NOT EXISTS index_users_screen_name ON users USING btree("screen_name");

CREATE INDEX IF NOT EXISTS index_users_created_at ON users USING btree("created_at");

CREATE TABLE IF NOT EXISTS followers(
  "followee_id" uuid NOT NULL REFERENCES users("id") ON DELETE CASCADE,
  "follower_id" uuid NOT NULL REFERENCES users("id") ON DELETE CASCADE,
  "created_at" timestamptz DEFAULT NOW(),
  PRIMARY KEY ("followee_id", "follower_id")
);

CREATE INDEX IF NOT EXISTS index_followers_created_at ON followers USING btree("created_at");

-- reference: https://github.com/shuber/postgres-twitter/blob/master/sql/005_behaviors.sql
-- increment_counter(table_name, column_name, pk_name, pk_value, step)
-- usage: increment_counter('users', 'followers_count', 'id', '00000000-0000-0000-0000-000000000000', 1)
-- description: increment/decrement counter column
CREATE FUNCTION increment_counter(table_name text, column_name text, pk_name text, pk_value uuid, step integer)
  RETURNS VOID
  AS $$
DECLARE
  table_name text := quote_ident(table_name);
  column_name text := quote_ident(column_name);
  conditions text := ' WHERE ' || quote_ident(pk_name) || ' = $1';
  updates text := column_name || '=' || column_name || '+' || step;
BEGIN
  EXECUTE 'UPDATE ' || table_name || ' SET ' || updates || conditions
  USING pk_value;
  END;
$$
LANGUAGE plpgsql;

-- counter_cache(table_name, counter_name, fk_name, pk_name)
-- usage: counter_cache('users', 'followers_count', 'followee_id', 'id')
-- description: increment/decrement counter column when insert/delete/update
CREATE FUNCTION counter_cache()
  RETURNS TRIGGER
  AS $$
DECLARE
  table_name text := quote_ident(TG_ARGV[0]);
  counter_name text := quote_ident(TG_ARGV[1]);
  fk_name text := quote_ident(TG_ARGV[2]);
  pk_name text := quote_ident(TG_ARGV[3]);
  fk_changed boolean := FALSE;
  fk_value uuid;
  record record;
BEGIN
  IF TG_OP = 'UPDATE' THEN
    record := NEW;
    EXECUTE 'SELECT ($1).' || fk_name || ' != ' || '($2).' || fk_name INTO fk_changed
    USING OLD, NEW;
  END IF;
    IF TG_OP = 'DELETE' OR fk_changed THEN
      record := OLD;
      EXECUTE 'SELECT ($1).' || fk_name INTO fk_value
      USING record;
      PERFORM
        increment_counter(table_name, counter_name, pk_name, fk_value, -1);
    END IF;
      IF TG_OP = 'INSERT' OR fk_changed THEN
        record := NEW;
        EXECUTE 'SELECT ($1).' || fk_name INTO fk_value
        USING record;
        PERFORM
          increment_counter(table_name, counter_name, pk_name, fk_value, 1);
      END IF;
        RETURN record;
END;
$$
LANGUAGE plpgsql;

CREATE TRIGGER update_follower_following
  AFTER INSERT OR UPDATE OR DELETE ON followers
  FOR EACH ROW
  EXECUTE PROCEDURE counter_cache('users', 'following', 'follower_id', 'id');

CREATE TRIGGER update_user_followers
  AFTER INSERT OR UPDATE OR DELETE ON followers
  FOR EACH ROW
  EXECUTE PROCEDURE counter_cache('users', 'followers', 'user_id', 'id');

COMMIT;

