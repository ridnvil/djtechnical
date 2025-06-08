package routes

import (
	"DeallsJobsTest/controllers"
	"DeallsJobsTest/middlewares"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB, redis *redis.Client) {

	authController := controllers.NewAuthController(db)
	attendanceController := controllers.NewAttendanceController(db)
	overtimeController := controllers.NewOvertimeController(db)
	reimbursementController := controllers.NewReimbursementController(db)
	payslipController := controllers.NewPayslipController(db)
	payrollController := controllers.NewPayrollController(db, redis)
	middleWareHandler := middlewares.NewJWTMiddlewareHandler(db)
	logHandler := middlewares.NewAuditLog(db)

	auth := app.Group("/api/auth")
	auth.Post("/login", authController.Login)

	employee := app.Group("/api/employee")
	employee.Post("/attendance", middleWareHandler.JWTMiddleware, middleWareHandler.ValidateUserRole, logHandler.LogAction, attendanceController.SubmitAttendance)
	employee.Post("/overtime", middleWareHandler.JWTMiddleware, middleWareHandler.ValidateUserRole, logHandler.LogAction, overtimeController.SubmitOvertime)
	employee.Post("/reimbursement", middleWareHandler.JWTMiddleware, middleWareHandler.ValidateUserRole, logHandler.LogAction, reimbursementController.SubmitReimbursement)
	employee.Post("/reimbursement/uploads/:id", middleWareHandler.JWTMiddleware, middleWareHandler.ValidateUserRole, logHandler.LogAction, reimbursementController.UploadAttcahments)
	employee.Get("/payslip", middleWareHandler.JWTMiddleware, middleWareHandler.ValidateUserRole, logHandler.LogAction, payslipController.GeneratePayslip)

	admin := app.Group("/api/admin")
	admin.Post("/period", middleWareHandler.JWTMiddleware, middleWareHandler.ValidateUserRole, logHandler.LogAction, payrollController.CreatePeriod)
	admin.Post("/payroll", middleWareHandler.JWTMiddleware, middleWareHandler.ValidateUserRole, logHandler.LogAction, payrollController.RunPayroll)
	admin.Get("/summary/:period_id", middleWareHandler.JWTMiddleware, middleWareHandler.ValidateUserRole, logHandler.LogAction, payrollController.GetPayrollSummary)
}
