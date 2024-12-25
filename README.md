# Project Title
Search-With-Go

## Demo Link
This is the [deployment](https://jocular-bubblegum-a47d9d.netlify.app/) for the project. The live version of the project contains smaller scraped data for web search due to limited compute to store large amounts of data. The demo show below contains larger amounts of data since data is stored locally, therefore presenting a more accurate description of the project.

![demo-release](https://github.com/user-attachments/assets/1377f4fc-bb20-43f5-b77d-48b04cde73c7)

## Table of Content:
- About
- System Design
- Technologies
- Setup
- Approach
- Next Steps

## About
This is a hobby project, showcasing my interests in search engine technologies. I researched how search engines, then implemented techniques in Natural Language Processing such as TF-IDF, Information Retreival such as Inverted Index, and web scraping such as politeness. The result is a search engine that ranks web pages based on their relevance to the user's query.

## System Design

### High Level Architecture
![image](https://github.com/user-attachments/assets/7f91ff05-0f6b-4559-9e7c-e5c6771677b0)

### Concurrency Model
![image](https://github.com/user-attachments/assets/ab7d8e91-4f43-47a6-a454-635eb41f8213)

### Database Modelling
![image](https://github.com/user-attachments/assets/8ce1a0a4-55fe-4cf3-9745-8a02e08b3ace)

## Technologies
- Frontend: This project uses [Vite](https://vite.dev/) for frontend tooling and [React](https://react.dev/) as a frontend library. Also, it uses [Tailwind CSS](https://tailwindcss.com/) for styling.
- Backend: The project uses [Go](https://go.dev/) because Go's concurrency paradigm is simple to work with, enabling a performant web crawler. It uses [Chromedp](https://github.com/chromedp/chromedp) to render dynamic websites in headless chrome. To perform TF-IDF computation, it used a [Go library](https://github.com/agonopol/go-stem?tab=readme-ov-file) as an implementation of the [porter stemming algorithm](https://tartarus.org/martin/PorterStemmer/index.html). When scraping websites, the crawler must respect the restrictions specified in the website's robots.txt file. To parse the robots.txt file, the project used this [Go library](https://github.com/benjaminestes/robots?tab=readme-ov-file).







