package utils

import (
	"fmt"
	"net/http"
)

// HandleHTTPResponse provides consistent HTTP status code handling across the CLI
// Returns true if the request was successful (200 or 201), false otherwise
func HandleHTTPResponse(resp *http.Response, context string) bool {
	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated:
		return true
	case http.StatusUnauthorized:
		fmt.Printf("❌ Authentication failed for %s\n", context)
		fmt.Println("💡 Your token may have expired. Run 'cool auth' to reconfigure")
		return false
	case http.StatusNotFound:
		fmt.Printf("❌ Resource not found for %s\n", context)
		fmt.Println("💡 The resource may have been deleted or moved")
		return false
	case http.StatusInternalServerError:
		fmt.Printf("❌ Server error occurred for %s\n", context)
		fmt.Println("💡 There may be an issue with the Coolify server")
		return false
	default:
		fmt.Printf("❌ %s failed with status %d: %s\n", context, resp.StatusCode, resp.Status)
		return false
	}
}
