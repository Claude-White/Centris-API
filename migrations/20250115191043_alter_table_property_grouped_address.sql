-- +goose Up
-- +goose StatementBegin
ALTER TABLE property
DROP COLUMN civic_number;

ALTER TABLE property
DROP COLUMN street_name;

ALTER TABLE property
DROP COLUMN apartment_number;

ALTER TABLE property
ADD address varchar(500) NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE property
ADD civic_number varchar(500) NOT NULL;

ALTER TABLE property
ADD street_name varchar(500) NOT NULL;

ALTER TABLE property
ADD apartment_number varchar(500) NOT NULL;

ALTER TABLE property
DROP COLUMN address;
-- +goose StatementEnd
