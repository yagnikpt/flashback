# ⚡ Flashback

A powerful command-line tool that serves as your second memory, intelligently storing and retrieving your notes using AI-powered semantic search. Never lose track of important information again!

## ✨ Features

### 🧠 AI-Powered Note Management
- **Smart Storage**: Automatically generates embeddings for your notes using Google's Gemini AI
- **Semantic Search**: Find notes using natural language queries, not just keywords
- **Intelligent Recall**: Get contextually relevant responses based on your stored notes
- **Timestamp-Aware**: Automatically tracks when notes were created with human-friendly time formatting

### 🎯 Three Powerful Modes
1. **Note Mode**: Quickly capture thoughts, ideas, and information
2. **Recall Mode**: Query your notes using natural language to get intelligent summaries
3. **Delete Mode**: Browse and manage your existing notes with easy deletion

### 🖥️ Beautiful Terminal Interface
- Modern TUI (Terminal User Interface) built with Bubble Tea
- Responsive design that adapts to your terminal size
- Visual mode indicators and loading states
- Keyboard-driven navigation for maximum efficiency

### 🔒 Privacy & Performance
- **Local Storage**: All notes stored locally in SQLite database
- **Fast Retrieval**: Vector similarity search for lightning-fast note matching
- **Offline First**: Works without internet connection (after initial AI setup)
- **Cross-Platform**: Runs on Linux, macOS, and Windows

## 🚀 Installation

### Prerequisites
- Go 1.24.4 or later
- Google AI API key (for Gemini embeddings)

### Option 1: Build from Source
```bash
# Clone the repository
git clone https://github.com/yagnik-patel-47/flashback.git
cd flashback

# Build the application
make build

# Run the application
./flashback
```

### Option 2: Direct Go Install
```bash
# Install directly from source
go install github.com/yagnik-patel-47/flashback/cmd/flashback@latest

# Run the application
flashback
```

### Option 3: Development Setup
```bash
# Clone and run in development mode
git clone https://github.com/yagnik-patel-47/flashback.git
cd flashback

# Install dependencies
go mod tidy

# Run without building
make run
```

## ⚙️ Configuration

### Setting up Google AI API
1. Get your API key from [Google AI Studio](https://aistudio.google.com/apikey)
2. Set the environment variable:
   - Linux, Mac
   ```bash
   export GOOGLE_API_KEY="your-api-key-here"
   ```
   - Powershell
   ```pwsh
   $Env:GOOGLE_API_KEY = "your-api-key-here"
   ```
3. Set permannent environment variable:
   - Linux, Mac
   edit `.bashrc` or `.zshrc` or your terminal config and add the above line
   - Powershell
   run `notepad $PROFILE` and add the above line


## 🎮 Usage

### Basic Navigation
- **Tab**: Switch between modes (Note → Recall → Delete → Note)
- **Enter**: Submit input or select item
- **Ctrl+C**: Exit the application
- **↑/↓ or j/k**: Navigate through notes in delete mode

### Note Mode
1. Type your note in the textarea
2. Press **Enter** to save
3. Get confirmation feedback

### Recall Mode
1. Enter a natural language query
2. Press **Enter** to search
3. Get AI-generated summary of relevant notes

### Delete Mode
1. Browse through all your notes
2. Use arrow keys or j/k to navigate
3. Press **Enter** to delete selected note

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
├── cmd/flashback/          # Main application entry point
├── internal/
│   ├── app/               # Core application logic and TUI models
│   ├── components/        # Reusable UI components
│   │   ├── notelist/      # Notes list component
│   │   ├── spinner/       # Loading spinner component
│   │   └── textarea/      # Text input component
│   ├── migration/         # Database migrations
│   └── notes/            # Notes business logic and AI integration
├── Makefile              # Build and development commands
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

## 🤝 Contributing

We welcome contributions from the community! Here's how you can help:

### Getting Started
1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/yagnik-patel-47/flashback.git
   cd flashback
   ```
3. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

### Development Guidelines

#### Code Style
- Follow standard Go conventions (`go fmt`, `go vet`)
- Use meaningful variable and function names
- Add comments for exported functions and complex logic
- Keep functions small and focused

#### Testing
- Write unit tests for new functionality
- Ensure existing tests pass: `go test ./...`
- Test the TUI manually across different terminal sizes

#### Commit Messages
Use conventional commit format:
```
type(scope): description

feat(notes): add export functionality
fix(ui): resolve textarea resize issue
docs(readme): update installation instructions
```

### Types of Contributions

#### 🐛 Bug Reports
- Use the issue template
- Include steps to reproduce
- Provide system information (OS, terminal, Go version)
- Include logs from `debug.log` if applicable

#### ✨ Feature Requests
- Describe the problem you're solving
- Provide clear use cases
- Consider backward compatibility
- Discuss implementation approach

#### 🔧 Code Contributions
- **UI/UX Improvements**: Better styling, new components, accessibility
- **Performance**: Query optimization, faster embeddings, memory usage
- **Features**: Export/import, note categories, search filters
- **Integrations**: Other AI providers, cloud storage, sync

#### 📚 Documentation
- README improvements
- Code comments and documentation
- Usage examples and tutorials
- API documentation

### Pull Request Process

1. **Update documentation** if needed
2. **Add or update tests** for your changes
3. **Ensure all tests pass**
4. **Update CHANGELOG.md** with your changes
5. **Submit pull request** with clear description

#### PR Template
```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Documentation update
- [ ] Performance improvement

## Testing
- [ ] Unit tests added/updated
- [ ] Manual testing completed
- [ ] No regressions introduced

## Screenshots (if applicable)
Add screenshots for UI changes
```

### Development Environment

#### Required Tools
- Go 1.24.4+
- Git
- Your favorite terminal
- Text editor with Go support
