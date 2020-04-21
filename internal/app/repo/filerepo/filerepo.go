package filerepo

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"github.com/khanhpdt/bookmark-api/internal/app/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

// SaveUploadedFiles saves files to disk and database.
func SaveUploadedFiles(fs []*multipart.FileHeader) []error {
	errs := make([]error, 0, len(fs))

	for _, f := range fs {
		fn := strings.ToLower(strings.ReplaceAll(f.Filename, " ", "_"))
		filePath := fmt.Sprintf("/tmp/%s", fn)

		if err := saveFileToDisk(f, filePath); err != nil {
			log.Printf("Error saving file %s to disk.", f.Filename)
			errs = append(errs, err)
			continue
		}

		if err := saveFileDocument(f.Filename, filePath); err != nil {
			log.Printf("Error saving file %s to database.", f.Filename)
			errs = append(errs, err)
			continue
		}

		log.Printf("Saved file %s to %s.", f.Filename, filePath)
	}

	return errs
}

func saveFileToDisk(f *multipart.FileHeader, filePath string) error {
	src, err := f.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		return err
	}

	return nil
}

func saveFileDocument(fileName, filePath string) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	_, err := mongo.FileColl().InsertOne(ctx, bson.M{"name": fileName, "path": filePath})
	return err
}
