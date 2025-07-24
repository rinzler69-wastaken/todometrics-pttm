package models

import (
	"time"
	//"gorm.io/gorm"
)

type PriorityLevel string

const (
	Low    PriorityLevel = "Low"
	Medium PriorityLevel = "Medium"
	High   PriorityLevel = "High"
)

type Task struct {
	ID          uint          `gorm:"primaryKey" json:"id"`
	Name        string        `gorm:"not null" json:"name"`
	Description string        `gorm:"type:text" json:"description"`
	DueDate     time.Time     `gorm:"not null" json:"due_date"`
	Priority    PriorityLevel `gorm:"type:varchar(10);not null" json:"priority"`
	Completed   bool          `gorm:"default:false" json:"completed"`
	CreatedAt   time.Time     `json:"created_at"`
	CompletedAt time.Time     `gorm:"column:completed_at"`
}

// Optional helper to compute CompletionStatus (not stored in DB)
func (t *Task) CompletionStatus() string {
	if !t.Completed {
		return "Incomplete"
	}

	now := time.Now()
	if t.DueDate.After(now) {
		return "Ahead of time"
	} else if t.DueDate.Equal(now.Truncate(24 * time.Hour)) {
		return "Just in time"
	} else {
		return "Overdue"
	}
}

// Optional method to check if task is urgent
func (t Task) IsUrgent() bool {
	today := time.Now().Truncate(24 * time.Hour)
	due := t.DueDate.Truncate(24 * time.Hour)
	return !t.Completed && (due.Equal(today) || due.Before(today))
}

func (t Task) IsDueToday() bool {
	location, _ := time.LoadLocation("Asia/Jakarta") // or "Asia/Bangkok", etc.
	today := time.Now().In(location).Truncate(24 * time.Hour)
	due := t.DueDate.In(location).Truncate(24 * time.Hour)
	return !t.Completed && today.Equal(due)
}

func (t Task) IsOverdue() bool {
	location, _ := time.LoadLocation("Asia/Jakarta")
	today := time.Now().In(location).Truncate(24 * time.Hour)
	due := t.DueDate.In(location).Truncate(24 * time.Hour)
	return !t.Completed && due.Before(today)
}
