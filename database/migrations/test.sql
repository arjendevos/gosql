CREATE TABLE account (
		id SERIAL NOT NULL,
		name VARCHAR(255) NOT NULL,
		profile_picture_url TEXT NOT NULL,
		email VARCHAR(255) NOT NULL,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
		CONSTRAINT account_ak_1 UNIQUE (id) NOT DEFERRABLE INITIALLY IMMEDIATE,
		CONSTRAINT account_pk PRIMARY KEY (id),
		CONSTRAINT account_ak_2 UNIQUE (email) NOT DEFERRABLE INITIALLY IMMEDIATE
);
CREATE INDEX account_idx_1 ON account (id);
CREATE INDEX account_idx_2 ON account (email);

CREATE TABLE proxy (
		id SERIAL NOT NULL,
		host VARCHAR(255) NOT NULL,
		scheme VARCHAR(5) NOT NULL DEFAULT 'http',
		port INTEGER NOT NULL DEFAULT 8080,
		username VARCHAR(255) NOT NULL DEFAULT 'PRIVATEVPN_USERNAME',
		password VARCHAR(255) NOT NULL DEFAULT 'PRIVATEVPN_PASSWORD',
		is_custom_provider BOOLEAN NOT NULL DEFAULT false,
		account_id INTEGER NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
		CONSTRAINT proxy_ak_1 UNIQUE (id) NOT DEFERRABLE INITIALLY IMMEDIATE,
		CONSTRAINT proxy_pk PRIMARY KEY (id),
		CONSTRAINT proxy_ak_2 UNIQUE (host) NOT DEFERRABLE INITIALLY IMMEDIATE,
		CONSTRAINT proxy_account_id_fk_1 FOREIGN KEY (account_id) REFERENCES account (id) NOT DEFERRABLE INITIALLY IMMEDIATE
);
CREATE INDEX proxy_idx_1 ON proxy (id);
CREATE INDEX proxy_idx_2 ON proxy (account_id);

CREATE TABLE credential (
		id SERIAL NOT NULL,
		csfr_token VARCHAR(255) NOT NULL,
		user_agent VARCHAR(255) NOT NULL,
		app_id VARCHAR(255) NOT NULL,
		cookie TEXT NOT NULL,
		instagram_id VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
		CONSTRAINT credential_ak_1 UNIQUE (id) NOT DEFERRABLE INITIALLY IMMEDIATE,
		CONSTRAINT credential_pk PRIMARY KEY (id),
		CONSTRAINT credential_ak_2 UNIQUE (instagram_id) NOT DEFERRABLE INITIALLY IMMEDIATE
);
CREATE INDEX credential_idx_1 ON credential (id);
CREATE INDEX credential_idx_2 ON credential (instagram_id);

CREATE TABLE job (
		id SERIAL NOT NULL,
		type VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
		CONSTRAINT job_ak_1 UNIQUE (id) NOT DEFERRABLE INITIALLY IMMEDIATE,
		CONSTRAINT job_pk PRIMARY KEY (id)
);
CREATE INDEX job_idx_1 ON job (id);
CREATE INDEX job_idx_2 ON job (type);

CREATE TABLE instagram_account (
		id SERIAL NOT NULL,
		username VARCHAR(255) NOT NULL,
		password VARCHAR(255) NOT NULL,
		country VARCHAR(255) NOT NULL,
		active BOOLEAN NOT NULL DEFAULT false,
		total_actions_done INTEGER NOT NULL DEFAULT 0,
		needs_to_be_checked BOOLEAN NOT NULL DEFAULT false,
		login_failed BOOLEAN NOT NULL DEFAULT false,
		is_setup BOOLEAN NOT NULL DEFAULT false,
		setup_target_account_username VARCHAR(255) NOT NULL,
		proxy_id INTEGER NULL,
		credential_id INTEGER NULL,
		account_id INTEGER NULL,
		job_id INTEGER NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
		CONSTRAINT instagram_account_ak_1 UNIQUE (id) NOT DEFERRABLE INITIALLY IMMEDIATE,
		CONSTRAINT instagram_account_pk PRIMARY KEY (id),
		CONSTRAINT instagram_account_proxy_id_fk_1 FOREIGN KEY (proxy_id) REFERENCES proxy (id) NOT DEFERRABLE INITIALLY IMMEDIATE,
		CONSTRAINT instagram_account_credential_id_fk_2 FOREIGN KEY (credential_id) REFERENCES credential (id) NOT DEFERRABLE INITIALLY IMMEDIATE,
		CONSTRAINT instagram_account_account_id_fk_3 FOREIGN KEY (account_id) REFERENCES account (id) NOT DEFERRABLE INITIALLY IMMEDIATE,
		CONSTRAINT instagram_account_job_id_fk_4 FOREIGN KEY (job_id) REFERENCES job (id) NOT DEFERRABLE INITIALLY IMMEDIATE
);
CREATE INDEX instagram_account_idx_1 ON instagram_account (id);
CREATE INDEX instagram_account_idx_2 ON instagram_account (proxy_id);
CREATE INDEX instagram_account_idx_3 ON instagram_account (credential_id);
CREATE INDEX instagram_account_idx_4 ON instagram_account (account_id);
CREATE INDEX instagram_account_idx_5 ON instagram_account (job_id);

