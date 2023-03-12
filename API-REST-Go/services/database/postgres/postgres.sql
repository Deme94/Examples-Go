DROP TABLE IF EXISTS users_roles;
DROP TABLE IF EXISTS roles_permissions;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS permissions;

-- Users
CREATE TABLE users (
	id serial PRIMARY KEY,
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
    verified_mail BOOLEAN NOT NULL DEFAULT FALSE,
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
  user_id INT NOT NULL,
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

