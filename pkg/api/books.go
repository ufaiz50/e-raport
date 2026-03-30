package api

import (
	"context"
	"encoding/json"
	"golang-rest-api-template/pkg/cache"
	"golang-rest-api-template/pkg/database"
	"golang-rest-api-template/pkg/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type BookRepository interface {
	Healthcheck(c *gin.Context)
	FindBooks(c *gin.Context)
	CreateBook(c *gin.Context)
	FindBook(c *gin.Context)
	UpdateBook(c *gin.Context)
	DeleteBook(c *gin.Context)
}

type bookRepository struct {
	DB          database.Database
	RedisClient cache.Cache
	Ctx         *context.Context
}

func NewBookRepository(db database.Database, redisClient cache.Cache, ctx *context.Context) *bookRepository {
	return &bookRepository{
		DB:          db,
		RedisClient: redisClient,
		Ctx:         ctx,
	}
}

func (r *bookRepository) Healthcheck(c *gin.Context) {
	c.JSON(http.StatusOK, "ok")
}

func (r *bookRepository) FindBooks(c *gin.Context) {
	schoolID, _, ok := requireTenant(c)
	if !ok {
		return
	}

	var books []models.Book
	type booksCachePayload struct {
		Data []models.Book `json:"data"`
		Meta struct {
			Offset int `json:"offset"`
			Limit  int `json:"limit"`
			Total  int `json:"total"`
			Count  int `json:"count"`
		} `json:"meta"`
	}

	offset, limit, ok := parsePagination(c)
	if !ok {
		return
	}

	scope := "all"
	if schoolID != nil {
		scope = *schoolID
	}
	cacheKey := "books_school_" + scope + "_offset_" + c.DefaultQuery("offset", "0") + "_limit_" + c.DefaultQuery("limit", "10")

	cachedBooks, err := r.RedisClient.Get(*r.Ctx, cacheKey).Result()
	if err == nil {
		var cachedPayload booksCachePayload
		if err := json.Unmarshal([]byte(cachedBooks), &cachedPayload); err == nil && len(cachedPayload.Data) > 0 {
			c.JSON(http.StatusOK, gin.H{
				"data": cachedPayload.Data,
				"meta": gin.H{
					"offset": cachedPayload.Meta.Offset,
					"limit":  cachedPayload.Meta.Limit,
					"total":  cachedPayload.Meta.Total,
					"count":  cachedPayload.Meta.Count,
				},
			})
			return
		}

		if err := json.Unmarshal([]byte(cachedBooks), &books); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"data": books,
				"meta": gin.H{
					"offset": offset,
					"limit":  limit,
					"total":  len(books),
					"count":  len(books),
				},
			})
			return
		}
	}

	var allBooks []models.Book
	countQuery := r.DB
	if schoolID != nil {
		countQuery = countQuery.Where("school_id = ?", *schoolID)
	}
	countQuery.Find(&allBooks)
	total := len(allBooks)

	dataQuery := r.DB.Offset(offset).Limit(limit)
	if schoolID != nil {
		dataQuery = dataQuery.Where("school_id = ?", *schoolID)
	}
	dataQuery.Find(&books)
	for i := range books {
		if books[i].StudentID != nil {
			var student models.Student
			studentQuery := r.DB.Where("id = ?", *books[i].StudentID)
			if schoolID != nil {
				studentQuery = studentQuery.Where("school_id = ?", *schoolID)
			}
			if err := studentQuery.First(&student).Error(); err == nil {
				books[i].Student = &student
			}
		}
	}

	cachePayload := booksCachePayload{Data: books}
	cachePayload.Meta.Offset = offset
	cachePayload.Meta.Limit = limit
	cachePayload.Meta.Total = total
	cachePayload.Meta.Count = len(books)

	serializedBooks, err := json.Marshal(cachePayload)
	if err == nil {
		_ = r.RedisClient.Set(*r.Ctx, cacheKey, serializedBooks, time.Minute).Err()
	}

	c.JSON(http.StatusOK, gin.H{
		"data": books,
		"meta": gin.H{
			"offset": offset,
			"limit":  limit,
			"total":  total,
			"count":  len(books),
		},
	})
}

func (r *bookRepository) CreateBook(c *gin.Context) {
	appCtx, exists := c.MustGet("appCtx").(*bookRepository)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	var input models.CreateBook

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schoolID, _, ok := resolveWriteSchoolID(c, input.SchoolID)
	if !ok {
		return
	}

	if input.TeacherID != nil {
		var teacher models.User
		teacherQuery := appCtx.DB.Where("id = ? AND role = ?", *input.TeacherID, "guru")
		if schoolID != nil {
			teacherQuery = teacherQuery.Where("school_id = ?", *schoolID)
		}
		if err := teacherQuery.First(&teacher).Error(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "teacher not found"})
			return
		}
	}

	if input.StudentID != nil {
		var student models.Student
		if err := appCtx.DB.Where("id = ? AND school_id = ?", *input.StudentID, *schoolID).First(&student).Error(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "student not found"})
			return
		}
	}

	book := models.Book{Title: input.Title, Author: input.Author, SchoolID: schoolID, TeacherID: input.TeacherID, StudentID: input.StudentID}
	appCtx.DB.Create(&book)

	keysPattern := "books_school_" + *schoolID + "_*"
	keys, err := appCtx.RedisClient.Keys(*appCtx.Ctx, keysPattern).Result()
	if err == nil {
		for _, key := range keys {
			appCtx.RedisClient.Del(*appCtx.Ctx, key)
		}
	}

	c.JSON(http.StatusCreated, gin.H{"data": book})
}

func (r *bookRepository) FindBook(c *gin.Context) {
	schoolID, _, ok := requireTenant(c)
	if !ok {
		return
	}

	var book models.Book
	if err := whereByIDOrUUID(r.DB, c.Param("id"), schoolID).First(&book).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	if book.StudentID != nil {
		var student models.Student
		if err := r.DB.Where("id = ? AND school_id = ?", *book.StudentID, *schoolID).First(&student).Error(); err == nil {
			book.Student = &student
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": book})
}

func (r *bookRepository) UpdateBook(c *gin.Context) {
	var input models.UpdateBook

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schoolID, _, ok := resolveWriteSchoolID(c, input.SchoolID)
	if !ok {
		return
	}

	var book models.Book
	if err := whereByIDOrUUID(r.DB, c.Param("id"), schoolID).First(&book).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	if input.StudentID != nil {
		var student models.Student
		if err := r.DB.Where("id = ? AND school_id = ?", *input.StudentID, *schoolID).First(&student).Error(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "student not found"})
			return
		}
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
	}

	r.DB.Model(&book).Updates(models.Book{Title: input.Title, Author: input.Author, SchoolID: schoolID, TeacherID: input.TeacherID, StudentID: input.StudentID})
	c.JSON(http.StatusOK, gin.H{"data": book})
}

func (r *bookRepository) DeleteBook(c *gin.Context) {
	schoolID, _, ok := requireTenant(c)
	if !ok {
		return
	}

	var book models.Book
	if err := whereByIDOrUUID(r.DB, c.Param("id"), schoolID).First(&book).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	r.DB.Delete(&book)
	c.JSON(http.StatusNoContent, gin.H{"data": true})
}
