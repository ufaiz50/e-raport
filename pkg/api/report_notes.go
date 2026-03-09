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
	schoolID, _, ok := requireTenant(c)
	if !ok {
		return
	}

	offset, limit, ok := parsePagination(c)
	if !ok {
		return
	}

	var notes []models.ReportNote
	query := r.DB.Model(&models.ReportNote{})
	if schoolID != nil {
		query = query.Where("school_id = ?", *schoolID)
	}
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count report notes"})
		return
	}

	query.Offset(offset).Limit(limit).Order("id desc").Find(&notes)
	c.JSON(http.StatusOK, gin.H{
		"data": notes,
		"meta": gin.H{
			"offset": offset,
			"limit":  limit,
			"total":  total,
			"count":  len(notes),
		},
	})
}

func (r *reportNoteRepository) UpsertReportNote(c *gin.Context) {
	schoolID, _, ok := requireTenant(c)
	if !ok {
		return
	}

	var input models.UpsertReportNote
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var student models.Student
	if err := r.DB.Where("id = ? AND school_id = ?", input.StudentID, *schoolID).First(&student).Error(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "student not found"})
		return
	}

	var note models.ReportNote
	if err := r.DB.Where("school_id = ? AND student_id = ? AND semester = ? AND academic_year = ?", *schoolID, input.StudentID, input.Semester, input.AcademicYear).First(&note).Error(); err != nil {
		note = models.ReportNote{
			SchoolID:        schoolID,
			StudentID:       input.StudentID,
			Semester:        input.Semester,
			AcademicYear:    input.AcademicYear,
			HomeroomComment: input.HomeroomComment,
		}
		r.DB.Create(&note)
		c.JSON(http.StatusCreated, gin.H{"data": note})
		return
	}

	r.DB.Model(&note).Updates(models.ReportNote{SchoolID: schoolID, HomeroomComment: input.HomeroomComment})
	c.JSON(http.StatusOK, gin.H{"data": note})
}
