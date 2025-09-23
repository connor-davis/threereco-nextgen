package common

import "github.com/connor-davis/threereco-nextgen/internal/models"

var Permissions = []models.PermissionGroup{
	{
		Name: "Materials",
		Permissions: []models.Permission{
			{
				Label:       "All Materials",
				Value:       "materials.*",
				Description: "Allows the user to perform any action on materials.",
			},
			{
				Label:       "Access Materials",
				Value:       "materials.access",
				Description: "Allows the user to access the materials module.",
			},
			{
				Label:       "View Materials",
				Value:       "materials.view",
				Description: "Allows the user to view materials.",
			},
			{
				Label:       "Create Material",
				Value:       "materials.create",
				Description: "Allows the user to create new materials.",
			},
			{
				Label:       "Update Material",
				Value:       "materials.update",
				Description: "Allows the user to update existing materials.",
			},
			{
				Label:       "Delete Material",
				Value:       "materials.delete",
				Description: "Allows the user to delete existing materials.",
			},
		},
	},
	{
		Name: "Collections",
		Permissions: []models.Permission{
			{
				Label:       "All Collections",
				Value:       "collections.*",
				Description: "Allows the user to perform any action on collections.",
			},
			{
				Label:       "Access Collections",
				Value:       "collections.access",
				Description: "Allows the user to access the collections module.",
			},
			{
				Label:       "View Collections",
				Value:       "collections.view",
				Description: "Allows the user to view collections.",
			},
			{
				Label:       "Create Collection",
				Value:       "collections.create",
				Description: "Allows the user to create new collections.",
			},
			{
				Label:       "Update Collection",
				Value:       "collections.update",
				Description: "Allows the user to update existing collections.",
			},
			{
				Label:       "Delete Collection",
				Value:       "collections.delete",
				Description: "Allows the user to delete existing collections.",
			},
		},
		SubGroups: []models.PermissionGroup{
			{
				Name: "Collection Materials",
				Permissions: []models.Permission{
					{
						Label:       "All Collection Materials",
						Value:       "collections.materials.*",
						Description: "Allows the user to perform any action on collection materials.",
					},
					{
						Label:       "Access Collection Materials",
						Value:       "collections.materials.access",
						Description: "Allows the user to access the collection materials module.",
					},
					{
						Label:       "View Collection Materials",
						Value:       "collections.materials.view",
						Description: "Allows the user to view collection materials.",
					},
					{
						Label:       "Create Collection Material",
						Value:       "collections.materials.create",
						Description: "Allows the user to create new collection materials.",
					},
					{
						Label:       "Update Collection Material",
						Value:       "collections.materials.update",
						Description: "Allows the user to update existing collection materials.",
					},
					{
						Label:       "Delete Collection Material",
						Value:       "collections.materials.delete",
						Description: "Allows the user to delete existing collection materials.",
					},
					{
						Label:       "Assign Materials to Collection",
						Value:       "collections.materials.assign",
						Description: "Allows the user to assign materials to collections.",
					},
					{
						Label:       "Unassign Materials from Collection",
						Value:       "collections.materials.unassign",
						Description: "Allows the user to unassign materials from collections.",
					},
				},
			},
		},
	},
	{
		Name: "Transactions",
		Permissions: []models.Permission{
			{
				Label:       "All Transactions",
				Value:       "transactions.*",
				Description: "Allows the user to perform any action on transactions.",
			},
			{
				Label:       "Access Transactions",
				Value:       "transactions.access",
				Description: "Allows the user to access the transactions module.",
			},
			{
				Label:       "View Transactions",
				Value:       "transactions.view",
				Description: "Allows the user to view transactions.",
			},
			{
				Label:       "Create Transaction",
				Value:       "transactions.create",
				Description: "Allows the user to create new transactions.",
			},
			{
				Label:       "Update Transaction",
				Value:       "transactions.update",
				Description: "Allows the user to update existing transactions.",
			},
			{
				Label:       "Delete Transaction",
				Value:       "transactions.delete",
				Description: "Allows the user to delete existing transactions.",
			},
		},
		SubGroups: []models.PermissionGroup{
			{
				Name: "Transaction Materials",
				Permissions: []models.Permission{
					{
						Label:       "All Transaction Materials",
						Value:       "transactions.materials.*",
						Description: "Allows the user to perform any action on transaction materials.",
					},
					{
						Label:       "Access Transaction Materials",
						Value:       "transactions.materials.access",
						Description: "Allows the user to access the transaction materials module.",
					},
					{
						Label:       "View Transaction Materials",
						Value:       "transactions.materials.view",
						Description: "Allows the user to view transaction materials.",
					},
					{
						Label:       "Create Transaction Material",
						Value:       "transactions.materials.create",
						Description: "Allows the user to create new transaction materials.",
					},
					{
						Label:       "Update Transaction Material",
						Value:       "transactions.materials.update",
						Description: "Allows the user to update existing transaction materials.",
					},
					{
						Label:       "Delete Transaction Material",
						Value:       "transactions.materials.delete",
						Description: "Allows the user to delete existing transaction materials.",
					},
					{
						Label:       "Assign Materials to Transaction",
						Value:       "transactions.materials.assign",
						Description: "Allows the user to assign materials to transactions.",
					},
					{
						Label:       "Unassign Materials from Transaction",
						Value:       "transactions.materials.unassign",
						Description: "Allows the user to unassign materials from transactions.",
					},
				},
			},
		},
	},
	{
		Name: "Users",
		Permissions: []models.Permission{
			{
				Label:       "All Users",
				Value:       "users.*",
				Description: "Allows the user to perform any action on users.",
			},
			{
				Label:       "Access Users",
				Value:       "users.access",
				Description: "Allows the user to access the users module.",
			},
			{
				Label:       "View User",
				Value:       "users.view",
				Description: "Allows the user to view a user.",
			},
			{
				Label:       "Create User",
				Value:       "users.create",
				Description: "Allows the user to create new users.",
			},
			{
				Label:       "Update Any User",
				Value:       "users.update.any",
				Description: "Allows the user to update any user.",
			},
			{
				Label:       "Update Self",
				Value:       "users.update.self",
				Description: "Allows the user to update their own user details.",
			},
			{
				Label:       "Delete Any User",
				Value:       "users.delete.any",
				Description: "Allows the user to delete any user.",
			},
			{
				Label:       "Delete Self",
				Value:       "users.delete.self",
				Description: "Allows the user to delete their own user account.",
			},
			{
				Label:       "Manage MFA",
				Value:       "users.mfa.manage",
				Description: "Allows the user to manage multi-factor authentication settings.",
			},
		},
	},
	{
		Name: "Businesses",
		Permissions: []models.Permission{
			{
				Label:       "All Businesses",
				Value:       "businesses.*",
				Description: "Allows the user to perform any action on businesses.",
			},
			{
				Label:       "Access Businesses",
				Value:       "businesses.access",
				Description: "Allows the user to access the businesses module.",
			},
			{
				Label:       "View Business",
				Value:       "businesses.view",
				Description: "Allows the user to view a business.",
			},
			{
				Label:       "Create Business",
				Value:       "businesses.create",
				Description: "Allows the user to create a new business.",
			},
			{
				Label:       "Update Business",
				Value:       "businesses.update",
				Description: "Allows the user to update a business.",
			},
			{
				Label:       "Delete Business",
				Value:       "businesses.delete",
				Description: "Allows the user to delete a business.",
			},
		},
		SubGroups: []models.PermissionGroup{
			{
				Name: "Business Users",
				Permissions: []models.Permission{
					{
						Label:       "All Business Users",
						Value:       "businesses.users.*",
						Description: "Allows the user to perform any action on business users.",
					},
					{
						Label:       "Access Business Users",
						Value:       "businesses.users.access",
						Description: "Allows the user to access the business users module.",
					},
					{
						Label:       "View Business Users",
						Value:       "businesses.users.view",
						Description: "Allows the user to view business users.",
					},
					{
						Label:       "Assign Business User",
						Value:       "businesses.users.assign",
						Description: "Allows the user to assign users to businesses.",
					},
					{
						Label:       "Unassign Business User",
						Value:       "businesses.users.unassign",
						Description: "Allows the user to unassign users from businesses.",
					},
				},
			},
			{
				Name: "Business Roles",
				Permissions: []models.Permission{
					{
						Label:       "All Business Roles",
						Value:       "businesses.roles.*",
						Description: "Allows the user to perform any action on business roles.",
					},
					{
						Label:       "Access Business Roles",
						Value:       "businesses.roles.access",
						Description: "Allows the user to access the business roles module.",
					},
					{
						Label:       "View Business Roles",
						Value:       "businesses.roles.view",
						Description: "Allows the user to view business roles.",
					},
					{
						Label:       "Assign Business Role",
						Value:       "businesses.roles.assign",
						Description: "Allows the user to assign roles to businesses.",
					},
					{
						Label:       "Unassign Business Role",
						Value:       "businesses.roles.unassign",
						Description: "Allows the user to unassign roles from businesses.",
					},
				},
			},
		},
	},
	{
		Name: "Roles",
		Permissions: []models.Permission{
			{
				Label:       "All Roles",
				Value:       "roles.*",
				Description: "Allows the user to perform any action on roles.",
			},
			{
				Label:       "Access Roles",
				Value:       "roles.access",
				Description: "Allows the user to access the roles module.",
			},
			{
				Label:       "View Roles",
				Value:       "roles.view",
				Description: "Allows the user to view roles.",
			},
			{
				Label:       "Create Role",
				Value:       "roles.create",
				Description: "Allows the user to create new roles.",
			},
			{
				Label:       "Update Role",
				Value:       "roles.update",
				Description: "Allows the user to update existing roles.",
			},
			{
				Label:       "Delete Role",
				Value:       "roles.delete",
				Description: "Allows the user to delete existing roles.",
			},
		},
	},
	{
		Name: "Permissions",
		Permissions: []models.Permission{
			{
				Label:       "View Permissions",
				Value:       "permissions.view",
				Description: "Allows the user to view available permissions.",
			},
		},
	},
}
