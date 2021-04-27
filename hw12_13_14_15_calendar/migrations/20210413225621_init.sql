-- +goose Up
-- +goose StatementBegin
CREATE TABLE events
(
    id          serial       NOT NULL,
    title       varchar(255) NOT NULL,
    date_start  timestamp    NOT NULL,
    date_end    timestamp    NOT NULL,
    description text         NULL,
    user_id     int          NOT NULL,
    created_at  timestamp    NOT NULL DEFAULT now(),
    updated_at  timestamp    NULL,
    CONSTRAINT events_pk PRIMARY KEY (id)
);
CREATE INDEX events_date_period_idx ON events USING btree (date_start, date_end);
CREATE UNIQUE INDEX events_user_id_date_start_date_end_idx ON events USING btree (user_id, date_start, date_end);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
-- +goose StatementEnd