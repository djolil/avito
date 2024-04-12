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

var UserAccount = newUserAccountTable("public", "user_account", "")

type userAccountTable struct {
	postgres.Table

	// Columns
	ID          postgres.ColumnInteger
	FirstName   postgres.ColumnString
	LastName    postgres.ColumnString
	Email       postgres.ColumnString
	PhoneNumber postgres.ColumnString
	Password    postgres.ColumnString

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type UserAccountTable struct {
	userAccountTable

	EXCLUDED userAccountTable
}

// AS creates new UserAccountTable with assigned alias
func (a UserAccountTable) AS(alias string) *UserAccountTable {
	return newUserAccountTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new UserAccountTable with assigned schema name
func (a UserAccountTable) FromSchema(schemaName string) *UserAccountTable {
	return newUserAccountTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new UserAccountTable with assigned table prefix
func (a UserAccountTable) WithPrefix(prefix string) *UserAccountTable {
	return newUserAccountTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new UserAccountTable with assigned table suffix
func (a UserAccountTable) WithSuffix(suffix string) *UserAccountTable {
	return newUserAccountTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newUserAccountTable(schemaName, tableName, alias string) *UserAccountTable {
	return &UserAccountTable{
		userAccountTable: newUserAccountTableImpl(schemaName, tableName, alias),
		EXCLUDED:         newUserAccountTableImpl("", "excluded", ""),
	}
}

func newUserAccountTableImpl(schemaName, tableName, alias string) userAccountTable {
	var (
		IDColumn          = postgres.IntegerColumn("id")
		FirstNameColumn   = postgres.StringColumn("first_name")
		LastNameColumn    = postgres.StringColumn("last_name")
		EmailColumn       = postgres.StringColumn("email")
		PhoneNumberColumn = postgres.StringColumn("phone_number")
		PasswordColumn    = postgres.StringColumn("password")
		allColumns        = postgres.ColumnList{IDColumn, FirstNameColumn, LastNameColumn, EmailColumn, PhoneNumberColumn, PasswordColumn}
		mutableColumns    = postgres.ColumnList{FirstNameColumn, LastNameColumn, EmailColumn, PhoneNumberColumn, PasswordColumn}
	)

	return userAccountTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:          IDColumn,
		FirstName:   FirstNameColumn,
		LastName:    LastNameColumn,
		Email:       EmailColumn,
		PhoneNumber: PhoneNumberColumn,
		Password:    PasswordColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
