package warp

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/peakefficiency/warp-diag-checker/wdc"
)

type CheckResult struct {
	CheckName    string
	CheckPass    bool
	IssueType    string
	Evidence     string
	ReplyMessage string
}

type Printer struct {
	Output     io.Writer
	HasPrinted bool
}

func NewPrinter() *Printer {
	return &Printer{Output: os.Stdout}
}

func (zipContent FileContentMap) DumpFiles(filename string) {

	if filename != "" {
		if content, ok := zipContent[filename]; ok {
			fmt.Println(filename)
			fmt.Println(string(content.Data))
		} else {
			fmt.Printf("File %s not found in zip\n", filename)
		}
	} else {
		fmt.Println("# File Contents")

		for name, content := range zipContent {
			fmt.Printf("## %s\n", name)
			fmt.Println(string(content.Data))
		}
	}

}

func ReportLogSearch(results map[string]LogSearchResult) (string, error) {
	var markdown strings.Builder

	if len(results) == 0 {
		return "", nil
	}
	markdown.WriteString("## Log Search Results\n")

	for issueType, result := range results {
		reply := wdc.WdcConf.ReplyByIssueType[issueType]

		markdown.WriteString(fmt.Sprintf("### %s\n", issueType))
		markdown.WriteString(fmt.Sprintf("%s\n", reply.Message))
		markdown.WriteString(fmt.Sprintf("- Evidence: \n\n```\n%s\n```\n\n", result.Evidence))
		markdown.WriteString("\n")
	}

	if wdc.Plain {
		return markdown.String(), nil
	}

	return glamour.Render(markdown.String(), "dark")
}

func (result CheckResult) MarkdownCheckResult() (string, error) {
	var markdown strings.Builder

	if !result.CheckPass && result.Evidence != "" {
		replyMsg := wdc.WdcConf.ReplyByIssueType[result.IssueType].Message

		markdown.WriteString(fmt.Sprintf("## %s\n", result.CheckName))

		markdown.WriteString(fmt.Sprintf("%s\n", replyMsg))

		markdown.WriteString(fmt.Sprintf("- Evidence: \n\n```\n%s\n```\n\n", result.Evidence))

		if wdc.Plain {
			return markdown.String(), nil
		}

		return glamour.Render(markdown.String(), "dark")
	}
	return "", nil
}

func (p *Printer) PrintCheckResult(result CheckResult, err error) {

	if err != nil {
		fmt.Fprintf(p.Output, "Error generating check result: %s", err)

	}
	markdown, err := result.MarkdownCheckResult()
	if err != nil {
		fmt.Fprintf(p.Output, "Error generating markdown: %s", err)
		return
	}
	if markdown != "" {
		p.HasPrinted = true
		fmt.Fprintf(p.Output, "%s", markdown)
	}
}

func (p *Printer) PrintString(s string, err error) {

	if err != nil {
		fmt.Fprintf(p.Output, "Error generating check result: %s", err)

	}

	if s != "" {
		p.HasPrinted = true
		fmt.Fprintf(p.Output, "%s", s)
	}
}
