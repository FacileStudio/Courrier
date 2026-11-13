-- +goose Up
-- baseline: the production schema as GORM AutoMigrate and schemas.MigrateWithIssuer built it, captured 2026-11-13 with pg_dump --schema-only --no-owner --no-privileges --no-comments --no-tablespaces --schema=public, public. stripped.
-- pg_trgm is pinned to public and its operator class stays qualified: an extension lands in the first schema on the search_path, and IF NOT EXISTS then no-ops for every other schema of the same database, which is the layout tronc/testdb uses.

CREATE EXTENSION IF NOT EXISTS pg_trgm WITH SCHEMA public;

CREATE TABLE accounts (
    id bigint NOT NULL,
    user_id bigint,
    name text,
    email text,
    imap_host text,
    imap_port bigint DEFAULT 993,
    imap_user text,
    imap_password text,
    smtp_host text,
    smtp_port bigint DEFAULT 587,
    smtp_user text,
    smtp_password text,
    signature text,
    is_default boolean DEFAULT false,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    space_id text
);

CREATE SEQUENCE accounts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE accounts_id_seq OWNED BY accounts.id;

CREATE TABLE app_settings (
    id bigint NOT NULL,
    encryption_key text DEFAULT ''::text NOT NULL
);

CREATE SEQUENCE app_settings_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE app_settings_id_seq OWNED BY app_settings.id;

CREATE TABLE attachments (
    id bigint NOT NULL,
    email_id bigint,
    filename text,
    mime_type text,
    size bigint,
    cid text,
    part_id text
);

CREATE SEQUENCE attachments_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE attachments_id_seq OWNED BY attachments.id;

CREATE TABLE contacts (
    id bigint NOT NULL,
    account_id bigint,
    space_id text,
    email text,
    name text,
    sent_count bigint DEFAULT 0,
    last_sent_at timestamp with time zone
);

CREATE SEQUENCE contacts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE contacts_id_seq OWNED BY contacts.id;

CREATE TABLE directory_deletion_state (
    id smallint NOT NULL,
    state bytea NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT directory_deletion_state_id_check CHECK ((id = 1))
);

CREATE TABLE email_templates (
    id bigint NOT NULL,
    user_id bigint,
    name text,
    subject text,
    body_html text,
    body_text text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    space_id text
);

CREATE SEQUENCE email_templates_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE email_templates_id_seq OWNED BY email_templates.id;

CREATE TABLE emails (
    id bigint NOT NULL,
    account_id bigint,
    folder_id bigint,
    message_id text,
    thread_id text,
    subject text,
    from_address text,
    from_name text,
    to_addresses text,
    cc_addresses text,
    bcc_addresses text,
    date timestamp with time zone,
    body_text text,
    body_html text,
    is_read boolean DEFAULT false,
    is_starred boolean DEFAULT false,
    has_attachments boolean DEFAULT false,
    in_reply_to text,
    "references" text,
    imap_uid bigint,
    facile_id text,
    cached_at timestamp with time zone,
    space_id text,
    category text DEFAULT ''::text,
    bulk_reason text DEFAULT ''::text
);

CREATE SEQUENCE emails_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE emails_id_seq OWNED BY emails.id;

CREATE TABLE folders (
    id bigint NOT NULL,
    account_id bigint,
    path text,
    name text,
    type text DEFAULT 'custom'::text,
    unread_count bigint DEFAULT 0,
    total_count bigint DEFAULT 0,
    uid_validity bigint DEFAULT 0,
    space_id text
);

CREATE SEQUENCE folders_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE folders_id_seq OWNED BY folders.id;

CREATE TABLE forfait_state (
    id smallint DEFAULT 1 NOT NULL,
    document text NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT forfait_state_id_check CHECK ((id = 1))
);

CREATE TABLE mcp_actions (
    id text NOT NULL,
    user_id text NOT NULL,
    kind text NOT NULL,
    target text NOT NULL,
    digest text NOT NULL,
    summary text DEFAULT ''::text NOT NULL,
    client_id text DEFAULT ''::text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    resolved_at timestamp with time zone,
    outcome text DEFAULT ''::text NOT NULL
);

