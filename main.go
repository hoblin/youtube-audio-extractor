package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kkdai/youtube/v2"
	"github.com/kkdai/youtube/v2/downloader"
)

// downloadAudio downloads audio from a YouTube video
func downloadAudio(videoURL, outputPath string) error {
	dl := downloader.Downloader{
		OutputDir: filepath.Dir(outputPath),
	}

	video, err := dl.GetVideo(videoURL)
	if err != nil {
		return fmt.Errorf("failed to get video info: %w", err)
	}

	fmt.Printf("Title: %s\n", video.Title)
	fmt.Printf("Duration: %s\n", video.Duration)

	// Get audio-only formats
	formats := video.Formats.Type("audio")
	if len(formats) == 0 {
		return fmt.Errorf("no audio formats found")
	}

	fmt.Printf("Found %d audio formats\n", len(formats))

	// Select the best audio format (highest bitrate)
	var selectedFormat *youtube.Format
	var maxBitrate int
	for i := range formats {
		format := &formats[i]
		if format.Bitrate > maxBitrate {
			maxBitrate = format.Bitrate
			selectedFormat = format
		}
	}

	if selectedFormat == nil {
		selectedFormat = &formats[0]
	}

	fmt.Printf("Downloading audio: %s (Bitrate: %d)\n", selectedFormat.MimeType, selectedFormat.Bitrate)

	err = dl.Download(context.Background(), video, selectedFormat, filepath.Base(outputPath))
	if err != nil {
		return fmt.Errorf("failed to download audio: %w", err)
	}

	fmt.Println("Download complete!")
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
	outputPath := filepath.Join(".", "audio.m4a")

	fmt.Printf("Downloading audio from: %s\n", videoURL)
	err := downloadAudio(videoURL, outputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Audio downloaded successfully to: %s\n", outputPath)
}
