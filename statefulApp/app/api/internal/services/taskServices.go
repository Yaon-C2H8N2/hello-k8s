package services

import (
	"api/internal/models"
	"api/pkg/utils"
	"context"
	"github.com/gin-gonic/gin"
)

func GetTaskForUser(c *gin.Context) {
	user, err := utils.NeedsAuth(c)
	if err != nil {
		return
	}

	conn := utils.GetConnections()
	defer conn.Close(context.Background())

	sql := `
		SELECT tasks.id, tasks.name, tasks.description, tasks.due_date
		FROM tasks
		JOIN user_tasks ON tasks.id = user_tasks.task_id
		WHERE user_tasks.user_id = $1
	`

	rows := utils.DoRequest(conn, sql, user.ID)
	tasks := make([]models.Task, 0)
	for rows.Next() {
		var task = models.Task{}
		err = rows.Scan(&task.ID, &task.Name, &task.Description, &task.DueDate)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   "Failed to scan task",
				"message": err.Error(),
			})
			return
		}
		tasks = append(tasks, task)
	}

	c.JSON(200, gin.H{
		"tasks": tasks,
	})
}

func CreateTaskForUser(c *gin.Context) {
	user, err := utils.NeedsAuth(c)
	if err != nil {
		return
	}

	var task models.Task
	err = c.BindJSON(&task)
	if err != nil {
		c.JSON(400, gin.H{
			"error":   "Failed to bind task",
			"message": err.Error(),
		})
		return
	}

	conn := utils.GetConnections()
	defer conn.Close(context.Background())

	sql := `
		INSERT INTO tasks (name, description, due_date)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	rows := utils.DoRequest(conn, sql, task.Name, task.Description, task.DueDate)
	if !rows.Next() {
		c.JSON(500, gin.H{
			"error":   "Failed to insert task",
			"message": "No task id returned",
		})
		return
	}
	var taskId int
	err = rows.Scan(&taskId)
	if err != nil {
		c.JSON(500, gin.H{
			"error":   "Failed to scan task id",
			"message": err.Error(),
		})
		return
	}

	conn2 := utils.GetConnections()
	defer conn2.Close(context.Background())

	sql = `
		INSERT INTO user_tasks (user_id, task_id)
		VALUES ($1, $2)
	`

	utils.DoRequest(conn2, sql, user.ID, taskId)

	c.JSON(200, gin.H{
		"message": "Task created",
		"task": gin.H{
			"id":          taskId,
			"name":        task.Name,
			"description": task.Description,
			"due_date":    task.DueDate,
		},
	})
}

func DeleteTaskForUser(c *gin.Context) {
	user, err := utils.NeedsAuth(c)
	if err != nil {
		return
	}

	taskId := c.Param("id")

	conn := utils.GetConnections()
	defer conn.Close(context.Background())

	sql := `
		DELETE FROM user_tasks
		WHERE user_id = $1 AND task_id = $2
	`

	utils.DoRequest(conn, sql, user.ID, taskId)

	c.JSON(200, gin.H{
		"message": "Task deleted",
	})
}
