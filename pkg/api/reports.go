package api

import (
	"bytes"
	"fmt"
	"golang-rest-api-template/pkg/database"
	"golang-rest-api-template/pkg/models"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
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
	Notes          string
}

type reportView struct {
	SchoolName    string
	StudentName   string
	StudentEmail  string
	StudentType   string
	Semester      int
	AcademicYear  string
	Rows          []reportRow
	Average       float64
	HomeroomTeach string
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
    <p class="muted">Laporan Hasil Belajar Siswa</p>
  </div>

  <div class="info">
    <p><strong>Nama:</strong> {{.StudentName}}</p>
    <p><strong>Email:</strong> {{.StudentEmail}}</p>
    <p><strong>Jenjang:</strong> {{.StudentType}}</p>
    <p><strong>Semester/Tahun Ajaran:</strong> {{.Semester}} / {{.AcademicYear}}</p>
  </div>

  <table>
    <thead>
      <tr>
        <th style="width: 40px;">No</th>
        <th>Mata Pelajaran</th>
        <th class="right">Pengetahuan</th>
        <th class="right">Keterampilan</th>
        <th class="right">Nilai Akhir</th>
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
        <td>{{.Notes}}</td>
      </tr>
      {{end}}
      <tr>
        <td colspan="4" class="right"><strong>Rata-rata</strong></td>
        <td class="right"><strong>{{printf "%.2f" .Average}}</strong></td>
        <td></td>
      </tr>
    </tbody>
  </table>

  <div class="footer">
    <div>
      <p class="muted">Wali Kelas</p>
      <br/><br/><p>{{.HomeroomTeach}}</p>
    </div>
    <div>
      <p class="muted">Orang Tua/Wali</p>
      <br/><br/><p>__________________</p>
    </div>
  </div>
</body>
</html>`))

// PrintReportCard godoc
// @Summary Print report card by student and term
// @Description Render printable report card (HTML) that can be converted to PDF by frontend/browser print.
// @Tags reports
// @Security ApiKeyAuth
// @Produce html
// @Param student_id path int true "Student ID"
// @Param semester query int true "Semester"
// @Param academic_year query string true "Academic year"
// @Success 200 {string} string "Printable HTML"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Student/report data not found"
// @Router /reports/students/{student_id}/print [get]
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

// PrintReportCardPDF godoc
// @Summary Export report card as PDF
// @Description Export ready-to-print report card PDF per student per term
// @Tags reports
// @Security ApiKeyAuth
// @Produce application/pdf
// @Param student_id path int true "Student ID"
// @Param semester query int true "Semester"
// @Param academic_year query string true "Academic year"
// @Success 200 {file} binary "PDF report card"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Student/report data not found"
// @Router /reports/students/{student_id}/pdf [get]
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
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 8, view.SchoolName)
	pdf.Ln(7)
	pdf.SetFont("Arial", "", 11)
	pdf.Cell(0, 7, "Laporan Hasil Belajar Siswa")
	pdf.Ln(10)

	pdf.Cell(0, 6, "Nama: "+view.StudentName)
	pdf.Ln(6)
	pdf.Cell(0, 6, "Email: "+view.StudentEmail)
	pdf.Ln(6)
	pdf.Cell(0, 6, "Jenjang: "+view.StudentType)
	pdf.Ln(6)
	pdf.Cell(0, 6, fmt.Sprintf("Semester/Tahun Ajaran: %d / %s", view.Semester, view.AcademicYear))
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 10)
	headers := []string{"No", "Mata Pelajaran", "Pengetahuan", "Keterampilan", "Nilai Akhir", "Catatan"}
	widths := []float64{10, 50, 28, 28, 25, 49}
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
		pdf.CellFormat(widths[5], 8, row.Notes, "1", 0, "L", false, 0, "")
		pdf.Ln(-1)
	}

	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(widths[0]+widths[1]+widths[2]+widths[3], 8, "Rata-rata", "1", 0, "R", false, 0, "")
	pdf.CellFormat(widths[4], 8, fmt.Sprintf("%.2f", view.Average), "1", 0, "R", false, 0, "")
	pdf.CellFormat(widths[5], 8, "", "1", 0, "L", false, 0, "")

	filename := fmt.Sprintf("raport_%d_s%d_%s.pdf", studentID, view.Semester, view.AcademicYear)
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%q", filename))

	if err := pdf.Output(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate pdf"})
		return
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
		var book models.Book
		subject := fmt.Sprintf("Mapel #%d", g.BookID)
		if err := r.DB.Where("id = ?", g.BookID).First(&book).Error(); err == nil {
			subject = book.Title
		}
		rows = append(rows, reportRow{
			No:             i + 1,
			Subject:        subject,
			KnowledgeScore: g.KnowledgeScore,
			SkillScore:     g.SkillScore,
			FinalScore:     g.FinalScore,
			Notes:          g.Notes,
		})
		total += g.FinalScore
	}

	view := reportView{
		SchoolName:    "E-Raport Internal School",
		StudentName:   student.Name,
		StudentEmail:  student.Email,
		StudentType:   student.Type,
		Semester:      semester,
		AcademicYear:  academicYear,
		Rows:          rows,
		Average:       total / float64(len(rows)),
		HomeroomTeach: "__________________",
	}

	return view, http.StatusOK, nil
}
