package Models

import (
	"strings"
)

type FileTreeItem struct {
	Size     int
	Name     string
	Icon     string         `json:"icon"`
	Children []FileTreeItem `json:"Children,omitempty"`
}

func (fti *FileTreeItem) AddPath(path string, size int, icon string, folderIcon string) bool {
	pathParts := strings.Split(path, "/")
	pathPartsLen := len(pathParts)
	fti.Icon = folderIcon
	if fti.Name != pathParts[0] {
		return false
	}
	if len(fti.Children) == 0 {
		var currentItem = fti
		for i := 1; i < pathPartsLen; i++ {
			item := new(FileTreeItem)
			item.Name = pathParts[i]
			item.Size = size

			currentItem.Children = append(currentItem.Children, *item)
			currentItem = item
			if i == pathPartsLen-1 {
				currentItem.Icon = icon
			} else {
				currentItem.Icon = folderIcon
			}
		}
	} else {
		currentItem := fti
		for i := 1; i < pathPartsLen; i++ {
			childIndex := FindChild(&currentItem.Children, pathParts[i])
			if childIndex > -1 {
				tmpItem := &currentItem.Children[childIndex]
				tmpItem.Size += size
				currentItem = tmpItem
				continue
			}
			item := new(FileTreeItem)
			item.Name = pathParts[i]
			item.Size = size

			currentItem.Children = append(currentItem.Children, *item)

			currentItem = item
			if i == pathPartsLen-1 {
				currentItem.Icon = icon
			} else {
				currentItem.Icon = folderIcon
			}
		}
	}
	return true
}

func initFileTreeItem(path string, size int) FileTreeItem {
	pathParts := strings.Split(path, "/")
	var result = FileTreeItem{
		Name: pathParts[0],
		Size: size,
	}
	if len(pathParts) > 1 {
		result.AddPath(path, size, "insert_drive_file", "folder")
	} else {

		result.Icon = "insert_drive_file"

	}

	return result
}

func FindChild(items *[]FileTreeItem, name string) int {
	for index := 0; index < len(*items); index++ {
		if (*items)[index].Name == name {
			return index
		}

	}
	return -1
}

func BuildTree(fileList []File) []FileTreeItem {
	var result []FileTreeItem

	fileListLen := len(fileList)
	for i := 0; i < fileListLen; i++ {
		file := fileList[i]
		pathParts := strings.Split(file.Path, "/")

		childIndex := FindChild(&result, pathParts[0])
		if childIndex == -1 {
			result = append(result, initFileTreeItem(file.Path, file.Size))
		} else {
			result[childIndex].Size += file.Size
			result[childIndex].AddPath(file.Path, file.Size, "insert_drive_file", "folder")
		}
	}

	return result
}
