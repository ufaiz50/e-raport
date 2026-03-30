package api

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

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

var academicYearPattern = regexp.MustCompile(`^(\d{4})/(\d{4})$`)

func validateAcademicYearFormat(year string) error {
	year = strings.TrimSpace(year)
	matches := academicYearPattern.FindStringSubmatch(year)
	if len(matches) != 3 {
		return fmt.Errorf("academic year must use format YYYY/YYYY")
	}
	if matches[2] != fmt.Sprintf("%04d", atoiSafe(matches[1])+1) {
		return fmt.Errorf("academic year range is invalid")
	}
	return nil
}

func atoiSafe(s string) int {
	var n int
	fmt.Sscanf(s, "%d", &n)
	return n
}

func (r *academicRepository) deactivateOtherAcademicYears(schoolID *string, keepID string) {
	r.DB.Model(&models.AcademicYear{}).Where("school_id = ? AND id <> ?", *schoolID, keepID).Updates(map[string]interface{}{"is_active": false})
}

func (r *academicRepository) deactivateOtherSemesters(schoolID *string, keepID string) {
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
	if err := query.Offset(offset).Limit(limit).Order("created_at desc").Find(&rows).Error; err != nil {
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
	if err := validateAcademicYearFormat(input.Year); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var existing models.AcademicYear
	if err := r.DB.Where("school_id = ? AND year = ?", *schoolID, input.Year).First(&existing).Error(); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "academic year already exists"})
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
	if err := validateAcademicYearFormat(input.Year); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var row models.AcademicYear
	if err := whereByIDOrUUID(r.DB, c.Param("id"), schoolID).First(&row).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "academic year not found"})
		return
	}
	var existing models.AcademicYear
	if err := r.DB.Where("school_id = ? AND year = ? AND id <> ?", *schoolID, input.Year, row.ID).First(&existing).Error(); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "academic year already exists"})
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
	if err := whereByIDOrUUID(r.DB, c.Param("id"), schoolID).First(&row).Error(); err != nil {
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
	query.Offset(offset).Limit(limit).Order("created_at desc").Find(&rows)
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
	if strings.TrimSpace(input.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "semester name is required"})
		return
	}
	if input.OrderNo < 1 || input.OrderNo > 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "semester order must be 1 or 2"})
		return
	}
	var year models.AcademicYear
	if err := r.DB.Where("id = ? AND school_id = ?", input.AcademicYearID, *schoolID).First(&year).Error(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "academic year not found"})
		return
	}
	var existing models.Semester
	if err := r.DB.Where("school_id = ? AND academic_year_id = ? AND order_no = ?", *schoolID, input.AcademicYearID, input.OrderNo).First(&existing).Error(); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "semester already exists for this academic year"})
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
	if strings.TrimSpace(input.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "semester name is required"})
		return
	}
	if input.OrderNo < 1 || input.OrderNo > 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "semester order must be 1 or 2"})
		return
	}
	var row models.Semester
	if err := whereByIDOrUUID(r.DB, c.Param("id"), schoolID).First(&row).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "semester not found"})
		return
	}
	var year models.AcademicYear
	if err := r.DB.Where("id = ? AND school_id = ?", input.AcademicYearID, *schoolID).First(&year).Error(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "academic year not found"})
		return
	}
	var existing models.Semester
	if err := r.DB.Where("school_id = ? AND academic_year_id = ? AND order_no = ? AND id <> ?", *schoolID, input.AcademicYearID, input.OrderNo, row.ID).First(&existing).Error(); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "semester already exists for this academic year"})
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
	if err := whereByIDOrUUID(r.DB, c.Param("id"), schoolID).First(&row).Error(); err != nil {
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
	query.Offset(offset).Limit(limit).Order("created_at desc").Find(&rows)
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
	if strings.TrimSpace(input.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "curriculum name is required"})
		return
	}
	input.SchoolID = schoolID
	var existing models.Curriculum
	if err := r.DB.Where("school_id = ? AND name = ? AND year = ?", *schoolID, input.Name, input.Year).First(&existing).Error(); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "curriculum already exists"})
		return
	}
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
	if strings.TrimSpace(input.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "curriculum name is required"})
		return
	}
	var row models.Curriculum
	if err := whereByIDOrUUID(r.DB, c.Param("id"), schoolID).First(&row).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "curriculum not found"})
		return
	}
	var existing models.Curriculum
	if err := r.DB.Where("school_id = ? AND name = ? AND year = ? AND id <> ?", *schoolID, input.Name, input.Year, row.ID).First(&existing).Error(); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "curriculum already exists"})
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
	if err := whereByIDOrUUID(r.DB, c.Param("id"), schoolID).First(&row).Error(); err != nil {
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
	query.Offset(offset).Limit(limit).Order("created_at desc").Find(&rows)
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
	var teacher models.User
	if err := r.DB.Where("id = ? AND school_id = ? AND role = ?", input.TeacherID, *schoolID, "guru").First(&teacher).Error(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "teacher not found"})
		return
	}
	var class models.Class
	if err := r.DB.Where("id = ? AND school_id = ?", input.ClassID, *schoolID).First(&class).Error(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "class not found"})
		return
	}
	var subject models.Book
	if err := r.DB.Where("id = ? AND school_id = ?", input.SubjectID, *schoolID).First(&subject).Error(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "subject not found"})
		return
	}
	if input.SemesterID != nil {
		var semester models.Semester
		if err := r.DB.Where("id = ? AND school_id = ?", *input.SemesterID, *schoolID).First(&semester).Error(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "semester not found"})
			return
		}
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
	var teacher models.User
	if err := r.DB.Where("id = ? AND school_id = ? AND role = ?", input.TeacherID, *schoolID, "guru").First(&teacher).Error(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "teacher not found"})
		return
	}
	var class models.Class
	if err := r.DB.Where("id = ? AND school_id = ?", input.ClassID, *schoolID).First(&class).Error(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "class not found"})
		return
	}
	var subject models.Book
	if err := r.DB.Where("id = ? AND school_id = ?", input.SubjectID, *schoolID).First(&subject).Error(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "subject not found"})
		return
	}
	if input.SemesterID != nil {
		var semester models.Semester
		if err := r.DB.Where("id = ? AND school_id = ?", *input.SemesterID, *schoolID).First(&semester).Error(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "semester not found"})
			return
		}
	}
	var row models.Teaching
	if err := whereByIDOrUUID(r.DB, c.Param("id"), schoolID).First(&row).Error(); err != nil {
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
	if err := whereByIDOrUUID(r.DB, c.Param("id"), schoolID).First(&row).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "teaching not found"})
		return
	}
	r.DB.Delete(&row)
	c.JSON(http.StatusNoContent, gin.H{"data": true})
}
