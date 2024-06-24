SELECT column_name as name
	, CASE udt_name
		WHEN 'int4' THEN 'int'
		WHEN 'integer' THEN 'int'
		WHEN 'text' THEN 'string'
		WHEN 'numeric' THEN 'float64'
		WHEN 'timestamp' THEN 'time.Time'
		WHEN 'bool' THEN 'bool'
	  END as data_type
	, is_nullable = 'YES' as is_nullable
FROM information_schema.columns
WHERE table_schema = $1
AND table_name = $2
ORDER BY ordinal_position