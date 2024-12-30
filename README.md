# Project Title
Search-With-Go

## Demo Link
This is the [deployment](https://jocular-bubblegum-a47d9d.netlify.app/) for the project. The live version of the project contains smaller scraped data due to limited compute to store large amounts of data. The demo below contains more data since data is stored locally, therefore presenting a more accurate description of the project.

![demo-release](https://github.com/user-attachments/assets/1377f4fc-bb20-43f5-b77d-48b04cde73c7)

## Table of Content:
- [About](https://github.com/Ikenna-Okpala/search-with-go/edit/main/README.md#about)
- [System Design](https://github.com/Ikenna-Okpala/search-with-go/edit/main/README.md#system-design)
- [Technologies](https://github.com/Ikenna-Okpala/search-with-go/edit/main/README.md#system-design)
- [Dev Environmet](https://github.com/Ikenna-Okpala/search-with-go/edit/main/README.md#system-design)
- [Approach](https://github.com/Ikenna-Okpala/search-with-go/edit/main/README.md#system-design)
- [Next Steps](https://github.com/Ikenna-Okpala/search-with-go/edit/main/README.md#system-design)

## About
This is a hobby project, showcasing my interests in search engine technologies. I researched how search engines work, then implemented techniques in natural language processing such as TF-IDF computation, information retreival such as inverted indices, and web scraping such as politeness. The result is a search engine that ranks web pages based on their relevance to the user's query.

## System Design

### High Level Architecture
![image](https://github.com/user-attachments/assets/7f91ff05-0f6b-4559-9e7c-e5c6771677b0)

### Concurrency Model
![image](https://github.com/user-attachments/assets/919345f0-d712-466c-8b1b-ea6e04ce390c)


### Database Modelling
![image](https://github.com/user-attachments/assets/8ce1a0a4-55fe-4cf3-9745-8a02e08b3ace)

## Technologies
- Frontend: This project uses [Vite](https://vite.dev/) for frontend tooling and [React](https://react.dev/) as a frontend library. Also, it uses [Tailwind CSS](https://tailwindcss.com/) for styling.
- Backend: The project uses [Go](https://go.dev/) because Go's concurrency paradigm makes developing multithreaded web crawlers easier. It uses [Chromedp](https://github.com/chromedp/chromedp) to render dynamic websites in headless chrome. To perform TF-IDF computation, it uses a [Go library](https://github.com/agonopol/go-stem?tab=readme-ov-file) as an implementation of the [porter stemming algorithm](https://tartarus.org/martin/PorterStemmer/index.html) to convert words to their root form. When scraping websites, the crawler must respect the restrictions specified in the website's robots.txt file. To parse the robots.txt file, the project uses this [library](https://github.com/benjaminestes/robots?tab=readme-ov-file).

## Dev Envorionment

### Requirements
- [Go 1.23+](https://go.dev/dl/)
- [PostgresSQL 15.3+](https://www.postgresql.org/download/)
- GCP Project (must have Cloud Storage enabled with a bucket created)

### Steps
- Clone this repository:
  - ```git clone https://github.com/Ikenna-Okpala/search-with-go.git```
 
- Enviroment Variables:
  - In the terminal session, set the following environment variables:
    - ```export SERVER_PORT=<your-server-port>```
    -  ```export SERVER_IP=localhost```
    -  ```export BUCKET_NAME=<your-gcp-cloud-bucket-name>```
    -  ```export DB_LOCAL_URL=postgresql://<postgres-user>:<password>@localhost/<database-name>?sslmode=disable```

- Dependencies:
  - In the root of the project, run:
     - ```go mod download```
     - For GCP client SDK, ensure the [default credentials](https://cloud.google.com/docs/authentication/application-default-credentials) is configured on your machine before moving on to the next steps
   
- Database:
  - Start a session:
      - ```psql -U username -d database_name```
  - Load Schema:
    - ```\i ${root}/internal/database/schema.sql``` (only works for psql client)

- Build:
  - Run crawler:
    - ```cd cmd/crawler``` 
    - ``` go run . 10``` (10 is an arbitrary number of websites to scrape)
  - Run server:
    - ``` cd cmd/app```
    - ```go run .``` (start the web server to serve incoming requests)
  - Run client:
    - ```cd client/search-with-go```
    - ```npm i```
    - ```npm run dev```

## Approach
To sort search results based on a user's search query, this project used TF-IDF.

Term Frequency (TF) captures how frequent a word shows up in a website. Furthermore, it tells us how relevant a word is in a website (when a word shows up multiple times, the TF for the word inceases).

```TF(w,d) = Total number of w in d / Total number of words in d, where w is a word, and d is a website```

Inverse Document Frequency (IDF) captures how common a word in an array of websites. The most common words are expected to have a lower IDF computation.

```IDF(w) = log2(Total number of websites / (Total number of websites that contain the word w + 1)) + 1, the first + 1 prevents negative results after computing log2, while the second + 1 prevents 0 as an IDF value```

TF-IDF assigns weight to each word in a website by merging TF and IDF computation. The words in a website with higher TF-IDF are more important to that website.

```TFIDF(w, d) = TF(w, d) * IDF(w), where w is a word in website d```

After crawling websites, the TF-IDF computation is performed and the result is stored in a relational [database](https://github.com/Ikenna-Okpala/search-with-go/edit/main/README.md#database-modelling). An [inverted index](https://en.wikipedia.org/wiki/Inverted_index) on words in websites is maintained, for blazingly fast retreival. Note that words are stored as root forms in the database using the [porter stemming algorithm](https://tartarus.org/martin/PorterStemmer/index.html).

At query time, the user's query consist of words reduced to their root form using the stemming algorithm. Using the inverted index, the database retrieves websites that contain words that matches the user's query. After, the program computes the TF-IDF of the query using the IDF retrieved from the database and the term frequency defined by:

```TF(w) = Frequency of w / Frequency of q, where w is a word, and q is the word with the highest frequency```

How can we know which websites are more relevant to a query? We represent each word in the query as its own axis in a vector space. Then, we use a technique in vector mathematics called Cosine Similarity to compute how closely related the websites are to the query. Finally, the websites are sorted based on the Cosine Similarity score, leading to a ranked list of websites where the top websites are the most relevant to the user's query.

## Next Steps
- Page Rank: Currently, the search engine ranks search results only using TF-IDF (text matching). As a result, search results are not completely ordered based on importance to the user's query. The solution is to implement a [page rank algorithm](https://en.wikipedia.org/wiki/PageRank) that counts the number of links that points to a website, indicating website importance. The result is a sorted list of search results based on importance.
  - ```Rank Score = w1 * Cosine Similarity + w2 * Page Rank, where w1 and w2 are arbitrary weights```
- Scaling:
  - Crawlers: Currently, the crawlers run in parallel on a single machine. When crawling the internet, the project needs distributed web crawlers to improve performance. To share crawling state (visited links) among crawler nodes, [Redis data structures](https://redis.io/technology/data-structures/) will be used.
  - Indexes: Currently, the search engine maintains a single PostgresSQL instance as an index. This current solution is not scalable as it introduces a Single Point of Failure (SPOF), and the search engine will lose performance as the number of users increase. Moving forward, the project will adopt horizontal sharding to partition the database into several instances, improving the latency of queries.







