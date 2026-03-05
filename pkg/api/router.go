package api

import (
	"context"
	"golang-rest-api-template/pkg/cache"
	"golang-rest-api-template/pkg/database"
	"golang-rest-api-template/pkg/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"golang.org/x/time/rate"
)

func ContextMiddleware(bookRepository BookRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("appCtx", bookRepository)
		c.Next()
	}
}

func NewRouter(logger *zap.Logger, mongoCollection *mongo.Collection, db database.Database, redisClient cache.Cache, ctx *context.Context) *gin.Engine {
	bookRepository := NewBookRepository(db, redisClient, ctx)
	studentRepository := NewStudentRepository(db, ctx)
	classRepository := NewClassRepository(db, ctx)
	gradeRepository := NewGradeRepository(db, ctx)
	reportRepository := NewReportRepository(db)
	schoolProfileRepository := NewSchoolProfileRepository(db)
	attendanceRepository := NewAttendanceRepository(db, ctx)
	reportNoteRepository := NewReportNoteRepository(db, ctx)
	userRepository := NewUserRepository(db, ctx)

	r := gin.Default()
	r.Use(ContextMiddleware(bookRepository))

	//r.Use(gin.Logger())
	r.Use(middleware.Logger(logger, mongoCollection))
	if gin.Mode() == gin.ReleaseMode {
		r.Use(middleware.Security())
		r.Use(middleware.Xss())
	}
	r.Use(middleware.Cors())
	r.Use(middleware.RateLimiter(rate.Every(1*time.Minute), 60)) // 60 requests per minute

	v1 := r.Group("/api/v1")
	{
		v1.GET("/", bookRepository.Healthcheck)
		v1.GET("/books", middleware.APIKeyAuth(), bookRepository.FindBooks)
		v1.POST("/books", middleware.APIKeyAuth(), middleware.JWTAuth(), bookRepository.CreateBook)
		v1.GET("/books/:id", middleware.APIKeyAuth(), bookRepository.FindBook)
		v1.PUT("/books/:id", middleware.APIKeyAuth(), bookRepository.UpdateBook)
		v1.DELETE("/books/:id", middleware.APIKeyAuth(), bookRepository.DeleteBook)

		v1.GET("/students", middleware.APIKeyAuth(), studentRepository.FindStudents)
		v1.POST("/students", middleware.APIKeyAuth(), middleware.JWTAuth(), studentRepository.CreateStudent)
		v1.GET("/students/:id", middleware.APIKeyAuth(), studentRepository.FindStudent)
		v1.PUT("/students/:id", middleware.APIKeyAuth(), studentRepository.UpdateStudent)
		v1.DELETE("/students/:id", middleware.APIKeyAuth(), studentRepository.DeleteStudent)

		v1.GET("/classes", middleware.APIKeyAuth(), classRepository.FindClasses)
		v1.POST("/classes", middleware.APIKeyAuth(), middleware.JWTAuth(), classRepository.CreateClass)
		v1.PUT("/classes/:id", middleware.APIKeyAuth(), classRepository.UpdateClass)
		v1.DELETE("/classes/:id", middleware.APIKeyAuth(), classRepository.DeleteClass)

		v1.GET("/grades", middleware.APIKeyAuth(), gradeRepository.FindGrades)
		v1.POST("/grades", middleware.APIKeyAuth(), middleware.JWTAuth(), gradeRepository.CreateGrade)
		v1.PUT("/grades/:id", middleware.APIKeyAuth(), gradeRepository.UpdateGrade)
		v1.DELETE("/grades/:id", middleware.APIKeyAuth(), gradeRepository.DeleteGrade)

		v1.GET("/attendances", middleware.APIKeyAuth(), attendanceRepository.FindAttendances)
		v1.PUT("/attendances", middleware.APIKeyAuth(), middleware.JWTAuth(), attendanceRepository.UpsertAttendance)

		v1.GET("/report-notes", middleware.APIKeyAuth(), reportNoteRepository.FindReportNotes)
		v1.PUT("/report-notes", middleware.APIKeyAuth(), middleware.JWTAuth(), reportNoteRepository.UpsertReportNote)

		v1.GET("/reports/students/:student_id/print", middleware.APIKeyAuth(), reportRepository.PrintReportCard)
		v1.GET("/reports/students/:student_id/pdf", middleware.APIKeyAuth(), reportRepository.PrintReportCardPDF)
		v1.GET("/reports/classes/:class_id/pdf", middleware.APIKeyAuth(), reportRepository.PrintReportCardClassPDF)
		v1.POST("/reports/students/:student_id/finalize", middleware.APIKeyAuth(), middleware.JWTAuth(), reportRepository.FinalizeReportCard)

		v1.GET("/school-profile", middleware.APIKeyAuth(), schoolProfileRepository.Get)
		v1.PUT("/school-profile", middleware.APIKeyAuth(), middleware.JWTAuth(), schoolProfileRepository.Upsert)

		v1.POST("/login", middleware.APIKeyAuth(), userRepository.LoginHandler)
		v1.POST("/register", middleware.APIKeyAuth(), userRepository.RegisterHandler)
	}
	r.GET("/docs", ScalarDocs)
	r.StaticFile("/openapi.yaml", "./docs/openapi.yaml")

	return r
}
