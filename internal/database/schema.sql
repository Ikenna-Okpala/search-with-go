CREATE TABLE websites(
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    url text NOT NULL,
    title text,
    icon text,
    name text NOT NULL,
    description TEXT
);

CREATE TABLE keywords(
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    word VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE website_keywords(
    keyword_id UUID NOT NULL,
    website_id UUID NOT NULL,
    tf_idf NUMERIC NOT NULL,
    idf NUMERIC NOT NULL,
    PRIMARY KEY (keyword_id, website_id),
    CONSTRAINT fk_websites
        FOREIGN KEY(website_id)
        REFERENCES websites(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_keywords
        FOREIGN KEY(keyword_id)
        REFERENCES keywords(id)
        ON DELETE CASCADE
);

CREATE UNIQUE INDEX word_index ON keywords (word);
CREATE UNIQUE INDEX website_index ON websites (url);