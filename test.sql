CREATE TABLE proxy (
		id SERIAL NOT NULL,
		host VARCHAR(255) NOT NULL,
		scheme VARCHAR(5) NOT NULL DEFAULT 'http',
		port INTEGER NOT NULL DEFAULT 8080,
		username VARCHAR(255) NOT NULL DEFAULT 'PRIVATEVPN_USERNAME',
		password VARCHAR(255) NOT NULL DEFAULT 'PRIVATEVPN_PASSWORD',
		is_custom_provider BOOLEAN NOT NULL DEFAULT false,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
		CONSTRAINT proxy_ak_1 UNIQUE (id) NOT DEFERRABLE INITIALLY IMMEDIATE,
		CONSTRAINT proxy_pk PRIMARY KEY (id),
		CONSTRAINT proxy_ak_2 UNIQUE (host) NOT DEFERRABLE INITIALLY IMMEDIATE,
		CONSTRAINT proxy_ak_3 UNIQUE (username) NOT DEFERRABLE INITIALLY IMMEDIATE,
		CONSTRAINT proxy_ak_4 UNIQUE (password) NOT DEFERRABLE INITIALLY IMMEDIATE
);
CREATE INDEX proxy_idx_1 ON proxy (id);

CREATE TABLE credential (
		id SERIAL NOT NULL,
		csfr_token VARCHAR(255) NOT NULL,
		user_agent VARCHAR(255) NOT NULL,
		app_id VARCHAR(255) NOT NULL,
		cookie TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
		CONSTRAINT credential_ak_1 UNIQUE (id) NOT DEFERRABLE INITIALLY IMMEDIATE,
		CONSTRAINT credential_pk PRIMARY KEY (id)
);
CREATE INDEX credential_idx_1 ON credential (id);

CREATE TABLE instagram_account (
		id SERIAL NOT NULL,
		instagram_id VARCHAR(255) NOT NULL,
		username VARCHAR(255) NOT NULL,
		password VARCHAR(255) NOT NULL,
		country VARCHAR(255) NOT NULL,
		active BOOLEAN NOT NULL DEFAULT false,
		total_actions_done INTEGER NOT NULL DEFAULT 0,
		needs_to_be_checked BOOLEAN NOT NULL DEFAULT false,
		proxy_id SERIAL NOT NULL,
		credential_id SERIAL NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
		CONSTRAINT instagram_account_ak_1 UNIQUE (id) NOT DEFERRABLE INITIALLY IMMEDIATE,
		CONSTRAINT instagram_account_pk PRIMARY KEY (id),
		CONSTRAINT instagram_account_ak_2 UNIQUE (instagram_id) NOT DEFERRABLE INITIALLY IMMEDIATE,
		CONSTRAINT instagram_account_ak_3 UNIQUE (username) NOT DEFERRABLE INITIALLY IMMEDIATE,
		CONSTRAINT instagram_account_proxy_id_fk_1 FOREIGN KEY (proxy_id) REFERENCES proxy (id) NOT DEFERRABLE INITIALLY IMMEDIATE,
		CONSTRAINT instagram_account_credential_id_fk_2 FOREIGN KEY (credential_id) REFERENCES credential (id) NOT DEFERRABLE INITIALLY IMMEDIATE
);
CREATE INDEX instagram_account_idx_1 ON instagram_account (id);
CREATE INDEX instagram_account_idx_2 ON instagram_account (instagram_id);
CREATE INDEX instagram_account_idx_3 ON instagram_account (username);
CREATE INDEX instagram_account_idx_4 ON instagram_account (proxy_id);
CREATE INDEX instagram_account_idx_5 ON instagram_account (credential_id);

