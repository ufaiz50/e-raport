-- +migrate Up
CREATE TABLE IF NOT EXISTS students (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

ALTER TABLE books
    ADD COLUMN IF NOT EXISTS student_id BIGINT NULL;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'fk_books_student'
    ) THEN
        ALTER TABLE books
            ADD CONSTRAINT fk_books_student
            FOREIGN KEY (student_id) REFERENCES students(id)
            ON UPDATE CASCADE ON DELETE SET NULL;
    END IF;
END$$;

CREATE INDEX IF NOT EXISTS idx_books_student_id ON books(student_id);

-- +migrate Down
ALTER TABLE books DROP CONSTRAINT IF EXISTS fk_books_student;
DROP INDEX IF EXISTS idx_books_student_id;
ALTER TABLE books DROP COLUMN IF EXISTS student_id;
DROP TABLE IF EXISTS students;
