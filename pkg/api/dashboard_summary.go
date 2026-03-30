package api

import (
	"net/http"
	"sort"

	"golang-rest-api-template/pkg/database"
	"golang-rest-api-template/pkg/models"

	"github.com/gin-gonic/gin"
)

type dashboardRepository struct {
	DB database.Database
}

func NewDashboardRepository(db database.Database) *dashboardRepository {
	return &dashboardRepository{DB: db}
}

type dashboardSummaryResponse struct {
	Totals            dashboardTotals         `json:"totals"`
	SemesterTrends    []semesterTrendItem     `json:"semester_trends"`
	ClassCompleteness []classCompletenessItem `json:"class_completeness"`
	Recommendations   []string                `json:"recommendations"`
}

type dashboardTotals struct {
	Students               int64 `json:"students"`
	Classes                int64 `json:"classes"`
	Grades                 int64 `json:"grades"`
	Attendances            int64 `json:"attendances"`
	ReportNotes            int64 `json:"report_notes"`
	SubjectsWithoutTeacher int64 `json:"subjects_without_teacher"`
	IncompleteStudents     int64 `json:"incomplete_students"`
}

type semesterTrendItem struct {
	AcademicYear string  `json:"academic_year"`
	Semester     int     `json:"semester"`
	AvgFinal     float64 `json:"avg_final"`
	Count        int     `json:"count"`
}

type classCompletenessItem struct {
	ClassID         string `json:"class_id"`
	ClassName       string `json:"class_name"`
	Level           string `json:"level"`
	TotalStudents   int    `json:"total_students"`
	GradePct        int    `json:"grade_pct"`
	AttendancePct   int    `json:"attendance_pct"`
	ReportNotePct   int    `json:"report_note_pct"`
	CompletenessPct int    `json:"completeness_pct"`
}

func (r *dashboardRepository) Summary(c *gin.Context) {
	schoolID, _, ok := requireTenant(c)
	if !ok {
		return
	}

	resp := dashboardSummaryResponse{}

	studentsCount := r.DB.Model(&models.Student{})
	classesCount := r.DB.Model(&models.Class{})
	gradesCount := r.DB.Model(&models.Grade{})
	attCount := r.DB.Model(&models.Attendance{})
	notesCount := r.DB.Model(&models.ReportNote{})
	subjectsWithoutTeacherCount := r.DB.Model(&models.Book{}).Where("teacher_id IS NULL")
	if schoolID != nil {
		studentsCount = studentsCount.Where("school_id = ?", *schoolID)
		classesCount = classesCount.Where("school_id = ?", *schoolID)
		gradesCount = gradesCount.Where("school_id = ?", *schoolID)
		attCount = attCount.Where("school_id = ?", *schoolID)
		notesCount = notesCount.Where("school_id = ?", *schoolID)
		subjectsWithoutTeacherCount = subjectsWithoutTeacherCount.Where("school_id = ?", *schoolID)
	}

	studentsCount.Count(&resp.Totals.Students)
	classesCount.Count(&resp.Totals.Classes)
	gradesCount.Count(&resp.Totals.Grades)
	attCount.Count(&resp.Totals.Attendances)
	notesCount.Count(&resp.Totals.ReportNotes)
	subjectsWithoutTeacherCount.Count(&resp.Totals.SubjectsWithoutTeacher)

	var students []models.Student
	studentQuery := r.DB.Order("created_at asc")
	if schoolID != nil {
		studentQuery = studentQuery.Where("school_id = ?", *schoolID)
	}
	studentQuery.Find(&students)

	var classes []models.Class
	classQuery := r.DB.Order("created_at asc")
	if schoolID != nil {
		classQuery = classQuery.Where("school_id = ?", *schoolID)
	}
	classQuery.Find(&classes)

	var grades []models.Grade
	gradeQuery := r.DB.Order("created_at asc")
	if schoolID != nil {
		gradeQuery = gradeQuery.Where("school_id = ?", *schoolID)
	}
	gradeQuery.Find(&grades)

	var attendances []models.Attendance
	attendanceQuery := r.DB.Order("created_at asc")
	if schoolID != nil {
		attendanceQuery = attendanceQuery.Where("school_id = ?", *schoolID)
	}
	attendanceQuery.Find(&attendances)

	var notes []models.ReportNote
	noteQuery := r.DB.Order("created_at asc")
	if schoolID != nil {
		noteQuery = noteQuery.Where("school_id = ?", *schoolID)
	}
	noteQuery.Find(&notes)

	resp.SemesterTrends = buildSemesterTrends(grades)
	resp.ClassCompleteness = buildClassCompleteness(classes, students, grades, attendances, notes)
	resp.Totals.IncompleteStudents = countIncompleteStudents(students, grades, attendances, notes)
	resp.Recommendations = buildRecommendations(resp.Totals, resp.ClassCompleteness)

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func buildSemesterTrends(grades []models.Grade) []semesterTrendItem {
	type key struct {
		academicYear string
		semester     int
	}
	type agg struct {
		total float64
		count int
	}
	bucket := map[key]agg{}
	for _, g := range grades {
		k := key{academicYear: g.AcademicYear, semester: g.Semester}
		a := bucket[k]
		a.total += g.FinalScore
		a.count++
		bucket[k] = a
	}

	out := make([]semesterTrendItem, 0, len(bucket))
	for k, a := range bucket {
		if a.count == 0 {
			continue
		}
		out = append(out, semesterTrendItem{
			AcademicYear: k.academicYear,
			Semester:     k.semester,
			AvgFinal:     a.total / float64(a.count),
			Count:        a.count,
		})
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].AcademicYear == out[j].AcademicYear {
			return out[i].Semester < out[j].Semester
		}
		return out[i].AcademicYear < out[j].AcademicYear
	})
	return out
}

