# ⚡ Flashback

A powerful command-line tool that serves as your second memory, intelligently storing and retrieving your notes using AI-powered semantic search. Never lose track of important information again!

## Demo
![Demo GIF](demo.gif)

## ✨ Features

### 🧠 AI-Powered Note Management
- **Smart Storage**: Automatically generates embeddings for your notes using Google's Gemini AI
- **Semantic Search**: Find notes using natural language queries, not just keywords
- **Intelligent Recall**: Get contextually relevant responses based on your stored notes
- **Timestamp-Aware**: Automatically tracks when notes were created with human-friendly time formatting

### 🎯 Three Modes
1. **Note Mode**: Quickly capture thoughts, ideas, and information
    - **URL Mode**: Fetch content from URLs. usage -> web:https://example.com
    - **File Mode**: Use content from files (only text files for now). usage -> file:/path/to/file.txt
    - **Clipboard Mode**: Use content from clipboard. usage -> #clipboard
2. **Recall Mode**: Query your notes using natural language to get intelligent summaries
3. **Delete Mode**: Browse and manage your existing notes with easy deletion

### 🔒 Privacy & Performance
- **Local Storage**: All notes stored locally in SQLite database
- **Fast Retrieval**: Vector similarity search for lightning-fast note matching
- **Cross-Platform**: Runs on Linux, macOS, and Windows

## 🚀 Installation

### Limitations
- I'm using **go-libsql** (i'll switch to something else in future) which uses CGO. You need a C compiler to build it.
- This package comes with a precompiled native libraries. Currently only linux amd64, linux arm64, darwin amd64 and darwin arm64 are supported. For windows you manually need to compile the native libraries yourself. https://github.com/tursodatabase/libsql/releases

### Prerequisites
- Go 1.24.4 or later
- Google AI API key (for Gemini embeddings)
- C compiler (go-libsql uses CGO)

### Option 1: Build from Source
```bash
# Clone the repository
git clone https://github.com/yagnikpt/flashback.git
cd flashback

# Build the application
make build

# Run the application
./flashback
```

### Option 2: Direct Go Install
```bash
# Install directly from source
go install github.com/yagnikpt/flashback/cmd/flashback@latest

# Run the application
flashback
```

### Option 3: Development Setup
```bash
# Clone and run in development mode
git clone https://github.com/yagnikpt/flashback.git
cd flashback

# Install dependencies
make tidy

# Run without building
make run
```

### Notification Daemon Setup
This will create a user level systemd service or a user level launchd service for the notification daemon.
```bash
# build binary
make build-daemon

# setup daemon
./setup-daemon.sh
```

## ⚙️ Configuration

### Setting up Google AI API
1. Get your API key from [Google AI Studio](https://aistudio.google.com/apikey)
2. Set the API key in the initial screen of the application.

## 🎮 Usage

### Basic Navigation
- **Alt+?**: Toggle keys helper visibility

### Example Workflows

**Capturing Information:**
```
Mode: note
> "Met with John about Q3 project timeline. Key deliverables due by July 15th"
```

**Recalling Information:**
```
Mode: recall
> "What did John say about Q3?"
Response: "John discussed Q3 project timeline with key deliverables due by July 15th"
```

## 🛠️ Development

### Project Structure
```
flashback/
├── cmd/flashback/         # Main application entry point
├── internal/
│   ├── app/               # Core application logic and TUI models
│   ├── components/        # Reusable UI components
│   ├── migration/         # Database migrations
│   ├── notes/             # Notes business logic and AI integration
│   └── config/            # Config save and load helpers
├── Makefile               # Build and development commands
└── README.md
```

### Key Technologies
- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)**: Modern TUI framework
- **[LibSQL](https://github.com/tursodatabase/go-libsql)**: SQLite-compatible database
- **[Google Generative AI](https://pkg.go.dev/google.golang.org/genai)**: Embeddings and text generation
- **[Goose](https://github.com/pressly/goose)**: Database migrations

### Available Make Commands
```bash
make build    # Build the binary
make run      # Run in development mode
make tidy     # Clean up dependencies
```

## 🔮 Future Features

### 🌐 Web Content Integration
- **Content Scraping**: Extract information directly from web pages to load context into notes
- **URL Processing**: Automatically detect URLs and generate summaries

### 📁 Attach Local Files
- **File Parsing**: Import and extract information from local files (PDF, TXT, DOCX, etc.)
- **Directory Indexing**: Recursively scan directories to build a knowledge base from your files

### ⏰ Smart Notifications
- **Time & Date Extraction**: Automatically identify dates and times mentioned in your notes
- **Notification System**: Run as a daemon to alert you about upcoming events extracted from notes
- **Custom Reminders**: Set specific notification preferences for different types of information

### Development Environment

#### Required Tools
- Go 1.24.4+
- Git
- C compiler
- Your favorite terminal
- Text editor with Go support
