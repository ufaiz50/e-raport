-- +migrate Up
ALTER TABLE students
    ADD COLUMN IF NOT EXISTS nama VARCHAR(100),
    ADD COLUMN IF NOT EXISTS nama_panggilan VARCHAR(50),
    ADD COLUMN IF NOT EXISTS tempat_lahir VARCHAR(100),
    ADD COLUMN IF NOT EXISTS tanggal_lahir DATE,
    ADD COLUMN IF NOT EXISTS agama VARCHAR(50),
    ADD COLUMN IF NOT EXISTS anak_ke INT,
    ADD COLUMN IF NOT EXISTS jenis_kelamin VARCHAR(20),
    ADD COLUMN IF NOT EXISTS nama_ayah VARCHAR(100),
    ADD COLUMN IF NOT EXISTS pekerjaan_ayah VARCHAR(100),
    ADD COLUMN IF NOT EXISTS nama_ibu VARCHAR(100),
    ADD COLUMN IF NOT EXISTS pekerjaan_ibu VARCHAR(100),
    ADD COLUMN IF NOT EXISTS no_hp_orangtua VARCHAR(20),
    ADD COLUMN IF NOT EXISTS alamat_orangtua_jalan TEXT,
    ADD COLUMN IF NOT EXISTS alamat_orangtua_kecamatan VARCHAR(100),
    ADD COLUMN IF NOT EXISTS alamat_orangtua_kabupaten VARCHAR(100),
    ADD COLUMN IF NOT EXISTS alamat_orangtua_provinsi VARCHAR(100),
    ADD COLUMN IF NOT EXISTS nama_wali VARCHAR(100),
    ADD COLUMN IF NOT EXISTS pekerjaan_wali VARCHAR(100),
    ADD COLUMN IF NOT EXISTS no_hp_wali VARCHAR(20),
    ADD COLUMN IF NOT EXISTS alamat_wali_jalan TEXT,
    ADD COLUMN IF NOT EXISTS alamat_wali_kecamatan VARCHAR(100),
    ADD COLUMN IF NOT EXISTS alamat_wali_kabupaten VARCHAR(100),
    ADD COLUMN IF NOT EXISTS alamat_wali_provinsi VARCHAR(100),
    ADD COLUMN IF NOT EXISTS tanggal_diterima DATE,
    ADD COLUMN IF NOT EXISTS catatan_guru TEXT;

UPDATE students
SET
    nama = COALESCE(NULLIF(TRIM(CONCAT(COALESCE(first_name, ''), ' ', COALESCE(last_name, ''))), ''), nama),
    tempat_lahir = COALESCE(NULLIF(birth_place, ''), tempat_lahir),
    tanggal_lahir = COALESCE(birth_date::date, tanggal_lahir),
    agama = COALESCE(NULLIF(religion, ''), agama),
    no_hp_orangtua = COALESCE(NULLIF(parent_phone, ''), no_hp_orangtua),
    alamat_orangtua_jalan = COALESCE(NULLIF(address, ''), alamat_orangtua_jalan),
    nama_ayah = COALESCE(NULLIF(parent_name, ''), nama_ayah)
WHERE
    nama IS NULL
    OR tempat_lahir IS NULL
    OR tanggal_lahir IS NULL
    OR agama IS NULL
    OR no_hp_orangtua IS NULL
    OR alamat_orangtua_jalan IS NULL
    OR nama_ayah IS NULL;

ALTER TABLE students
    ALTER COLUMN email DROP NOT NULL;

DROP INDEX IF EXISTS idx_student_email_school;
CREATE UNIQUE INDEX IF NOT EXISTS idx_student_email_school
    ON students (school_id, email)
    WHERE email IS NOT NULL AND email <> '';

ALTER TABLE students DROP COLUMN IF EXISTS first_name;
ALTER TABLE students DROP COLUMN IF EXISTS last_name;
ALTER TABLE students DROP COLUMN IF EXISTS birth_place;
ALTER TABLE students DROP COLUMN IF EXISTS birth_date;
ALTER TABLE students DROP COLUMN IF EXISTS address;
ALTER TABLE students DROP COLUMN IF EXISTS phone;
ALTER TABLE students DROP COLUMN IF EXISTS religion;
ALTER TABLE students DROP COLUMN IF EXISTS parent_name;
ALTER TABLE students DROP COLUMN IF EXISTS parent_phone;

