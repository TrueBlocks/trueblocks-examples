package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

type TestCase struct {
	Directory string
	Command   string
	Group     int
}

func main() {
	// Group test cases by Group number
	groups := make(map[int][]TestCase)
	for _, tc := range testCases {
		groups[tc.Group] = append(groups[tc.Group], tc)
	}

	// Run groups concurrently
	var wg sync.WaitGroup
	for groupNum, groupTests := range groups {
		wg.Add(1)
		go func(groupNum int, tests []TestCase) {
			defer wg.Done()
			// Run tests within group sequentially
			for _, tc := range tests {
				if err := processTestCase(tc); err != nil {
					fmt.Printf("Error in group %d, test %s: %v\n", groupNum, tc.Directory, err)
				}
			}
		}(groupNum, groupTests)
	}

	// Wait for all groups to complete
	wg.Wait()
}

func processTestCase(tc TestCase) error {
	defer fmt.Println("Processed:", tc.Directory, "(Group", tc.Group, ")")

	originalFile := filepath.Join(tc.Directory, "transfers.csv")
	tempFile := filepath.Join(tc.Directory, "transfers_new.csv")

	originalContent, err := os.ReadFile(originalFile)
	if err != nil {
		return fmt.Errorf("failed to read original transfers.csv: %v", err)
	}

	header := fmt.Sprintf("# %s\n", tc.Command)
	if err := os.WriteFile(tempFile, []byte(header), 0644); err != nil {
		return fmt.Errorf("failed to write command header: %v", err)
	}

	cmd := exec.Command("sh", "-c", tc.Command+" --append --output "+tempFile)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command execution failed: %v\nStderr: %s", err, stderr.String())
	}

	newContent, err := os.ReadFile(tempFile)
	if err != nil {
		return fmt.Errorf("failed to read new transfers.csv: %v", err)
	}

	compareFiles(originalContent, newContent, tc.Directory)

	if err := os.Remove(originalFile); err != nil {
		return fmt.Errorf("failed to remove original file: %v", err)
	}
	if err := os.Rename(tempFile, originalFile); err != nil {
		return fmt.Errorf("failed to rename new file to original: %v", err)
	}

	return nil
}

func compareFiles(original, new []byte, directory string) bool {
	originalLines := strings.Split(string(original), "\n")
	newLines := strings.Split(string(new), "\n")

	originalLines = trimEmptyLines(originalLines)
	newLines = trimEmptyLines(newLines)

	hasMismatches := false

	if len(originalLines) != len(newLines) {
		hasMismatches = true
		fmt.Printf("Line count mismatch in %s: original=%d, new=%d\n", directory, len(originalLines), len(newLines))
	}

	maxLines := len(originalLines)
	if len(newLines) < maxLines {
		maxLines = len(newLines)
	}

	for i := 0; i < maxLines; i++ {
		if originalLines[i] != newLines[i] {
			hasMismatches = true
			fmt.Printf("Difference at line %d in %s:\n", i+1, directory)
			fmt.Printf("  Original: %s\n", originalLines[i])
			fmt.Printf("  New:      %s\n", newLines[i])
		}
	}

	if len(originalLines) > len(newLines) {
		hasMismatches = true
		fmt.Printf("Original has %d extra lines in %s\n", len(originalLines)-len(newLines), directory)
	} else if len(newLines) > len(originalLines) {
		hasMismatches = true
		fmt.Printf("New has %d extra lines in %s\n", len(newLines)-len(originalLines), directory)
	}

	return hasMismatches
}

