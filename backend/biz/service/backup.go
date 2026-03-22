package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"nanamiku-blog/backend/query"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BackupService struct {
	q  *query.Queries
	db *pgxpool.Pool
}

func NewBackupService(db *pgxpool.Pool) *BackupService {
	return &BackupService{q: query.New(db), db: db}
}

var backupTableOrder = []string{
	"admin_users",
	"admin_refresh_tokens",
	"visitors",
	"posts",
	"tags",
	"post_tags",
	"post_revisions",
	"post_likes",
	"post_view_daily",
	"post_comments",
	"guestbook_messages",
	"guestbook_votes",
	"moments",
	"moment_likes",
	"moment_reposts",
	"moment_comments",
	"moment_comment_likes",
	"friend_links",
	"friend_link_health_logs",
	"sensitive_words",
	"blocked_ips",
	"audit_logs",
	"analytics_sessions",
	"analytics_pageviews",
}

func (s *BackupService) BuildJSON(ctx context.Context) ([]byte, error) {
	payload, err := s.q.ExportBackupPayload(ctx)
	if err != nil {
		return nil, fmt.Errorf("export backup payload: %w", err)
	}

	data := make(map[string]json.RawMessage)
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, fmt.Errorf("decode backup payload: %w", err)
	}

	tables := make([]string, 0, len(data))
	for table := range data {
		tables = append(tables, table)
	}
	sort.Strings(tables)

	doc := map[string]interface{}{
		"meta": map[string]interface{}{
			"format":       "json",
			"generated_at": time.Now().UTC().Format(time.RFC3339),
			"table_count":  len(tables),
			"tables":       tables,
		},
		"data": data,
	}

	bytes, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal backup json: %w", err)
	}
	return bytes, nil
}

func (s *BackupService) BuildSQL(ctx context.Context) ([]byte, error) {
	payload, err := s.q.ExportBackupPayload(ctx)
	if err != nil {
		return nil, fmt.Errorf("export backup payload: %w", err)
	}

	data := make(map[string]json.RawMessage)
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, fmt.Errorf("decode backup payload: %w", err)
	}

	tableColumns, err := s.loadBackupTableColumns(ctx)
	if err != nil {
		return nil, fmt.Errorf("load backup table columns: %w", err)
	}

	lines := make([]string, 0, 4+len(backupTableOrder))
	lines = append(lines,
		"-- NanaMiku Blog SQL Backup",
		"-- Generated at "+time.Now().UTC().Format(time.RFC3339),
		"BEGIN;",
		"TRUNCATE TABLE "+strings.Join(backupTableOrder, ", ")+" RESTART IDENTITY CASCADE;",
	)

	for _, table := range backupTableOrder {
		raw, presentColumns, err := normalizeBackupRows(table, data[table])
		if err != nil {
			return nil, fmt.Errorf("normalize backup rows for %s: %w", table, err)
		}

		insertColumns := pickInsertColumns(tableColumns[table], presentColumns)
		if len(insertColumns) == 0 {
			continue
		}

		columnList := strings.Join(insertColumns, ", ")
		escaped := escapeSQLLiteral(string(raw))
		lines = append(lines,
			fmt.Sprintf(
				"INSERT INTO %s (%s) SELECT %s FROM jsonb_populate_recordset(NULL::%s, '%s'::jsonb);",
				table,
				columnList,
				columnList,
				table,
				escaped,
			),
		)
	}

	lines = append(lines, buildMomentsCompatSQL())
	lines = append(lines, "COMMIT;")
	return []byte(strings.Join(lines, "\n") + "\n"), nil
}

func (s *BackupService) loadBackupTableColumns(ctx context.Context) (map[string][]string, error) {
	rows, err := s.db.Query(ctx, `
SELECT table_name, column_name
FROM information_schema.columns
WHERE table_schema = 'public'
  AND table_name = ANY($1::text[])
ORDER BY table_name, ordinal_position
`, backupTableOrder)
	if err != nil {
		return nil, fmt.Errorf("query information_schema.columns: %w", err)
	}
	defer rows.Close()

	result := make(map[string][]string, len(backupTableOrder))
	for _, table := range backupTableOrder {
		result[table] = []string{}
	}

	for rows.Next() {
		var tableName string
		var columnName string
		if err := rows.Scan(&tableName, &columnName); err != nil {
			return nil, fmt.Errorf("scan table columns: %w", err)
		}
		result[tableName] = append(result[tableName], columnName)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate table columns: %w", err)
	}

	for _, table := range backupTableOrder {
		if len(result[table]) == 0 {
			return nil, fmt.Errorf("table %s has no columns in public schema", table)
		}
	}

	return result, nil
}

func normalizeBackupRows(table string, raw json.RawMessage) (json.RawMessage, map[string]struct{}, error) {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		raw = json.RawMessage("[]")
	}

	var rows []map[string]json.RawMessage
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, nil, err
	}

	presentColumns := make(map[string]struct{})
	for i := range rows {
		for key := range rows[i] {
			presentColumns[key] = struct{}{}
		}

		// jsonb `null` in this column cannot be restored to a NOT NULL jsonb field.
		if table == "audit_logs" && isJSONNull(rows[i]["detail"]) {
			rows[i]["detail"] = json.RawMessage("{}")
		}
	}

	normalizedRaw, err := json.Marshal(rows)
	if err != nil {
		return nil, nil, err
	}

	return normalizedRaw, presentColumns, nil
}

func isJSONNull(raw json.RawMessage) bool {
	return bytes.Equal(bytes.TrimSpace(raw), []byte("null"))
}

func pickInsertColumns(tableColumns []string, presentColumns map[string]struct{}) []string {
	columns := make([]string, 0, len(tableColumns))
	for _, column := range tableColumns {
		if _, ok := presentColumns[column]; ok {
			columns = append(columns, column)
		}
	}
	return columns
}

func buildMomentsCompatSQL() string {
	return strings.TrimSpace(`
DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'public' AND table_name = 'moments' AND column_name = 'publish_status'
    ) AND EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'public' AND table_name = 'moments' AND column_name = 'published_at'
    ) THEN
        EXECUTE '
            UPDATE moments
            SET published_at = created_at
            WHERE publish_status = ''published''
              AND published_at IS NULL
        ';
    END IF;
END $$;
`)
}

func escapeSQLLiteral(input string) string {
	return strings.ReplaceAll(input, "'", "''")
}
