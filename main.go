package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kkdai/youtube/v2"
	"github.com/kkdai/youtube/v2/downloader"
)

// getDefaultDownloadDir returns ~/Downloads if it exists, otherwise current directory
func getDefaultDownloadDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "."
	}

	downloadsDir := filepath.Join(homeDir, "Downloads")
	if _, err := os.Stat(downloadsDir); err == nil {
		return downloadsDir
	}

	return "."
}

// isYouTubeURL validates if the URL is a YouTube URL
func isYouTubeURL(url string) bool {
	youtubeRegex := regexp.MustCompile(`^(https?://)?(www\.)?(youtube\.com|youtu\.be)/.+$`)
	return youtubeRegex.MatchString(url)
}

// cleanYouTubeURL removes tracking parameters from YouTube URLs
func cleanYouTubeURL(url string) string {
	// Remove all query parameters except v= for /watch URLs
	if strings.Contains(url, "?") {
		parts := strings.Split(url, "?")
		baseURL := parts[0]

		// If it's a /watch URL, preserve only the v= parameter
		if strings.Contains(url, "/watch") && len(parts) > 1 {
			params := strings.Split(parts[1], "&")
			for _, param := range params {
				if strings.HasPrefix(param, "v=") {
					return baseURL + "?" + param
				}
			}
		}

		// For youtu.be links, just return the base URL
		return baseURL
	}
	return url
}

// sanitizeFilename creates a safe filename from a string
func sanitizeFilename(name string) string {
	// Remove invalid filename characters
	reg := regexp.MustCompile(`[<>:"/\\|?*]`)
	safe := reg.ReplaceAllString(name, "_")
	// Limit length
	if len(safe) > 200 {
		safe = safe[:200]
	}
	return strings.TrimSpace(safe)
}


// downloadAudio downloads audio from a YouTube video and returns the output path
func downloadAudio(videoURL, outputDir string) (string, error) {
	dl := downloader.Downloader{
		OutputDir: outputDir,
	}

	video, err := dl.GetVideo(videoURL)
	if err != nil {
		return "", err
	}

	fmt.Printf("Title: %s\n", video.Title)
	fmt.Printf("Duration: %s\n", video.Duration)

	// Get audio-only formats
	formats := video.Formats.Type("audio")
	if len(formats) == 0 {
		return "", fmt.Errorf("no audio formats found")
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

	// Create safe filename from channel and video title
	safeAuthor := sanitizeFilename(video.Author)
	safeTitle := sanitizeFilename(video.Title)
	filename := fmt.Sprintf("%s - %s.m4a", safeAuthor, safeTitle)

	err = dl.Download(context.Background(), video, selectedFormat, filename)
	outputPath := filepath.Join(outputDir, filename)

	if err != nil {
		// Clean up failed download file if it exists
		os.Remove(outputPath)
		return "", err
	}

	// Verify the file was actually downloaded and has content
	fileInfo, err := os.Stat(outputPath)
	if err != nil || fileInfo.Size() == 0 {
		os.Remove(outputPath)
		return "", fmt.Errorf("download failed: file is empty or missing")
	}

	fmt.Println("Download complete!")
	return outputPath, nil
}

func main() {
	a := app.New()
	w := a.NewWindow("YouTube Audio Extractor")
	w.Resize(fyne.NewSize(650, 280))

	// Output directory (default to ~/Downloads if it exists, otherwise current directory)
	outputDir := getDefaultDownloadDir()

	// URL input - make it multiline for better visibility
	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("Paste YouTube URL here...")

	// Output directory section
	dirLabel := widget.NewLabel("📁 " + outputDir)
	dirLabel.Wrapping = fyne.TextWrapWord

	// Choose directory button
	chooseDirBtn := widget.NewButton("Choose Folder...", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if uri != nil {
				outputDir = uri.Path()
				dirLabel.SetText("📁 " + outputDir)
			}
		}, w)
	})
	chooseDirBtn.Importance = widget.LowImportance

	// Status label
	statusLabel := widget.NewLabel("Ready to download")
	statusLabel.Wrapping = fyne.TextWrapWord

	// Download button - declare first so it can be referenced in callback
	var downloadBtn *widget.Button
	downloadBtn = widget.NewButton("⬇ Download Audio", func() {
		url := urlEntry.Text
		if url == "" {
			statusLabel.SetText("Error: Please enter a YouTube URL")
			return
		}

		// Validate URL
		if !isYouTubeURL(url) {
			statusLabel.SetText("Error: Please enter a valid YouTube URL")
			return
		}

		// Disable button during download
		downloadBtn.Disable()
		chooseDirBtn.Disable()
		statusLabel.SetText("Downloading...")

		// Download in goroutine to keep UI responsive
		go func() {
			outputPath, err := downloadAudio(url, outputDir)

			// If 403 error, retry with cleaned URL
			if err != nil && strings.Contains(err.Error(), "403") {
				cleanedURL := cleanYouTubeURL(url)
				if cleanedURL != url {
					fyne.Do(func() {
						statusLabel.SetText("Got 403 error, retrying with sanitized URL...")
					})
					outputPath, err = downloadAudio(cleanedURL, outputDir)
				}
			}

			// Update UI in main thread using fyne.Do
			if err != nil {
				fyne.Do(func() {
					statusLabel.SetText(err.Error())
					downloadBtn.Enable()
					chooseDirBtn.Enable()
				})
			} else {
				fyne.Do(func() {
					statusLabel.SetText(fmt.Sprintf("✓ Downloaded successfully to:\n%s", outputPath))
					downloadBtn.Enable()
					chooseDirBtn.Enable()
				})
			}
		}()
	})
	downloadBtn.Importance = widget.HighImportance

	// Layout with better spacing and hierarchy
	content := container.NewVBox(
		widget.NewLabel("YouTube Audio Extractor"),
		widget.NewSeparator(),
		urlEntry,
		downloadBtn,
		widget.NewSeparator(),
		container.NewHBox(widget.NewLabel("Save to:"), chooseDirBtn),
		dirLabel,
		widget.NewSeparator(),
		statusLabel,
	)

	w.SetContent(content)
	w.ShowAndRun()
}
