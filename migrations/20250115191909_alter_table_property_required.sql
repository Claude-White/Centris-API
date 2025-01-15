-- +goose Up
-- +goose StatementBegin
ALTER TABLE property
DROP COLUMN neighbourhood_name;

ALTER TABLE property
ALTER COLUMN city_name SET DATA TYPE varchar(500);

ALTER TABLE property
ALTER COLUMN city_name SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE property
ADD neighbourhood_name varchar(500) NULL;

ALTER TABLE property
ALTER COLUMN city_name DROP NOT NULL;

ALTER TABLE property
ALTER COLUMN city_name SET DATA TYPE varchar(255); -- Adjust as needed
-- +goose StatementEnd
