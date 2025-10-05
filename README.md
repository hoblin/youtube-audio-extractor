# YouTube Audio Extractor

A self-contained cross-platform GUI application for downloading YouTube videos and extracting audio.

## Features

- Download YouTube videos
- Extract audio to separate file (MP3 format)
- Simple GUI interface
- Cross-platform (Linux and macOS)
- Zero external dependencies (ffmpeg embedded)
- Single standalone executable

## Requirements

None! All dependencies are bundled into the executable.

## Installation

Download the latest release for your platform from the [Releases](https://github.com/hoblin/youtube-audio-extractor/releases) page.

## Development

### Prerequisites

- Go 1.21 or later
- mise (for tool management)

### Setup

```bash
mise trust
mise install
go mod download
```

### Build

```bash
go build -o youtube-audio-extractor
```

### Run

```bash
./youtube-audio-extractor
```

## Usage

1. Launch the application
2. Paste a YouTube video URL into the input field
3. Click the Download button
4. The app will download the video and extract the audio to the current directory

## License

MIT
