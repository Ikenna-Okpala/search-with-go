package database

const (
	InsertWebsite = `INSERT INTO websites (url)
	VALUES ($1)
	ON CONFLICT (url) DO NOTHING
	RETURNING id
	`
	InsertKeyword = `INSERT INTO keywords (word)
	VALUES ($1)
	ON CONFLICT (word) DO NOTHING
	RETURNING id
	`

	InsertWebsiteKeyword = `INSERT into website_keywords (website_id, keyword_id, tf_idf, idf)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (website_id, keyword_id) DO UPDATE
	SET tf_idf = EXCLUDED.tf_idf, idf = EXCLUDED.idf
	`

	SearchWord = `SELECT keywords.word, websites.url, website_keywords.idf, website_keywords.tf_idf FROM website_keywords
	INNER JOIN keywords ON website_keywords.keyword_id = keywords.id
	INNER JOIN websites ON website_keywords.website_id = websites.id
	WHERE keywords.word IN (%s)`
)