package service

import (
	"encoding/json"
	"testing"
)

func TestNormalizeBackupRowsAuditLogsDetailNull(t *testing.T) {
	raw := json.RawMessage(`[
		{"id":"1","action":"approve","detail":null},
		{"id":"2","action":"reject","detail":{"reason":"manual"}}
	]`)

	normalized, presentColumns, err := normalizeBackupRows("audit_logs", raw)
	if err != nil {
		t.Fatalf("normalize backup rows: %v", err)
	}

	var rows []map[string]interface{}
	if err := json.Unmarshal(normalized, &rows); err != nil {
		t.Fatalf("decode normalized rows: %v", err)
	}

	firstDetail, ok := rows[0]["detail"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected first detail to be object, got %T", rows[0]["detail"])
	}
	if len(firstDetail) != 0 {
		t.Fatalf("expected first detail to be empty object, got %#v", firstDetail)
	}

	if _, ok := presentColumns["detail"]; !ok {
		t.Fatalf("expected present columns to include detail")
	}
	if _, ok := presentColumns["action"]; !ok {
		t.Fatalf("expected present columns to include action")
	}
}

func TestNormalizeBackupRowsEmptyInput(t *testing.T) {
	normalized, presentColumns, err := normalizeBackupRows("posts", nil)
	if err != nil {
		t.Fatalf("normalize backup rows: %v", err)
	}
	if string(normalized) != "[]" {
		t.Fatalf("expected empty json array, got %s", string(normalized))
	}
	if len(presentColumns) != 0 {
		t.Fatalf("expected no columns for empty payload, got %d", len(presentColumns))
	}
}

func TestPickInsertColumnsPreservesSchemaOrder(t *testing.T) {
	tableColumns := []string{"id", "content", "publish_status", "scheduled_at"}
	presentColumns := map[string]struct{}{
		"id":           {},
		"content":      {},
		"scheduled_at": {},
		"extra_field":  {},
	}

	columns := pickInsertColumns(tableColumns, presentColumns)
	if len(columns) != 3 {
		t.Fatalf("expected 3 columns, got %d", len(columns))
	}
	if columns[0] != "id" || columns[1] != "content" || columns[2] != "scheduled_at" {
		t.Fatalf("unexpected column order: %#v", columns)
	}
}
