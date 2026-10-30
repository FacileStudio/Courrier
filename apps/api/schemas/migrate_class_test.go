package schemas

import (
	"testing"

	"gorm.io/gorm"
)

func TestMigrateTurnsNullClassColumnsIntoEmptyStrings(t *testing.T) {
	orm := openTestDatabase(t)
	acc := Account{UserID: 1, Name: "A", Email: "a@example.test"}
	if err := orm.Create(&acc).Error; err != nil {
		t.Fatal(err)
	}
	row := Email{AccountID: acc.ID, MessageID: "<m@x.test>"}
	if err := orm.Create(&row).Error; err != nil {
		t.Fatal(err)
	}
	if err := orm.Exec("UPDATE emails SET category = NULL, bulk_reason = NULL").Error; err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		if err := Migrate(orm); err != nil {
			t.Fatal(err)
		}
	}
	fresh := Email{AccountID: acc.ID, MessageID: "<n@x.test>"}
	orm.Create(&fresh)
	if n := countNullClass(orm); n != 0 {
		t.Fatalf("%d rows hold NULL after migrate or on insert", n)
	}
}

func countNullClass(orm *gorm.DB) int64 {
	var n int64
	orm.Model(&Email{}).Where("category IS NULL OR bulk_reason IS NULL").Count(&n)
	return n
}
