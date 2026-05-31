CREATE TABLE assignments(
    id TEXT PRIMARY KEY,
    course_id TEXT NOT NULL,
    slug TEXT NOT NULL,
    title TEXT NOT NULL,
    description_markdown TEXT NOT NULL DEFAULT '',
    deadline TIMESTAMP,
    hard_deadline TIMESTAMP,
    max_score INT NOT NULL DEFAULT 0,
    template_repo TEXT NOT NULL DEFAULT '',
    repo_naming_pattern TEXT NOT NULL DEFAULT '',
    auto_request_review_on_pass BOOLEAN NOT NULL DEFAULT FALSE,
    requires_defense BOOLEAN NOT NULL DEFAULT FALSE,
    weight_tests DOUBLE PRECISION NOT NULL DEFAULT 0,
    weight_quality DOUBLE PRECISION NOT NULL DEFAULT 0,
    defence_multiplier BOOLEAN NOT NULL DEFAULT FALSE,
    custom_formula TEXT NOT NULL DEFAULT '',
    runner TEXT NOT NULL DEFAULT 'external_ci',
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT assignments_course_slug_key UNIQUE (course_id, slug)
);

CREATE INDEX assignments_course_idx ON assignments (course_id);

CREATE TRIGGER assignments_updated_at
    BEFORE UPDATE ON assignments
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();
