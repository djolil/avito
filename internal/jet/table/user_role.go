//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var UserRole = newUserRoleTable("public", "user_role", "")

type userRoleTable struct {
	postgres.Table

	// Columns
	UserID postgres.ColumnInteger
	RoleID postgres.ColumnInteger

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type UserRoleTable struct {
	userRoleTable

	EXCLUDED userRoleTable
}

// AS creates new UserRoleTable with assigned alias
func (a UserRoleTable) AS(alias string) *UserRoleTable {
	return newUserRoleTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new UserRoleTable with assigned schema name
func (a UserRoleTable) FromSchema(schemaName string) *UserRoleTable {
	return newUserRoleTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new UserRoleTable with assigned table prefix
func (a UserRoleTable) WithPrefix(prefix string) *UserRoleTable {
	return newUserRoleTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new UserRoleTable with assigned table suffix
func (a UserRoleTable) WithSuffix(suffix string) *UserRoleTable {
	return newUserRoleTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newUserRoleTable(schemaName, tableName, alias string) *UserRoleTable {
	return &UserRoleTable{
		userRoleTable: newUserRoleTableImpl(schemaName, tableName, alias),
		EXCLUDED:      newUserRoleTableImpl("", "excluded", ""),
	}
}

func newUserRoleTableImpl(schemaName, tableName, alias string) userRoleTable {
	var (
		UserIDColumn   = postgres.IntegerColumn("user_id")
		RoleIDColumn   = postgres.IntegerColumn("role_id")
		allColumns     = postgres.ColumnList{UserIDColumn, RoleIDColumn}
		mutableColumns = postgres.ColumnList{}
	)

	return userRoleTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		UserID: UserIDColumn,
		RoleID: RoleIDColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
