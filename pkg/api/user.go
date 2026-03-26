package api

import (
	"context"
	"errors"
	"fmt"
	"golang-rest-api-template/pkg/auth"
	"golang-rest-api-template/pkg/database"
	"golang-rest-api-template/pkg/models"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	LoginHandler(c *gin.Context)
	RefreshTokenHandler(c *gin.Context)
	LogoutHandler(c *gin.Context)
	RegisterHandler(c *gin.Context)
	MeHandler(c *gin.Context)
}

// bookRepository holds shared resources like database and Redis client
type userRepository struct {
	DB  database.Database
	Ctx *context.Context
}

func NewUserRepository(db database.Database, ctx *context.Context) *userRepository {
	return &userRepository{
		DB:  db,
		Ctx: ctx,
	}
}

func refreshCookieOptions() (maxAge int, secure bool, sameSite http.SameSite) {
	maxAge = int((7 * 24 * time.Hour).Seconds())
	secure = os.Getenv("COOKIE_SECURE") == "true"
	sameSite = http.SameSiteLaxMode
	if os.Getenv("COOKIE_SAMESITE") == "none" {
		sameSite = http.SameSiteNoneMode
	}
	return
}

func setRefreshCookie(c *gin.Context, token string) {
	maxAge, secure, sameSite := refreshCookieOptions()
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSite,
		MaxAge:   maxAge,
	})
}

func clearRefreshCookie(c *gin.Context) {
	_, secure, sameSite := refreshCookieOptions()
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSite,
		MaxAge:   -1,
	})
}

// @BasePath /api/v1

// LoginHandler godoc
// @Summary Authenticate a user
// @Schemes
// @Description Authenticates a user using username and password, returns a JWT token if successful
// @Tags user
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   user     body    models.LoginUser     true        "User login object"
// @Success 200 {string} string "JWT Token"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /login [post]
func (r *userRepository) LoginHandler(c *gin.Context) {
	var incomingUser models.LoginUser
	var dbUser models.User

	// Get JSON body
	if err := c.ShouldBindJSON(&incomingUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	// Fetch the user from the database
	if err := r.DB.Where("username = ?", incomingUser.Username).First(&dbUser).Error(); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		}
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(incomingUser.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	accessToken, err := auth.GenerateAccessToken(dbUser.Username, dbUser.Role, dbUser.SchoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	refreshToken, err := auth.GenerateRefreshToken(dbUser.Username, dbUser.Role, dbUser.SchoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating refresh token"})
		return
	}

	setRefreshCookie(c, refreshToken)
	c.JSON(http.StatusOK, gin.H{"token": accessToken, "access_token": accessToken})
}

func (r *userRepository) RefreshTokenHandler(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	refreshToken := req.RefreshToken
	if refreshToken == "" {
		if cookie, err := c.Cookie("refresh_token"); err == nil {
			refreshToken = cookie
		}
	}
	if refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing refresh token"})
		return
	}

	claims, err := auth.ParseRefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	var dbUser models.User
	if err := r.DB.Where("username = ?", claims.Username).First(&dbUser).Error(); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	accessToken, err := auth.GenerateAccessToken(dbUser.Username, dbUser.Role, dbUser.SchoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	refreshToken, err = auth.GenerateRefreshToken(dbUser.Username, dbUser.Role, dbUser.SchoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating refresh token"})
		return
	}

	setRefreshCookie(c, refreshToken)
	c.JSON(http.StatusOK, gin.H{"token": accessToken, "access_token": accessToken})
}

func (r *userRepository) LogoutHandler(c *gin.Context) {
	clearRefreshCookie(c)
	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}

// RegisterHandler godoc
// @Summary Register a new user
// @Schemes http
// @Description Registers a new user with the given username and password
// @Tags user
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   user     body    models.LoginUser     true        "User registration object"
// @Success 201 {string} string	"Successfully registered"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /register [post]
func (r *userRepository) RegisterHandler(c *gin.Context) {
	var user models.RegisterUser

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.Username = strings.TrimSpace(user.Username)
	if user.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}

	// Create new user
	role := user.Role
	if role == "" {
		role = "guru"
	}

	if role != "super_admin" && user.SchoolID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "school_id is required for non-super_admin users"})
		return
	}

	if role == "super_admin" {
		// Keep super admin global (non-tenant) to avoid accidental tenant scoping.
		user.SchoolID = nil
	}

	var existingUser models.User
	if err := r.DB.Where("username = ?", user.Username).First(&existingUser).Error(); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	newUser := models.User{Username: user.Username, Password: hashedPassword, Role: role, SchoolID: user.SchoolID}

	// Save the user to the database
	if err := r.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Could not save user: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registration successful"})
}

func (r *userRepository) MeHandler(c *gin.Context) {
	claimsAny, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	claims, ok := claimsAny.(map[string]any)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid auth context"})
		return
	}
	username, _ := claims["username"].(string)
	var user models.User
	if err := r.DB.Where("username = ?", username).First(&user).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	user.Password = ""
	c.JSON(http.StatusOK, gin.H{"data": user})
}
