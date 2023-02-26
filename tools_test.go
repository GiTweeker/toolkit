package toolkit_test

import (
	"fmt"
	"github.com/GiTweeker/toolkit"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
)

func TestTools_RandomString(t *testing.T) {
	var testTools toolkit.Tools
	s := testTools.RandomString(10)

	if len(s) != 10 {
		t.Error("wrong length of random string returned")
	}
}

var uploadTests = []struct {
	name          string
	allowedTypes  []string
	renameFile    bool
	errorExpected bool
}{
	{
		name:          "allowed no rename",
		allowedTypes:  []string{"image/jpeg", "image/png"},
		renameFile:    false,
		errorExpected: false,
	},
	{
		name:          "allowed renamed",
		allowedTypes:  []string{"image/jpeg", "image/png"},
		renameFile:    true,
		errorExpected: false,
	},
	{
		name:          "not allowed file type",
		allowedTypes:  []string{"image/jpeg"},
		renameFile:    false,
		errorExpected: true,
	},
}

func TestTools_UploadFiles(t *testing.T) {

	for _, test := range uploadTests {
		t.Run(test.name, func(t *testing.T) {
			pr, pw := io.Pipe()
			writer := multipart.NewWriter(pw)

			wg := sync.WaitGroup{}

			wg.Add(1)

			go func() {
				defer writer.Close()
				defer wg.Done()
				const filePath = "./testdata/img.png"
				part, err := writer.CreateFormFile("file", filePath)

				if err != nil {
					t.Error(err)
				}

				f, err := os.Open(filePath)

				if err != nil {
					t.Error(err)
				}

				defer f.Close()

				img, _, err := image.Decode(f)

				if err != nil {
					t.Error("error decoding image ", err)
				}

				err = png.Encode(part, img)

				if err != nil {
					t.Error(err)
				}

			}()

			request := httptest.NewRequest("POST", "/", pr)
			request.Header.Add("Content-Type", writer.FormDataContentType())

			var testTools toolkit.Tools

			testTools.AllowFileTypes = test.allowedTypes
			uploadedFiles, err := testTools.UploadFiles(request, "./testdata/uploads/", test.renameFile)

			if err != nil && !test.errorExpected {
				t.Error(err)
			}

			if !test.errorExpected {
				fileName := fmt.Sprintf("./testdata/uploads/%s", uploadedFiles[0].NewFileName)
				if _, err := os.Stat(fileName); os.IsNotExist(err) {
					t.Errorf("%s: expected file to exist: %s", test.name, err.Error())
				}

				//clean up

				_ = os.Remove(fileName)
			}

			if !test.errorExpected && err != nil {
				t.Errorf("%s: error expected but none received", test.name)
			}

			wg.Wait()
		})

	}
}

func TestTools_UploadOneFile(t *testing.T) {
	for _, test := range uploadTests {
		t.Run(test.name, func(t *testing.T) {
			pr, pw := io.Pipe()
			writer := multipart.NewWriter(pw)

			wg := sync.WaitGroup{}

			wg.Add(1)

			go func() {
				defer writer.Close()
				defer wg.Done()
				const filePath = "./testdata/img.png"
				part, err := writer.CreateFormFile("file", filePath)

				if err != nil {
					t.Error(err)
				}

				f, err := os.Open(filePath)

				if err != nil {
					t.Error(err)
				}

				defer f.Close()

				img, _, err := image.Decode(f)

				if err != nil {
					t.Error("error decoding image ", err)
				}

				err = png.Encode(part, img)

				if err != nil {
					t.Error(err)
				}

			}()

			request := httptest.NewRequest("POST", "/", pr)
			request.Header.Add("Content-Type", writer.FormDataContentType())

			var testTools toolkit.Tools

			testTools.AllowFileTypes = test.allowedTypes
			uploadedFile, err := testTools.UploadOneFile(request, "./testdata/uploads/", test.renameFile)

			if err != nil && !test.errorExpected {
				t.Error(err)
			}

			if !test.errorExpected {
				fileName := fmt.Sprintf("./testdata/uploads/%s", uploadedFile.NewFileName)
				if _, err := os.Stat(fileName); os.IsNotExist(err) {
					t.Errorf("%s: expected file to exist: %s", test.name, err.Error())
				}

				//clean up

				_ = os.Remove(fileName)
			}

			if !test.errorExpected && err != nil {
				t.Errorf("%s: error expected but none received", test.name)
			}

			wg.Wait()
		})

	}
}

func TestTools_CreateDirIfNotExist(t *testing.T) {
	var testTools toolkit.Tools
	const path = "./testdata/mytestdir"
	err := testTools.CreateDirIfNotExist(path)

	if err != nil {
		t.Error(err)
	}

	err = testTools.CreateDirIfNotExist(path)

	if err != nil {
		t.Error(err)
	}

	_ = os.Remove(path)
}

var slugTests = []struct {
	name          string
	s             string
	expected      string
	errorExpected bool
}{
	{name: "valid string", s: "now is the time", expected: "now-is-the-time", errorExpected: false},
	{name: "empty string", s: "", expected: "", errorExpected: true},
	{name: "complex string", s: "now is the time for all Good men! + fish & such &^123", expected: "now-is-the-time-for-all-good-men-fish-such-123", errorExpected: false},
	{name: "japanese string", s: "トラブルが好きな人もいる", expected: "", errorExpected: true},
	{name: "japanese and roman characters string", s: "トラブルが好きな人もいる+hello+world", expected: "hello-world", errorExpected: false},
}

func TestTools_Sluggify(t *testing.T) {
	var testTool toolkit.Tools

	for _, test := range slugTests {
		t.Run(test.name, func(t *testing.T) {
			slug, err := testTool.Sluggify(test.s)

			if err != nil && !test.errorExpected {
				t.Errorf("%s: error received when none expected: %s", test.name, err.Error())
			}

			if !test.errorExpected && slug != test.expected {
				t.Errorf("%s: wrong slug returned; expected %s but got %s", test.name, test.expected, slug)
			}
		})

	}
}
