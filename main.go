package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

func GetRepositoryPath(repositoryName string) (string, error) {
	repositoryPaths := map[string]string{
		"appstream":        "AppStream/$basearch/os/",
		"baseos":           "BaseOS/$basearch/os/",
		"crb":              "CRB/$basearch/os/",
		"devel":            "devel/$basearch/os/",
		"extras":           "extras/$basearch/os/",
		"ha":               "HighAvailability/$basearch/os/",
		"highavailability": "HighAvailability/$basearch/os/",
		"nfv":              "NFV/$basearch/os/",
		"plus":             "plus/$basearch/os/",
		"powertools":       "PowerTools/$basearch/os/",
		"resilientstorage": "ResilientStorage/$basearch/os/",
		"rt":               "RT/$basearch/os/",
		"sap":              "SAP/$basearch/os/",
		"saphana":          "SAPHANA/$basearch/os/",
		"synergy":          "synergy/$basearch/os/",
	}

	path, exists := repositoryPaths[repositoryName]
	if !exists {
		return "", fmt.Errorf("'%s' not found", repositoryName)
	}

	return path, nil
}

func GetVersion(version string) (string, error) {
	versions := map[string]string{
		"9": "9.4",
		"8": "8.10",
	}

	if strings.Contains(version, ".") {
		return version, nil
	} else {
		path, exists := versions[version]
		if !exists {
			return "", fmt.Errorf("'%s' not found", version)
		}

		return path, nil
	}
}

func AccessLog(clientIp string, proxyIP string, mirror string, userAgent string) {
	data := map[string]interface{}{
		"clientIp":  clientIp,
		"proxyIp":   proxyIP,
		"mirror":    mirror,
		"userAgent": userAgent,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(jsonData))
}

func Home(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://mirrors.almalinux.org", http.StatusTemporaryRedirect)
}

func MirrorList(w http.ResponseWriter, r *http.Request) {
	mirrorDC01, _ := os.LookupEnv("DC01_MIRROR")
	prefixDC01, _ := os.LookupEnv("DC01_PREFIX")
	mirrorDC02, _ := os.LookupEnv("DC02_MIRROR")
	prefixDC02, _ := os.LookupEnv("DC02_PREFIX")
	mirrorDefault, _ := os.LookupEnv("DEFAULT_MIRROR")

	urlPath := path.Clean(r.URL.Path)
	parts := strings.Split(urlPath, "/")
	clientIP := r.RemoteAddr
	userAgent := r.Header.Get("User-Agent")
	proxyIP := r.Header.Get("X-Forwarded-For")

	if len(parts) < 3 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	version := parts[2]
	repo := parts[3]

	repoPath, pathErr := GetRepositoryPath(repo)
	versionPath, versionErr := GetVersion(version)

	if pathErr != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Not Found")
		return
	} else if versionErr != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Not Found")
		return
	} else {
		if strings.HasPrefix(clientIP, prefixDC01) {
			mirror := fmt.Sprintf("%s/%s/%s", mirrorDC01, versionPath, repoPath)
			AccessLog(clientIP, proxyIP, mirror, userAgent)
			fmt.Fprintf(w, "%s\n", mirror)
		} else if strings.HasPrefix(clientIP, prefixDC02) {
			mirror := fmt.Sprintf("%s/%s/%s", mirrorDC02, versionPath, repoPath)
			AccessLog(clientIP, proxyIP, mirror, userAgent)
			fmt.Fprintf(w, "%s\n", mirror)
		} else {
			mirror := fmt.Sprintf("%s/%s/%s", mirrorDefault, versionPath, repoPath)
			AccessLog(clientIP, proxyIP, mirror, userAgent)
			fmt.Fprintf(w, "%s\n", mirror)
		}
	}
}

func main() {
	portNum := ":8080"

	log.SetFlags(0)

	http.HandleFunc("/", Home)
	http.HandleFunc("/mirrorlist/", MirrorList)

	err := http.ListenAndServe(portNum, nil)
	if err != nil {
		log.Fatal(err)
	}
}
