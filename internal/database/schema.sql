CREATE SEQUENCE IF NOT EXISTS property_photo_id_seq;
CREATE SEQUENCE IF NOT EXISTS property_features_id_seq;
CREATE SEQUENCE IF NOT EXISTS broker_property_id_seq;
CREATE SEQUENCE IF NOT EXISTS broker_exteral_links_id_seq;
CREATE SEQUENCE IF NOT EXISTS broker_phone_id_seq;
CREATE SEQUENCE IF NOT EXISTS property_expenses_id_seq;

CREATE TABLE IF NOT EXISTS property_photo (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  property_id bigint NOT NULL,
  link varchar(500) NOT NULL,
  description varchar(500)
);

CREATE TABLE IF NOT EXISTS broker (
  id bigint NOT NULL PRIMARY KEY,
  first_name varchar(500) NOT NULL,
  middle_name varchar(500),
  last_name varchar(500) NOT NULL,
  title varchar(500) NOT NULL,
  profile_photo varchar(500),
  complementary_info varchar(500) NOT NULL,
  served_areas varchar(500),
  presentation varchar(500),
  corporation_name varchar(500),
  agency_name varchar(500) NOT NULL,
  agency_address varchar(500) NOT NULL,
  agency_logo varchar(500)
);

CREATE TABLE IF NOT EXISTS property_features (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  property_id bigint NOT NULL,
  title varchar(500) NOT NULL,
  value varchar(500) NOT NULL
);

CREATE TABLE IF NOT EXISTS broker_property (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  broker_id bigint NOT NULL,
  property_id bigint NOT NULL
);

CREATE TABLE IF NOT EXISTS broker_exteral_links (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  broker_id bigint NOT NULL,
  type varchar(500) NOT NULL,
  link varchar(500) NOT NULL
);

CREATE TABLE IF NOT EXISTS broker_phone (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  broker_id bigint NOT NULL,
  type varchar(500) NOT NULL,
  number varchar(500) NOT NULL
);

CREATE TABLE IF NOT EXISTS property_expenses (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  property_id bigint NOT NULL,
  type varchar(500) NOT NULL,
  annual_price varchar(500) NOT NULL,
  monthly_price varchar(500) NOT NULL
);

CREATE TABLE IF NOT EXISTS property (
  id bigint NOT NULL PRIMARY KEY,
  title varchar(500) NOT NULL,
  category varchar(500) NOT NULL,
  civic_number varchar(500),
  street_name varchar(500),
  appartment_number varchar(500),
  city_name varchar(500),
  neighbourhood_name varchar(500),
  price varchar(500) NOT NULL,
  description varchar(500),
  bedroom_number varchar(500),
  room_number varchar(500),
  bathroom_number varchar(500),
  longitude varchar(500) NOT NULL,
  latitude varchar(500) NOT NULL
);

COMMENT ON COLUMN property.id IS 'MLS number';

ALTER TABLE broker ADD CONSTRAINT Courtier_id_fk FOREIGN KEY (id) REFERENCES broker_exteral_links (broker_id);
ALTER TABLE broker ADD CONSTRAINT Courtier_id_fk FOREIGN KEY (id) REFERENCES broker_property (broker_id);
ALTER TABLE broker ADD CONSTRAINT Courtier_id_fk FOREIGN KEY (id) REFERENCES broker_phone (broker_id);
ALTER TABLE property ADD CONSTRAINT Propriété_id_fk FOREIGN KEY (id) REFERENCES property_expenses (property_id);
ALTER TABLE property ADD CONSTRAINT Propriété_id_fk FOREIGN KEY (id) REFERENCES property_photo (property_id);
ALTER TABLE property ADD CONSTRAINT Propriété_id_fk FOREIGN KEY (id) REFERENCES property_features (property_id);
ALTER TABLE property ADD CONSTRAINT Propriété_id_fk FOREIGN KEY (id) REFERENCES broker_property (property_id);
