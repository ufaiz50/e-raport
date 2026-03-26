package api

import (
	"net/http"
	"time"

	"golang-rest-api-template/pkg/models"

	"github.com/gin-gonic/gin"
)

type subjectResponse struct {
	ID        uint         `json:"id"`
	Name      string       `json:"name"`
	Title     string       `json:"title"`
	TeacherID *uint        `json:"teacher_id,omitempty"`
	Teacher   *models.User `json:"teacher,omitempty"`
	SchoolID  *uint        `json:"school_id,omitempty"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type subjectPayload struct {
	Name      string `json:"name"`
	Title     string `json:"title"`
	TeacherID *uint  `json:"teacher_id"`
	SchoolID  *uint  `json:"school_id"`
}

func buildSubjectResponse(book models.Book) subjectResponse {
	return subjectResponse{
		ID:        book.ID,
		Name:      book.Title,
		Title:     book.Title,
		TeacherID: book.TeacherID,
		Teacher:   book.Teacher,
		SchoolID:  book.SchoolID,
		CreatedAt: book.CreatedAt,
		UpdatedAt: book.UpdatedAt,
	}
}

func (r *bookRepository) FindSubjects(c *gin.Context) {
	schoolID, _, ok := requireTenant(c)
	if !ok {
		return
	}
	offset, limit, ok := parsePagination(c)
	if !ok {
		return
	}

	query := r.DB.Model(&models.Book{})
	if schoolID != nil {
		query = query.Where("school_id = ?", *schoolID)
	}
	if uuid := c.Query("uuid"); uuid != "" {
		query = query.Where("uuid = ?", uuid)
	}
	if teacherID := c.Query("teacher_id"); teacherID != "" {
		query = query.Where("teacher_id = ?", teacherID)
	}
	if q := c.Query("q"); q != "" {
		query = query.Where("title ILIKE ?", "%"+q+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count subjects"})
		return
	}

	var books []models.Book
	if err := query.Offset(offset).Limit(limit).Order("title asc").Find(&books).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch subjects"})
		return
	}

	for i := range books {
		if books[i].TeacherID != nil {
			var teacher models.User
			teacherQuery := r.DB.Where("id = ?", *books[i].TeacherID)
			if schoolID != nil {
				teacherQuery = teacherQuery.Where("school_id = ?", *schoolID)
			}
			if err := teacherQuery.First(&teacher).Error(); err == nil {
				teacher.Password = ""
				books[i].Teacher = &teacher
			}
		}
	}

	items := make([]subjectResponse, 0, len(books))
	for _, book := range books {
		items = append(items, buildSubjectResponse(book))
	}

	c.JSON(http.StatusOK, gin.H{"data": items, "meta": gin.H{"offset": offset, "limit": limit, "total": total, "count": len(items)}})
}

func (r *bookRepository) CreateSubject(c *gin.Context) {
	var input subjectPayload
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	name := input.Name
	if name == "" {
		name = input.Title
	}
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}
	schoolID, _, ok := resolveWriteSchoolID(c, input.SchoolID)
	if !ok {
		return
	}
	author := "-"
	if input.TeacherID != nil {
		var teacher models.User
		teacherQuery := r.DB.Where("id = ? AND role = ?", *input.TeacherID, "guru")
		if schoolID != nil {
			teacherQuery = teacherQuery.Where("school_id = ?", *schoolID)
		}
		if err := teacherQuery.First(&teacher).Error(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "teacher not found"})
			return
		}
		author = teacher.Username
	}
	book := models.Book{Title: name, Author: author, SchoolID: schoolID, TeacherID: input.TeacherID}
	if err := r.DB.Create(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create subject"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": buildSubjectResponse(book)})
}

func (r *bookRepository) FindSubject(c *gin.Context) {
	schoolID, _, ok := requireTenant(c)
	if !ok {
		return
	}
	var book models.Book
	if err := whereByIDOrUUID(r.DB, c.Param("id"), schoolID).First(&book).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "subject not found"})
		return
	}
	if book.TeacherID != nil {
		var teacher models.User
		if err := r.DB.Where("id = ? AND school_id = ?", *book.TeacherID, *schoolID).First(&teacher).Error(); err == nil {
			teacher.Password = ""
			book.Teacher = &teacher
		}
	}
	c.JSON(http.StatusOK, gin.H{"data": buildSubjectResponse(book)})
}

func (r *bookRepository) UpdateSubject(c *gin.Context) {
	var input subjectPayload
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	name := input.Name
	if name == "" {
		name = input.Title
	}
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}
	schoolID, _, ok := resolveWriteSchoolID(c, input.SchoolID)
	if !ok {
		return
	}
	var book models.Book
	if err := whereByIDOrUUID(r.DB, c.Param("id"), schoolID).First(&book).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "subject not found"})
		return
	}
	if input.TeacherID != nil {
		var teacher models.User
		teacherQuery := r.DB.Where("id = ? AND role = ?", *input.TeacherID, "guru")
		if schoolID != nil {
			teacherQuery = teacherQuery.Where("school_id = ?", *schoolID)
		}
		if err := teacherQuery.First(&teacher).Error(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "teacher not found"})
			return
		}
		book.Author = teacher.Username
	}
	if err := r.DB.Model(&book).Updates(models.Book{Title: name, Author: book.Author, SchoolID: schoolID, TeacherID: input.TeacherID}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update subject"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": buildSubjectResponse(book)})
}

func (r *bookRepository) DeleteSubject(c *gin.Context) { r.DeleteBook(c) }
