package data

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
)

func (data Call) GetBodyAsTempFile() (*os.File, error) {
	bodyStr := strings.Join(data.Body, "\n")

	// TODO Make this configurable so it can be inspected
	tmpFile, err := ioutil.TempFile("", "ain-body")
	if err != nil {
		return nil, errors.Wrap(err, "Could not create tempfile")
	}

	if _, err := tmpFile.Write([]byte(bodyStr)); err != nil {
		// This also returns an error, but the first is more significant
		// so ignore this, it's only a temp-file that will be deleted eventually
		_ = tmpFile.Close()

		return nil, errors.Wrap(err, "Could not write to tempfile")
	}

	return tmpFile, nil
}
