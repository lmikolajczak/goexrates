CREATE TABLE IF NOT EXISTS currencies (
    id serial PRIMARY KEY,
    code text NOT NULL,
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS rates (
    id serial PRIMARY KEY,
    currency_id int NOT NULL REFERENCES currencies(id) ON DELETE CASCADE,
    rate numeric NOT NULL,
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);
