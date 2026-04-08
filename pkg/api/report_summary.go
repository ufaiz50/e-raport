package api

import (
	"net/http"
	"sort"

	"golang-rest-api-template/pkg/database"
	"golang-rest-api-template/pkg/models"

	"github.com/gin-gonic/gin"
)

type reportSummaryRepository struct {
	DB database.Database
}

func NewReportSummaryRepository(db database.Database) *reportSummaryRepository {
	return &reportSummaryRepository{DB: db}
}

type reportSummaryStudentItem struct {
	StudentID       string `json:"student_id"`
	StudentName     string `json:"student_name"`
	ClassID         string `json:"class_id"`
	ClassName       string `json:"class_name"`
	HasGrades       bool   `json:"has_grades"`
	HasAttendance   bool   `json:"has_attendance"`
	HasReportNote   bool   `json:"has_report_note"`
	Finalized       bool   `json:"finalized"`
	CompletenessPct int    `json:"completeness_pct"`
}

type reportSummaryClassItem struct {
	ClassID           string `json:"class_id"`
	ClassName         string `json:"class_name"`
	Level             string `json:"level"`
	TotalStudents     int    `json:"total_students"`
	ReadyStudents     int    `json:"ready_students"`
	FinalizedStudents int    `json:"finalized_students"`
	CompletenessPct   int    `json:"completeness_pct"`
}

type reportSummaryResponse struct {
	AcademicYear string                     `json:"academic_year"`
	Semester     int                        `json:"semester"`
	Totals       map[string]int             `json:"totals"`
	Classes      []reportSummaryClassItem   `json:"classes"`
	Students     []reportSummaryStudentItem `json:"students"`
}

func studentDisplayName(s models.Student) string {
	if s.Nama != "" {
		return s.Nama
	}
	if s.NamaPanggilan != "" {
		return s.NamaPanggilan
	}
	if s.Email != "" {
		return s.Email
	}
	if s.NIS != "" {
		return s.NIS
	}
	return "Siswa"
}

func (r *reportSummaryRepository) Summary(c *gin.Context) {
	schoolID, _, ok := requireTenant(c)
	if !ok {
		return
	}

	semester, err := parseRequiredInt(c, "semester")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	academicYear := c.Query("academic_year")
	if academicYear == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "academic_year is required"})
		return
	}

	var enrollments []models.StudentEnrollment
	enrollmentQuery := r.DB.Where("school_id = ? AND academic_year = ? AND semester = ?", *schoolID, academicYear, semester).Order("created_at asc")
	enrollmentQuery.Find(&enrollments)

	var students []models.Student
	studentQuery := r.DB.Where("school_id = ?", *schoolID).Order("created_at asc")
	studentQuery.Find(&students)
	studentMap := map[string]models.Student{}
	for _, s := range students {
		studentMap[s.ID] = s
	}

	var classes []models.Class
	classQuery := r.DB.Where("school_id = ?", *schoolID).Order("created_at asc")
	classQuery.Find(&classes)
	classMap := map[string]models.Class{}
	for _, cls := range classes {
		classMap[cls.ID] = cls
	}

	var grades []models.Grade
	r.DB.Where("school_id = ? AND academic_year = ? AND semester = ?", *schoolID, academicYear, semester).Find(&grades)
	gradeStudent := map[string]bool{}
	for _, g := range grades {
		gradeStudent[g.StudentID] = true
	}

	var attendances []models.Attendance
	r.DB.Where("school_id = ? AND academic_year = ? AND semester = ?", *schoolID, academicYear, semester).Find(&attendances)
	attendanceStudent := map[string]bool{}
	for _, a := range attendances {
		attendanceStudent[a.StudentID] = true
	}

	var notes []models.ReportNote
	r.DB.Where("school_id = ? AND academic_year = ? AND semester = ?", *schoolID, academicYear, semester).Find(&notes)
	noteStudent := map[string]bool{}
	for _, n := range notes {
		noteStudent[n.StudentID] = true
	}

	var cards []models.ReportCard
	r.DB.Where("school_id = ? AND academic_year = ? AND semester = ?", *schoolID, academicYear, semester).Find(&cards)
	finalizedStudent := map[string]bool{}
	for _, rc := range cards {
		if rc.Status == models.ReportCardFinalized {
			finalizedStudent[rc.StudentID] = true
		}
	}

	studentItems := make([]reportSummaryStudentItem, 0, len(enrollments))
	classAgg := map[string]*reportSummaryClassItem{}
	for _, e := range enrollments {
		s := studentMap[e.StudentID]
		cls := classMap[e.ClassID]
		studentName := studentDisplayName(s)
		hasGrades := gradeStudent[e.StudentID]
		hasAttendance := attendanceStudent[e.StudentID]
		hasNote := noteStudent[e.StudentID]
		finalized := finalizedStudent[e.StudentID]
		completeness := 0
		if hasGrades {
			completeness += 34
		}
		if hasAttendance {
			completeness += 33
		}
		if hasNote {
			completeness += 33
		}
		item := reportSummaryStudentItem{
			StudentID:       e.StudentID,
			StudentName:     studentName,
			ClassID:         e.ClassID,
			ClassName:       cls.Name,
			HasGrades:       hasGrades,
			HasAttendance:   hasAttendance,
			HasReportNote:   hasNote,
			Finalized:       finalized,
			CompletenessPct: completeness,
		}
		studentItems = append(studentItems, item)

		if classAgg[e.ClassID] == nil {
			classAgg[e.ClassID] = &reportSummaryClassItem{ClassID: e.ClassID, ClassName: cls.Name, Level: cls.Level}
		}
		agg := classAgg[e.ClassID]
		agg.TotalStudents++
		if hasGrades && hasAttendance && hasNote {
			agg.ReadyStudents++
		}
		if finalized {
			agg.FinalizedStudents++
		}
	}

	classItems := make([]reportSummaryClassItem, 0, len(classAgg))
	for _, agg := range classAgg {
		if agg.TotalStudents > 0 {
			agg.CompletenessPct = (agg.ReadyStudents * 100) / agg.TotalStudents
		}
		classItems = append(classItems, *agg)
	}

	sort.Slice(classItems, func(i, j int) bool { return classItems[i].CompletenessPct < classItems[j].CompletenessPct })
	sort.Slice(studentItems, func(i, j int) bool {
		if studentItems[i].ClassName == studentItems[j].ClassName {
			return studentItems[i].StudentName < studentItems[j].StudentName
		}
		return studentItems[i].ClassName < studentItems[j].ClassName
	})

	resp := reportSummaryResponse{
		AcademicYear: academicYear,
		Semester:     semester,
		Totals: map[string]int{
			"students":           len(studentItems),
			"ready_students":     countReadyStudents(studentItems),
			"finalized_students": countFinalizedStudents(studentItems),
			"classes":            len(classItems),
		},
		Classes:  classItems,
		Students: studentItems,
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func countReadyStudents(items []reportSummaryStudentItem) int {
	total := 0
	for _, item := range items {
		if item.HasGrades && item.HasAttendance && item.HasReportNote {
			total++
		}
	}
	return total
}

func countFinalizedStudents(items []reportSummaryStudentItem) int {
	total := 0
	for _, item := range items {
		if item.Finalized {
			total++
		}
	}
	return total
}
