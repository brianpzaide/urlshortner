CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    email citext UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    created_at timestamp without time zone DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS urls (
    url_key character(7) PRIMARY KEY,
    target_url VARCHAR(255) NOT NULL,
    visits bigint default 0,
    user_id bigint NOT NULL,
    created_at timestamp without time zone DEFAULT NOW(),
    updated_at timestamp without time zone DEFAULT NOW(),
    CONSTRAINT fk_user
      FOREIGN KEY(user_id) 
	  REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS tokens (
    token_hash bytea NOT NULL,
    user_id bigint NOT NULL,
    expiry timestamp without time zone DEFAULT NOW(),
    scope text NOT NULL,
    CONSTRAINT fk_user_token
      FOREIGN KEY(user_id) 
	  REFERENCES users(id)
);

DROP TABLE urls; 
DROP TABLE tokens;
DROP TABLE users; 