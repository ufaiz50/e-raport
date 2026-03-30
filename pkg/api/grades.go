package api

import (
	"context"
	"fmt"
	"golang-rest-api-template/pkg/database"
	"golang-rest-api-template/pkg/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GradeRepository interface {
	FindGrades(c *gin.Context)
	CreateGrade(c *gin.Context)
	UpdateGrade(c *gin.Context)
	DeleteGrade(c *gin.Context)
}

type gradeRepository struct {
	DB  database.Database
	Ctx *context.Context
}

type gradeResponse struct {
	models.Grade
	SubjectID   string `json:"subject_id"`
	SubjectName string `json:"subject_name,omitempty"`
}

func NewGradeRepository(db database.Database, ctx *context.Context) *gradeRepository {
	return &gradeRepository{DB: db, Ctx: ctx}
}

func computeFinalScore(knowledge, skill float64) float64 {
	return (knowledge * 0.6) + (skill * 0.4)
}

func resolveTermBySemesterID(r *gradeRepository, schoolID *string, semesterID *string) (*models.Semester, *models.AcademicYear, error) {
	if semesterID == nil || *semesterID == "" {
		return nil, nil, nil
	}
	var sem models.Semester
	if err := r.DB.Where("id = ? AND school_id = ?", *semesterID, *schoolID).First(&sem).Error(); err != nil {
		return nil, nil, err
	}
	var ay models.AcademicYear
	if err := r.DB.Where("id = ? AND school_id = ?", sem.AcademicYearID, *schoolID).First(&ay).Error(); err != nil {
		return nil, nil, err
	}
	return &sem, &ay, nil
}

func (r *gradeRepository) FindGrades(c *gin.Context) {
	schoolID, _, ok := requireTenant(c)
	if !ok {
		return
	}

	offset, limit, ok := parsePagination(c)
	if !ok {
		return
	}

	query := r.DB.Model(&models.Grade{})
	if schoolID != nil {
		query = query.Where("school_id = ?", *schoolID)
	}

	if studentID := c.Query("student_id"); studentID != "" {
		resolvedStudentID, err := resolveStudentID(r.DB, schoolID, studentID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		query = query.Where("student_id = ?", resolvedStudentID)
	}
	if semester := c.Query("semester"); semester != "" {
		query = query.Where("semester = ?", semester)
	}
	if academicYear := c.Query("academic_year"); academicYear != "" {
		query = query.Where("academic_year = ?", academicYear)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count grades"})
		return
	}

	var grades []models.Grade
	if err := query.Offset(offset).Limit(limit).Order("book_id asc").Find(&grades).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch grades"})
		return
	}

	responses := make([]gradeResponse, 0, len(grades))
	for _, grade := range grades {
		resp := gradeResponse{Grade: grade, SubjectID: grade.BookID}
		var book models.Book
		if err := r.DB.Where("id = ?", grade.BookID).First(&book).Error(); err == nil {
			resp.SubjectName = book.Title
		}
		responses = append(responses, resp)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": responses,
		"meta": gin.H{
			"offset": offset,
			"limit":  limit,
			"total":  total,
			"count":  len(responses),
		},
	})
}

