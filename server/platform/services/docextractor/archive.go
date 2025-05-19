// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package docextractor

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/mholt/archives"
)

type archiveExtractor struct {
	SubExtractor Extractor
}

func (ae *archiveExtractor) Name() string {
	return "archiveExtractor"
}

func (ae *archiveExtractor) Match(filename string) bool {
	format, _, err := archives.Identify(context.Background(), filename, nil)
	return err == nil && format != nil
}

func (ae *archiveExtractor) Extract(name string, r io.ReadSeeker) (string, error) {
	dir, err := os.MkdirTemp(os.TempDir(), "archiver")
	if err != nil {
		return "", fmt.Errorf("error creating temporary file: %v", err)
	}
	defer os.RemoveAll(dir)

	f, err := os.Create(filepath.Join(dir, name))
	if err != nil {
		return "", fmt.Errorf("error copying data into temporary file: %v", err)
	}
	_, err = io.Copy(f, r)
	f.Close()
	if err != nil {
		return "", fmt.Errorf("error copying data into temporary file: %v", err)
	}

	var text strings.Builder
	fsys, err := archives.FileSystem(context.Background(), f.Name(), nil)
	if err != nil {
		return "", fmt.Errorf("error creating file system: %v", err)
	}

	err = fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		text.WriteString(path + " ")
		if ae.SubExtractor != nil {
			filename := filepath.Base(path)
			filename = strings.ReplaceAll(filename, "-", " ")
			filename = strings.ReplaceAll(filename, ".", " ")
			filename = strings.ReplaceAll(filename, ",", " ")

			file, err := fsys.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			data, err := io.ReadAll(file)
			if err != nil {
				return err
			}

			subtext, extractErr := ae.SubExtractor.Extract(filename, bytes.NewReader(data))
			if extractErr == nil {
				text.WriteString(subtext + " ")
			}
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	return text.String(), nil
}
