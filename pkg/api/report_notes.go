package api

import (
	"context"
	"golang-rest-api-template/pkg/database"
	"golang-rest-api-template/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type reportNoteRepository struct {
	DB  database.Database
	Ctx *context.Context
}

func NewReportNoteRepository(db database.Database, ctx *context.Context) *reportNoteRepository {
	return &reportNoteRepository{DB: db, Ctx: ctx}
}

func (r *reportNoteRepository) FindReportNotes(c *gin.Context) {
	var notes []models.ReportNote
	query := r.DB
	if studentID := c.Query("student_id"); studentID != "" {
		query = query.Where("student_id = ?", studentID)
	}
	if semester := c.Query("semester"); semester != "" {
		query = query.Where("semester = ?", semester)
	}
	if academicYear := c.Query("academic_year"); academicYear != "" {
		query = query.Where("academic_year = ?", academicYear)
	}
	query.Order("id desc").Find(&notes)
	c.JSON(http.StatusOK, gin.H{"data": notes})
}

func (r *reportNoteRepository) UpsertReportNote(c *gin.Context) {
	var input models.UpsertReportNote
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var student models.Student
	if err := r.DB.Where("id = ?", input.StudentID).First(&student).Error(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "student not found"})
		return
	}

	var note models.ReportNote
	if err := r.DB.Where("student_id = ? AND semester = ? AND academic_year = ?", input.StudentID, input.Semester, input.AcademicYear).First(&note).Error(); err != nil {
		note = models.ReportNote{
			StudentID:       input.StudentID,
			Semester:        input.Semester,
			AcademicYear:    input.AcademicYear,
			HomeroomComment: input.HomeroomComment,
		}
		r.DB.Create(&note)
		c.JSON(http.StatusCreated, gin.H{"data": note})
		return
	}

	r.DB.Model(&note).Updates(models.ReportNote{HomeroomComment: input.HomeroomComment})
	c.JSON(http.StatusOK, gin.H{"data": note})
}
