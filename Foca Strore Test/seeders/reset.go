package seeders

import "gorm.io/gorm"

func ResetDatabase(db *gorm.DB) error {
	if err := DropAllTable(db); err != nil {
		return err
	}

	if err := MigrateAll(db); err != nil {
		return err
	}

	if err := SeedRoles(db); err != nil {
			return err
		}

	if err := SeedAdmin(db); err != nil {
		return err
	}

	if err := SeedCategories(db); err != nil {
		return err
	}

	if err := SeedUsers(db); err != nil {
		return err
	}

	if err := SeedCoupons(db); err != nil {
		return err
	}
	
	return nil
}

func ResetDatabaseWithProduct(db *gorm.DB) error {
	if err := DropAllTable(db); err != nil {
		return err
	}

	if err := MigrateAll(db); err != nil {
		return err
	}

	if err := SeedRoles(db); err != nil {
			return err
		}

	if err := SeedAdmin(db); err != nil {
		return err
	}

	if err := SeedCategories(db); err != nil {
		return err
	}

	if err := SeedUsers(db); err != nil {
		return err
	}

	if err := SeedCoupons(db); err != nil {
		return err
	}

	if err := SeedProducts(db); err != nil {
		return err
	}
	
	return nil
}