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

func NewGradeRepository(db database.Database, ctx *context.Context) *gradeRepository {
	return &gradeRepository{DB: db, Ctx: ctx}
}

func computeFinalScore(knowledge, skill float64) float64 {
	return (knowledge * 0.6) + (skill * 0.4)
}

// FindGrades godoc
// @Summary Get grades
// @Description Get grade list with optional filters: student_id, semester, academic_year
// @Tags grades
// @Security ApiKeyAuth
// @Produce json
// @Param offset query int false "Offset for pagination" default(0)
// @Param limit query int false "Limit for pagination" default(10)
// @Param student_id query int false "Student ID"
// @Param semester query int false "Semester"
// @Param academic_year query string false "Academic year"
// @Success 200 {array} models.Grade "Successfully retrieved list of grades"
// @Router /grades [get]
func (r *gradeRepository) FindGrades(c *gin.Context) {
	offset, limit, ok := parsePagination(c)
	if !ok {
		return
	}

	query := r.DB.Model(&models.Grade{})

	if studentID := c.Query("student_id"); studentID != "" {
		query = query.Where("student_id = ?", studentID)
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

	c.JSON(http.StatusOK, gin.H{
		"data": grades,
		"meta": gin.H{
			"offset": offset,
			"limit":  limit,
			"total":  total,
			"count":  len(grades),
		},
	})
}

// CreateGrade godoc
// @Summary Create a grade
// @Description Create grade for a student and subject(book)
// @Tags grades
// @Security ApiKeyAuth
// @Security JwtAuth
// @Accept json
// @Produce json
// @Param input body models.CreateGrade true "Create grade object"
// @Success 201 {object} models.Grade "Successfully created grade"
// @Failure 400 {string} string "Bad Request"
// @Router /grades [post]
func (r *gradeRepository) CreateGrade(c *gin.Context) {
	var input models.CreateGrade
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var student models.Student
	if err := r.DB.Where("id = ?", input.StudentID).First(&student).Error(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "student not found"})
		return
	}

	var book models.Book
	if err := r.DB.Where("id = ?", input.BookID).First(&book).Error(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "subject/book not found"})
		return
	}

	grade := models.Grade{
		SchoolID:       input.SchoolID,
		StudentID:      input.StudentID,
		BookID:         input.BookID,
		Semester:       input.Semester,
		AcademicYear:   input.AcademicYear,
		KnowledgeScore: input.KnowledgeScore,
		SkillScore:     input.SkillScore,
		FinalScore:     computeFinalScore(input.KnowledgeScore, input.SkillScore),
		Notes:          input.Notes,
	}

	if err := r.DB.Create(&grade).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create grade"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": grade})
}

// UpdateGrade godoc
// @Summary Update grade by ID
// @Tags grades
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "Grade ID"
// @Param input body models.UpdateGrade true "Update grade object"
// @Success 200 {object} models.Grade "Successfully updated grade"
// @Router /grades/{id} [put]
func (r *gradeRepository) UpdateGrade(c *gin.Context) {
	var grade models.Grade
	if err := r.DB.Where("id = ?", c.Param("id")).First(&grade).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "grade not found"})
		return
	}

	var input models.UpdateGrade
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
	if input.SchoolID != nil {
		grade.SchoolID = input.SchoolID
	}
	grade.FinalScore = computeFinalScore(grade.KnowledgeScore, grade.SkillScore)

	if err := r.DB.Model(&grade).Updates(grade).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update grade"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": grade})
}

// DeleteGrade godoc
// @Summary Delete grade by ID
// @Tags grades
// @Security ApiKeyAuth
// @Produce json
// @Param id path string true "Grade ID"
// @Success 204 {string} string "Successfully deleted grade"
// @Router /grades/{id} [delete]
func (r *gradeRepository) DeleteGrade(c *gin.Context) {
	var grade models.Grade
	if err := r.DB.Where("id = ?", c.Param("id")).First(&grade).Error(); err != nil {
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
