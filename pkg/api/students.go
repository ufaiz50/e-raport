package api

import (
	"context"
	"golang-rest-api-template/pkg/database"
	"golang-rest-api-template/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StudentRepository interface {
	FindStudents(c *gin.Context)
	CreateStudent(c *gin.Context)
	FindStudent(c *gin.Context)
	UpdateStudent(c *gin.Context)
	DeleteStudent(c *gin.Context)
}

type studentRepository struct {
	DB  database.Database
	Ctx *context.Context
}

func NewStudentRepository(db database.Database, ctx *context.Context) *studentRepository {
	return &studentRepository{DB: db, Ctx: ctx}
}

func (r *studentRepository) FindStudents(c *gin.Context) {
	schoolID, _, ok := requireTenant(c)
	if !ok {
		return
	}

	offset, limit, ok := parsePagination(c)
	if !ok {
		return
	}

	query := r.DB.Model(&models.Student{})
	if schoolID != nil {
		query = query.Where("school_id = ?", *schoolID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count students"})
		return
	}

	var students []models.Student
	dataQuery := r.DB.Offset(offset).Limit(limit).Order("created_at asc")
	if schoolID != nil {
		dataQuery = dataQuery.Where("school_id = ?", *schoolID)
	}
	dataQuery.Find(&students)
	// Enrich ClassID from active enrollment (backward compat)
	for i := range students {
		if students[i].SchoolID == nil {
			continue
		}
		var active models.StudentEnrollment
		if err := r.DB.Where("student_id = ? AND school_id = ? AND is_active = ?", students[i].ID, *students[i].SchoolID, true).Order("created_at desc").First(&active).Error; err == nil {
			students[i].ClassID = &active.ClassID
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"data": students,
		"meta": gin.H{
			"offset": offset,
			"limit":  limit,
			"total":  total,
			"count":  len(students),
		},
	})
}

func (r *studentRepository) CreateStudent(c *gin.Context) {
	var input models.CreateStudent
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schoolID, _, ok := resolveWriteSchoolID(c, input.SchoolID)
	if !ok {
		return
	}

	// Validate ClassID if provided (optional — student is master data)
	if input.ClassID != nil && *input.ClassID != "" {
		var class models.Class
		if err := r.DB.Where("id = ? AND school_id = ?", *input.ClassID, *schoolID).First(&class).Error(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "class not found"})
			return
		}
	}

	student := models.Student{
		Nama:                    input.Nama,
		NamaPanggilan:           input.NamaPanggilan,
		Email:                   input.Email,
		NIS:                     input.NIS,
		NISN:                    input.NISN,
		TempatLahir:             input.TempatLahir,
		TanggalLahir:            input.TanggalLahir,
		Agama:                   input.Agama,
		AnakKe:                  input.AnakKe,
		JenisKelamin:            input.JenisKelamin,
		NamaAyah:                input.NamaAyah,
		PekerjaanAyah:           input.PekerjaanAyah,
		NamaIbu:                 input.NamaIbu,
		PekerjaanIbu:            input.PekerjaanIbu,
		NoHPOrangtua:            input.NoHPOrangtua,
		AlamatOrangtuaJalan:     input.AlamatOrangtuaJalan,
		AlamatOrangtuaKecamatan: input.AlamatOrangtuaKecamatan,
		AlamatOrangtuaKabupaten: input.AlamatOrangtuaKabupaten,
		AlamatOrangtuaProvinsi:  input.AlamatOrangtuaProvinsi,
		NamaWali:                input.NamaWali,
		PekerjaanWali:           input.PekerjaanWali,
		NoHPWali:                input.NoHPWali,
		AlamatWaliJalan:         input.AlamatWaliJalan,
		AlamatWaliKecamatan:     input.AlamatWaliKecamatan,
		AlamatWaliKabupaten:     input.AlamatWaliKabupaten,
		AlamatWaliProvinsi:      input.AlamatWaliProvinsi,
		TanggalDiterima:         input.TanggalDiterima,
		CatatanGuru:             input.CatatanGuru,
		Status:                  input.Status,
		SchoolID:                schoolID,
		ClassID:                 input.ClassID,
		Type:                    "junior",
	}

	if err := r.DB.Create(&student).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create student"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": student})
}

func (r *studentRepository) FindStudent(c *gin.Context) {
	schoolID, _, ok := requireTenant(c)
	if !ok {
		return
	}

	var student models.Student
	if err := whereByIDOrUUID(r.DB, c.Param("id"), schoolID).First(&student).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": student})
}

func (r *studentRepository) UpdateStudent(c *gin.Context) {
	var input models.UpdateStudent

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schoolID, _, ok := resolveWriteSchoolID(c, input.SchoolID)
	if !ok {
		return
	}

	var student models.Student
	if err := whereByIDOrUUID(r.DB, c.Param("id"), schoolID).First(&student).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
		return
	}

	// Validate ClassID if provided (optional)
	if input.ClassID != nil && *input.ClassID != "" {
		var class models.Class
		if err := r.DB.Where("id = ? AND school_id = ?", *input.ClassID, *schoolID).First(&class).Error(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "class not found"})
			return
		}
	}

	if err := r.DB.Model(&student).Updates(models.Student{
		Nama:                    input.Nama,
		NamaPanggilan:           input.NamaPanggilan,
		Email:                   input.Email,
		NIS:                     input.NIS,
		NISN:                    input.NISN,
		TempatLahir:             input.TempatLahir,
		TanggalLahir:            input.TanggalLahir,
		Agama:                   input.Agama,
		AnakKe:                  input.AnakKe,
		JenisKelamin:            input.JenisKelamin,
		NamaAyah:                input.NamaAyah,
		PekerjaanAyah:           input.PekerjaanAyah,
		NamaIbu:                 input.NamaIbu,
		PekerjaanIbu:            input.PekerjaanIbu,
		NoHPOrangtua:            input.NoHPOrangtua,
		AlamatOrangtuaJalan:     input.AlamatOrangtuaJalan,
		AlamatOrangtuaKecamatan: input.AlamatOrangtuaKecamatan,
		AlamatOrangtuaKabupaten: input.AlamatOrangtuaKabupaten,
		AlamatOrangtuaProvinsi:  input.AlamatOrangtuaProvinsi,
		NamaWali:                input.NamaWali,
		PekerjaanWali:           input.PekerjaanWali,
		NoHPWali:                input.NoHPWali,
		AlamatWaliJalan:         input.AlamatWaliJalan,
		AlamatWaliKecamatan:     input.AlamatWaliKecamatan,
		AlamatWaliKabupaten:     input.AlamatWaliKabupaten,
		AlamatWaliProvinsi:      input.AlamatWaliProvinsi,
		TanggalDiterima:         input.TanggalDiterima,
		CatatanGuru:             input.CatatanGuru,
		Status:                  input.Status,
		SchoolID:                schoolID,
		ClassID:                 input.ClassID,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update student"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": student})
}

func (r *studentRepository) DeleteStudent(c *gin.Context) {
	schoolID, _, ok := requireTenant(c)
	if !ok {
		return
	}

	var student models.Student
	if err := whereByIDOrUUID(r.DB, c.Param("id"), schoolID).First(&student).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
		return
	}

	r.DB.Delete(&student)
	c.JSON(http.StatusNoContent, gin.H{"data": true})
}
