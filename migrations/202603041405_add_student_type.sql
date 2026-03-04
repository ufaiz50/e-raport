-- +migrate Up
ALTER TABLE students
    ADD COLUMN IF NOT EXISTS student_type VARCHAR(20);

UPDATE students
SET student_type = 'junior'
WHERE student_type IS NULL;

ALTER TABLE students
    ALTER COLUMN student_type SET NOT NULL;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'chk_students_student_type'
    ) THEN
        ALTER TABLE students
            ADD CONSTRAINT chk_students_student_type
            CHECK (student_type IN ('junior', 'senior'));
    END IF;
END$$;

-- +migrate Down
ALTER TABLE students DROP CONSTRAINT IF EXISTS chk_students_student_type;
ALTER TABLE students DROP COLUMN IF EXISTS student_type;