func trimEmptyLines(lines []string) []string {
	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

var testCases = []TestCase{
	// Group 1: 0x08166f02313feae18bb044e7877c808b55b5bf58
	{Directory: "acctExport_statement_2_ether", Command: "chifra export 0x08166f02313feae18bb044e7877c808b55b5bf58 --accounting --transfers --last_block 4000000 --fmt csv", Group: 1},
	{Directory: "acctExport_statement_2", Command: "chifra export 0x08166f02313feae18bb044e7877c808b55b5bf58 --accounting --transfers --last_block 4000000 --fmt csv", Group: 1},
	{Directory: "acctExport_statement_3_accounting", Command: "chifra export 0x08166f02313feae18bb044e7877c808b55b5bf58 --accounting --transfers --last_block 4000000 --fmt csv", Group: 1},
	{Directory: "acctExport_statement_3", Command: "chifra export 0x08166f02313feae18bb044e7877c808b55b5bf58 --accounting --transfers --last_block 4000000 --fmt csv", Group: 1},
	{Directory: "acctExport_transfer_2_ether", Command: "chifra export 0x08166f02313feae18bb044e7877c808b55b5bf58 --accounting --transfers --last_block 4000000 --fmt csv", Group: 1},
	{Directory: "acctExport_transfer_2", Command: "chifra export 0x08166f02313feae18bb044e7877c808b55b5bf58 --accounting --transfers --last_block 4000000 --fmt csv", Group: 1},
	{Directory: "acctExport_transfer_3_accounting", Command: "chifra export 0x08166f02313feae18bb044e7877c808b55b5bf58 --accounting --transfers --last_block 4000000 --fmt csv", Group: 1},
	{Directory: "acctExport_transfer_3", Command: "chifra export 0x08166f02313feae18bb044e7877c808b55b5bf58 --accounting --transfers --last_block 4000000 --fmt csv", Group: 1},

	// Group 2: 0x054993ab0f2b1acc0fdc65405ee203b4271bebe6
	{Directory: "acctExport_statement_failed_2572_1", Command: "chifra export 0x054993ab0f2b1acc0fdc65405ee203b4271bebe6 --accounting --transfers --asset 0xf5b2c59f6db42ffcdfc1625999c81fdf17953384 --first_block 15549162 --last_block 15549163 --max_records 1 --fmt csv", Group: 2},
	{Directory: "acctExport_transfer_failed_2572_1", Command: "chifra export 0x054993ab0f2b1acc0fdc65405ee203b4271bebe6 --accounting --transfers --asset 0xf5b2c59f6db42ffcdfc1625999c81fdf17953384 --last_block 15549163 --max_records 40 --fmt csv", Group: 2},

	// Group 3: 0x65b0d5e1dc0dee0704f53f660aa865c72e986fc7
	{Directory: "acctExport_statement_failed_2572_2", Command: "chifra export 0x65b0d5e1dc0dee0704f53f660aa865c72e986fc7 --accounting --transfers --asset 0xc713e5e149d5d0715dcd1c156a020976e7e56b88 --first_block 11670418 --last_block 11670420 --max_records 40 --fmt csv", Group: 3},
	{Directory: "acctExport_transfer_failed_2572_2", Command: "chifra export 0x65b0d5e1dc0dee0704f53f660aa865c72e986fc7 --accounting --transfers --asset 0xc713e5e149d5d0715dcd1c156a020976e7e56b88 --first_block 11670418 --last_block 11670420 --max_records 40 --fmt csv", Group: 3},

	// Group 4: 0xf503017d7baf7fbc0fff7492b751025c6a78179b
	{Directory: "acctExport_statement_filter_traces", Command: "chifra export 0xf503017d7baf7fbc0fff7492b751025c6a78179b --accounting --transfers --traces --first_block 8860513 --last_block 8860531 --asset 0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359 --fmt csv", Group: 4},
	{Directory: "acctExport_statement_filtered", Command: "chifra export 0xf503017d7baf7fbc0fff7492b751025c6a78179b --accounting --transfers --first_block 8860513 --last_block 8860531 --asset 0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359 --fmt csv", Group: 4},
	{Directory: "acctExport_statement_nft", Command: "chifra export 0xf503017d7baf7fbc0fff7492b751025c6a78179b --accounting --transfers --first_block 8876230 --last_block 9024186 --fmt csv", Group: 4},
	{Directory: "acctExport_statement_token_ibt_2", Command: "chifra export 0xf503017d7baf7fbc0fff7492b751025c6a78179b --accounting --transfers --first_block 12704455 --last_block 12705893 --fmt csv", Group: 4},
	{Directory: "acctExport_statement_unfiltered", Command: "chifra export 0xf503017d7baf7fbc0fff7492b751025c6a78179b --accounting --transfers --first_block 8860513 --last_block 8860531 --fmt csv", Group: 4},
	{Directory: "acctExport_transfer_2_asset_filt", Command: "chifra export 0xf503017d7baf7fbc0fff7492b751025c6a78179b --accounting --transfers --asset 0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee --asset 0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359 --first_block 8856476 --last_block 9193814 --fmt csv", Group: 4},
	{Directory: "acctExport_transfer_filter_traces", Command: "chifra export 0xf503017d7baf7fbc0fff7492b751025c6a78179b --accounting --transfers --traces --first_block 8860513 --last_block 8860531 --asset 0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359 --fmt csv", Group: 4},
	{Directory: "acctExport_transfer_filtered", Command: "chifra export 0xf503017d7baf7fbc0fff7492b751025c6a78179b --accounting --transfers --first_block 8860513 --last_block 8860531 --asset 0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359 --fmt csv", Group: 4},
	{Directory: "acctExport_transfer_nft", Command: "chifra export 0xf503017d7baf7fbc0fff7492b751025c6a78179b --accounting --transfers --first_block 8876230 --last_block 9024186 --fmt csv", Group: 4},
	{Directory: "acctExport_transfer_token_ibt_2", Command: "chifra export 0xf503017d7baf7fbc0fff7492b751025c6a78179b --accounting --transfers --first_block 12704455 --last_block 12705893 --fmt csv", Group: 4},
	{Directory: "acctExport_transfer_unfiltered", Command: "chifra export 0xf503017d7baf7fbc0fff7492b751025c6a78179b --accounting --transfers --first_block 8860513 --last_block 8860531 --fmt csv", Group: 4},

	// Group 5: 0x868b8fd259abfcfdf9634c343593b34ef359641d
	{Directory: "acctExport_statement_forward", Command: "chifra export 0x868b8fd259abfcfdf9634c343593b34ef359641d --accounting --transfers --traces --last_block 8769141 --fmt csv", Group: 5},
	{Directory: "acctExport_statement_tributes", Command: "chifra export 0x868b8fd259abfcfdf9634c343593b34ef359641d --accounting --transfers --first_block 8769018 --last_block 8769053 --asset 0x0ba45a8b5d5575935b8158a88c631e9f9c95a2e5 --fmt csv", Group: 5},
	{Directory: "acctExport_transfer_forward", Command: "chifra export 0x868b8fd259abfcfdf9634c343593b34ef359641d --accounting --transfers --traces --last_block 8769141 --fmt csv", Group: 5},
	{Directory: "acctExport_transfer_tributes", Command: "chifra export 0x868b8fd259abfcfdf9634c343593b34ef359641d --accounting --transfers --first_block 8769018 --last_block 8769053 --asset 0x0ba45a8b5d5575935b8158a88c631e9f9c95a2e5 --fmt csv", Group: 5},

	// Group 6: 0xec3ef464bf821c3b10a18adf9ac7177a628e87cc
	{Directory: "acctExport_statement_token_ibt", Command: "chifra export 0xec3ef464bf821c3b10a18adf9ac7177a628e87cc --accounting --transfers --first_block 7005600 --last_block 7005780 --fmt csv", Group: 6},
	{Directory: "acctExport_transfer_token_ibt", Command: "chifra export 0xec3ef464bf821c3b10a18adf9ac7177a628e87cc --accounting --transfers --first_block 7005600 --last_block 7005780 --fmt csv", Group: 6},

	// Group 7: 0x05a56e2d52c817161883f50c441c3228cfe54d9f
	{Directory: "acctExport_statement_wei_2_1", Command: "chifra export 0x05a56e2d52c817161883f50c441c3228cfe54d9f --accounting --transfers --first_record 0 --max_records 15 --fmt csv", Group: 7},
	{Directory: "acctExport_transfer_wei_2_1", Command: "chifra export 0x05a56e2d52c817161883f50c441c3228cfe54d9f --accounting --transfers --first_record 0 --max_records 15 --fmt csv", Group: 7},
}
