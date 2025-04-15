CREATE TABLE IF NOT EXISTS locations(
    id bigserial,
    city character varying(255) NOT NULL,
    state character varying(255) NOT NULL,
    country character varying(255) NOT NULL,
    code character varying(3),
    lat numeric(9, 6),
    lng numeric(9, 6),
    CONSTRAINT locations_pkey PRIMARY KEY (id)
);

