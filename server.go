// Reference: https://www.reddit.com/r/golang/comments/apf6l5/multiple_files_upload_using_gos_standard_library/
//
// Original author: https://www.reddit.com/user/teizz
//
// handles multiple files being uploaded
// reads them in blocks of 4K
// writes them to a temporary file in $TMPDIR
// calculates and logs the SHA256 sum
// proceeds to remove the temporary file through a defer statement
package main

import (
	"crypto/sha256"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var indexPage = `
<html>
<body>
	<form
		enctype="multipart/form-data"
		action="http://localhost:8080/upload"
		method="post"
	>
		<input type="file" name="files" multiple />
		<input type="submit" value="upload" />
	</form>
</body>
</html>
`

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(indexPage))
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	mr, err := r.MultipartReader()
	if err != nil {
		log.Printf("Hit error while opening multipart reader: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var fileSize int // bytes number of uploaded files

	// buffer to be used for reading bytes from files.
	chunk := make([]byte, 4096) // 4k size byte slice
	tempDir := os.TempDir()     // temp dir for chunk files
	tempFile, err := ioutil.TempFile(tempDir, "example-temp-file")
	if err != nil {
		log.Printf("Hit error while creating temp file: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// continue looping through all parts, *multipart.Reader.NextPart() will
	// return an End of File when all parts have been read.
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			// err is io.EOF, files upload completes.
			log.Printf("Hit last part of multipart upload")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Files upload complete"))
			break
		}
		if err != nil {
			// A normal error occurred
			log.Printf("Hit error while fetching next part: %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// at this point the filename and the mimetype is known
		log.Printf("Uploaded filename: %s\n", p.FileName())
		log.Printf("Uploaded mimetype: %s\n", p.Header)

		// continue reading the part stream of this loop until either done or err.
		for {
			n, err := p.Read(chunk)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("Hit error while reading chunk: %s", err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if _, err = tempFile.Write(chunk[:n]); err != nil {
				log.Printf("Hit error while writing chunk: %s", err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fileSize += n
			log.Printf("Uploaded filesize: %d bytes\n", fileSize)

			// Log file sum.
			sum(tempFile)
		}
	}
}

func sum(f *os.File) {
	if n, err := f.Seek(0, 0); err != nil || n != 0 {
		log.Printf("unable to seek to beginning of file: %q\n", f.Name())
		h := sha256.New()
		if _, err := io.Copy(h, f); err != nil {
			log.Printf("unable to hash %q: %s\n", f.Name(), err.Error())
			return
		}
		log.Printf("SHA256 sum of %q: %x\n", f.Name(), h.Sum(nil))
	}
	defer f.Truncate(0)
}

func main() {
	log.Println("Gopher files upload service started!")
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/upload", uploadHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
