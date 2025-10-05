package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kkdai/youtube/v2"
)

// downloadVideo downloads a YouTube video to the specified output path
func downloadVideo(videoURL, outputPath string) error {
	client := youtube.Client{}

	video, err := client.GetVideo(videoURL)
	if err != nil {
		return fmt.Errorf("failed to get video info: %w", err)
	}

	// Get the highest quality format with both video and audio
	formats := video.Formats.WithAudioChannels()
	if len(formats) == 0 {
		return fmt.Errorf("no formats with audio found")
	}

	// Select the first format (usually highest quality)
	format := formats[0]

	stream, _, err := client.GetStream(video, &format)
	if err != nil {
		return fmt.Errorf("failed to get stream: %w", err)
	}
	defer stream.Close()

	// Create output file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Copy stream to file
	_, err = io.Copy(file, stream)
	if err != nil {
		return fmt.Errorf("failed to download video: %w", err)
	}

	return nil
}

func main() {
	fmt.Println("YouTube Audio Extractor")

	// Example usage (will be replaced with GUI)
	if len(os.Args) < 2 {
		fmt.Println("Usage: youtube-audio-extractor <video-url>")
		os.Exit(1)
	}

	videoURL := os.Args[1]
	outputPath := filepath.Join(".", "video.mp4")

	fmt.Printf("Downloading video from: %s\n", videoURL)
	err := downloadVideo(videoURL, outputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Video downloaded successfully to: %s\n", outputPath)
}
