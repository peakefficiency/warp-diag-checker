package warp

import (
	"strings"
)

type LogSearchResult struct {
	Filename     string
	SearchTerm   string
	SearchStatus bool
	IssueType    string
	Evidence     string
}

type IssueStruct struct {
	SearchTerms []string `yaml:"search_term"`
}

var LogSearchOutput = map[string]LogSearchResult{}

func (zipContent FileContentMap) LogSearch(info ParsedDiag) map[string]LogSearchResult {

	for _, logPattern := range WdcConf.LogPatternsByIssue {
		// Split the comma-separated string into a slice of strings
		searchFilenames := strings.Split(logPattern.SearchFile, ",")

		for _, filename := range searchFilenames {
			filename = strings.TrimSpace(filename) // Remove any leading/trailing whitespace

			content, found := zipContent[filename]
			if !found {
				continue
			}

			fileContent := string(content.Data)
			lines := strings.Split(fileContent, "\n")

			// Reverse the lines slice
			for i := len(lines)/2 - 1; i >= 0; i-- {
				opp := len(lines) - 1 - i
				lines[i], lines[opp] = lines[opp], lines[i]
			}

			for issueType, issue := range logPattern.Issue {
				evidence := []string{}
				numEntries := 5

				for _, searchTerm := range issue.SearchTerms {
					for i := len(lines) - 1; i >= 0; i-- {
						line := lines[i]

						if strings.Contains(line, searchTerm) {
							evidence = append(evidence, line)
							if len(evidence) >= numEntries {
								break
							}
						}
					}
				}

				// Reverse the evidence slice
				for i := len(evidence)/2 - 1; i >= 0; i-- {
					opp := len(evidence) - 1 - i
					evidence[i], evidence[opp] = evidence[opp], evidence[i]
				}

				if len(evidence) > 0 {
					LogSearchOutput[issueType] = LogSearchResult{
						IssueType: issueType,
						Evidence:  strings.Join(evidence, "\n"),
					}
				}
			}
		}
	}

	return LogSearchOutput
}
