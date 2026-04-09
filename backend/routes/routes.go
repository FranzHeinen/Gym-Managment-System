package routes

import (
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/handlers"
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// --- 1. Vistas HTML ---
	r.GET("/", func(c *gin.Context) { c.HTML(http.StatusOK, "index.html", gin.H{"TemplateName": "index"}) })
	r.GET("/home", func(c *gin.Context) { c.HTML(http.StatusOK, "home.html", gin.H{"TemplateName": "home"}) })
	r.GET("/login", func(c *gin.Context) { c.HTML(http.StatusOK, "login.html", gin.H{"TemplateName": "login"}) })
	r.GET("/register", func(c *gin.Context) { c.HTML(http.StatusOK, "register.html", gin.H{"TemplateName": "register"}) })
	r.GET("/profile", func(c *gin.Context) { c.HTML(http.StatusOK, "profile.html", gin.H{"TemplateName": "profile"}) })
	r.GET("/rutinas", func(c *gin.Context) { c.HTML(http.StatusOK, "rutinas.html", gin.H{"TemplateName": "rutinas"}) })
	r.GET("/rutina-detalle", func(c *gin.Context) {
		c.HTML(http.StatusOK, "rutina-detalle.html", gin.H{"TemplateName": "rutina-detalle"})
	})
	r.GET("/dashboard", func(c *gin.Context) { c.HTML(http.StatusOK, "dashboard.html", gin.H{"TemplateName": "dashboard"}) })
	r.GET("/admin/ejercicios", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin_ejercicios.html", gin.H{"TemplateName": "admin_ejercicios"})
	})
	r.GET("/admin/dashboard", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin_dashboard.html", gin.H{"TemplateName": "admin_dashboard"})
	})
	r.GET("/admin/users", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin_usuarios.html", gin.H{"TemplateName": "admin_usuarios"})
	})
	r.GET("/admin/logs", handlers.GetLogsPage)

	// --- 2. Grupo de API ---
	api := r.Group("/api")
	{
		// -- 2.1 Endpoints Públicos (Autenticación y Mocks) --
		api.POST("/auth/register", handlers.Register)
		api.POST("/auth/login", handlers.Login)

		// -- 2.2 Endpoints Privados (Requieren autenticación de cualquier rol) --
		private := api.Group("/")
		private.Use(middleware.AuthMiddleware())
		{
			private.POST("/auth/logout", handlers.Logout)

			// Perfil de Usuario
			private.GET("/profile", handlers.GetProfile)
			private.PUT("/profile", handlers.UpdateProfile)
			private.POST("/profile/change-password", handlers.ChangePassword)

			// Estadísticas Personales
			private.GET("/profile/stats", handlers.GetUserStats)

			// Catálogo de Ejercicios (solo lectura para usuarios)
			private.GET("/ejercicios", handlers.GetExercises)

			// Gestión de Rutinas del Usuario
			private.GET("/rutinas", handlers.GetRoutines)
			private.POST("/rutinas", handlers.CreateRoutine)
			private.GET("/rutinas/:id", handlers.GetRoutineByID)
			private.PUT("/rutinas/:id", handlers.UpdateRoutine)
			private.DELETE("/rutinas/:id", handlers.DeleteRoutine)
			private.POST("/rutinas/:id/duplicate", handlers.DuplicateRoutine)

			private.POST("/rutinas/:id/ejercicios", handlers.AddExerciseToRoutine)
			private.PUT("/rutinas/:id/ejercicios/:ejercicioId", handlers.UpdateExerciseInRoutine)
			private.DELETE("/rutinas/:id/ejercicios/:ejercicioId", handlers.RemoveExerciseFromRoutine)
			private.GET("/rutinas/:id/ejercicios", handlers.GetExercisesFromRoutine)
			private.GET("/rutinas/:id/ejercicios/:ejercicioId", handlers.GetExerciseFromRoutine)

			// Seguimiento de Progreso (Workouts)
			private.POST("/workouts", handlers.RegisterWorkout)
			private.GET("/workouts", handlers.GetWorkoutHistory)
		}

		// -- 2.3 Endpoints de Administrador (Requieren rol ADMIN) --
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(), middleware.CheckAdmin())
		{
			// Gestión del Catálogo de Ejercicios
			admin.POST("/ejercicios", handlers.CreateExercise)
			admin.PUT("/ejercicios/:id", handlers.UpdateExercise)
			admin.DELETE("/ejercicios/:id", handlers.DeleteExercise)
			admin.GET("/users", handlers.GetAllUsers)
			admin.GET("/stats", handlers.GetGlobalStats)
		}
	}
}
