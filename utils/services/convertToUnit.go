package services

import "fmt"

func ConvertToAppropriateUnit(bytes int64) string {
	if bytes < 1024*1024 {
		// If size is less than 1 MB, convert to KB
		sizeInKB := float64(bytes) / 1024
		return fmt.Sprintf("%.2f KB", sizeInKB)
	} else {
		// If size is 1 MB or more, convert to MB
		sizeInMB := float64(bytes) / (1024 * 1024)
		return fmt.Sprintf("%.2f MB", sizeInMB)
	}
}
