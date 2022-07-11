package util

import (
	"github.com/saintfish/chardet"
	"golang.org/x/text/encoding/simplifiedchinese"
	"runtime"
)

func Convert2Utf8Str(str string) string {

	rawBytes := []byte(str)
	detector := chardet.NewTextDetector()
	charset, err := detector.DetectBest(rawBytes)
	if err != nil {
		panic(err)
	}

	outStr := ""

	//log.Print(rawBytes, charset)

	if charset.Charset != "UTF-8" && runtime.GOOS == "windows" { // en... maybe  need other ways to detect chinese
		charset.Charset = "GB18030"
	}

	switch charset.Charset {
	case "ISO-8859-1", "GB18030", "GB-18030":
		var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(rawBytes)
		outStr = string(decodeBytes)
	case "UTF-8":
		fallthrough
	default:
		outStr = string(rawBytes)
	}

	return outStr
}