CREATE TABLE mcp_pending_actions (
    id text NOT NULL,
    user_id bigint NOT NULL,
    account_id bigint NOT NULL,
    email_id bigint NOT NULL,
    kind text NOT NULL,
    digest text NOT NULL,
    client_id text DEFAULT ''::text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    resolved_at timestamp with time zone,
    outcome text DEFAULT ''::text NOT NULL
);

CREATE TABLE porte_identities (
    user_id bigint NOT NULL,
    provider text NOT NULL,
    subject text NOT NULL,
    password_hash text DEFAULT ''::text NOT NULL,
    access_token text DEFAULT ''::text NOT NULL,
    refresh_token text DEFAULT ''::text NOT NULL,
    token_expiry timestamp with time zone,
    roles jsonb,
    roles_synced_at timestamp with time zone,
    synced_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now()
);

CREATE TABLE porte_login_codes (
    code_hash text NOT NULL,
    user_id bigint NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    consumed_at timestamp with time zone,
    provider text DEFAULT ''::text NOT NULL,
    subject text DEFAULT ''::text NOT NULL
);

CREATE TABLE porte_sessions (
    id bigint NOT NULL,
    token_hash text NOT NULL,
    user_id bigint NOT NULL,
    label text DEFAULT ''::text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    last_used_at timestamp with time zone DEFAULT now() NOT NULL,
    expires_at timestamp with time zone
);

CREATE SEQUENCE porte_sessions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE porte_sessions_id_seq OWNED BY porte_sessions.id;

CREATE TABLE push_devices (
    token text NOT NULL,
    user_id bigint NOT NULL,
    kind text NOT NULL,
    p256dh text DEFAULT ''::text NOT NULL,
    auth text DEFAULT ''::text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);

CREATE TABLE sender_rules (
    id bigint NOT NULL,
    account_id bigint NOT NULL,
    address text NOT NULL,
    category text DEFAULT ''::text NOT NULL,
    updated_at timestamp with time zone
);

CREATE SEQUENCE sender_rules_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE sender_rules_id_seq OWNED BY sender_rules.id;

CREATE TABLE space_members (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    space_id uuid NOT NULL,
    user_id bigint NOT NULL,
    role text DEFAULT 'member'::text NOT NULL,
    joined_at timestamp with time zone
);

CREATE TABLE spaces (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name text NOT NULL,
    description text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    registre_space_id text
);

CREATE TABLE thread_links (
    account_id bigint NOT NULL,
    message_id text NOT NULL,
    thread_id text
);

CREATE TABLE users (
    id bigint NOT NULL,
    email text,
    name text,
    avatar_url text,
    avatar_source text,
    oidc_picture_url text,
    oidc_access_token text,
    oidc_refresh_token text,
    oidc_token_expiry timestamp with time zone,
    profile_synced_at timestamp with time zone,
    password_hash text,
    created_at timestamp with time zone,
    oidc_subject text,
    avatar_upload_path text,
    locale text DEFAULT 'en'::text NOT NULL
);

CREATE SEQUENCE users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE users_id_seq OWNED BY users.id;

ALTER TABLE ONLY accounts ALTER COLUMN id SET DEFAULT nextval('accounts_id_seq'::regclass);

ALTER TABLE ONLY app_settings ALTER COLUMN id SET DEFAULT nextval('app_settings_id_seq'::regclass);

ALTER TABLE ONLY attachments ALTER COLUMN id SET DEFAULT nextval('attachments_id_seq'::regclass);

ALTER TABLE ONLY contacts ALTER COLUMN id SET DEFAULT nextval('contacts_id_seq'::regclass);

ALTER TABLE ONLY email_templates ALTER COLUMN id SET DEFAULT nextval('email_templates_id_seq'::regclass);

