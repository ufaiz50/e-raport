package api

import (
	"bytes"
	"fmt"
	"golang-rest-api-template/pkg/database"
	"golang-rest-api-template/pkg/models"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	qrcode "github.com/skip2/go-qrcode"
)

type reportRepository struct {
	DB database.Database
}

func NewReportRepository(db database.Database) *reportRepository {
	return &reportRepository{DB: db}
}

type reportRow struct {
	No             int
	Subject        string
	KnowledgeScore float64
	SkillScore     float64
	FinalScore     float64
	Predicate      string
	Notes          string
}

type reportView struct {
	SchoolName     string
	SchoolAddress  string
	SchoolNPSN     string
	PrincipalName  string
	StudentName    string
	StudentEmail   string
	StudentType    string
	ClassName      string
	Semester       int
	AcademicYear   string
	Rows           []reportRow
	Average        float64
	Rank           int
	ReportStatus   string
	HomeroomTeach  string
	FinalizedAtStr string
	SickDays       int
	PermissionDays int
	AbsentDays     int
	HomeroomNote   string
	VerificationID string
}

var reportTemplate = template.Must(template.New("report_card").Parse(`<!doctype html>
<html>
<head>
  <meta charset="utf-8"/>
  <title>Raport - {{.StudentName}}</title>
  <style>
    body { font-family: Arial, sans-serif; margin: 24px; color: #111; }
    h1,h2,h3,p { margin: 0; }
    .header { margin-bottom: 16px; }
    .muted { color: #666; font-size: 12px; }
    .info { margin: 14px 0; }
    table { width: 100%; border-collapse: collapse; margin-top: 10px; }
    th, td { border: 1px solid #333; padding: 8px; font-size: 13px; }
    th { background: #f4f4f4; text-align: left; }
    .right { text-align: right; }
    .footer { margin-top: 28px; display:flex; justify-content:space-between; }
  </style>
</head>
<body>
  <div class="header">
    <h2>{{.SchoolName}}</h2>
    <p class="muted">{{.SchoolAddress}} | NPSN: {{.SchoolNPSN}}</p>
    <p class="muted">Laporan Hasil Belajar Siswa</p>
  </div>

  <div class="info">
    <p><strong>Nama:</strong> {{.StudentName}}</p>
    <p><strong>Email:</strong> {{.StudentEmail}}</p>
    <p><strong>Kelas:</strong> {{.ClassName}}</p>
    <p><strong>Jenjang:</strong> {{.StudentType}}</p>
    <p><strong>Semester/Tahun Ajaran:</strong> {{.Semester}} / {{.AcademicYear}}</p>
    <p><strong>Rangking Kelas:</strong> {{.Rank}}</p>
    <p><strong>Status Raport:</strong> {{.ReportStatus}} {{if .FinalizedAtStr}}({{.FinalizedAtStr}}){{end}}</p>
  </div>

  <table>
    <thead>
      <tr>
        <th style="width: 40px;">No</th>
        <th>Mata Pelajaran</th>
        <th class="right">Pengetahuan</th>
        <th class="right">Keterampilan</th>
        <th class="right">Nilai Akhir</th>
        <th>Predikat</th>
        <th>Catatan</th>
      </tr>
    </thead>
    <tbody>
      {{range .Rows}}
      <tr>
        <td>{{.No}}</td>
        <td>{{.Subject}}</td>
        <td class="right">{{printf "%.2f" .KnowledgeScore}}</td>
        <td class="right">{{printf "%.2f" .SkillScore}}</td>
        <td class="right">{{printf "%.2f" .FinalScore}}</td>
        <td>{{.Predicate}}</td>
        <td>{{.Notes}}</td>
      </tr>
      {{end}}
      <tr>
        <td colspan="4" class="right"><strong>Rata-rata</strong></td>
        <td class="right"><strong>{{printf "%.2f" .Average}}</strong></td>
        <td colspan="2"></td>
      </tr>
    </tbody>
  </table>

  <div class="info">
    <p><strong>Absensi:</strong> Sakit {{.SickDays}} hari, Izin {{.PermissionDays}} hari, Alpha {{.AbsentDays}} hari</p>
    <p><strong>Catatan Wali Kelas:</strong> {{.HomeroomNote}}</p>
    <p><strong>Verifikasi:</strong> {{.VerificationID}}</p>
  </div>

  <div class="footer">
    <div>
      <p class="muted">Wali Kelas</p>
      <br/><br/><p>{{.HomeroomTeach}}</p>
    </div>
    <div>
      <p class="muted">Kepala Sekolah</p>
      <br/><br/><p>{{.PrincipalName}}</p>
    </div>
  </div>
</body>
</html>`))

