# YouTube Audio Extractor

A simple cross-platform GUI application for downloading audio from YouTube videos.

## Features

- 🎵 Download audio directly from YouTube videos
- 🖥️ Clean, intuitive GUI interface
- 📁 Choose custom output folder (defaults to ~/Downloads)
- 🏷️ Smart file naming: "Channel Name - Video Title.m4a"
- ✅ URL validation
- 🧹 Automatic cleanup of failed downloads
- 🌍 Cross-platform (Linux and macOS)

## Screenshots

[Coming soon]

## Installation

### From Releases

Download the latest release for your platform from the [Releases](https://github.com/hoblin/youtube-audio-extractor/releases) page.

### From Source

**Prerequisites:**
- Go 1.21 or later (or use [mise](https://mise.jdx.dev/) for tool management)

**Build:**
```bash
git clone https://github.com/hoblin/youtube-audio-extractor.git
cd youtube-audio-extractor

# If using mise
mise trust
mise install

# Build the application
go build -o youtube-audio-extractor
```

## Usage

1. Launch the application: `./youtube-audio-extractor`
2. Paste a YouTube video URL into the input field
3. (Optional) Click "Choose Folder..." to select a custom output directory
4. Click "⬇ Download Audio" button
5. Wait for the download to complete

Audio files are saved as high-quality m4a/webm format with the naming pattern:
```
Channel Name - Video Title.m4a
```

## Known Issues

- Occasional 403 errors due to YouTube's anti-bot measures (retry usually works)
- Some videos may not be available depending on regional restrictions

## Development

### Project Structure

```
youtube-audio-extractor/
├── main.go           # Main application code
├── go.mod            # Go dependencies
├── .mise.toml        # Tool version management
├── .gitignore        # Git ignore rules
└── README.md         # This file
```

### Tech Stack

- **Language:** Go
- **GUI Framework:** [Fyne](https://fyne.io/)
- **YouTube Library:** [kkdai/youtube](https://github.com/kkdai/youtube)

### Running in Development

```bash
go run main.go
```

## License

MIT
