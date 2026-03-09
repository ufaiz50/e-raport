package api

import (
	"context"
	"encoding/json"
	"golang-rest-api-template/pkg/cache"
	"golang-rest-api-template/pkg/database"
	"golang-rest-api-template/pkg/models"
	"net/http"
	"strconv"
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

// bookRepository holds shared resources like database and Redis client
type bookRepository struct {
	DB          database.Database
	RedisClient cache.Cache
	Ctx         *context.Context
}

// NewAppContext creates a new AppContext
func NewBookRepository(db database.Database, redisClient cache.Cache, ctx *context.Context) *bookRepository {
	return &bookRepository{
		DB:          db,
		RedisClient: redisClient,
		Ctx:         ctx,
	}
}

// @BasePath /api/v1

// Healthcheck godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} ok
// @Router / [get]
func (r *bookRepository) Healthcheck(c *gin.Context) {
	c.JSON(http.StatusOK, "ok")
}

// FindBooks godoc
// @Summary Get all books with pagination
// @Description Get a list of all books with optional pagination
// @Tags books
// @Security ApiKeyAuth
// @Produce json
// @Param offset query int false "Offset for pagination" default(0)
// @Param limit query int false "Limit for pagination" default(10)
// @Success 200 {array} models.Book "Successfully retrieved list of books"
// @Router /books [get]
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

	// Get query params
	offsetQuery := c.DefaultQuery("offset", "0")
	limitQuery := c.DefaultQuery("limit", "10")

	// Convert query params to integers
	offset, err := strconv.Atoi(offsetQuery)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset format"})
		return
	}

	limit, err := strconv.Atoi(limitQuery)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit format"})
		return
	}

	if offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset format"})
		return
	}

	if limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit format"})
		return
	}

	// Create a cache key based on query params and effective school scope
	scope := "all"
	if schoolID != nil {
		scope = strconv.Itoa(int(*schoolID))
	}
	cacheKey := "books_school_" + scope + "_offset_" + offsetQuery + "_limit_" + limitQuery

	// Try fetching the data from Redis first
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

		// Backward compatibility for older cache format ([]Book only)
		if err := json.Unmarshal([]byte(cachedBooks), &books); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal cached data"})
			return
		}

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

	var allBooks []models.Book
	countQuery := r.DB
	if schoolID != nil {
		countQuery = countQuery.Where("school_id = ?", *schoolID)
	}
	countQuery.Find(&allBooks)
	total := len(allBooks)

	// If cache missed, fetch data from the database
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

	// Serialize books object and store it in Redis
	cachePayload := booksCachePayload{Data: books}
	cachePayload.Meta.Offset = offset
	cachePayload.Meta.Limit = limit
	cachePayload.Meta.Total = total
	cachePayload.Meta.Count = len(books)

	serializedBooks, err := json.Marshal(cachePayload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal data"})
		return
	}
	err = r.RedisClient.Set(*r.Ctx, cacheKey, serializedBooks, time.Minute).Err() // Here TTL is set to one hour
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set cache"})
		return
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

// CreateBook godoc
// @Summary Create a new book
// @Description Create a new book with the given input data
// @Tags books
// @Security ApiKeyAuth
// @Security JwtAuth
// @Accept  json
// @Produce  json
// @Param   input     body   models.CreateBook   true   "Create book object"
// @Success 201 {object} models.Book "Successfully created book"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Router /books [post]
func (r *bookRepository) CreateBook(c *gin.Context) {
	schoolID, _, ok := requireTenant(c)
	if !ok {
		return
	}

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

	if input.StudentID != nil {
		var student models.Student
		if err := appCtx.DB.Where("id = ? AND school_id = ?", *input.StudentID, *schoolID).First(&student).Error(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "student not found"})
			return
		}
	}

	book := models.Book{Title: input.Title, Author: input.Author, SchoolID: schoolID, StudentID: input.StudentID}

	appCtx.DB.Create(&book)

	// Invalidate cache
	keysPattern := "books_school_" + strconv.Itoa(int(*schoolID)) + "_*"
	keys, err := appCtx.RedisClient.Keys(*appCtx.Ctx, keysPattern).Result()
	if err == nil {
		for _, key := range keys {
			appCtx.RedisClient.Del(*appCtx.Ctx, key)
		}
	}

	c.JSON(http.StatusCreated, gin.H{"data": book})
}

// FindBook godoc
// @Summary Find a book by ID
// @Description Get details of a book by its ID
// @Tags books
// @Security ApiKeyAuth
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} models.Book "Successfully retrieved book"
// @Failure 404 {string} string "Book not found"
// @Router /books/{id} [get]
func (r *bookRepository) FindBook(c *gin.Context) {
	schoolID, _, ok := requireTenant(c)
	if !ok {
		return
	}

	var book models.Book

	if err := r.DB.Where("id = ? AND school_id = ?", c.Param("id"), *schoolID).First(&book).Error(); err != nil {
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

// UpdateBook godoc
// @Summary Update a book by ID
// @Description Update the book details for the given ID
// @Tags books
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param id path string true "Book ID"
// @Param input body models.UpdateBook true "Update book object"
// @Success 200 {object} models.Book "Successfully updated book"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "book not found"
// @Router /books/{id} [put]
func (r *bookRepository) UpdateBook(c *gin.Context) {
	schoolID, _, ok := requireTenant(c)
	if !ok {
		return
	}

	var book models.Book
	var input models.UpdateBook

	if err := r.DB.Where("id = ? AND school_id = ?", c.Param("id"), *schoolID).First(&book).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.StudentID != nil {
		var student models.Student
		if err := r.DB.Where("id = ? AND school_id = ?", *input.StudentID, *schoolID).First(&student).Error(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "student not found"})
			return
		}
	}

	r.DB.Model(&book).Updates(models.Book{Title: input.Title, Author: input.Author, SchoolID: schoolID, StudentID: input.StudentID})

	c.JSON(http.StatusOK, gin.H{"data": book})
}

// DeleteBook godoc
// @Summary Delete a book by ID
// @Description Delete the book with the given ID
// @Tags books
// @Security ApiKeyAuth
// @Produce json
// @Param id path string true "Book ID"
// @Success 204 {string} string "Successfully deleted book"
// @Failure 404 {string} string "book not found"
// @Router /books/{id} [delete]
func (r *bookRepository) DeleteBook(c *gin.Context) {
	schoolID, _, ok := requireTenant(c)
	if !ok {
		return
	}

	var book models.Book

	if err := r.DB.Where("id = ? AND school_id = ?", c.Param("id"), *schoolID).First(&book).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	r.DB.Delete(&book)

	c.JSON(http.StatusNoContent, gin.H{"data": true})
}