func predicate(score float64) string {
	switch {
	case score >= 90:
		return "A"
	case score >= 80:
		return "B"
	case score >= 70:
		return "C"
	default:
		return "D"
	}
}

// PrintReportCard godoc
func (r *reportRepository) PrintReportCard(c *gin.Context) {
	studentID, err := strconv.Atoi(c.Param("student_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student_id"})
		return
	}

	view, statusCode, err := r.buildReportView(studentID, c)
	if err != nil {
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	var html bytes.Buffer
	if err := reportTemplate.Execute(&html, view); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to render report template"})
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html.String())
}

func (r *reportRepository) PrintReportCardPDF(c *gin.Context) {
	studentID, err := strconv.Atoi(c.Param("student_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student_id"})
		return
	}

	view, statusCode, err := r.buildReportView(studentID, c)
	if err != nil {
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	r.renderStudentReportPDF(pdf, view)

	filename := fmt.Sprintf("raport_%d_s%d_%s.pdf", studentID, view.Semester, view.AcademicYear)
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%q", filename))
	_ = pdf.Output(c.Writer)
}

func (r *reportRepository) PrintReportCardClassPDF(c *gin.Context) {
	classID, err := strconv.Atoi(c.Param("class_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid class_id"})
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

	var students []models.Student
	if err := r.DB.Where("class_id = ?", classID).Order("name asc").Find(&students).Error; err != nil || len(students) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "class students not found"})
		return
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	for _, s := range students {
		clone := *c
		clone.Params = append(clone.Params[:0], gin.Param{Key: "student_id", Value: fmt.Sprintf("%d", s.ID)})
		clone.Request = c.Request
		clone.Request.URL.RawQuery = fmt.Sprintf("semester=%d&academic_year=%s", semester, academicYear)
		view, _, err := r.buildReportView(int(s.ID), &clone)
		if err != nil {
			continue
		}
		r.renderStudentReportPDF(pdf, view)
	}

	filename := fmt.Sprintf("raport_kelas_%d_s%d_%s.pdf", classID, semester, academicYear)
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%q", filename))
	_ = pdf.Output(c.Writer)
}

func (r *reportRepository) FinalizeReportCard(c *gin.Context) {
	studentID, err := strconv.Atoi(c.Param("student_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student_id"})
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

	var grades []models.Grade
	r.DB.Where("student_id = ? AND semester = ? AND academic_year = ?", studentID, semester, academicYear).Find(&grades)
	if len(grades) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot finalize without grades"})
		return
	}

	now := time.Now()
	var reportCard models.ReportCard
	if err := r.DB.Where("student_id = ? AND semester = ? AND academic_year = ?", studentID, semester, academicYear).First(&reportCard).Error(); err != nil {
		reportCard = models.ReportCard{
			StudentID:    uint(studentID),
			Semester:     semester,
			AcademicYear: academicYear,
			Status:       models.ReportCardFinalized,
			FinalizedAt:  &now,
		}
		r.DB.Create(&reportCard)
		c.JSON(http.StatusCreated, gin.H{"data": reportCard})
		return
	}

	r.DB.Model(&reportCard).Updates(models.ReportCard{Status: models.ReportCardFinalized, FinalizedAt: &now})
	c.JSON(http.StatusOK, gin.H{"data": reportCard})
}

