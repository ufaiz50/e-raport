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
	reportSummaryRepository := NewReportSummaryRepository(db)
	dashboardRepository := NewDashboardRepository(db)
	academicRepository := NewAcademicRepository(db)
	schoolRepository := NewSchoolRepository(db)
	schoolProfileRepository := NewSchoolProfileRepository(db)
	attendanceRepository := NewAttendanceRepository(db, ctx)
	reportNoteRepository := NewReportNoteRepository(db, ctx)
	enrollmentRepository := NewEnrollmentRepository(db, ctx)
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
		v1.GET("/books", middleware.APIKeyAuth(), middleware.JWTAuth(), bookRepository.FindBooks)
		v1.POST("/books", middleware.APIKeyAuth(), middleware.JWTAuth(), bookRepository.CreateBook)
		v1.GET("/books/:id", middleware.APIKeyAuth(), middleware.JWTAuth(), bookRepository.FindBook)
		v1.PUT("/books/:id", middleware.APIKeyAuth(), middleware.JWTAuth(), bookRepository.UpdateBook)
		v1.DELETE("/books/:id", middleware.APIKeyAuth(), middleware.JWTAuth(), bookRepository.DeleteBook)

		// Phase-2 alias endpoints (domain-friendly naming)
		v1.GET("/subjects", middleware.APIKeyAuth(), middleware.JWTAuth(), bookRepository.FindSubjects)
		v1.POST("/subjects", middleware.APIKeyAuth(), middleware.JWTAuth(), bookRepository.CreateSubject)
		v1.GET("/subjects/:id", middleware.APIKeyAuth(), middleware.JWTAuth(), bookRepository.FindSubject)
		v1.PUT("/subjects/:id", middleware.APIKeyAuth(), middleware.JWTAuth(), bookRepository.UpdateSubject)
		v1.DELETE("/subjects/:id", middleware.APIKeyAuth(), middleware.JWTAuth(), bookRepository.DeleteSubject)

		v1.GET("/students", middleware.APIKeyAuth(), middleware.JWTAuth(), studentRepository.FindStudents)
		v1.POST("/students", middleware.APIKeyAuth(), middleware.JWTAuth(), studentRepository.CreateStudent)
		v1.GET("/students/:id", middleware.APIKeyAuth(), middleware.JWTAuth(), studentRepository.FindStudent)
		v1.PUT("/students/:id", middleware.APIKeyAuth(), middleware.JWTAuth(), studentRepository.UpdateStudent)
		v1.DELETE("/students/:id", middleware.APIKeyAuth(), middleware.JWTAuth(), studentRepository.DeleteStudent)

		v1.GET("/classes", middleware.APIKeyAuth(), middleware.JWTAuth(), classRepository.FindClasses)
		v1.POST("/classes", middleware.APIKeyAuth(), middleware.JWTAuth(), classRepository.CreateClass)
		v1.PUT("/classes/:id", middleware.APIKeyAuth(), middleware.JWTAuth(), classRepository.UpdateClass)
		v1.DELETE("/classes/:id", middleware.APIKeyAuth(), middleware.JWTAuth(), classRepository.DeleteClass)

		v1.GET("/grades", middleware.APIKeyAuth(), middleware.JWTAuth(), gradeRepository.FindGrades)
		v1.POST("/grades", middleware.APIKeyAuth(), middleware.JWTAuth(), gradeRepository.CreateGrade)
		v1.PUT("/grades/:id", middleware.APIKeyAuth(), middleware.JWTAuth(), gradeRepository.UpdateGrade)
		v1.DELETE("/grades/:id", middleware.APIKeyAuth(), middleware.JWTAuth(), gradeRepository.DeleteGrade)

		v1.GET("/attendances", middleware.APIKeyAuth(), middleware.JWTAuth(), attendanceRepository.FindAttendances)
		v1.PUT("/attendances", middleware.APIKeyAuth(), middleware.JWTAuth(), attendanceRepository.UpsertAttendance)

		v1.GET("/report-notes", middleware.APIKeyAuth(), middleware.JWTAuth(), reportNoteRepository.FindReportNotes)
		v1.PUT("/report-notes", middleware.APIKeyAuth(), middleware.JWTAuth(), reportNoteRepository.UpsertReportNote)

		v1.GET("/enrollments", middleware.APIKeyAuth(), middleware.JWTAuth(), enrollmentRepository.FindEnrollments)
		v1.POST("/enrollments", middleware.APIKeyAuth(), middleware.JWTAuth(), enrollmentRepository.CreateEnrollment)
		v1.PUT("/enrollments/:id/close", middleware.APIKeyAuth(), middleware.JWTAuth(), enrollmentRepository.CloseEnrollment)

		v1.GET("/reports/students/:student_id/print", middleware.APIKeyAuth(), middleware.JWTAuth(), reportRepository.PrintReportCard)
		v1.GET("/reports/students/:student_id/pdf", middleware.APIKeyAuth(), middleware.JWTAuth(), reportRepository.PrintReportCardPDF)
		v1.GET("/reports/classes/:class_id/pdf", middleware.APIKeyAuth(), middleware.JWTAuth(), reportRepository.PrintReportCardClassPDF)
		v1.GET("/reports/summary", middleware.APIKeyAuth(), middleware.JWTAuth(), reportSummaryRepository.Summary)
		v1.POST("/reports/students/:student_id/finalize", middleware.APIKeyAuth(), middleware.JWTAuth(), reportRepository.FinalizeReportCard)
		v1.GET("/dashboard/summary", middleware.APIKeyAuth(), middleware.JWTAuth(), dashboardRepository.Summary)

		v1.GET("/academic-years", middleware.APIKeyAuth(), middleware.JWTAuth(), academicRepository.ListAcademicYears)
		v1.POST("/academic-years", middleware.APIKeyAuth(), middleware.JWTAuth(), academicRepository.CreateAcademicYear)
		v1.PUT("/academic-years/:id", middleware.APIKeyAuth(), middleware.JWTAuth(), academicRepository.UpdateAcademicYear)
		v1.DELETE("/academic-years/:id", middleware.APIKeyAuth(), middleware.JWTAuth(), academicRepository.DeleteAcademicYear)
		v1.GET("/semesters", middleware.APIKeyAuth(), middleware.JWTAuth(), academicRepository.ListSemesters)
		v1.POST("/semesters", middleware.APIKeyAuth(), middleware.JWTAuth(), academicRepository.CreateSemester)
		v1.PUT("/semesters/:id", middleware.APIKeyAuth(), middleware.JWTAuth(), academicRepository.UpdateSemester)
		v1.DELETE("/semesters/:id", middleware.APIKeyAuth(), middleware.JWTAuth(), academicRepository.DeleteSemester)
		v1.GET("/curriculums", middleware.APIKeyAuth(), middleware.JWTAuth(), academicRepository.ListCurriculums)
		v1.POST("/curriculums", middleware.APIKeyAuth(), middleware.JWTAuth(), academicRepository.CreateCurriculum)
		v1.PUT("/curriculums/:id", middleware.APIKeyAuth(), middleware.JWTAuth(), academicRepository.UpdateCurriculum)
		v1.DELETE("/curriculums/:id", middleware.APIKeyAuth(), middleware.JWTAuth(), academicRepository.DeleteCurriculum)
		v1.GET("/teachings", middleware.APIKeyAuth(), middleware.JWTAuth(), academicRepository.ListTeachings)
		v1.POST("/teachings", middleware.APIKeyAuth(), middleware.JWTAuth(), academicRepository.CreateTeaching)
		v1.PUT("/teachings/:id", middleware.APIKeyAuth(), middleware.JWTAuth(), academicRepository.UpdateTeaching)
		v1.DELETE("/teachings/:id", middleware.APIKeyAuth(), middleware.JWTAuth(), academicRepository.DeleteTeaching)

		v1.GET("/schools", middleware.APIKeyAuth(), middleware.JWTAuth(), middleware.RequireRoles("super_admin"), schoolRepository.ListSchools)
		v1.POST("/schools", middleware.APIKeyAuth(), middleware.JWTAuth(), middleware.RequireRoles("super_admin"), schoolRepository.CreateSchool)
		v1.PUT("/schools/:id", middleware.APIKeyAuth(), middleware.JWTAuth(), middleware.RequireRoles("super_admin"), schoolRepository.UpdateSchool)
		v1.DELETE("/schools/:id", middleware.APIKeyAuth(), middleware.JWTAuth(), middleware.RequireRoles("super_admin"), schoolRepository.DeleteSchool)

		// Preferred school profile endpoints (backed by schools table)
		v1.GET("/schools/profile", middleware.APIKeyAuth(), middleware.JWTAuth(), schoolProfileRepository.Get)
		v1.PUT("/schools/profile", middleware.APIKeyAuth(), middleware.JWTAuth(), schoolProfileRepository.Upsert)

		// Backward-compat legacy endpoints (deprecated)
		v1.GET("/school-profile", middleware.APIKeyAuth(), middleware.JWTAuth(), schoolProfileRepository.Get)
		v1.PUT("/school-profile", middleware.APIKeyAuth(), middleware.JWTAuth(), schoolProfileRepository.Upsert)

		v1.POST("/login", middleware.APIKeyAuth(), userRepository.LoginHandler)
		v1.POST("/refresh", middleware.APIKeyAuth(), userRepository.RefreshTokenHandler)
		v1.POST("/logout", middleware.APIKeyAuth(), userRepository.LogoutHandler)
		v1.POST("/register", middleware.APIKeyAuth(), userRepository.RegisterHandler)
		v1.GET("/teachers", middleware.APIKeyAuth(), middleware.JWTAuth(), userRepository.ListTeachers)
	}
	r.GET("/docs", ScalarDocs)
	r.GET("/openapi.yaml", OpenAPIDoc)

	return r
}