ALTER TABLE ONLY emails ALTER COLUMN id SET DEFAULT nextval('emails_id_seq'::regclass);

ALTER TABLE ONLY folders ALTER COLUMN id SET DEFAULT nextval('folders_id_seq'::regclass);

ALTER TABLE ONLY porte_sessions ALTER COLUMN id SET DEFAULT nextval('porte_sessions_id_seq'::regclass);

ALTER TABLE ONLY sender_rules ALTER COLUMN id SET DEFAULT nextval('sender_rules_id_seq'::regclass);

ALTER TABLE ONLY users ALTER COLUMN id SET DEFAULT nextval('users_id_seq'::regclass);

ALTER TABLE ONLY accounts
    ADD CONSTRAINT accounts_pkey PRIMARY KEY (id);

ALTER TABLE ONLY app_settings
    ADD CONSTRAINT app_settings_pkey PRIMARY KEY (id);

ALTER TABLE ONLY attachments
    ADD CONSTRAINT attachments_pkey PRIMARY KEY (id);

ALTER TABLE ONLY contacts
    ADD CONSTRAINT contacts_pkey PRIMARY KEY (id);

ALTER TABLE ONLY directory_deletion_state
    ADD CONSTRAINT directory_deletion_state_pkey PRIMARY KEY (id);

ALTER TABLE ONLY email_templates
    ADD CONSTRAINT email_templates_pkey PRIMARY KEY (id);

ALTER TABLE ONLY emails
    ADD CONSTRAINT emails_pkey PRIMARY KEY (id);

ALTER TABLE ONLY folders
    ADD CONSTRAINT folders_pkey PRIMARY KEY (id);

ALTER TABLE ONLY forfait_state
    ADD CONSTRAINT forfait_state_pkey PRIMARY KEY (id);

ALTER TABLE ONLY mcp_actions
    ADD CONSTRAINT mcp_actions_pkey PRIMARY KEY (id);

ALTER TABLE ONLY mcp_pending_actions
    ADD CONSTRAINT mcp_pending_actions_pkey PRIMARY KEY (id);

ALTER TABLE ONLY porte_identities
    ADD CONSTRAINT porte_identities_pkey PRIMARY KEY (provider, subject);

ALTER TABLE ONLY porte_login_codes
    ADD CONSTRAINT porte_login_codes_pkey PRIMARY KEY (code_hash);

ALTER TABLE ONLY porte_sessions
    ADD CONSTRAINT porte_sessions_pkey PRIMARY KEY (id);

ALTER TABLE ONLY porte_sessions
    ADD CONSTRAINT porte_sessions_token_hash_key UNIQUE (token_hash);

ALTER TABLE ONLY push_devices
    ADD CONSTRAINT push_devices_pkey PRIMARY KEY (token);

ALTER TABLE ONLY sender_rules
    ADD CONSTRAINT sender_rules_pkey PRIMARY KEY (id);

ALTER TABLE ONLY space_members
    ADD CONSTRAINT space_members_pkey PRIMARY KEY (id);

ALTER TABLE ONLY spaces
    ADD CONSTRAINT spaces_pkey PRIMARY KEY (id);

ALTER TABLE ONLY thread_links
    ADD CONSTRAINT thread_links_pkey PRIMARY KEY (account_id, message_id);

ALTER TABLE ONLY users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);

CREATE INDEX idx_accounts_space_id ON accounts USING btree (space_id);

CREATE INDEX idx_accounts_user_id ON accounts USING btree (user_id);

CREATE INDEX idx_attachments_email_id ON attachments USING btree (email_id);

CREATE UNIQUE INDEX idx_contacts_account_email ON contacts USING btree (account_id, email);

CREATE INDEX idx_contacts_last_sent_at ON contacts USING btree (last_sent_at);

CREATE INDEX idx_contacts_space_id ON contacts USING btree (space_id);

CREATE INDEX idx_email_templates_space_id ON email_templates USING btree (space_id);

