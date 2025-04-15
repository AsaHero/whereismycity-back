CREATE TABLE IF NOT EXISTS location_geoname_ids(
    location_id bigint NOT NULL,
    geoname_id bigint NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT location_geoname_ids_pkey PRIMARY KEY (location_id, geoname_id),
    CONSTRAINT fk_location_geoname_ids_location_id FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_location_geoname_ids_location_id ON location_geoname_ids(location_id);

CREATE INDEX IF NOT EXISTS idx_location_geoname_ids_geoname_id ON location_geoname_ids(geoname_id);

