package worker

import (
	"DeallsJobsTest/models"
	"DeallsJobsTest/services"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"log"
	"time"
)

func BackgroundTaskPayroll(rdb *redis.Client, ctx context.Context, db *gorm.DB) {
	log.Println("Payroll background task is running...")
	defer rdb.Close()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Background task shutting down...")
			return
		default:
			go ProcessPayrollData(rdb, ctx, "payroll_channel_employees", db)
			time.Sleep(5 * time.Second)
		}
	}
}

func ProcessPayrollData(client *redis.Client, ctx context.Context, channel string, db *gorm.DB) {
	pubsub := client.Subscribe(ctx, channel)
	defer func(pubsub *redis.PubSub) {
		err := pubsub.Close()
		if err != nil {
			log.Println("failed to close pubsub: ", err)
		}
	}(pubsub)

	_, err := pubsub.Receive(ctx)
	if err != nil {
		log.Println("failed to subscribe: ", err)
		return
	}
	ch := pubsub.Channel()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Stopping subscriber...")
			return
		case msg := <-ch:
			if msg == nil {
				log.Println("Channel closed")
				return
			}

			var userData struct {
				PeriodID  uint          `json:"period_id"`
				IP        string        `json:"ip"`
				Users     []models.User `json:"users"`
				RequestID string        `json:"request_id"`
			}
			if err := json.Unmarshal([]byte(msg.Payload), &userData); err != nil {
				log.Println("Error unmarshalling message payload:", err)
				return
			}

			for _, user := range userData.Users {
				res, errupdate := services.GeneratePaySlipByEmployeeID(db, user.ID, userData.IP, userData.RequestID)
				if errupdate != nil {
					log.Printf("Error generating payslip for user ID %d: %v\n", user.ID, errupdate)
					return
				}
				log.Printf("Processing user: ID=%d, Name=%s, PaySlipID=%d", user.ID, user.FullName, res.ID)
			}

			if errupdate := db.Table("attendance_periods").Where("id = ?", userData.PeriodID).Update("is_locked", true).Error; errupdate != nil {
				log.Printf("Error locking attendance period %d: %v\n", userData.PeriodID, errupdate)
				return
			}
		}
	}
}
