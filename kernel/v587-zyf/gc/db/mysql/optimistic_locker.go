package mysql

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// OptimisticLocker struct defines the optimistic locking plugin.
type OptimisticLocker struct{}

// Name returns the name of the plugin.
func (o *OptimisticLocker) Name() string {
	return "optimisticLocker"
}

// Initialize sets up the plugin by registering callbacks before and after update operations.
func (o *OptimisticLocker) Initialize(db *gorm.DB) error {
	// Register a callback before the update operation to increment the version and add a where clause for the current version.
	if err := db.Callback().Update().Before("gorm:update").Register("optimistic_locking:before_update", o.beforeUpdate); err != nil {
		return err
	}

	// Register a callback after the update operation to check if any rows were affected.
	if err := db.Callback().Update().After("gorm:update").Register("optimistic_locking:after_update", o.afterUpdate); err != nil {
		return err
	}

	return nil
}

// beforeUpdate increments the version field and adds a where clause to ensure the version matches the previous value.
func (o *OptimisticLocker) beforeUpdate(db *gorm.DB) {
	if versionField, ok := db.Statement.Schema.FieldsByName["Version"]; ok {
		// Get the current value of the Version field from the model instance
		versionValue := db.Statement.ReflectValue.FieldByName(versionField.Name)
		if !versionValue.IsValid() || !versionValue.CanSet() {
			db.AddError(errors.New("optimistic lock error: cannot access Version field"))
			return
		}

		curVersion := versionValue.Int()
		// Increment the version
		versionValue.SetInt(curVersion + 1)

		// Add a where clause to ensure the version matches the current value
		db.Statement.AddClause(clause.Where{Exprs: []clause.Expression{clause.Eq{Column: versionField.DBName, Value: curVersion}}})
	} else {
		db.AddError(errors.New("optimistic lock error: no Version field found"))
	}
}

// afterUpdate checks if any rows were affected by the update operation.
func (o *OptimisticLocker) afterUpdate(db *gorm.DB) {
	// Check if any rows were affected
	if db.RowsAffected == 0 {
		db.AddError(errors.New("optimistic lock error: no rows were updated, indicating a concurrent modification"))
	}
}

// ExampleModel represents a simple model with an optimistic locking Version field.
type ExampleModel struct {
	gorm.Model
	Name  string
	Value int
	// Version field for optimistic locking
	Version int
}