func (r *gradeRepository) CreateGrade(c *gin.Context) {
	var input models.CreateGrade
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schoolID, _, ok := resolveWriteSchoolID(c, input.SchoolID)
	if !ok {
		return
	}

	var student models.Student
	if err := r.DB.Where("id = ? AND school_id = ?", input.StudentID, *schoolID).First(&student).Error(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "student not found"})
		return
	}

	var book models.Book
	if err := r.DB.Where("id = ? AND school_id = ?", input.BookID, *schoolID).First(&book).Error(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "subject/book not found"})
		return
	}

	sem, ay, err := resolveTermBySemesterID(r, schoolID, input.SemesterID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid semester_id"})
		return
	}
	if sem != nil && ay != nil {
		input.Semester = sem.OrderNo
		input.AcademicYear = ay.Year
	}

	grade := models.Grade{
		SchoolID:       schoolID,
		SemesterID:     input.SemesterID,
		TeachingID:     input.TeachingID,
		EnrollmentID:   input.EnrollmentID,
		StudentID:      input.StudentID,
		BookID:         input.BookID,
		Semester:       input.Semester,
		AcademicYear:   input.AcademicYear,
		KnowledgeScore: input.KnowledgeScore,
		SkillScore:     input.SkillScore,
		FinalScore:     computeFinalScore(input.KnowledgeScore, input.SkillScore),
		Notes:          input.Notes,
	}

	enrollment, err := resolveEnrollmentForTerm(r.DB, schoolID, input.EnrollmentID, input.StudentID, input.AcademicYear, input.Semester)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	grade.EnrollmentID = &enrollment.ID
	if input.TeachingID != nil {
		var teaching models.Teaching
		if err := r.DB.Where("id = ? AND school_id = ?", *input.TeachingID, *schoolID).First(&teaching).Error(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "teaching not found"})
			return
		}
		if teaching.SubjectID != input.BookID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "teaching subject mismatch"})
			return
		}
	}

	if err := r.DB.Create(&grade).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create grade"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": grade})
}

func (r *gradeRepository) UpdateGrade(c *gin.Context) {
	var input models.UpdateGrade
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schoolID, _, ok := resolveWriteSchoolID(c, input.SchoolID)
	if !ok {
		return
	}

	var grade models.Grade
	if err := whereByIDOrUUID(r.DB, c.Param("id"), schoolID).First(&grade).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "grade not found"})
		return
	}

	if input.KnowledgeScore != nil {
		grade.KnowledgeScore = *input.KnowledgeScore
	}
	if input.SkillScore != nil {
		grade.SkillScore = *input.SkillScore
	}
	if input.Notes != nil {
		grade.Notes = *input.Notes
	}
	if input.SemesterID != nil {
		grade.SemesterID = input.SemesterID
		sem, ay, err := resolveTermBySemesterID(r, schoolID, input.SemesterID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid semester_id"})
			return
		}
		if sem != nil && ay != nil {
			grade.Semester = sem.OrderNo
			grade.AcademicYear = ay.Year
		}
	}
	if input.TeachingID != nil {
		grade.TeachingID = input.TeachingID
	}
	if input.EnrollmentID != nil {
		grade.EnrollmentID = input.EnrollmentID
	}
	grade.SchoolID = schoolID
	enrollment, err := resolveEnrollmentForTerm(r.DB, schoolID, grade.EnrollmentID, grade.StudentID, grade.AcademicYear, grade.Semester)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	grade.EnrollmentID = &enrollment.ID
	if grade.TeachingID != nil {
		var teaching models.Teaching
		if err := r.DB.Where("id = ? AND school_id = ?", *grade.TeachingID, *schoolID).First(&teaching).Error(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "teaching not found"})
			return
		}
		if teaching.SubjectID != grade.BookID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "teaching subject mismatch"})
			return
		}
	}
	grade.FinalScore = computeFinalScore(grade.KnowledgeScore, grade.SkillScore)

	if err := r.DB.Model(&grade).Updates(grade).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update grade"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": grade})
}

func (r *gradeRepository) DeleteGrade(c *gin.Context) {
	schoolID, _, ok := requireTenant(c)
	if !ok {
		return
	}

	var grade models.Grade
	if err := whereByIDOrUUID(r.DB, c.Param("id"), schoolID).First(&grade).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "grade not found"})
		return
	}

	r.DB.Delete(&grade)
	c.JSON(http.StatusNoContent, gin.H{"data": true})
}

func parseRequiredInt(c *gin.Context, key string) (int, error) {
	value := c.Query(key)
	if value == "" {
		return 0, fmt.Errorf("%s is required", key)
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s is invalid", key)
	}
	return parsed, nil
}
