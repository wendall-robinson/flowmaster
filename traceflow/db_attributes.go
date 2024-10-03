package traceflow

import "go.opentelemetry.io/otel/attribute"

// AddDBQuery adds database query information to the trace.
func (t *Trace) AddDBQuery(query, dbType string) *Trace {
	t.attrs = append(t.attrs,
		attribute.String("db.query", query),
		attribute.String("db.system", dbType),
	)

	return t
}

// AddDBInfo adds database-related attributes like database name and version.
func (t *Trace) AddDBInfo(dbName, dbVersion string) *Trace {
	t.attrs = append(t.attrs,
		attribute.String("db.name", dbName),
		attribute.String("db.version", dbVersion),
	)

	return t
}

// AddDBConnectionInfo adds database connection-related attributes like connection string and connection count.
func (t *Trace) AddDBConnectionInfo(connectionString string, connectionCount int) *Trace {
	t.attrs = append(t.attrs,
		attribute.String("db.connection_string", connectionString),
		attribute.Int("db.connection_count", connectionCount),
	)

	return t
}

// AddDBTableInfo adds database table-related attributes like table name and row count.
func (t *Trace) AddDBTableInfo(tableName string, rowCount int) *Trace {
	t.attrs = append(t.attrs,
		attribute.String("db.table_name", tableName),
		attribute.Int("db.row_count", rowCount),
	)

	return t
}

// AddDBIndexInfo adds database index-related attributes like index name and index count.
func (t *Trace) AddDBIndexInfo(indexName string, indexCount int) *Trace {
	t.attrs = append(t.attrs,
		attribute.String("db.index_name", indexName),
		attribute.Int("db.index_count", indexCount),
	)

	return t
}

// AddDBColumnInfo adds database column-related attributes like column name and column count.
func (t *Trace) AddDBColumnInfo(columnName string, columnCount int) *Trace {
	t.attrs = append(t.attrs,
		attribute.String("db.column_name", columnName),
		attribute.Int("db.column_count", columnCount),
	)

	return t
}

// AddDBTransactionInfo adds database transaction-related attributes like transaction ID and status.
func (t *Trace) AddDBTransactionInfo(transactionID, status string) *Trace {
	t.attrs = append(t.attrs,
		attribute.String("db.transaction_id", transactionID),
		attribute.String("db.transaction_status", status),
	)

	return t
}

// AddDBErrorInfo adds database error-related attributes like error message and error code.
func (t *Trace) AddDBErrorInfo(errorMessage, errorCode string) *Trace {
	t.attrs = append(t.attrs,
		attribute.String("db.error_message", errorMessage),
		attribute.String("db.error_code", errorCode),
	)

	return t
}
