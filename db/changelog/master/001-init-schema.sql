SET TIME ZONE 'Europe/Moscow';

-- Table spaces
ALTER TABLESPACE pg_global
    OWNER TO postgres;
ALTER TABLESPACE pg_default
    OWNER TO postgres;

-- Records
create table if not exists records (
    nickname varchar(255) not null unique,
    score bigint not null
) tablespace pg_default;