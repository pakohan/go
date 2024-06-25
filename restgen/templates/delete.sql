DELETE FROM "{{.Schema}}"."{{.Model.TableName}}"
WHERE id = $1