func buildClassCompleteness(classes []models.Class, students []models.Student, grades []models.Grade, attendances []models.Attendance, notes []models.ReportNote) []classCompletenessItem {
	gradeStudent := map[string]bool{}
	for _, g := range grades {
		gradeStudent[g.StudentID] = true
	}
	attendanceStudent := map[string]bool{}
	for _, a := range attendances {
		attendanceStudent[a.StudentID] = true
	}
	noteStudent := map[string]bool{}
	for _, n := range notes {
		noteStudent[n.StudentID] = true
	}

	out := make([]classCompletenessItem, 0, len(classes))
	for _, cls := range classes {
		var classStudents []models.Student
		for _, s := range students {
			if s.ClassID != nil && *s.ClassID == cls.ID {
				classStudents = append(classStudents, s)
			}
		}

		total := len(classStudents)
		if total == 0 {
			out = append(out, classCompletenessItem{ClassID: cls.ID, ClassName: cls.Name, Level: cls.Level, TotalStudents: 0})
			continue
		}

		gradeHits := 0
		attendanceHits := 0
		noteHits := 0
		for _, s := range classStudents {
			if gradeStudent[s.ID] {
				gradeHits++
			}
			if attendanceStudent[s.ID] {
				attendanceHits++
			}
			if noteStudent[s.ID] {
				noteHits++
			}
		}

		gradePct := (gradeHits * 100) / total
		attendancePct := (attendanceHits * 100) / total
		notePct := (noteHits * 100) / total
		completeness := (gradePct + attendancePct + notePct) / 3

		out = append(out, classCompletenessItem{
			ClassID:         cls.ID,
			ClassName:       cls.Name,
			Level:           cls.Level,
			TotalStudents:   total,
			GradePct:        gradePct,
			AttendancePct:   attendancePct,
			ReportNotePct:   notePct,
			CompletenessPct: completeness,
		})
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].CompletenessPct < out[j].CompletenessPct
	})
	return out
}

func buildRecommendations(t dashboardTotals, classes []classCompletenessItem) []string {
	tips := make([]string, 0, 4)
	if t.Students == 0 {
		tips = append(tips, "Mulai dari input data siswa agar proses akademik bisa berjalan.")
	}
	if t.Classes == 0 {
		tips = append(tips, "Tambahkan kelas dan tahun ajaran agar struktur rapor siap dipakai.")
	}
	if t.Grades < t.Students {
		tips = append(tips, "Fokus input nilai untuk meningkatkan coverage rapor.")
	}
	if t.Attendances < t.Students {
		tips = append(tips, "Lengkapi absensi agar ringkasan rapor lebih akurat.")
	}
	if t.ReportNotes < t.Students {
		tips = append(tips, "Isi catatan wali kelas untuk personalisasi evaluasi siswa.")
	}
	if t.SubjectsWithoutTeacher > 0 {
		tips = append(tips, "Masih ada mata pelajaran tanpa guru pengampu. Lengkapi agar pengajaran konsisten.")
	}
	if len(classes) > 0 && classes[0].TotalStudents > 0 && classes[0].CompletenessPct < 70 {
		tips = append(tips, "Prioritaskan perbaikan data pada kelas dengan completeness terendah.")
	}
	if len(tips) == 0 {
		tips = append(tips, "Data utama sudah terisi rapi. Lanjutkan quality check sebelum cetak rapor.")
	}
	if len(tips) > 4 {
		return tips[:4]
	}
	return tips
}

func countIncompleteStudents(students []models.Student, grades []models.Grade, attendances []models.Attendance, notes []models.ReportNote) int64 {
	gradeStudent := map[string]bool{}
	for _, g := range grades {
		gradeStudent[g.StudentID] = true
	}
	attendanceStudent := map[string]bool{}
	for _, a := range attendances {
		attendanceStudent[a.StudentID] = true
	}
	noteStudent := map[string]bool{}
	for _, n := range notes {
		noteStudent[n.StudentID] = true
	}
	var total int64
	for _, s := range students {
		if !gradeStudent[s.ID] || !attendanceStudent[s.ID] || !noteStudent[s.ID] {
			total++
		}
	}
	return total
}
