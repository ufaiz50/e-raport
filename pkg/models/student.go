package models

import (
	"time"

	"gorm.io/gorm"
)

type Student struct {
	UUIDPrimaryKey
	Nama                    string     `json:"nama"`
	NamaPanggilan           string     `json:"nama_panggilan"`
	Email                   string     `json:"email,omitempty" gorm:"index:idx_student_email_school,unique,where:email IS NOT NULL"`
	NIS                     string     `json:"nis"`
	NISN                    string     `json:"nisn"`
	TempatLahir             string     `json:"tempat_lahir"`
	TanggalLahir            *time.Time `json:"tanggal_lahir"`
	Agama                   string     `json:"agama"`
	AnakKe                  *int       `json:"anak_ke"`
	JenisKelamin            string     `json:"jenis_kelamin"`
	NamaAyah                string     `json:"nama_ayah"`
	PekerjaanAyah           string     `json:"pekerjaan_ayah"`
	NamaIbu                 string     `json:"nama_ibu"`
	PekerjaanIbu            string     `json:"pekerjaan_ibu"`
	NoHPOrangtua            string     `json:"no_hp_orangtua"`
	AlamatOrangtuaJalan     string     `json:"alamat_orangtua_jalan"`
	AlamatOrangtuaKecamatan string     `json:"alamat_orangtua_kecamatan"`
	AlamatOrangtuaKabupaten string     `json:"alamat_orangtua_kabupaten"`
	AlamatOrangtuaProvinsi  string     `json:"alamat_orangtua_provinsi"`
	NamaWali                string     `json:"nama_wali"`
	PekerjaanWali           string     `json:"pekerjaan_wali"`
	NoHPWali                string     `json:"no_hp_wali"`
	AlamatWaliJalan         string     `json:"alamat_wali_jalan"`
	AlamatWaliKecamatan     string     `json:"alamat_wali_kecamatan"`
	AlamatWaliKabupaten     string     `json:"alamat_wali_kabupaten"`
	AlamatWaliProvinsi      string     `json:"alamat_wali_provinsi"`
	TanggalDiterima         *time.Time `json:"tanggal_diterima"`
	CatatanGuru             string     `json:"catatan_guru"`
	Status                  string     `json:"status" gorm:"type:varchar(20);default:active"`
	SchoolID                *string    `json:"school_id,omitempty" gorm:"type:uuid"`
	School                  *School    `json:"school,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:SchoolID;references:ID"`
	ClassID                 *string    `json:"class_id,omitempty" gorm:"type:uuid"`
	Class                   *Class     `json:"class,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:ClassID;references:ID"`
	Type                    string     `json:"-" gorm:"column:student_type;type:varchar(20);not null;default:junior"`
	CreatedAt               time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt               time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

func (s *Student) BeforeCreate(_ *gorm.DB) error { s.ID = ensureUUID(s.ID); return nil }

type CreateStudent struct {
	Nama                    string     `json:"nama" binding:"required"`
	NamaPanggilan           string     `json:"nama_panggilan"`
	Email                   string     `json:"email" binding:"omitempty,email"`
	NIS                     string     `json:"nis"`
	NISN                    string     `json:"nisn"`
	TempatLahir             string     `json:"tempat_lahir"`
	TanggalLahir            *time.Time `json:"tanggal_lahir"`
	Agama                   string     `json:"agama"`
	AnakKe                  *int       `json:"anak_ke"`
	JenisKelamin            string     `json:"jenis_kelamin"`
	NamaAyah                string     `json:"nama_ayah"`
	PekerjaanAyah           string     `json:"pekerjaan_ayah"`
	NamaIbu                 string     `json:"nama_ibu"`
	PekerjaanIbu            string     `json:"pekerjaan_ibu"`
	NoHPOrangtua            string     `json:"no_hp_orangtua"`
	AlamatOrangtuaJalan     string     `json:"alamat_orangtua_jalan"`
	AlamatOrangtuaKecamatan string     `json:"alamat_orangtua_kecamatan"`
	AlamatOrangtuaKabupaten string     `json:"alamat_orangtua_kabupaten"`
	AlamatOrangtuaProvinsi  string     `json:"alamat_orangtua_provinsi"`
	NamaWali                string     `json:"nama_wali"`
	PekerjaanWali           string     `json:"pekerjaan_wali"`
	NoHPWali                string     `json:"no_hp_wali"`
	AlamatWaliJalan         string     `json:"alamat_wali_jalan"`
	AlamatWaliKecamatan     string     `json:"alamat_wali_kecamatan"`
	AlamatWaliKabupaten     string     `json:"alamat_wali_kabupaten"`
	AlamatWaliProvinsi      string     `json:"alamat_wali_provinsi"`
	TanggalDiterima         *time.Time `json:"tanggal_diterima"`
	CatatanGuru             string     `json:"catatan_guru"`
	Status                  string     `json:"status"`
	SchoolID                *string    `json:"school_id"`
	ClassID                 *string    `json:"class_id"`
}

type UpdateStudent struct {
	Nama                    string     `json:"nama"`
	NamaPanggilan           string     `json:"nama_panggilan"`
	Email                   string     `json:"email" binding:"omitempty,email"`
	NIS                     string     `json:"nis"`
	NISN                    string     `json:"nisn"`
	TempatLahir             string     `json:"tempat_lahir"`
	TanggalLahir            *time.Time `json:"tanggal_lahir"`
	Agama                   string     `json:"agama"`
	AnakKe                  *int       `json:"anak_ke"`
	JenisKelamin            string     `json:"jenis_kelamin"`
	NamaAyah                string     `json:"nama_ayah"`
	PekerjaanAyah           string     `json:"pekerjaan_ayah"`
	NamaIbu                 string     `json:"nama_ibu"`
	PekerjaanIbu            string     `json:"pekerjaan_ibu"`
	NoHPOrangtua            string     `json:"no_hp_orangtua"`
	AlamatOrangtuaJalan     string     `json:"alamat_orangtua_jalan"`
	AlamatOrangtuaKecamatan string     `json:"alamat_orangtua_kecamatan"`
	AlamatOrangtuaKabupaten string     `json:"alamat_orangtua_kabupaten"`
	AlamatOrangtuaProvinsi  string     `json:"alamat_orangtua_provinsi"`
	NamaWali                string     `json:"nama_wali"`
	PekerjaanWali           string     `json:"pekerjaan_wali"`
	NoHPWali                string     `json:"no_hp_wali"`
	AlamatWaliJalan         string     `json:"alamat_wali_jalan"`
	AlamatWaliKecamatan     string     `json:"alamat_wali_kecamatan"`
	AlamatWaliKabupaten     string     `json:"alamat_wali_kabupaten"`
	AlamatWaliProvinsi      string     `json:"alamat_wali_provinsi"`
	TanggalDiterima         *time.Time `json:"tanggal_diterima"`
	CatatanGuru             string     `json:"catatan_guru"`
	Status                  string     `json:"status"`
	SchoolID                *string    `json:"school_id"`
	ClassID                 *string    `json:"class_id"`
}
