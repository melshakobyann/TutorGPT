# TutorGPT

TutorGPT is an AI-powered interactive tutor designed to teach subjects such as mathematics, coding, and physics..

## Project Structure

```
tutorgpt/
├── backend/               # Go backend service
│   ├── cmd/               # Application entry points
│   │   └── main.go        # Main entry point
│   ├── pkg/               # Backend packages
│   │   ├── ai/            # AI modules (lecture builder, visualization, etc.)
│   │   ├── handler/       # HTTP request handlers
│   │   ├── models/        # Data models
│   │   ├── utils/         # Utility functions (config, logging)
│   ├── go.mod             # Go module definition
│   └── go.sum             # Go module checksums
├── frontend/              # Vanilla HTML/CSS/JS frontend
│   ├── index.html         # Main HTML file
│   ├── styles.css         # CSS styles
│   ├── app.js             # Main application logic
│   └── assets/            # Static assets
├── .env                   # Environment variables (not committed to version control)
├── .env.example           # Example environment file with dummy values
├── Makefile               # Build and run scripts
└── README.md              # This file
```

## Setup Instructions

1. Clone the repository:
```bash
git clone [repository-url]
cd tutorgpt
```

2. Configure environment variables:
```bash
cp .env.example .env
# Edit .env and add your actual API keys and configuration
```

3. Run the application:
```bash
tutorgpt # for unix based
# or
tutorgpt.exe # for windows
```

Optional
If you want to build the project run
```bash
make run
```

This will start both the backend server and frontend static file server.

## Development

- Backend: Written in Go, the backend provides a stateless API service that integrates with OpenAI to deliver tutoring content.
- Frontend: Uses vanilla HTML/CSS/JavaScript to provide a chat interface with support for interactive visualizations.

## Available Make Commands

- `make build-backend` - Build the Go backend binary
- `make run-backend` - Start the backend service
- `make run-frontend` - Serve the frontend files
- `make run` - Run both backend and frontend

## Demo
<video width="600" controls>
  <source src="assets/TutorGPT.mp4" type="video/mp4">
  Your browser does not support the video tag.
</video>