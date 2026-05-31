CREATE TABLE enrollments(
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    course_id TEXT NOT NULL,
    role TEXT NOT NULL,
    enrolled_at TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT enrollments_user_course_key UNIQUE (user_id, course_id)
);

CREATE INDEX enrollments_course_idx ON enrollments (course_id);
CREATE INDEX enrollments_user_idx ON enrollments (user_id);
