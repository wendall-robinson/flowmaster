package traceflow

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/attribute"
)

// TestAddDBQuery tests the AddDBQuery method.
func TestAddDBQuery(t *testing.T) {
	ctx := context.TODO()
	trace := New(ctx, "test-service")
	trace.AddDBQuery("SELECT * FROM users", "mysql")

	if len(trace.attrs) != 2 {
		t.Fatalf("Expected 2 attributes, got %d", len(trace.attrs))
	}
	if trace.attrs[0] != attribute.String("db.query", "SELECT * FROM users") {
		t.Errorf("Expected db.query to be 'SELECT * FROM users', got %v", trace.attrs[0])
	}
	if trace.attrs[1] != attribute.String("db.system", "mysql") {
		t.Errorf("Expected db.system to be 'mysql', got %v", trace.attrs[1])
	}
}

// TestAddDBInfo tests the AddDBInfo method.
func TestAddDBInfo(t *testing.T) {
	ctx := context.TODO()
	trace := New(ctx, "test-service")
	trace.AddDBInfo("test_db", "5.7")

	if len(trace.attrs) != 2 {
		t.Fatalf("Expected 2 attributes, got %d", len(trace.attrs))
	}
	if trace.attrs[0] != attribute.String("db.name", "test_db") {
		t.Errorf("Expected db.name to be 'test_db', got %v", trace.attrs[0])
	}
	if trace.attrs[1] != attribute.String("db.version", "5.7") {
		t.Errorf("Expected db.version to be '5.7', got %v", trace.attrs[1])
	}
}

// TestAddDBConnectionInfo tests the AddDBConnectionInfo method.
func TestAddDBConnectionInfo(t *testing.T) {
	ctx := context.TODO()
	trace := New(ctx, "test-service")
	trace.AddDBConnectionInfo("user:password@tcp(localhost:3306)/test_db", 10)

	if len(trace.attrs) != 2 {
		t.Fatalf("Expected 2 attributes, got %d", len(trace.attrs))
	}
	if trace.attrs[0] != attribute.String("db.connection_string", "user:password@tcp(localhost:3306)/test_db") {
		t.Errorf("Expected db.connection_string to be 'user:password@tcp(localhost:3306)/test_db', got %v", trace.attrs[0])
	}
	if trace.attrs[1] != attribute.Int("db.connection_count", 10) {
		t.Errorf("Expected db.connection_count to be 10, got %v", trace.attrs[1])
	}
}

// TestAddDBTableInfo tests the AddDBTableInfo method.
func TestAddDBTableInfo(t *testing.T) {
	ctx := context.TODO()
	trace := New(ctx, "test-service")
	trace.AddDBTableInfo("users", 100)

	if len(trace.attrs) != 2 {
		t.Fatalf("Expected 2 attributes, got %d", len(trace.attrs))
	}
	if trace.attrs[0] != attribute.String("db.table_name", "users") {
		t.Errorf("Expected db.table_name to be 'users', got %v", trace.attrs[0])
	}
	if trace.attrs[1] != attribute.Int("db.row_count", 100) {
		t.Errorf("Expected db.row_count to be 100, got %v", trace.attrs[1])
	}
}

// TestAddDBIndexInfo tests the AddDBIndexInfo method.
func TestAddDBIndexInfo(t *testing.T) {
	ctx := context.TODO()
	trace := New(ctx, "test-service")
	trace.AddDBIndexInfo("idx_users", 3)

	if len(trace.attrs) != 2 {
		t.Fatalf("Expected 2 attributes, got %d", len(trace.attrs))
	}
	if trace.attrs[0] != attribute.String("db.index_name", "idx_users") {
		t.Errorf("Expected db.index_name to be 'idx_users', got %v", trace.attrs[0])
	}
	if trace.attrs[1] != attribute.Int("db.index_count", 3) {
		t.Errorf("Expected db.index_count to be 3, got %v", trace.attrs[1])
	}
}

// TestAddDBColumnInfo tests the AddDBColumnInfo method.
func TestAddDBColumnInfo(t *testing.T) {
	ctx := context.TODO()
	trace := New(ctx, "test-service")
	trace.AddDBColumnInfo("user_id", 5)

	if len(trace.attrs) != 2 {
		t.Fatalf("Expected 2 attributes, got %d", len(trace.attrs))
	}
	if trace.attrs[0] != attribute.String("db.column_name", "user_id") {
		t.Errorf("Expected db.column_name to be 'user_id', got %v", trace.attrs[0])
	}
	if trace.attrs[1] != attribute.Int("db.column_count", 5) {
		t.Errorf("Expected db.column_count to be 5, got %v", trace.attrs[1])
	}
}

// TestAddDBTransactionInfo tests the AddDBTransactionInfo method.
func TestAddDBTransactionInfo(t *testing.T) {
	ctx := context.TODO()
	trace := New(ctx, "test-service")
	trace.AddDBTransactionInfo("txn_123", "committed")

	if len(trace.attrs) != 2 {
		t.Fatalf("Expected 2 attributes, got %d", len(trace.attrs))
	}
	if trace.attrs[0] != attribute.String("db.transaction_id", "txn_123") {
		t.Errorf("Expected db.transaction_id to be 'txn_123', got %v", trace.attrs[0])
	}
	if trace.attrs[1] != attribute.String("db.transaction_status", "committed") {
		t.Errorf("Expected db.transaction_status to be 'committed', got %v", trace.attrs[1])
	}
}

// TestAddDBErrorInfo tests the AddDBErrorInfo method.
func TestAddDBErrorInfo(t *testing.T) {
	ctx := context.TODO()
	trace := New(ctx, "test-service")
	trace.AddDBErrorInfo("Syntax error", "1234")

	if len(trace.attrs) != 2 {
		t.Fatalf("Expected 2 attributes, got %d", len(trace.attrs))
	}
	if trace.attrs[0] != attribute.String("db.error_message", "Syntax error") {
		t.Errorf("Expected db.error_message to be 'Syntax error', got %v", trace.attrs[0])
	}
	if trace.attrs[1] != attribute.String("db.error_code", "1234") {
		t.Errorf("Expected db.error_code to be '1234', got %v", trace.attrs[1])
	}
}
