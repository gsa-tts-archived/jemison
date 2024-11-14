

-- -- This is for working directly against the river_job table.
-- -- A horrible idea, really. But, perhaps it is OK for reading only.
CREATE TABLE river_job (
    id bigint NOT NULL,
    state public.river_job_state DEFAULT 'available'::public.river_job_state NOT NULL,
    attempt smallint DEFAULT 0 NOT NULL,
    max_attempts smallint NOT NULL,
    attempted_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    finalized_at timestamp with time zone,
    scheduled_at timestamp with time zone DEFAULT now() NOT NULL,
    priority smallint DEFAULT 1 NOT NULL,
    args jsonb NOT NULL,
    attempted_by text[],
    errors jsonb[],
    kind text NOT NULL,
    metadata jsonb DEFAULT '{}'::jsonb NOT NULL,
    queue text DEFAULT 'default'::text NOT NULL,
    tags character varying(255)[] DEFAULT '{}'::character varying[] NOT NULL,
    unique_key bytea,
    unique_states bit(8),
    CONSTRAINT finalized_or_finalized_at_null CHECK ((((finalized_at IS NULL) AND (state <> ALL (ARRAY['cancelled'::public.river_job_state, 'completed'::public.river_job_state, 'discarded'::public.river_job_state]))) OR ((finalized_at IS NOT NULL) AND (state = ANY (ARRAY['cancelled'::public.river_job_state, 'completed'::public.river_job_state, 'discarded'::public.river_job_state]))))),
    CONSTRAINT kind_length CHECK (((char_length(kind) > 0) AND (char_length(kind) < 128))),
    CONSTRAINT max_attempts_is_positive CHECK ((max_attempts > 0)),
    CONSTRAINT priority_in_range CHECK (((priority >= 1) AND (priority <= 4))),
    CONSTRAINT queue_length CHECK (((char_length(queue) > 0) AND (char_length(queue) < 128)))
);

