CREATE TABLE courses(
    id TEXT PRIMARY KEY,
    code TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    instructor_id TEXT NOT NULL,
    vcs_provider_kind INT,
    vcs_provider_instance TEXT,
    vcs_target_org TEXT,
    vcs_student_team TEXT,
    vcs_reviewer_team TEXT,
    vcs_reviewer_logins TEXT[],
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX courses_instructor_idx ON courses (instructor_id);

CREATE TRIGGER courses_updated_at
    BEFORE UPDATE ON courses
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();
