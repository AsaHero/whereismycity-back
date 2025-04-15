CREATE TABLE IF NOT EXISTS location_alternate_names(
    id bigserial PRIMARY KEY,
    location_id bigint NOT NULL,
    geoname_id bigint NOT NULL,
    alternate_name_id bigint NOT NULL,
    type varchar(50) NOT NULL,
    iso_language_code varchar(10),
    alternate_name text NOT NULL,
    is_preferred boolean DEFAULT FALSE,
    is_short boolean DEFAULT FALSE,
    is_colloquial boolean DEFAULT FALSE,
    is_historic boolean DEFAULT FALSE,
    created_at timestamp with time zone DEFAULT now(),
    FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_location_alternate_names_location_id ON location_alternate_names(location_id);

CREATE INDEX IF NOT EXISTS idx_location_alternate_names_geoname_id ON location_alternate_names(geoname_id);

CREATE INDEX IF NOT EXISTS idx_location_alternate_names_alternate_name_id ON location_alternate_names(alternate_name_id);

CREATE INDEX IF NOT EXISTS idx_location_alternate_names_type ON location_alternate_names(type);

