package handlers

import (
	"errors"
	"net/http"
	"time"
	"todo-api/internal/config"
	"todo-api/internal/models"
	repo "todo-api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type SignupRequest struct {
	Email string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

func SignupHandler(pool *pgxpool.Pool) gin.HandlerFunc {
    return func(c *gin.Context) {
        var payload SignupRequest

        if err := c.ShouldBindJSON(&payload); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        if len(payload.Password) < 8 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Password must contain at least 8 characters"})
            return
        }

        if payload.ConfirmPassword != payload.Password {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords don't match!"})
            return
        }

        hashedPwd, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process security credentials"})
            return
        }

        user := &models.User{
            Email:    payload.Email,
            Password: string(hashedPwd),
        }

        newUser, err := repo.CreateUser(pool, user)
        if err != nil {
            var pgErr *pgconn.PgError	// postgreSQL error
            
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
                c.JSON(http.StatusConflict, gin.H{"error": "Email provided is already registered"})
                return
            }

            c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
            return
        }

        c.JSON(http.StatusCreated, newUser)
    }
}