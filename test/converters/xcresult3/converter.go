package xcresult3

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-xcode/xcodeproject/serialized"
	"github.com/bitrise-steplib/bitrise-step-build-router-start/test/junit"
	"howett.net/plist"
)

// Converter ...
type Converter struct {
	xcresultPth string
}

func majorVersion(document serialized.Object) (int, error) {
	version, err := document.Object("version")
	if err != nil {
		return -1, err
	}

	major, err := version.Value("major")
	if err != nil {
		return -1, err
	}
	return int(major.(uint64)), nil
}

func documentMajorVersion(pth string) (int, error) {
	content, err := fileutil.ReadBytesFromFile(pth)
	if err != nil {
		return -1, err
	}

	var info serialized.Object
	if _, err := plist.Unmarshal(content, &info); err != nil {
		return -1, err
	}

	return majorVersion(info)
}

// Detect ...
func (c *Converter) Detect(files []string) bool {
	if !isXcresulttoolAvailable() {
		log.Debugf("xcresult tool is not available")
		return false
	}

	for _, file := range files {
		if filepath.Ext(file) != ".xcresult" {
			continue
		}

		infoPth := filepath.Join(file, "Info.plist")
		if exist, err := pathutil.IsPathExists(infoPth); err != nil {
			log.Debugf("Failed to find Info.plist at %s: %s", infoPth, err)
			continue
		} else if !exist {
			log.Debugf("No Info.plist found at %s", infoPth)
			continue
		}

		version, err := documentMajorVersion(infoPth)
		if err != nil {
			log.Debugf("failed to get document version: %s", err)
			continue
		}

		if version < 3 {
			log.Debugf("version < 3: %d", version)
			continue
		}

		c.xcresultPth = file
		return true
	}
	return false
}

// XML ...
func (c *Converter) XML() (junit.XML, error) {
	testResultDir := filepath.Dir(c.xcresultPth)

	record, summaries, err := Parse(c.xcresultPth)
	if err != nil {
		return junit.XML{}, err
	}

	var xmlData junit.XML
	for _, summary := range summaries {
		testsByName := summary.tests()

		for name, tests := range testsByName {
			testSuit := junit.TestSuite{
				Name:     name,
				Tests:    len(tests),
				Failures: summary.failuresCount(name),
				Time:     summary.totalTime(name),
			}

			for _, test := range tests {
				var duartion float64
				if test.Duration.Value != "" {
					duartion, err = strconv.ParseFloat(test.Duration.Value, 64)
					if err != nil {
						return junit.XML{}, err
					}
				}

				failureMessage := record.failure(test, testSuit)

				var failure *junit.Failure
				if len(failureMessage) > 0 {
					failure = &junit.Failure{
						Value: failureMessage,
					}
				}

				testSuit.TestCases = append(testSuit.TestCases, junit.TestCase{
					Name:      test.Name.Value,
					ClassName: strings.Split(test.Identifier.Value, "/")[0],
					Failure:   failure,
					Time:      duartion,
				})

				if err := test.exportScreenshots(c.xcresultPth, testResultDir); err != nil {
					return junit.XML{}, err
				}
			}

			xmlData.TestSuites = append(xmlData.TestSuites, testSuit)
		}
	}

	return xmlData, nil
}
