package utils

import (
	"bytes"
	"context"
	"fmt"
	"github.com/conku/webp"
	"github.com/pkg/errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"patreon/internal/app"
	repoFiles "patreon/internal/app/repository/files"
	"strings"
)

type ConverterToWebp struct {
}

// Convert Errors:
// 	app.GeneralError:
//		utils.ConvertErr
// 		utils.UnknownExtOfFileName
func (cv *ConverterToWebp) Convert(_ context.Context,
	file io.Reader, name repoFiles.FileName) (io.Reader, repoFiles.FileName, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, "", app.GeneralError{
			Err:         errors.Wrap(ConvertErr, "error in get image"),
			ExternalErr: err,
		}
	}

	buf, err := webp.EncodeExactLosslessRGBA(img)
	if err != nil {
		return nil, "", app.GeneralError{
			Err:         errors.Wrap(ConvertErr, "error in webp convertor"),
			ExternalErr: err,
		}
	}

	pos := strings.LastIndex(string(name), ".")
	if pos == -1 {
		return nil, "", errors.Wrap(UnknownExtOfFileName, fmt.Sprintf("error with %s: ", name))
	}

	name = name[:pos] + ".webp"

	res := bytes.NewReader(buf)
	return res, name, nil
}
