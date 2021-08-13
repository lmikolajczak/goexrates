CREATE TABLE IF NOT EXISTS currencies (
    id serial PRIMARY KEY,
    code text NOT NULL,
    rate numeric NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);
