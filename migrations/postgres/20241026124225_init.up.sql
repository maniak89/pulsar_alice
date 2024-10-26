CREATE TYPE log_level AS ENUM
    ('Error', 'Info');

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS meters
(
    id uuid NOT NULL default uuid_generate_v4(),
    user_id uuid NOT NULL,
    name character varying(45) NOT NULL,
    address character varying(45) NOT NULL,
    serail_number character varying(8) NOT NULL,
    period_check bigint NOT NULL,
    is_cold BOOLEAN NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NOT NULL DEFAULT now(),
    PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS user_id_idx
    ON meters USING btree
        (user_id ASC NULLS LAST);

CREATE TABLE IF NOT EXISTS logs (
    id uuid NOT NULL default uuid_generate_v4(),
    meter_id uuid NOT NULL,
    time timestamp without time zone NOT NULL DEFAULT now(),
    level log_level NOT NULL,
    message text NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (meter_id) REFERENCES meters(id) ON UPDATE CASCADE ON DELETE CASCADE
);