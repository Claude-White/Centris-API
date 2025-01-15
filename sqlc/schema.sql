CREATE TABLE IF NOT EXISTS property (
    id bigint NOT NULL PRIMARY KEY,
    title varchar(500) NOT NULL,
    category varchar(500) NOT NULL,
    address varchar(500) NOT NULL,
    city_name varchar(500) NOT NULL,
    neighbourhood_name varchar(500),
    price numeric(15,2) NOT NULL,
    description text,
    bedroom_number integer,
    room_number integer,
    bathroom_number integer,
    latitude decimal(9,6) NOT NULL,
    longitude decimal(9,6) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON COLUMN property.id IS 'MLS number';

CREATE TABLE IF NOT EXISTS broker (
    id bigint NOT NULL PRIMARY KEY,
    first_name varchar(100) NOT NULL,
    middle_name varchar(100),
    last_name varchar(100) NOT NULL,
    title varchar(200) NOT NULL,
    profile_photo varchar(500),
    complementary_info text,
    served_areas text,
    presentation text,
    corporation_name varchar(200),
    agency_name varchar(200) NOT NULL,
    agency_address text NOT NULL,
    agency_logo varchar(500),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS property_photo (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    property_id bigint NOT NULL,
    link varchar(500) NOT NULL,
    description text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_property_photo_property FOREIGN KEY (property_id)
        REFERENCES property(id) ON DELETE CASCADE
);
CREATE INDEX idx_property_photo_property_id ON property_photo(property_id);

CREATE TABLE IF NOT EXISTS property_features (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    property_id bigint NOT NULL,
    title varchar(200) NOT NULL,
    value text NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_property_features_property FOREIGN KEY (property_id)
        REFERENCES property(id) ON DELETE CASCADE
);
CREATE INDEX idx_property_features_property_id ON property_features(property_id);

CREATE TABLE IF NOT EXISTS broker_property (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    broker_id bigint NOT NULL,
    property_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT broker_property_broker_property_unique UNIQUE (broker_id, property_id),
    CONSTRAINT fk_broker_property_broker FOREIGN KEY (broker_id)
        REFERENCES broker(id) ON DELETE CASCADE,
    CONSTRAINT fk_broker_property_property FOREIGN KEY (property_id)
        REFERENCES property(id) ON DELETE CASCADE
);
CREATE INDEX idx_broker_property_broker_id ON broker_property(broker_id);
CREATE INDEX idx_broker_property_property_id ON broker_property(property_id);

CREATE TABLE IF NOT EXISTS broker_external_links (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    broker_id bigint NOT NULL,
    type varchar(50) NOT NULL,
    link varchar(500) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_broker_external_links_broker FOREIGN KEY (broker_id)
        REFERENCES broker(id) ON DELETE CASCADE
);
CREATE INDEX idx_broker_external_links_broker_id ON broker_external_links(broker_id);

CREATE TABLE IF NOT EXISTS broker_phone (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    broker_id bigint NOT NULL,
    type varchar(50) NOT NULL,
    number varchar(50) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_broker_phone_broker FOREIGN KEY (broker_id)
        REFERENCES broker(id) ON DELETE CASCADE
);
CREATE INDEX idx_broker_phone_broker_id ON broker_phone(broker_id);

CREATE TABLE IF NOT EXISTS property_expenses (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    property_id bigint NOT NULL,
    type varchar(100) NOT NULL,
    annual_price numeric(15,2) NOT NULL,
    monthly_price numeric(15,2) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_property_expenses_property FOREIGN KEY (property_id)
        REFERENCES property(id) ON DELETE CASCADE
);
CREATE INDEX idx_property_expenses_property_id ON property_expenses(property_id);

CREATE TYPE property_status AS ENUM ('available', 'pending', 'sold', 'rented', 'off_market');
CREATE TYPE property_category AS ENUM ('house', 'apartment', 'condo', 'land', 'commercial');
CREATE TYPE broker_link_type AS ENUM ('website', 'linkedin', 'facebook', 'twitter', 'instagram', 'youtube');
CREATE TYPE phone_type AS ENUM ('mobile', 'office', 'home', 'fax');
CREATE TYPE expense_type AS ENUM ('tax', 'insurance', 'maintenance', 'utilities', 'other');

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_property_updated_at
    BEFORE UPDATE ON property
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_broker_updated_at
    BEFORE UPDATE ON broker
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
