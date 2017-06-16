package files

import (
	"io"
	"io/ioutil"
	"mime/multipart"
	"strings"
	"testing"
)

func TestOutput(t *testing.T) {
	text := "Some text! :)"
	fileset := []File{
		NewReaderFile("file.txt", "file.txt", ioutil.NopCloser(strings.NewReader(text)), nil),
		NewSliceFile("boop", "boop", []File{
			NewReaderFile("boop/a.txt", "boop/a.txt", ioutil.NopCloser(strings.NewReader("bleep")), nil),
			NewReaderFile("boop/b.txt", "boop/b.txt", ioutil.NopCloser(strings.NewReader("bloop")), nil),
		}),
		NewReaderFile("beep.txt", "beep.txt", ioutil.NopCloser(strings.NewReader("beep")), nil),
	}
	sf := NewSliceFile("", "", fileset)
	buf := make([]byte, 20)

	// testing output by reading it with the go stdlib "mime/multipart" Reader
	mfr := NewMultiFileReader(sf, true)
	mpReader := multipart.NewReader(mfr, mfr.Boundary())

	part, err := mpReader.NextPart()
	if part == nil || err != nil {
		t.Error("Expected non-nil part, nil error")
	}
	mpf, err := NewFileFromPart(part)
	if mpf == nil || err != nil {
		t.Error("Expected non-nil MultipartFile, nil error")
	}
	if mpf.IsDirectory() {
		t.Error("Expected file to not be a directory")
	}
	if mpf.FileName() != "file.txt" {
		t.Error("Expected filename to be \"file.txt\"")
	}
	if n, err := mpf.Read(buf); n != len(text) || err != nil {
		t.Error("Expected to read from file", n, err)
	}
	if string(buf[:len(text)]) != text {
		t.Error("Data read was different than expected")
	}

	part, err = mpReader.NextPart()
	if part == nil || err != nil {
		t.Error("Expected non-nil part, nil error")
	}
	mpf, err = NewFileFromPart(part)
	if mpf == nil || err != nil {
		t.Error("Expected non-nil MultipartFile, nil error")
	}
	if !mpf.IsDirectory() {
		t.Error("Expected file to be a directory")
	}
	if mpf.FileName() != "boop" {
		t.Error("Expected filename to be \"boop\"")
	}

	part, err = mpReader.NextPart()
	if part == nil || err != nil {
		t.Error("Expected non-nil part, nil error")
	}
	mpf, err = NewFileFromPart(part)
	if mpf.IsDirectory() {
		t.Error("Expected file to not be a directory")
	}
	if mpf.FileName() != "boop/a.txt" {
		t.Error("Expected filename to be \"some/file/path\"")
	}

	part, err = mpReader.NextPart()
	if part == nil || err != nil {
		t.Error("Expected non-nil part, nil error")
	}
	mpf, err = NewFileFromPart(part)
	if mpf.IsDirectory() {
		t.Error("Expected file to not be a directory")
	}
	if mpf.FileName() != "boop/b.txt" {
		t.Error("Expected filename to be \"some/file/path\"")
	}

	child, err := mpf.NextFile()
	if child != nil || err != ErrNotDirectory {
		t.Error("Expected a nil file and ErrNotDirectory")
	}

	part, err = mpReader.NextPart()
	if part == nil || err != nil {
		t.Error("Expected non-nil part, nil error")
	}
	mpf, err = NewFileFromPart(part)
	if mpf == nil || err != nil {
		t.Error("Expected non-nil MultipartFile, nil error")
	}

	part, err = mpReader.NextPart()
	if part != nil || err != io.EOF {
		t.Error("Expected to get (nil, io.EOF)")
	}
}
