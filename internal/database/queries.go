package database

const (
	InsertWebsite = `INSERT INTO websites (url, title, icon, name, description)
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (url) DO UPDATE
	SET url = EXCLUDED.url, title = EXCLUDED.title, icon = EXCLUDED.icon, name = EXCLUDED.name, description = EXCLUDED.description
	RETURNING id
	`
	InsertKeyword = `INSERT INTO keywords (word)
	VALUES ($1)
	ON CONFLICT (word) DO UPDATE
	SET word = EXCLUDED.word
	RETURNING id
	`

	InsertWebsiteKeyword = `INSERT into website_keywords (website_id, keyword_id, tf_idf, idf)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (website_id, keyword_id) DO UPDATE
	SET tf_idf = EXCLUDED.tf_idf, idf = EXCLUDED.idf
	`

	SearchWord = `SELECT keywords.word, websites.url, websites.title, websites.icon, websites.name, websites.description, website_keywords.idf, website_keywords.tf_idf FROM website_keywords
	INNER JOIN keywords ON website_keywords.keyword_id = keywords.id
	INNER JOIN websites ON website_keywords.website_id = websites.id
	WHERE keywords.word IN (%s)`
)