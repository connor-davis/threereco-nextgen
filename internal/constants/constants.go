package constants

import "github.com/connor-davis/threereco-nextgen/internal/models"

const (
	InternalServerError        string = "Internal server error"
	InternalServerErrorDetails string = "An unexpected error occurred. Please try again later or contact support."
	UnauthorizedError          string = "Unauthorized"
	UnauthorizedErrorDetails   string = "You are not authorized to access this resource. Please log in or contact support."
	NotFoundError              string = "Not Found"
	NotFoundErrorDetails       string = "The requested resource could not be found. Please check the URL or contact support."
	BadRequestError            string = "Bad Request"
	BadRequestErrorDetails     string = "The request could not be understood or was missing required parameters."
	ConflictError              string = "Conflict"
	ConflictErrorDetails       string = "The request could not be completed due to a conflict with the current state of the resource."
	ForbiddenError             string = "Forbidden"
	ForbiddenErrorDetails      string = "You do not have permission to access this resource. Please check your permissions or contact support."
	Created                    string = "Created"
	CreatedDetails             string = "The resource has been successfully created."
	Success                    string = "Success"
	SuccessDetails             string = "The request was successful."
)

var AvailablePermissionsGroups = []models.AvailablePermissionsGroup{
	{
		Name: "Global",
		Permissions: []models.AvailablePermission{
			{
				Value:       "*",
				Description: "All permissions",
			},
		},
	},
	{
		Name: "Users",
		Permissions: []models.AvailablePermission{
			{
				Value:       "users.*",
				Description: "All user permissions",
			},
			{
				Value:       "users.create",
				Description: "Create user permission",
			},
			{
				Value:       "users.view",
				Description: "View user permission",
			},
			{
				Value:       "users.update",
				Description: "Update user permission",
			},
			{
				Value:       "users.delete",
				Description: "Delete user permission",
			},
		},
	},
	{
		Name: "Organizations",
		Permissions: []models.AvailablePermission{
			{
				Value:       "organizations.*",
				Description: "All organization permissions",
			},
			{
				Value:       "organizations.create",
				Description: "Create organization permission",
			},
			{
				Value:       "organizations.view",
				Description: "View organization permission",
			},
			{
				Value:       "organizations.update",
				Description: "Update organization permission",
			},
			{
				Value:       "organizations.delete",
				Description: "Delete organization permission",
			},
		},
	},
	{
		Name: "Roles",
		Permissions: []models.AvailablePermission{
			{
				Value:       "roles.*",
				Description: "All role permissions",
			},
			{
				Value:       "roles.create",
				Description: "Create role permission",
			},
			{
				Value:       "roles.view",
				Description: "View role permission",
			},
			{
				Value:       "roles.update",
				Description: "Update role permission",
			},
			{
				Value:       "roles.delete",
				Description: "Delete role permission",
			},
		},
	},
	{
		Name: "Audit Logs",
		Permissions: []models.AvailablePermission{
			{
				Value:       "audit_logs.*",
				Description: "All audit log permissions",
			},
			{
				Value:       "audit_logs.create",
				Description: "Create audit log permission",
			},
			{
				Value:       "audit_logs.view",
				Description: "View audit log permission",
			},
			{
				Value:       "audit_logs.update",
				Description: "Update audit log permission",
			},
			{
				Value:       "audit_logs.delete",
				Description: "Delete audit log permission",
			},
		},
	},
	{
		Name: "Materials",
		Permissions: []models.AvailablePermission{
			{
				Value:       "materials.*",
				Description: "All material permissions",
			},
			{
				Value:       "materials.create",
				Description: "Create material permission",
			},
			{
				Value:       "materials.view",
				Description: "View material permission",
			},
			{
				Value:       "materials.update",
				Description: "Update material permission",
			},
			{
				Value:       "materials.delete",
				Description: "Delete material permission",
			},
		},
	},
	{
		Name: "Products",
		Permissions: []models.AvailablePermission{
			{
				Value:       "products.*",
				Description: "All product permissions",
			},
			{
				Value:       "products.create",
				Description: "Create product permission",
			},
			{
				Value:       "products.view",
				Description: "View product permission",
			},
			{
				Value:       "products.update",
				Description: "Update product permission",
			},
			{
				Value:       "products.delete",
				Description: "Delete product permission",
			},
		},
	},
	{
		Name: "Transactions",
		Permissions: []models.AvailablePermission{
			{
				Value:       "transactions.*",
				Description: "All transaction permissions",
			},
			{
				Value:       "transactions.create",
				Description: "Create transaction permission",
			},
			{
				Value:       "transactions.view",
				Description: "View transaction permission",
			},
			{
				Value:       "transactions.update",
				Description: "Update transaction permission",
			},
			{
				Value:       "transactions.delete",
				Description: "Delete transaction permission",
			},
		},
	},
	{
		Name: "Notifications",
		Permissions: []models.AvailablePermission{
			{
				Value:       "notifications.*",
				Description: "All notification permissions",
			},
			{
				Value:       "notifications.create",
				Description: "Create notification permission",
			},
			{
				Value:       "notifications.view",
				Description: "View notification permission",
			},
			{
				Value:       "notifications.update",
				Description: "Update notification permission",
			},
			{
				Value:       "notifications.delete",
				Description: "Delete notification permission",
			},
		},
	},
}
