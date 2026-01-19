package mysql

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
)

// SaveModel saves a model instance to the database. It requires the model to have a primary ID.
// It supports transactional databases through the txDB parameter.
func SaveModel[T any](m *T, txDB *gorm.DB) error {
	if txDB != nil {
		if err := txDB.Save(m).Error; err != nil {
			return fmt.Errorf("failed to save model: %w", err)
		}
	} else {
		if err := GetDB().Save(m).Error; err != nil {
			return fmt.Errorf("failed to save model (using default DB): %w", err)
		}
	}
	return nil
}

// UpdateModel updates a model instance in the database. It requires the model to have a primary ID.
// It supports transactional databases through the txDB parameter.
func UpdateModel[T any](m *T, txDB *gorm.DB) error {
	if txDB != nil {
		if err := txDB.Updates(m).Error; err != nil {
			return fmt.Errorf("failed to update model: %w", err)
		}
	} else {
		if err := GetDB().Updates(m).Error; err != nil {
			return fmt.Errorf("failed to update model (using default DB): %w", err)
		}
	}
	return nil
}

// UpdateWhereModel updates a model instance based on a condition.
// It does not require the model to have a primary ID but uses the model's fields as conditions.
func UpdateWhereModel[T any](m *T) error {
	if err := GetDB().Model(m).Updates(m).Error; err != nil {
		return fmt.Errorf("failed to update model by conditions: %w", err)
	}
	return nil
}

// CreateModel creates a new model instance in the database.
// It does not require the model to have a primary ID.
func CreateModel[T any](m *T, txDB *gorm.DB) error {
	if txDB != nil {
		if err := txDB.Create(m).Error; err != nil {
			return fmt.Errorf("failed to create model: %w", err)
		}
	} else {
		if err := GetDB().Create(m).Error; err != nil {
			return fmt.Errorf("failed to create model (using default DB): %w", err)
		}
	}
	return nil
}

// CreateManyModel creates multiple model instances in the database.
func CreateManyModel[T any](m *[]T, txDB *gorm.DB) error {
	if txDB != nil {
		if err := txDB.Create(m).Error; err != nil {
			return fmt.Errorf("failed to create multiple models: %w", err)
		}
	} else {
		if err := GetDB().Create(m).Error; err != nil {
			return fmt.Errorf("failed to create multiple models (using default DB): %w", err)
		}
	}
	return nil
}

// LoadModel loads a model instance by its primary ID.
func LoadModel[T any](m *T) error {
	if err := GetDB().First(m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("record not found")
		}
		return fmt.Errorf("failed to load model: %w", err)
	}
	return nil
}

// SoftDeleteModel soft-deletes a model instance from the database.
// It requires the model to have a primary ID.
func SoftDeleteModel[T any](m *T, txDB *gorm.DB) error {
	if txDB != nil {
		if err := txDB.Unscoped().Delete(m).Error; err != nil {
			return fmt.Errorf("failed to soft-delete model: %w", err)
		}
	} else {
		if err := GetDB().Unscoped().Delete(m).Error; err != nil {
			return fmt.Errorf("failed to soft-delete model (using default DB): %w", err)
		}
	}
	return nil
}

// MustDeleteModel permanently deletes a model instance from the database.
// It requires the model to have a primary ID.
func MustDeleteModel[T any](m *T, txDB *gorm.DB) error {
	if txDB != nil {
		if err := txDB.Delete(m).Error; err != nil {
			return fmt.Errorf("failed to delete model: %w", err)
		}
	} else {
		if err := GetDB().Delete(m).Error; err != nil {
			return fmt.Errorf("failed to delete model (using default DB): %w", err)
		}
	}
	return nil
}
