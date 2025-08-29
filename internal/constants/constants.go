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
				Description: "All permissions.",
			},
		},
	},
	{
		Name: "Users",
		Permissions: []models.AvailablePermission{
			{
				Value:       "users.*",
				Description: "All permissions related to users.",
			},
			{
				Value:       "users.access",
				Description: "Permission to access users management.",
			},
			{
				Value:       "users.create",
				Description: "Permission to create users.",
			},
			{
				Value:       "users.view",
				Description: "Permission to view users.",
			},
			{
				Value:       "users.view.self",
				Description: "Permission to view own user.",
			},
			{
				Value:       "users.view.other",
				Description: "Permission to view other user.",
			},
			{
				Value:       "users.update",
				Description: "Permission to update users.",
			},
			{
				Value:       "users.update.self",
				Description: "Permission to update own user.",
			},
			{
				Value:       "users.view.other",
				Description: "Permission to view other user.",
			},
			{
				Value:       "users.delete",
				Description: "Permission to delete users.",
			},
			{
				Value:       "users.delete.self",
				Description: "Permission to delete own user.",
			},
			{
				Value:       "users.delete.other",
				Description: "Permission to delete other user.",
			},
		},
	},
	{
		Name: "Organizations",
		Permissions: []models.AvailablePermission{
			{
				Value:       "organizations.*",
				Description: "All permissions related to organizations.",
			},
			{
				Value:       "organizations.access",
				Description: "Permission to access organizations management.",
			},
			{
				Value:       "organizations.create",
				Description: "Permission to create organizations.",
			},
			{
				Value:       "organizations.view",
				Description: "Permission to view organizations.",
			},
			{
				Value:       "organizations.update",
				Description: "Permission to update organizations.",
			},
			{
				Value:       "organizations.delete",
				Description: "Permission to delete organizations.",
			},
		},
	},
	{
		Name: "Roles",
		Permissions: []models.AvailablePermission{
			{
				Value:       "roles.*",
				Description: "All permissions related to roles.",
			},
			{
				Value:       "roles.access",
				Description: "Permission to access roles management.",
			},
			{
				Value:       "roles.create",
				Description: "Permission to create roles.",
			},
			{
				Value:       "roles.view",
				Description: "Permission to view roles.",
			},
			{
				Value:       "roles.update",
				Description: "Permission to update roles.",
			},
			{
				Value:       "roles.delete",
				Description: "Permission to delete roles.",
			},
		},
	},
	{
		Name: "Audit Logs",
		Permissions: []models.AvailablePermission{
			{
				Value:       "audit_logs.*",
				Description: "All permissions related to audit logs.",
			},
			{
				Value:       "audit_logs.create",
				Description: "Permission to create audit logs.",
			},
			{
				Value:       "audit_logs.view",
				Description: "Permission to view audit logs.",
			},
			{
				Value:       "audit_logs.update",
				Description: "Permission to update audit logs.",
			},
			{
				Value:       "audit_logs.delete",
				Description: "Permission to delete audit logs.",
			},
		},
	},
	{
		Name: "Materials",
		Permissions: []models.AvailablePermission{
			{
				Value:       "materials.*",
				Description: "All permissions related to materials.",
			},
			{
				Value:       "materials.access",
				Description: "Permission to access materials management.",
			},
			{
				Value:       "materials.create",
				Description: "Permission to create materials.",
			},
			{
				Value:       "materials.view",
				Description: "Permission to view materials.",
			},
			{
				Value:       "materials.update",
				Description: "Permission to update materials.",
			},
			{
				Value:       "materials.delete",
				Description: "Permission to delete materials.",
			},
		},
	},
	{
		Name: "Products",
		Permissions: []models.AvailablePermission{
			{
				Value:       "products.*",
				Description: "All permissions related to products.",
			},
			{
				Value:       "products.access",
				Description: "Permission to access products management.",
			},
			{
				Value:       "products.create",
				Description: "Permission to create products.",
			},
			{
				Value:       "products.view",
				Description: "Permission to view products.",
			},
			{
				Value:       "products.update",
				Description: "Permission to update products.",
			},
			{
				Value:       "products.delete",
				Description: "Permission to delete products.",
			},
		},
	},
	{
		Name: "Transactions",
		Permissions: []models.AvailablePermission{
			{
				Value:       "transactions.*",
				Description: "All permissions related to transactions.",
			},
			{
				Value:       "transactions.access",
				Description: "Permission to access transactions management.",
			},
			{
				Value:       "transactions.create",
				Description: "Permission to create transactions.",
			},
			{
				Value:       "transactions.view",
				Description: "Permission to view transactions.",
			},
			{
				Value:       "transactions.update",
				Description: "Permission to update transactions.",
			},
			{
				Value:       "transactions.delete",
				Description: "Permission to delete transactions.",
			},
		},
	},
	{
		Name: "Notifications",
		Permissions: []models.AvailablePermission{
			{
				Value:       "notifications.*",
				Description: "All permissions related to notifications.",
			},
			{
				Value:       "notifications.access",
				Description: "Permission to access notifications management.",
			},
			{
				Value:       "notifications.create",
				Description: "Permission to create notifications.",
			},
			{
				Value:       "notifications.view",
				Description: "Permission to view notifications.",
			},
			{
				Value:       "notifications.update",
				Description: "Permission to update notifications.",
			},
			{
				Value:       "notifications.delete",
				Description: "Permission to delete notifications.",
			},
		},
	},
}
