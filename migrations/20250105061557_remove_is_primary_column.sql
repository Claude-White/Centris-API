-- +goose Up
-- +goose StatementBegin
ALTER TABLE property_photo 
    DROP COLUMN is_primary;

ALTER TABLE broker_phone 
    DROP COLUMN is_primary;

ALTER TABLE broker_property 
    DROP COLUMN is_primary_broker;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE property_photo 
    ADD COLUMN is_primary boolean DEFAULT false;

ALTER TABLE broker_phone 
    ADD COLUMN is_primary boolean DEFAULT false;

ALTER TABLE broker_property 
    ADD COLUMN is_primary_broker boolean DEFAULT false;
-- +goose StatementEnd