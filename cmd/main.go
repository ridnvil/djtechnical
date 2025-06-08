package main

import (
	"DeallsJobsTest/config"
	"DeallsJobsTest/models"
	"DeallsJobsTest/routes"
	"DeallsJobsTest/services"
	"DeallsJobsTest/utils"
	"DeallsJobsTest/worker"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var cntx = context.Background()

func main() {
	db := config.InitDatabase()
	redis := config.InitRedis(cntx)
	app := fiber.New(fiber.Config{AppName: "DeallsJobsTest by Ridwan"})
	app.Use(logger.New())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	if errmig := db.AutoMigrate(&models.User{},
		&models.OvertimePaid{},
		&models.Attendance{},
		&models.AttendancesPeriod{},
		&models.Attendance{},
		&models.Overtime{},
		&models.Reimbursement{},
		&models.Payslip{},
		&models.AuditLog{}); errmig != nil {
		log.Fatal("Failed to migrate database:", errmig)
	}

	if err := utils.AutoCreateFolder("uploads"); err != nil {
		log.Fatal(err)
	}

	if errgenadmin := services.GenerateUserAdmin(db); errgenadmin != nil {
		log.Fatal("Failed to generate admin user:", errgenadmin)
	}

	if errgenerateslip := services.GeneratePaySlipEmployees(db); errgenerateslip != nil {
		log.Fatal("Failed to generate payslips for employees:", errgenerateslip)
	}

	if errgenempployee := services.Generate100Employees(db); errgenempployee != nil {
		log.Fatal("Failed to generate 100 employees:", errgenempployee)
	}

	if errsetOvertimeOunt := services.GenerateOvertimeAmount(db); errsetOvertimeOunt != nil {
		log.Fatal("Failed to set overtime amount for employees:", errsetOvertimeOunt)
	}

	routes.SetupRoutes(app, db, redis)

	// for simulation payroll
	go worker.BackgroundTaskPayroll(redis, context.Background(), db)
	go func() {
		if err := app.Listen(":3000"); err != nil {
			log.Fatal(err)
		}
	}()

	<-quit
	fmt.Println("Gracefully shutting down...")

	ctx, cancel := context.WithTimeout(cntx, 2*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server exiting")
}
