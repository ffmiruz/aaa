# Get Data Insights Using AI

A web application that helps you gain insights from any dataset using natural language queries. Ask questions about the data and get AI-powered visualizations and analysis.

## Live Demo

**Deployed Version:** [https://2go.my](https://2go.my)
[![video]()](./video.mp4)

## Overview

This application allows you to:

- Query a database of 40+ Malaysian social enterprises
- Ask natural language questions about the data
- Generate visual representations of data trends
- Explore enterprises by sector, location, impact areas, and more
- Receive AI-powered insights about social enterprise patterns

## Features

- **AI-Powered Data Analysis**: Leverages Gemini Flash model to interpret your queries
- **Interactive Web Interface**: Easy-to-use frontend for submitting questions
- **Natural Language Processing**: Ask questions in plain English
- **Multiple Output Formats**: Get responses as formatted HTML, data tables, or visualizations
- **Rich Dataset**: Access information on company names, sectors, locations, impact areas, funding sources, and more

## Example Dataset Information

The application provides access to a comprehensive dataset of Malaysian social enterprises that includes:

- Company profiles and descriptions
- Geographic locations
- Social impact areas
- Problems addressed
- Funding sources
- Institutional supporters
- Target beneficiaries
- Revenue models
- Founding years
- Award recognitions

## Prerequisites

- Go 1.21 or higher
- Google Gemini API key

## Installation

1. Clone this repository:
   ```
   git clone <repository-url>
   cd aaa
   ```

2. Set up your Gemini API key:
   - Add your API key to `gemini.txt`

3. Install dependencies:
   ```
   go mod download
   ```

## Usage

1. Start the server:
   ```
   go run server.go websiter.go
   ```

2. Open your browser and navigate to:
   ```
   http://localhost:5000
   ```

3. Enter your query in the form. Example queries:
   - "Which sectors have the most social enterprises?"
   - "Show me all enterprises supporting refugee communities"
   - "What are the top impact areas in Kuala Lumpur?"
   - "Compare women empowerment initiatives across different regions"
   - "Which enterprises received MaGIC grants?"

4. The AI will process your request and return insights based on the data.

## API Endpoints

- `POST /api`: Submit data queries with the `input` form parameter
- `/hello`: Test endpoint that returns a simple greeting
- `/`: Home page with the web interface

## Sample Queries

Here are some example queries you can try:

- "List all enterprises in Selangor"
- "Which social enterprises focus on women empowerment?"
- "Show me enterprises founded after 2018"
- "What are the most common problems being addressed?"
- "Compare rural versus urban social enterprises"
- "Visualize distribution of enterprises by sector"

## Docker Support

This project includes Docker support for easy deployment:

```
docker build -t data-insight-generator .
docker run -p 5000:5000 data-insight-generator
```

## License

Unlicensed

## Acknowledgments

- Google Gemini API for powering the AI functionality
- Go standard library for web server implementation
