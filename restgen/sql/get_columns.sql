SELECT column_name as name
	, CASE data_type
		WHEN 'integer' THEN 'int'
		WHEN 'text' THEN 'string'
	  END as data_type
	, is_nullable = 'YES' as is_nullable
FROM information_schema.columns
WHERE table_schema = $1
AND table_name = $2
ORDER BY ordinal_position