func (r *reportRepository) renderStudentReportPDF(pdf *gofpdf.Fpdf, view reportView) {
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 8, view.SchoolName)
	pdf.Ln(7)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, fmt.Sprintf("NPSN: %s | %s", view.SchoolNPSN, view.SchoolAddress))
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 11)
	pdf.Cell(0, 6, "Nama: "+view.StudentName+" | Kelas: "+view.ClassName)
	pdf.Ln(6)
	pdf.Cell(0, 6, fmt.Sprintf("Semester/TA: %d / %s | Rank: %d", view.Semester, view.AcademicYear, view.Rank))
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 10)
	headers := []string{"No", "Mapel", "P", "K", "Akhir", "Pred", "Catatan"}
	widths := []float64{10, 45, 18, 18, 18, 15, 66}
	for i, h := range headers {
		pdf.CellFormat(widths[i], 8, h, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 10)
	for _, row := range view.Rows {
		pdf.CellFormat(widths[0], 8, fmt.Sprintf("%d", row.No), "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[1], 8, row.Subject, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[2], 8, fmt.Sprintf("%.2f", row.KnowledgeScore), "1", 0, "R", false, 0, "")
		pdf.CellFormat(widths[3], 8, fmt.Sprintf("%.2f", row.SkillScore), "1", 0, "R", false, 0, "")
		pdf.CellFormat(widths[4], 8, fmt.Sprintf("%.2f", row.FinalScore), "1", 0, "R", false, 0, "")
		pdf.CellFormat(widths[5], 8, row.Predicate, "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[6], 8, row.Notes, "1", 0, "L", false, 0, "")
		pdf.Ln(-1)
	}

	pdf.Ln(6)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, fmt.Sprintf("Absensi: Sakit %d | Izin %d | Alpha %d", view.SickDays, view.PermissionDays, view.AbsentDays))
	pdf.Ln(6)
	pdf.MultiCell(0, 6, "Catatan Wali Kelas: "+view.HomeroomNote, "", "L", false)
	pdf.Cell(0, 6, "Verifikasi: "+view.VerificationID)

	if png, err := qrcode.Encode(view.VerificationID, qrcode.Medium, 160); err == nil {
		name := "qr_" + fmt.Sprintf("%d", time.Now().UnixNano())
		opt := gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: false}
		pdf.RegisterImageOptionsReader(name, opt, bytes.NewReader(png))
		x := 170.0
		y := pdf.GetY() - 10
		if y < 20 {
			y = 20
		}
		pdf.ImageOptions(name, x, y, 25, 25, false, opt, 0, "")
	}
}

