package database

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func (dbCfg *DbConfig) saveImage(item ToplistItem, listID int) error {
	// Get the current directory path
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// Define the path where images will be stored
	imagesDir := filepath.Join(currentDir, "internal", "database", "images")
	toplistDir := filepath.Join(imagesDir, fmt.Sprintf("%d", listID))

	// Create the toplist directory if it doesn't exist
	err = os.MkdirAll(toplistDir, os.ModePerm)
	if err != nil {
		return err
	}

	// Define the image file path
	imagePath := filepath.Join(toplistDir, fmt.Sprintf("%d.png", item.Rank))

	// Check if item.Image is not present
	if len(item.Image) == 0 {
		// Check if the image file already exists
		if _, err := os.Stat(imagePath); err == nil {
			// Image file exists, delete it
			err := os.Remove(imagePath)
			if err != nil {
				return err
			}
		}
		return nil
	}

	// Create the image file
	file, err := os.Create(imagePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the image data to the file
	_, err = io.Copy(file, bytes.NewReader(item.Image))
	if err != nil {
		return err
	}

	return nil
}
