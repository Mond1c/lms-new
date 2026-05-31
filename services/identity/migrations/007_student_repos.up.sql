CREATE TABLE student_repos(
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    assignment_id TEXT NOT NULL,
    provider_kind INT NOT NULL,
    provider_instance TEXT NOT NULL DEFAULT '',
    full_name TEXT NOT NULL,
    external_id BIGINT NOT NULL DEFAULT 0,
    state TEXT NOT NULL DEFAULT 'pending',
    clone_url_https TEXT NOT NULL DEFAULT '',
    clone_url_ssh TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT student_repos_user_assignment_key UNIQUE (user_id, assignment_id)
);

CREATE INDEX student_repos_assignment_idx ON student_repos (assignment_id);

CREATE TRIGGER student_repos_updated_at
    BEFORE UPDATE ON student_repos
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();
