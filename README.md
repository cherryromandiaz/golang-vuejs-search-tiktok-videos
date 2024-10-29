# Custom TikTok Video Player Prototype

This project is a TikTok video search and display component built with Vue.js, styled for flexible integration and customization within other applications. It includes lazy loading, caching, and sorting features, with a Go backend that fetches TikTok videos based on user defined search terms.

## Features

- **Search and Display**: Search TikTok videos by keyword and display results in a responsive grid layout.
- **Embed Video Playback**: On clicking a video, fetch and display an embedded TikTok player.
- **Sorting and Filtering**: Sort videos by relevance or date.
- **Caching and Lazy Loading**: Cache search results to reduce network requests and improve performance.


## Technologies Used

- **Golang** for backend
- **Vue.js framework**: for frontend

## Getting Started

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/cherryromandiaz/golang-vuejs-search-tiktok-videos.git
   cd golang-vuejs-search-tiktok-videos

2. Install Frontend Dependencies:
   ```bash
   cd frontend
   npm install
   
3. Start the Go Backend
   The Go backend serves as the API endpoint for retrieving TikTok videos. Ensure Go is installed, then start the backend:
   ```bash
   cd golang-vuejs-search-tiktok-videos
   go run main.go

4. Configure and Run Vue Frontend
   With the backend running, return to the frontend directory and start the Vue app.
   ```bash
   npm run serve

## Frontend File Structure
- SearchComponent.vue: Main Vue js file.
- SearchComponentTemplate.html: HTML template for the Vue component.
- SearchComponentScript.js: Handles data properties, methods, and interaction with the Go backend.
- SearchComponentStyle.css: Customizable CSS to match various application styles.

## Go Backend
  The Go backend retrieves and processes TikTok search results. It provides:

- Search Endpoint: Supports querying TikTok videos by search term.
- Image and Video URLs: Returns URLs for displaying video thumbnails and embedding the video player.
  
### Contributing
Contributions are welcome! Please open an issue or submit a pull request for any improvements or suggestions.

### License
This project is licensed under the MIT License.
