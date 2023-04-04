DROP TABLE IF EXISTS users_roles;
DROP TABLE IF EXISTS roles_permissions;
DROP TABLE IF EXISTS features; -- Only if using timescaledb
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS permissions;

-- Users
CREATE TABLE users (
	id VARCHAR ( 50 ) PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP,
	deleted_at TIMESTAMP,

	username VARCHAR ( 50 ) UNIQUE NOT NULL,
	email VARCHAR ( 255 ) UNIQUE NOT NULL,
    password VARCHAR ( 255 ) NOT NULL,

	nick VARCHAR ( 50 ) NOT NULL, -- Default=<username>
	first_name VARCHAR ( 50 ),
	last_name VARCHAR ( 50 ),
	phone VARCHAR ( 50 ),
    address VARCHAR ( 255 ),
    
    last_login TIMESTAMP,
    last_password_change TIMESTAMP NOT NULL DEFAULT NOW(),
    verified_email BOOLEAN NOT NULL DEFAULT FALSE,
    verified_phone BOOLEAN NOT NULL DEFAULT FALSE,
    ban_date TIMESTAMP,
    ban_expire TIMESTAMP,

    -- Extra fields
    photo_name VARCHAR ( 50 ) UNIQUE,
    cv_name VARCHAR ( 50 ) UNIQUE
);

-- Roles
CREATE TABLE roles (
   id serial PRIMARY KEY,
   name VARCHAR (50) UNIQUE NOT NULL
);
INSERT INTO roles (name) VALUES ('superadmin'); -- Insert superadmin as default role

-- Permissions
CREATE TABLE permissions (
   id serial PRIMARY KEY,
   resource VARCHAR (50) NOT NULL,
   operation VARCHAR (50) NOT NULL,
   UNIQUE (resource, operation)
);

-- Users-Roles
CREATE TABLE users_roles (
  user_id VARCHAR ( 50 ) NOT NULL,
  role_id INT NOT NULL,
  grant_date TIMESTAMP NOT NULL DEFAULT NOW(),
  PRIMARY KEY (user_id, role_id),
  FOREIGN KEY (role_id)
      REFERENCES roles (id) ON DELETE CASCADE,
  FOREIGN KEY (user_id)
      REFERENCES users (id) ON DELETE CASCADE
);

-- Roles-Permissions
CREATE TABLE roles_permissions (
  role_id INT NOT NULL,
  permission_id INT NOT NULL,
  grant_date TIMESTAMP NOT NULL DEFAULT NOW(),
  PRIMARY KEY (role_id, permission_id),
  FOREIGN KEY (permission_id)
      REFERENCES permissions (id) ON DELETE CASCADE,
  FOREIGN KEY (role_id)
      REFERENCES roles (id) ON DELETE CASCADE
);

-- ONLY IF USING TIMESCALEDB
CREATE EXTENSION IF NOT EXISTS postgis;
-- Features (1:M)
CREATE TABLE features (
  geom geometry(Point, 4326) NOT NULL,
  timestamp TIMESTAMP NOT NULL DEFAULT NOW(),
  user_id VARCHAR ( 50 ) NOT NULL,
  CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id)
);
SELECT UpdateGeometrySRID('public', 'features', 'geom', '4326') ;
SELECT create_hypertable('features','timestamp');
CREATE INDEX ix_geom_timestamp ON features (geom, timestamp DESC); -- efficient queries on geom and timestamp
ALTER TABLE features SET ( -- enables data compression
  timescaledb.compress,
  timescaledb.compress_orderby = 'timestamp DESC',
  timescaledb.compress_segmentby = 'user_id'
);
SELECT add_compression_policy('features', INTERVAL '2 weeks'); -- add compression policy for 2 weeks old data
-- MapLibre Martin function
CREATE OR REPLACE FUNCTION function_zxy_query(z integer, x integer, y integer, query_params json) RETURNS bytea AS $$
DECLARE
  mvt bytea;
BEGIN
  IF json_typeof(query_params->'fromDate') = 'string' AND json_typeof(query_params->'toDate') = 'string' THEN -- if fromDate and toDate params exist
    SELECT INTO mvt ST_AsMVT(tile, 'function_zxy_query', 4096, 'geom') FROM (
      SELECT
        ST_AsMVTGeom(ST_Transform(ST_CurveToLine(geom), 3857), ST_TileEnvelope(z, x, y), 4096, 64, true) AS geom
      FROM features
      WHERE geom && ST_Transform(ST_TileEnvelope(z, x, y), 4326)
      AND user_id = (query_params->>'userId')::varchar -- cast varchar
      AND timestamp >= (query_params->>'fromDate')::timestamp -- cast timestamp
      AND timestamp <= (query_params->>'toDate')::timestamp -- cast timestamp
    ) as tile WHERE geom IS NOT NULL;
  ELSE -- if fromDate and toDate do not exist, get most recent records
    SELECT INTO mvt ST_AsMVT(tile, 'function_zxy_query', 4096, 'geom') FROM (
      SELECT
        ST_AsMVTGeom(ST_Transform(ST_CurveToLine(geom), 3857), ST_TileEnvelope(z, x, y), 4096, 64, true) AS geom
      FROM features
      WHERE geom && ST_Transform(ST_TileEnvelope(z, x, y), 4326)
      AND user_id = (query_params->>'userId')::varchar -- cast varchar
      ORDER BY timestamp DESC 
      LIMIT 1
    ) as tile WHERE geom IS NOT NULL;
  END IF;

  RETURN mvt;
END
$$ LANGUAGE plpgsql IMMUTABLE STRICT PARALLEL SAFE;
