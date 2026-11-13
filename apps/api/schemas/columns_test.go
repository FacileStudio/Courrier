package schemas

import (
	"errors"
	"reflect"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	portepg "github.com/FacileStudio/porte/pg"
	"github.com/FacileStudio/tronc/mcp"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var queryModels = []any{
	&User{}, &AppSetting{}, &Account{}, &Folder{}, &Email{}, &ThreadLink{}, &Attachment{},
	&Contact{}, &EmailTemplate{}, &Space{}, &SpaceMember{}, &SenderRule{},
}

var timeType = reflect.TypeFor[time.Time]()

func compatibleTypes(field reflect.Type) []string {
	for field.Kind() == reflect.Pointer {
		field = field.Elem()
	}
	if field == timeType {
		return []string{"timestamp with time zone"}
	}
	switch field.Kind() {
	case reflect.String:
		return []string{"text", "character varying", "uuid"}
	case reflect.Bool:
		return []string{"boolean"}
	case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint64:
		return []string{"bigint"}
	case reflect.Int32, reflect.Uint32:
		return []string{"integer", "bigint"}
	case reflect.Int16, reflect.Uint16:
		return []string{"smallint", "integer", "bigint"}
	default:
		return nil
	}
}

func migratedColumnTypes(t *testing.T, db *gorm.DB, table string) map[string]string {
	t.Helper()
	var rows []struct {
		ColumnName string
		DataType   string
	}
	err := db.Raw(`SELECT column_name, data_type FROM information_schema.columns
		WHERE table_schema = current_schema() AND table_name = ?`, table).Scan(&rows).Error
	if err != nil {
		t.Fatalf("read the columns of %s: %v", table, err)
	}
	columns := map[string]string{}
	for _, row := range rows {
		columns[row.ColumnName] = row.DataType
	}
	return columns
}

// TestGORMModelsMatchTheMigratedSchema checks every column a GORM model reads
// or writes exists in the schema the migrations build, with a type the field
// can hold. AutoMigrate no longer adds a missing column, so a field added to a
// struct without a migration fails here instead of at query time.
func TestGORMModelsMatchTheMigratedSchema(t *testing.T) {
	db := openTestDatabase(t)
	for _, model := range queryModels {
		parsed, err := schema.Parse(model, &sync.Map{}, db.NamingStrategy)
		if err != nil {
			t.Fatalf("parse %T: %v", model, err)
		}
		columns := migratedColumnTypes(t, db, parsed.Table)
		if len(columns) == 0 {
			t.Errorf("%T: table %s does not exist", model, parsed.Table)
			continue
		}
		for _, field := range parsed.Fields {
			if field.DBName == "" {
				continue
			}
			got, ok := columns[field.DBName]
			if !ok {
				t.Errorf("%s.%s: no such column", parsed.Table, field.DBName)
				continue
			}
			want := compatibleTypes(field.FieldType)
			if !slices.Contains(want, got) {
				t.Errorf("%s.%s: column is %s, field %s needs one of %v", parsed.Table, field.DBName, got, field.FieldType, want)
			}
		}
	}
}

func tableShape(t *testing.T, db *gorm.DB, table string, withIndexes bool) []string {
	t.Helper()
	query := `SELECT column_name || ' ' || data_type || ' ' || is_nullable || ' ' || coalesce(column_default, '')
		FROM information_schema.columns WHERE table_schema = current_schema() AND table_name = ?`
	args := []any{table}
	if withIndexes {
		query += ` UNION ALL SELECT indexdef FROM pg_indexes WHERE schemaname = current_schema() AND tablename = ?`
		args = append(args, table)
	}
	var shape []string
	if err := db.Raw(query+` ORDER BY 1`, args...).Scan(&shape).Error; err != nil {
		t.Fatalf("read the shape of %s: %v", table, err)
	}
	return shape
}

// TestLibrarySchemaIsInTheMigrations checks the migrations already hold
// tronc/mcp's PendingSchema, which Courrier used to execute at every boot. A
// tronc bump that changes it makes this fail, and the change then needs a
// numbered migration of its own.
func TestLibrarySchemaIsInTheMigrations(t *testing.T) {
	db := openTestDatabase(t)
	before := tableShape(t, db, "mcp_actions", true)
	if err := db.Exec(mcp.PendingSchema).Error; err != nil {
		t.Fatalf("run PendingSchema: %v", err)
	}
	if after := tableShape(t, db, "mcp_actions", true); !slices.Equal(before, after) {
		t.Fatalf("PendingSchema changed mcp_actions:\nbefore %v\nafter  %v", before, after)
	}
}

// TestPorteTablesMatchPorteSchema builds porte/pg's Schema in a throwaway
// schema and checks the three porte tables Courrier keeps, as a private copy
// pointed at its own users, have the same columns. A porte bump that adds or
// changes a column there fails here and needs a numbered migration.
func TestPorteTablesMatchPorteSchema(t *testing.T) {
	db := openTestDatabase(t)
	tables := []string{"porte_identities", "porte_login_codes", "porte_sessions"}
	migrated := map[string][]string{}
	for _, table := range tables {
		migrated[table] = tableShape(t, db, table, false)
	}
	library := map[string][]string{}
	err := db.Transaction(func(tx *gorm.DB) error {
		statements := []string{`CREATE SCHEMA courrier_porte_drift`, `SET LOCAL search_path TO courrier_porte_drift`, portepg.Schema}
		for _, statement := range statements {
			if err := tx.Exec(statement).Error; err != nil {
				return err
			}
		}
		for _, table := range tables {
			library[table] = tableShape(t, tx, table, false)
		}
		return errRollback
	})
	if !errors.Is(err, errRollback) {
		t.Fatalf("build porte/pg's Schema: %v", err)
	}
	for _, table := range tables {
		if !slices.Equal(migrated[table], library[table]) {
			t.Errorf("%s drifted from porte/pg:\nmigrations %s\nporte      %s", table,
				strings.Join(migrated[table], ", "), strings.Join(library[table], ", "))
		}
	}
}

var errRollback = errors.New("rollback")