func (r *reportRepository) buildReportView(studentID int, c *gin.Context) (reportView, int, error) {
	semester, err := parseRequiredInt(c, "semester")
	if err != nil {
		return reportView{}, http.StatusBadRequest, err
	}
	academicYear := c.Query("academic_year")
	if academicYear == "" {
		return reportView{}, http.StatusBadRequest, fmt.Errorf("academic_year is required")
	}

	var student models.Student
	if err := r.DB.Where("id = ?", studentID).First(&student).Error(); err != nil {
		return reportView{}, http.StatusNotFound, fmt.Errorf("student not found")
	}

	className := "-"
	homeroom := "__________________"
	if student.ClassID != nil {
		var class models.Class
		if err := r.DB.Where("id = ?", *student.ClassID).First(&class).Error(); err == nil {
			className = class.Name
			if class.Homeroom != "" {
				homeroom = class.Homeroom
			}
		}
	}

	var grades []models.Grade
	if err := r.DB.Where("student_id = ? AND semester = ? AND academic_year = ?", studentID, semester, academicYear).Order("book_id asc").Find(&grades).Error; err != nil {
		return reportView{}, http.StatusInternalServerError, fmt.Errorf("failed to fetch report data")
	}
	if len(grades) == 0 {
		return reportView{}, http.StatusNotFound, fmt.Errorf("report data not found for this term")
	}

	rows := make([]reportRow, 0, len(grades))
	var total float64
	for i, g := range grades {
		subject := fmt.Sprintf("Mapel #%d", g.BookID)
		var book models.Book
		if err := r.DB.Where("id = ?", g.BookID).First(&book).Error(); err == nil {
			subject = book.Title
		}
		rows = append(rows, reportRow{No: i + 1, Subject: subject, KnowledgeScore: g.KnowledgeScore, SkillScore: g.SkillScore, FinalScore: g.FinalScore, Predicate: predicate(g.FinalScore), Notes: g.Notes})
		total += g.FinalScore
	}
	avg := total / float64(len(rows))

	rank := 1
	if student.ClassID != nil {
		rank = r.computeRank(*student.ClassID, uint(studentID), semester, academicYear)
	}

	reportStatus := string(models.ReportCardDraft)
	finalizedAtStr := ""
	var rc models.ReportCard
	if err := r.DB.Where("student_id = ? AND semester = ? AND academic_year = ?", studentID, semester, academicYear).First(&rc).Error(); err == nil {
		reportStatus = string(rc.Status)
		if rc.FinalizedAt != nil {
			finalizedAtStr = rc.FinalizedAt.Format("2006-01-02 15:04")
		}
	}

	school := models.SchoolProfile{SchoolName: "E-Raport Internal School"}
	_ = r.DB.Order("id asc").First(&school).Error

	attendance := models.Attendance{}
	_ = r.DB.Where("student_id = ? AND semester = ? AND academic_year = ?", studentID, semester, academicYear).First(&attendance).Error()

	note := models.ReportNote{HomeroomComment: "-"}
	_ = r.DB.Where("student_id = ? AND semester = ? AND academic_year = ?", studentID, semester, academicYear).First(&note).Error()

	verificationID := fmt.Sprintf("ERAPORT-%d-%d-%s", studentID, semester, academicYear)

	return reportView{
		SchoolName:     school.SchoolName,
		SchoolAddress:  school.Address,
		SchoolNPSN:     school.NPSN,
		PrincipalName:  school.PrincipalName,
		StudentName:    student.Name,
		StudentEmail:   student.Email,
		StudentType:    student.Type,
		ClassName:      className,
		Semester:       semester,
		AcademicYear:   academicYear,
		Rows:           rows,
		Average:        avg,
		Rank:           rank,
		ReportStatus:   reportStatus,
		HomeroomTeach:  homeroom,
		FinalizedAtStr: finalizedAtStr,
		SickDays:       attendance.SickDays,
		PermissionDays: attendance.PermissionDays,
		AbsentDays:     attendance.AbsentDays,
		HomeroomNote:   note.HomeroomComment,
		VerificationID: verificationID,
	}, http.StatusOK, nil
}

func (r *reportRepository) computeRank(classID uint, studentID uint, semester int, academicYear string) int {
	var students []models.Student
	r.DB.Where("class_id = ?", classID).Find(&students)
	type score struct {
		studentID uint
		avg       float64
	}
	list := make([]score, 0, len(students))
	for _, s := range students {
		var grades []models.Grade
		r.DB.Where("student_id = ? AND semester = ? AND academic_year = ?", s.ID, semester, academicYear).Find(&grades)
		if len(grades) == 0 {
			continue
		}
		var total float64
		for _, g := range grades {
			total += g.FinalScore
		}
		list = append(list, score{studentID: s.ID, avg: total / float64(len(grades))})
	}
	if len(list) == 0 {
		return 1
	}
	sort.SliceStable(list, func(i, j int) bool { return list[i].avg > list[j].avg })
	for i, item := range list {
		if item.studentID == studentID {
			return i + 1
		}
	}
	return len(list)
}
