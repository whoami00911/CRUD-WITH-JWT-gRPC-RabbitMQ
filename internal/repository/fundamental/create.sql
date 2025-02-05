-- Table: public.AbuseEntity

-- DROP TABLE IF EXISTS public."AbuseEntity";

CREATE TABLE IF NOT EXISTS public."AbuseEntity"
(
    "ipAddress" character varying(15) COLLATE pg_catalog."default",
    "isPublic" boolean,
    "ipVersion" integer,
    "isWhitelisted" boolean,
    "abuseConfidenceScore" integer,
    "countryCode" character varying(5) COLLATE pg_catalog."default",
    "countryName" character varying(100) COLLATE pg_catalog."default",
    "usageType" character varying(100) COLLATE pg_catalog."default",
    isp character varying(200) COLLATE pg_catalog."default",
    "isTor" boolean,
    "isFromDB" boolean,
    CONSTRAINT "unique_ipAddress" UNIQUE ("ipAddress")
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public."AbuseEntity"
    OWNER to postgres;


-- Table: public.refresh_tokens

-- DROP TABLE IF EXISTS public.refresh_tokens;

CREATE TABLE IF NOT EXISTS public.refresh_tokens
(
    id integer NOT NULL DEFAULT nextval('refresh_tokens_id_seq'::regclass),
    "userId" integer NOT NULL,
    "refreshToken" character varying(255) COLLATE pg_catalog."default" NOT NULL,
    "expiresAt" timestamp without time zone NOT NULL,
    CONSTRAINT "refresh_tokens_refreshToken_key" UNIQUE ("refreshToken"),
    CONSTRAINT "refresh_tokens_userId_fkey" FOREIGN KEY ("userId")
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.refresh_tokens
    OWNER to postgres;


-- Table users
-- Table: public.users

-- DROP TABLE IF EXISTS public.users;
CREATE TABLE IF NOT EXISTS public.users
(
    id integer NOT NULL DEFAULT nextval('users_id_seq'::regclass),
    name character varying(255) COLLATE pg_catalog."default" NOT NULL,
    username character varying(255) COLLATE pg_catalog."default" NOT NULL,
    password_hash character varying(255) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT users_id_key UNIQUE (id),
    CONSTRAINT users_username_key UNIQUE (username)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.users
    OWNER to postgres;
