// This is data service for the go-index API and static dist. The raw data is based on the https://index.golang.org/ api.
package idx

// @FetchIdx - Given a properly configured postgresql instance it queries the https://index.golang.org api
// using the since param to ensure all packages are found and stored with their version information as well.
func FetchIdx() {}