CREATE INDEX idx_email_templates_user_id ON email_templates USING btree (user_id);

CREATE INDEX idx_emails_account_id ON emails USING btree (account_id);

CREATE INDEX idx_emails_account_thread ON emails USING btree (account_id, thread_id);

CREATE INDEX idx_emails_body_text_trgm ON emails USING gin (body_text public.gin_trgm_ops);

CREATE INDEX idx_emails_date ON emails USING btree (date);

CREATE UNIQUE INDEX idx_emails_facile_id ON emails USING btree (facile_id);

CREATE INDEX idx_emails_folder_id ON emails USING btree (folder_id);

CREATE INDEX idx_emails_folder_thread_date ON emails USING btree (account_id, folder_id, thread_id, date DESC, id DESC);

CREATE INDEX idx_emails_from_address_trgm ON emails USING gin (from_address public.gin_trgm_ops);

CREATE INDEX idx_emails_from_name_trgm ON emails USING gin (from_name public.gin_trgm_ops);

CREATE INDEX idx_emails_message_id ON emails USING btree (message_id);

CREATE INDEX idx_emails_space_id ON emails USING btree (space_id);

CREATE INDEX idx_emails_subject_trgm ON emails USING gin (subject public.gin_trgm_ops);

CREATE INDEX idx_emails_thread_id ON emails USING btree (thread_id);

CREATE INDEX idx_folders_account_id ON folders USING btree (account_id);

CREATE INDEX idx_folders_space_id ON folders USING btree (space_id);

CREATE UNIQUE INDEX idx_sender_rules_account_address ON sender_rules USING btree (account_id, address);

CREATE INDEX idx_space_members_space_id ON space_members USING btree (space_id);

CREATE INDEX idx_space_members_user_id ON space_members USING btree (user_id);

CREATE UNIQUE INDEX idx_spaces_registre_space_id ON spaces USING btree (registre_space_id);

CREATE INDEX idx_thread_links_account_thread ON thread_links USING btree (account_id, thread_id);

CREATE INDEX idx_thread_links_thread_id ON thread_links USING btree (thread_id);

CREATE UNIQUE INDEX idx_users_email ON users USING btree (email);

CREATE UNIQUE INDEX idx_users_o_id_c_subject ON users USING btree (oidc_subject);

CREATE INDEX mcp_actions_open_idx ON mcp_actions USING btree (user_id, created_at DESC) WHERE (resolved_at IS NULL);

CREATE INDEX mcp_pending_actions_open_idx ON mcp_pending_actions USING btree (user_id, created_at DESC) WHERE (resolved_at IS NULL);

CREATE INDEX porte_identities_user_idx ON porte_identities USING btree (user_id);

CREATE INDEX porte_sessions_expiry_idx ON porte_sessions USING btree (expires_at) WHERE (expires_at IS NOT NULL);

CREATE INDEX porte_sessions_user_idx ON porte_sessions USING btree (user_id);

CREATE INDEX push_devices_user_idx ON push_devices USING btree (user_id);

ALTER TABLE ONLY sender_rules
    ADD CONSTRAINT fk_sender_rules_account FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE;

ALTER TABLE ONLY space_members
    ADD CONSTRAINT fk_space_members_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE ONLY space_members
    ADD CONSTRAINT fk_spaces_members FOREIGN KEY (space_id) REFERENCES spaces(id);

ALTER TABLE ONLY mcp_pending_actions
    ADD CONSTRAINT mcp_pending_actions_account_id_fkey FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE;

ALTER TABLE ONLY mcp_pending_actions
    ADD CONSTRAINT mcp_pending_actions_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE ONLY porte_identities
    ADD CONSTRAINT porte_identities_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE ONLY porte_login_codes
    ADD CONSTRAINT porte_login_codes_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE ONLY porte_sessions
    ADD CONSTRAINT porte_sessions_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE ONLY push_devices
    ADD CONSTRAINT push_devices_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- +goose Down
SELECT 'baseline migration is not reversible' WHERE false;
