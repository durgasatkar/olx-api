CREATE TABLE listings (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    title text NOT NULL,
    description TEXT NOT NULL,
    price bigint NOT NULL,
    city text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW()
);