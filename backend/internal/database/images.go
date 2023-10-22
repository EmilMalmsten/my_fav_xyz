package database

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func getToplistDir(listID int) (string, error) {

	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	exeDir := filepath.Dir(exePath)

	imagesDir := filepath.Join(exeDir, "images")
	toplistDir := filepath.Join(imagesDir, fmt.Sprintf("%d", listID))
	return toplistDir, nil
}

func saveImage(item ToplistItem, listID int) error {

	toplistDir, err := getToplistDir(listID)
	if err != nil {
		return err
	}

	err = os.MkdirAll(toplistDir, os.ModePerm)
	if err != nil {
		return err
	}

	mimeType := http.DetectContentType(item.Image)
	ext := getExtensionFromMimeType(mimeType)
	if ext == "" {
		return errors.New("unsupported img type")
	}

	// Toplist item images are saved in a folder named after the toplist ID
	// The image file name is the rank of the item it belongs to.
	imagePath := filepath.Join(toplistDir, fmt.Sprintf("%d%s", item.Rank, ext))

	file, err := os.Create(filepath.Clean(imagePath))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, bytes.NewReader(item.Image))
	if err != nil {
		return err
	}

	return nil
}

func getExtensionFromMimeType(mimeType string) string {
	switch mimeType {
	case "image/jpeg":
		return ".jpg"
	default:
		return ""
	}
}

// This function is used for responding to a request with toplist items
// The filepath for each toplist item image gets put together and added to the item
func setImagePaths(items []ToplistItem, listID int) ([]ToplistItem, error) {
	imageDirPath, err := getToplistDir(listID)
	if err != nil {
		return []ToplistItem{}, err
	}
	if _, err := os.Stat(imageDirPath); err == nil {
		imgFiles, err := os.ReadDir(imageDirPath)
		if err != nil {
			return []ToplistItem{}, err
		}

		imageMap := make(map[int]string) // maps rank to image path

		for _, file := range imgFiles {
			fileBase := filepath.Base(file.Name())
			fileWithoutExt := strings.TrimSuffix(fileBase, filepath.Ext(fileBase))
			rank, err := strconv.Atoi(fileWithoutExt)
			if err == nil {
				imageMap[rank] = file.Name()
			}
		}

		// Set each items image path if it exists
		for i, item := range items {
			if path, exists := imageMap[item.Rank]; exists {
				items[i].ImagePath = path
			}
		}
	} else {
		if !os.IsNotExist(err) {
			return []ToplistItem{}, err
		}
	}

	return items, nil
}

func deleteToplistImages(listID int) error {
	toplistDir, err := getToplistDir(listID)
	if err != nil {
		return err
	}

	err = os.RemoveAll(toplistDir)
	if err != nil {
		return err
	}

	return nil
}

// Swap the image paths of item A and item B
func swapImageFiles(item, swapItem *ToplistItem, listID int) error {

	toplistDir, err := getToplistDir(listID)
	if err != nil {
		return err
	}

	swapItemFullPath := filepath.Join(toplistDir, swapItem.ImagePath)
	currItemFullPath := filepath.Join(toplistDir, item.ImagePath)
	tempPath := filepath.Join(toplistDir, "temp.png")

	// Set item B image path to a temporary path
	err = os.Rename(swapItemFullPath, tempPath)
	if err != nil {
		return err
	}

	// Set the path of Item A to its target path (Item B's initial path)
	err = os.Rename(currItemFullPath, swapItemFullPath)
	if err != nil {
		return err
	}

	// Set the path of Item B to its target path (Item A's initial path)
	err = os.Rename(tempPath, currItemFullPath)
	if err != nil {
		return err
	}

	// Swap the image paths on the item structs
	tmp := item.ImagePath
	item.ImagePath = swapItem.ImagePath
	swapItem.ImagePath = tmp

	return nil
}

func findSwapItem(items []ToplistItem, currentItem *ToplistItem) *ToplistItem {
	for i := range items {
		fileRank := getRankFromFileName(items[i].ImagePath)

		if currentItem.Rank == fileRank {
			// we found the item to swap with
			return &items[i]
		}
	}

	return nil
}

// Move image file from oldPath to newPath
func moveImage(currentFileName string, itemRank, listID int) error {
	toplistDir, err := getToplistDir(listID)
	if err != nil {
		return err
	}

	dotIndex := strings.LastIndex(currentFileName, ".")
	itemRankString := strconv.Itoa(itemRank)
	newFileName := itemRankString + currentFileName[dotIndex:]

	oldPath := filepath.Join(toplistDir, currentFileName)
	newPath := filepath.Join(toplistDir, newFileName)

	err = os.Rename(oldPath, newPath)
	if err != nil {
		return err
	}

	return nil
}

// When user updates the toplist this function checks for needed image changes
func handleImageChanges(items []ToplistItem, listID int) ([]ToplistItem, error) {
	for i := range items {

		if items[i].ImagePath == "" && len(items[i].Image) > 0 {
			// item did not have an image previously but now got one added
			err := saveImage(items[i], listID)
			if err != nil {
				fmt.Println(err)
				return []ToplistItem{}, err
			}
		} else if items[i].ImagePath != "" && len(items[i].Image) == 0 {
			// item did have an item previosly, and no new one got uploaded
			// check if the item rank matches with its img file path
			// if it doesnt match, it means the user swapped item ranks
			// and filepaths needs to be updated

			fileRank := getRankFromFileName(items[i].ImagePath)
			if fileRank == items[i].Rank {
				// no image path changes needed
				continue
			}

			item := &items[i]
			swapItem := findSwapItem(items, item)
			if swapItem == nil {
				// no swap items means that item got removed, so just need to change
				// the old image path to match the new rank
				err := moveImage(items[i].ImagePath, items[i].Rank, listID)
				if err != nil {
					return []ToplistItem{}, err
				}
				continue
			}

			err := swapImageFiles(item, swapItem, listID)
			if err != nil {
				return []ToplistItem{}, err
			}

		} else if items[i].ImagePath == "" && len(items[i].Image) == 0 {
			// Item should no longer have an image. Check if one exists and if so, delete it
			toplistDir, err := getToplistDir(listID)
			if err != nil {
				return []ToplistItem{}, err
			}

			imageFiles, err := os.ReadDir(toplistDir)
			if err != nil {
				continue
			}

			for _, imageFile := range imageFiles {
				rank := getRankFromFileName(imageFile.Name())
				if rank == items[i].Rank {
					filePath := filepath.Join(toplistDir, imageFile.Name())
					err := os.Remove(filePath)
					if err != nil {
						fmt.Println("Error deleting file:", err)
					}
				}
			}
		}
	}

	return items, nil
}

func getRankFromFileName(path string) int {
	splitPath := strings.Split(path, ".")
	if len(splitPath) < 2 {
		return -1
	}
	fileName := splitPath[0]
	fileRank, err := strconv.Atoi(fileName)
	if err != nil {
		return -1
	}
	return fileRank
}
