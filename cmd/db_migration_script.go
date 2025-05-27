package main

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/userblog/management/internal/models"
)

// initializeDatabaseScript creates default roles and permissions
func initializeDatabaseScript(ctx context.Context, db *gorm.DB) {
	db.Debug().AutoMigrate(&models.User{}, &models.Role{}, &models.Permission{}, &models.Blog{})

	// Create admin role if it doesn't exist
	var adminRole models.Role
	if db.Debug().Where("name = ?", "admin").First(&adminRole).RecordNotFound() {
		adminRole = models.Role{
			Name:        "admin",
			Description: "Administrator with all permissions",
		}
		db.Create(&adminRole)
	}

	// Create user role if it doesn't exist
	var userRole models.Role
	if db.Debug().Where("name = ?", "user").First(&userRole).RecordNotFound() {
		userRole = models.Role{
			Name:        "user",
			Description: "Regular user with limited permissions",
		}
		db.Debug().Create(&userRole)
	}

	// Create permissions
	permissions := []models.Permission{
		{Name: "create_blog", Description: "Can create blog posts", Resource: "blog", Action: "create"},
		{Name: "read_blog", Description: "Can read blog posts", Resource: "blog", Action: "read"},
		{Name: "update_blog", Description: "Can update blog posts", Resource: "blog", Action: "update"},
		{Name: "delete_blog", Description: "Can delete blog posts", Resource: "blog", Action: "delete"},
		{Name: "create_user", Description: "Can create users", Resource: "user", Action: "create"},
		{Name: "read_user", Description: "Can read user information", Resource: "user", Action: "read"},
		{Name: "update_user", Description: "Can update user information", Resource: "user", Action: "update"},
		{Name: "delete_user", Description: "Can delete users", Resource: "user", Action: "delete"},
	}

	// Create permissions if they don't exist
	for _, p := range permissions {
		var permission models.Permission
		if db.Debug().Where("name = ?", p.Name).First(&permission).RecordNotFound() {
			db.Debug().Create(&models.Permission{
				Name:        p.Name,
				Description: p.Description,
				Resource:    p.Resource,
				Action:      p.Action,
			})
		}
	}

	// Assign all permissions to admin role
	db.Debug().Exec("DELETE FROM role_permissions WHERE role_id = ?", adminRole.ID)
	for _, p := range permissions {
		var permission models.Permission
		db.Debug().Where("name = ?", p.Name).First(&permission)
		db.Debug().Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)", adminRole.ID, permission.ID)
	}

	// Assign only read permissions to user role
	db.Debug().Exec("DELETE FROM role_permissions WHERE role_id = ?", userRole.ID)
	for _, p := range permissions {
		if p.Action == "read" {
			var permission models.Permission
			db.Debug().Where("name = ?", p.Name).First(&permission)
			db.Debug().Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)", userRole.ID, permission.ID)
		}
	}
}
