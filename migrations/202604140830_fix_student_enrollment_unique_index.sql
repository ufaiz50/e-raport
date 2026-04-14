-- +migrate Up
-- Fix: the old idx_student_term unique index only covered (academic_year, semester),
-- which prevented more than one student from being enrolled in the same term.
-- The correct constraint should be unique per (student_id, class_id, academic_year, semester).

DROP INDEX IF EXISTS idx_student_term;

CREATE UNIQUE INDEX idx_student_term
    ON student_enrollments (student_id, class_id, academic_year, semester);

-- +migrate Down
DROP INDEX IF EXISTS idx_student_term;

CREATE UNIQUE INDEX idx_student_term
    ON student_enrollments (academic_year, semester);
