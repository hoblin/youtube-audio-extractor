package main

import (
	"context"
	"fmt"
	"path/filepath"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
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
	a := app.New()
	w := a.NewWindow("YouTube Audio Extractor")
	w.Resize(fyne.NewSize(600, 200))

	// URL input
	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("Enter YouTube URL here...")

	// Status label
	statusLabel := widget.NewLabel("Ready to download")
	statusLabel.Wrapping = fyne.TextWrapWord

	// Download button - declare first so it can be referenced in callback
	var downloadBtn *widget.Button
	downloadBtn = widget.NewButton("Download Audio", func() {
		url := urlEntry.Text
		if url == "" {
			statusLabel.SetText("Error: Please enter a YouTube URL")
			return
		}

		// Disable button during download
		downloadBtn.Disable()
		statusLabel.SetText("Downloading...")

		// Download in goroutine to keep UI responsive
		go func() {
			outputPath := filepath.Join(".", "audio.m4a")
			err := downloadAudio(url, outputPath)

			// Update UI in main thread using fyne.Do
			if err != nil {
				fyne.Do(func() {
					statusLabel.SetText(fmt.Sprintf("Error: %v", err))
					downloadBtn.Enable()
				})
			} else {
				fyne.Do(func() {
					statusLabel.SetText(fmt.Sprintf("Success! Downloaded to: %s", outputPath))
					downloadBtn.Enable()
				})
			}
		}()
	})

	// Layout
	content := container.NewVBox(
		widget.NewLabel("YouTube Audio Extractor"),
		widget.NewSeparator(),
		urlEntry,
		downloadBtn,
		widget.NewSeparator(),
		statusLabel,
	)

	w.SetContent(content)
	w.ShowAndRun()
}
