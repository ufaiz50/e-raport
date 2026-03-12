package api

import "github.com/gin-gonic/gin"

// Subject endpoints are domain aliases to existing books endpoints during migration.
func (r *bookRepository) FindSubjects(c *gin.Context)   { r.FindBooks(c) }
func (r *bookRepository) CreateSubject(c *gin.Context)  { r.CreateBook(c) }
func (r *bookRepository) FindSubject(c *gin.Context)    { r.FindBook(c) }
func (r *bookRepository) UpdateSubject(c *gin.Context)  { r.UpdateBook(c) }
func (r *bookRepository) DeleteSubject(c *gin.Context)  { r.DeleteBook(c) }
