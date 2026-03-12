package api

import (
	"net/http"

	"golang-rest-api-template/pkg/database"
	"golang-rest-api-template/pkg/models"

	"github.com/gin-gonic/gin"
)

type academicRepository struct {
	DB database.Database
}

func NewAcademicRepository(db database.Database) *academicRepository {
	return &academicRepository{DB: db}
}

func (r *academicRepository) ListAcademicYears(c *gin.Context) {
	schoolID, _, ok := requireTenant(c)
	if !ok {
		return
	}
	offset, limit, ok := parsePagination(c)
	if !ok {
		return
	}
	query := r.DB.Model(&models.AcademicYear{})
	if schoolID != nil {
		query = query.Where("school_id = ?", *schoolID)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count academic years"})
		return
	}
	var rows []models.AcademicYear
	if err := query.Offset(offset).Limit(limit).Order("id desc").Find(&rows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch academic years"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rows, "meta": gin.H{"offset": offset, "limit": limit, "total": total, "count": len(rows)}})
}

func (r *academicRepository) CreateAcademicYear(c *gin.Context) {
	var input models.AcademicYear
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	schoolID, _, ok := resolveWriteSchoolID(c, input.SchoolID)
	if !ok {
		return
	}
	input.SchoolID = schoolID
	if err := r.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create academic year"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": input})
}

func (r *academicRepository) ListSemesters(c *gin.Context) {
	schoolID, _, ok := requireTenant(c)
	if !ok {
		return
	}
	offset, limit, ok := parsePagination(c)
	if !ok {
		return
	}
	query := r.DB.Model(&models.Semester{})
	if schoolID != nil {
		query = query.Where("school_id = ?", *schoolID)
	}
	var total int64
	query.Count(&total)
	var rows []models.Semester
	query.Offset(offset).Limit(limit).Order("id desc").Find(&rows)
	c.JSON(http.StatusOK, gin.H{"data": rows, "meta": gin.H{"offset": offset, "limit": limit, "total": total, "count": len(rows)}})
}

func (r *academicRepository) CreateSemester(c *gin.Context) {
	var input models.Semester
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	schoolID, _, ok := resolveWriteSchoolID(c, input.SchoolID)
	if !ok {
		return
	}
	input.SchoolID = schoolID
	if err := r.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create semester"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": input})
}

func (r *academicRepository) ListCurriculums(c *gin.Context) {
	schoolID, _, ok := requireTenant(c)
	if !ok {
		return
	}
	offset, limit, ok := parsePagination(c)
	if !ok {
		return
	}
	query := r.DB.Model(&models.Curriculum{})
	if schoolID != nil {
		query = query.Where("school_id = ?", *schoolID)
	}
	var total int64
	query.Count(&total)
	var rows []models.Curriculum
	query.Offset(offset).Limit(limit).Order("id desc").Find(&rows)
	c.JSON(http.StatusOK, gin.H{"data": rows, "meta": gin.H{"offset": offset, "limit": limit, "total": total, "count": len(rows)}})
}

func (r *academicRepository) CreateCurriculum(c *gin.Context) {
	var input models.Curriculum
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	schoolID, _, ok := resolveWriteSchoolID(c, input.SchoolID)
	if !ok {
		return
	}
	input.SchoolID = schoolID
	if err := r.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create curriculum"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": input})
}

func (r *academicRepository) ListTeachings(c *gin.Context) {
	schoolID, _, ok := requireTenant(c)
	if !ok {
		return
	}
	offset, limit, ok := parsePagination(c)
	if !ok {
		return
	}
	query := r.DB.Model(&models.Teaching{})
	if schoolID != nil {
		query = query.Where("school_id = ?", *schoolID)
	}
	var total int64
	query.Count(&total)
	var rows []models.Teaching
	query.Offset(offset).Limit(limit).Order("id desc").Find(&rows)
	c.JSON(http.StatusOK, gin.H{"data": rows, "meta": gin.H{"offset": offset, "limit": limit, "total": total, "count": len(rows)}})
}

func (r *academicRepository) CreateTeaching(c *gin.Context) {
	var input models.Teaching
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	schoolID, _, ok := resolveWriteSchoolID(c, input.SchoolID)
	if !ok {
		return
	}
	input.SchoolID = schoolID
	if err := r.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create teaching"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": input})
}