-- +migrate Down
ALTER TABLE students
    ADD COLUMN IF NOT EXISTS first_name VARCHAR(255),
    ADD COLUMN IF NOT EXISTS last_name VARCHAR(255),
    ADD COLUMN IF NOT EXISTS birth_place VARCHAR(255),
    ADD COLUMN IF NOT EXISTS birth_date TIMESTAMP WITH TIME ZONE,
    ADD COLUMN IF NOT EXISTS address TEXT,
    ADD COLUMN IF NOT EXISTS phone VARCHAR(255),
    ADD COLUMN IF NOT EXISTS religion VARCHAR(255),
    ADD COLUMN IF NOT EXISTS parent_name VARCHAR(255),
    ADD COLUMN IF NOT EXISTS parent_phone VARCHAR(255);

UPDATE students
SET
    first_name = COALESCE(first_name, nama),
    last_name = COALESCE(last_name, ''),
    birth_place = COALESCE(birth_place, tempat_lahir),
    birth_date = COALESCE(birth_date, tanggal_lahir::timestamp with time zone),
    address = COALESCE(address, alamat_orangtua_jalan),
    religion = COALESCE(religion, agama),
    parent_name = COALESCE(parent_name, nama_ayah),
    parent_phone = COALESCE(parent_phone, no_hp_orangtua);

DROP INDEX IF EXISTS idx_student_email_school;
CREATE UNIQUE INDEX IF NOT EXISTS idx_student_email_school
    ON students (school_id, email)
    WHERE email IS NOT NULL;

ALTER TABLE students
    ADD COLUMN IF NOT EXISTS phone VARCHAR(255);

ALTER TABLE students DROP COLUMN IF EXISTS nama;
ALTER TABLE students DROP COLUMN IF EXISTS nama_panggilan;
ALTER TABLE students DROP COLUMN IF EXISTS tempat_lahir;
ALTER TABLE students DROP COLUMN IF EXISTS tanggal_lahir;
ALTER TABLE students DROP COLUMN IF EXISTS agama;
ALTER TABLE students DROP COLUMN IF EXISTS anak_ke;
ALTER TABLE students DROP COLUMN IF EXISTS jenis_kelamin;
ALTER TABLE students DROP COLUMN IF EXISTS nama_ayah;
ALTER TABLE students DROP COLUMN IF EXISTS pekerjaan_ayah;
ALTER TABLE students DROP COLUMN IF EXISTS nama_ibu;
ALTER TABLE students DROP COLUMN IF EXISTS pekerjaan_ibu;
ALTER TABLE students DROP COLUMN IF EXISTS no_hp_orangtua;
ALTER TABLE students DROP COLUMN IF EXISTS alamat_orangtua_jalan;
ALTER TABLE students DROP COLUMN IF EXISTS alamat_orangtua_kecamatan;
ALTER TABLE students DROP COLUMN IF EXISTS alamat_orangtua_kabupaten;
ALTER TABLE students DROP COLUMN IF EXISTS alamat_orangtua_provinsi;
ALTER TABLE students DROP COLUMN IF EXISTS nama_wali;
ALTER TABLE students DROP COLUMN IF EXISTS pekerjaan_wali;
ALTER TABLE students DROP COLUMN IF EXISTS no_hp_wali;
ALTER TABLE students DROP COLUMN IF EXISTS alamat_wali_jalan;
ALTER TABLE students DROP COLUMN IF EXISTS alamat_wali_kecamatan;
ALTER TABLE students DROP COLUMN IF EXISTS alamat_wali_kabupaten;
ALTER TABLE students DROP COLUMN IF EXISTS alamat_wali_provinsi;
ALTER TABLE students DROP COLUMN IF EXISTS tanggal_diterima;
ALTER TABLE students DROP COLUMN IF EXISTS catatan_guru;
