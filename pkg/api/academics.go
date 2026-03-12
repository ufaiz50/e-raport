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

func (r *academicRepository) deactivateOtherAcademicYears(schoolID *uint, keepID uint) {
	r.DB.Model(&models.AcademicYear{}).Where("school_id = ? AND id <> ?", *schoolID, keepID).Updates(map[string]interface{}{"is_active": false})
}

func (r *academicRepository) deactivateOtherSemesters(schoolID *uint, keepID uint) {
	r.DB.Model(&models.Semester{}).Where("school_id = ? AND id <> ?", *schoolID, keepID).Updates(map[string]interface{}{"is_active": false})
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
	if input.IsActive {
		r.deactivateOtherAcademicYears(schoolID, input.ID)
	}
	c.JSON(http.StatusCreated, gin.H{"data": input})
}

func (r *academicRepository) UpdateAcademicYear(c *gin.Context) {
	var input models.AcademicYear
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	schoolID, _, ok := resolveWriteSchoolID(c, input.SchoolID)
	if !ok {
		return
	}
	var row models.AcademicYear
	if err := r.DB.Where("id = ? AND school_id = ?", c.Param("id"), *schoolID).First(&row).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "academic year not found"})
		return
	}
	r.DB.Model(&row).Updates(models.AcademicYear{Year: input.Year, IsActive: input.IsActive, SchoolID: schoolID})
	if input.IsActive {
		r.deactivateOtherAcademicYears(schoolID, row.ID)
	}
	c.JSON(http.StatusOK, gin.H{"data": row})
}

func (r *academicRepository) DeleteAcademicYear(c *gin.Context) {
	schoolID, _, ok := requireTenant(c)
	if !ok {
		return
	}
	var row models.AcademicYear
	if err := r.DB.Where("id = ? AND school_id = ?", c.Param("id"), *schoolID).First(&row).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "academic year not found"})
		return
	}
	r.DB.Delete(&row)
	c.JSON(http.StatusNoContent, gin.H{"data": true})
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
	if input.IsActive {
		r.deactivateOtherSemesters(schoolID, input.ID)
	}
	c.JSON(http.StatusCreated, gin.H{"data": input})
}

func (r *academicRepository) UpdateSemester(c *gin.Context) {
	var input models.Semester
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	schoolID, _, ok := resolveWriteSchoolID(c, input.SchoolID)
	if !ok {
		return
	}
	var row models.Semester
	if err := r.DB.Where("id = ? AND school_id = ?", c.Param("id"), *schoolID).First(&row).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "semester not found"})
		return
	}
	r.DB.Model(&row).Updates(models.Semester{AcademicYearID: input.AcademicYearID, Name: input.Name, OrderNo: input.OrderNo, IsActive: input.IsActive, SchoolID: schoolID})
	if input.IsActive {
		r.deactivateOtherSemesters(schoolID, row.ID)
	}
	c.JSON(http.StatusOK, gin.H{"data": row})
}

func (r *academicRepository) DeleteSemester(c *gin.Context) {
	schoolID, _, ok := requireTenant(c)
	if !ok {
		return
	}
	var row models.Semester
	if err := r.DB.Where("id = ? AND school_id = ?", c.Param("id"), *schoolID).First(&row).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "semester not found"})
		return
	}
	r.DB.Delete(&row)
	c.JSON(http.StatusNoContent, gin.H{"data": true})
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

func (r *academicRepository) UpdateCurriculum(c *gin.Context) {
	var input models.Curriculum
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	schoolID, _, ok := resolveWriteSchoolID(c, input.SchoolID)
	if !ok {
		return
	}
	var row models.Curriculum
	if err := r.DB.Where("id = ? AND school_id = ?", c.Param("id"), *schoolID).First(&row).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "curriculum not found"})
		return
	}
	r.DB.Model(&row).Updates(models.Curriculum{Name: input.Name, Year: input.Year, Description: input.Description, SchoolID: schoolID})
	c.JSON(http.StatusOK, gin.H{"data": row})
}

func (r *academicRepository) DeleteCurriculum(c *gin.Context) {
	schoolID, _, ok := requireTenant(c)
	if !ok {
		return
	}
	var row models.Curriculum
	if err := r.DB.Where("id = ? AND school_id = ?", c.Param("id"), *schoolID).First(&row).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "curriculum not found"})
		return
	}
	r.DB.Delete(&row)
	c.JSON(http.StatusNoContent, gin.H{"data": true})
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
	var existing models.Teaching
	dup := r.DB.Where("school_id = ? AND teacher_id = ? AND class_id = ? AND subject_id = ?", *schoolID, input.TeacherID, input.ClassID, input.SubjectID)
	if input.SemesterID == nil {
		dup = dup.Where("semester_id IS NULL")
	} else {
		dup = dup.Where("semester_id = ?", *input.SemesterID)
	}
	if err := dup.First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "teaching already exists"})
		return
	}
	if err := r.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create teaching"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": input})
}

func (r *academicRepository) UpdateTeaching(c *gin.Context) {
	var input models.Teaching
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	schoolID, _, ok := resolveWriteSchoolID(c, input.SchoolID)
	if !ok {
		return
	}
	var row models.Teaching
	if err := r.DB.Where("id = ? AND school_id = ?", c.Param("id"), *schoolID).First(&row).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "teaching not found"})
		return
	}

	var existing models.Teaching
	dup := r.DB.Where("school_id = ? AND teacher_id = ? AND class_id = ? AND subject_id = ? AND id <> ?", *schoolID, input.TeacherID, input.ClassID, input.SubjectID, row.ID)
	if input.SemesterID == nil {
		dup = dup.Where("semester_id IS NULL")
	} else {
		dup = dup.Where("semester_id = ?", *input.SemesterID)
	}
	if err := dup.First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "teaching already exists"})
		return
	}

	r.DB.Model(&row).Updates(models.Teaching{SchoolID: schoolID, TeacherID: input.TeacherID, ClassID: input.ClassID, SubjectID: input.SubjectID, SemesterID: input.SemesterID})
	c.JSON(http.StatusOK, gin.H{"data": row})
}

func (r *academicRepository) DeleteTeaching(c *gin.Context) {
	schoolID, _, ok := requireTenant(c)
	if !ok {
		return
	}
	var row models.Teaching
	if err := r.DB.Where("id = ? AND school_id = ?", c.Param("id"), *schoolID).First(&row).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "teaching not found"})
		return
	}
	r.DB.Delete(&row)
	c.JSON(http.StatusNoContent, gin.H{"data": true})
}