CREATE TABLE scraped_instagram_account (
		id SERIAL NOT NULL,
		instagram_id VARCHAR(255) NOT NULL,
		username VARCHAR(255) NOT NULL,
		full_name VARCHAR(255) NOT NULL,
		biography TEXT NOT NULL,
		facebook_id VARCHAR(255) NOT NULL,
		is_verified BOOLEAN NOT NULL,
		is_private BOOLEAN NOT NULL,
		profile_picture_url TEXT NOT NULL,
		external_url TEXT NOT NULL,
		business_category_name VARCHAR(255) NOT NULL,
		category_name VARCHAR(255) NOT NULL,
		is_business_account BOOLEAN NOT NULL,
		is_professional_account BOOLEAN NOT NULL,
		followed_by_count INTEGER NOT NULL,
		follow_count INTEGER NOT NULL,
		media_count INTEGER NOT NULL,
		highlight_count INTEGER NOT NULL,
		has_reels BOOLEAN NOT NULL,
		extracted_email VARCHAR(255) NOT NULL,
		has_management BOOLEAN NOT NULL,
		account_id INTEGER NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
		CONSTRAINT scraped_instagram_account_ak_1 UNIQUE (id) NOT DEFERRABLE INITIALLY IMMEDIATE,
		CONSTRAINT scraped_instagram_account_pk PRIMARY KEY (id),
		CONSTRAINT scraped_instagram_account_account_id_fk_1 FOREIGN KEY (account_id) REFERENCES account (id) NOT DEFERRABLE INITIALLY IMMEDIATE
);
CREATE INDEX scraped_instagram_account_idx_1 ON scraped_instagram_account (id);
CREATE INDEX scraped_instagram_account_idx_2 ON scraped_instagram_account (account_id);

CREATE TABLE scraped_instagram_account_chained (
		id SERIAL NOT NULL,
		instagram_id VARCHAR(255) NOT NULL,
		username VARCHAR(255) NOT NULL,
		full_name VARCHAR(255) NOT NULL,
		profile_picture_url TEXT NOT NULL,
		scraped_instagram_account_id INTEGER NOT NULL,
		account_id INTEGER NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
		CONSTRAINT scraped_instagram_account_chained_ak_1 UNIQUE (id) NOT DEFERRABLE INITIALLY IMMEDIATE,
		CONSTRAINT scraped_instagram_account_chained_pk PRIMARY KEY (id),
		CONSTRAINT scraped_instagram_account_chained_scraped_instagram_account_id_fk_1 FOREIGN KEY (scraped_instagram_account_id) REFERENCES scraped_instagram_account (id) NOT DEFERRABLE INITIALLY IMMEDIATE,
		CONSTRAINT scraped_instagram_account_chained_account_id_fk_2 FOREIGN KEY (account_id) REFERENCES account (id) NOT DEFERRABLE INITIALLY IMMEDIATE
);
CREATE INDEX scraped_instagram_account_chained_idx_1 ON scraped_instagram_account_chained (id);
CREATE INDEX scraped_instagram_account_chained_idx_2 ON scraped_instagram_account_chained (scraped_instagram_account_id);
CREATE INDEX scraped_instagram_account_chained_idx_3 ON scraped_instagram_account_chained (account_id);

