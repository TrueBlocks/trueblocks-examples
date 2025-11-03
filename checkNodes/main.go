package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/TrueBlocks/trueblocks-chifra/v6/pkg/base"
	"github.com/TrueBlocks/trueblocks-chifra/v6/pkg/colors"
	"github.com/TrueBlocks/trueblocks-chifra/v6/pkg/file"
	"github.com/TrueBlocks/trueblocks-chifra/v6/pkg/logger"
	"github.com/TrueBlocks/trueblocks-chifra/v6/pkg/tslib"
	"github.com/TrueBlocks/trueblocks-chifra/v6/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v6"
)

func main() {
	verbose := false
	for i, arg := range os.Args {
		if i == 0 {
			continue
		}
		if arg == "--verbose" {
			verbose = true
		} else {
			logger.Fatal("unknown argument:", arg)
		}
	}

	writer := logger.GetLoggerWriter()
	defer logger.SetLoggerWriter(writer)
	if !verbose {
		logger.SetLoggerWriter(io.Discard)
	}

	fmt.Println("")

	var buffer bytes.Buffer
	w := tabwriter.NewWriter(&buffer, 0, 0, 2, ' ', 0)

	fmt.Fprintf(w, colors.Yellow+"\tSource\tChain\tBlock\tTimestamp\tDate\tTime\tBehind"+colors.Off+"\n")
	fmt.Fprintf(w, colors.Yellow+"\t------\t---------\t-----------\t---------\t----\t----\t------"+colors.Off+"\n")

	modes := []string{"latest", "stage", "final"}
	chains, err := getChains()
	if err != nil {
		logger.Error("error getting chains:", err)
		os.Exit(1)
	}
	for i, mode := range modes {
		if i > 0 {
			fmt.Fprintf(w, colors.Yellow+"\t\t\t\t\t\t\t"+colors.Off+"\n")
		}
		for _, chain := range chains {
			if meta, err := sdk.GetMetaData(chain); err != nil {
				logger.Fatal("error from GetMetaData:", mode, chain, err)
			} else {
				ReportOne(w, meta, mode)
			}
		}
	}

	w.Flush()
	text := buffer.String()
	for _, chain := range chains {
		text = strings.ReplaceAll(text, chain, colors.Cyan+chain+colors.Green)
	}
	fmt.Print(text)
	fmt.Println(colors.Off)
}

func ReportOne(w *tabwriter.Writer, meta *types.MetaData, mode string) {
	report := Report{
		Source: mode,
		Meta:   *meta,
	}

	var blockNumber base.Blknum
	var timestamp base.Timestamp
	var err error

	switch mode {
	case "latest":
		blockNumber = meta.Latest
		timestamp, err = sdk.TsFromBlock(meta.Chain, blockNumber)
	case "stage":
		blockNumber = meta.Staging
		timestamp, err = tslib.FromBnToTs(meta.Chain, blockNumber)
	case "final":
		blockNumber = meta.Finalized
		timestamp, err = tslib.FromBnToTs(meta.Chain, blockNumber)
	}

	if err != nil {
		fmt.Println(err.Error())
	}

	report.Block = types.NamedBlock{
		BlockNumber: blockNumber,
		Timestamp:   timestamp,
	}

	fmt.Fprintf(w, colors.Green+"%s\n", report.String())
}

type Report struct {
	Block  types.NamedBlock
	Meta   types.MetaData
	Source string
}

func (t *Report) String() string {
	now := time.Now().Unix()
	duration := time.Unix(now, 0).Sub(time.Unix(t.Block.Timestamp.Int64(), 0))
	formattedDuration := formatDuration(duration)
	dParts := strings.Split(t.Block.Date(), " ")

	return fmt.Sprintf("\t%s\t%s\t% 10d\t%d\t%s\t%s\t%s",
		t.Source,
		t.Meta.Chain,
		t.Block.BlockNumber,
		t.Block.Timestamp,
		dParts[0],
		dParts[1],
		formattedDuration)
}

func formatDuration(duration time.Duration) string {
	hours := int(duration / time.Hour)
	minutes := int((duration % time.Hour) / time.Minute)
	seconds := duration.Seconds() - float64(hours*3600+minutes*60)

	switch {
	case hours > 0:
		return fmt.Sprintf("%2dh %2dm %5.2fs", hours, minutes, seconds)
	case minutes > 0:
		return fmt.Sprintf("%2dm %5.2fs", minutes, seconds)
	default:
		return fmt.Sprintf("%5.2fs", seconds)
	}
}

func getChains() ([]string, error) {
	chainsFn := "./chains.csv"
	if !file.FileExists(chainsFn) {
		return nil, fmt.Errorf("chains file does not exist")
	}

	chains := file.AsciiFileToLines(chainsFn)
	uniqueChains := make(map[string]bool)
	for _, chain := range chains {
		if !strings.HasPrefix(chain, "#") {
			uniqueChains[chain] = true
		}
	}

	result := make([]string, 0, len(uniqueChains))
	for chain := range uniqueChains {
		result = append(result, chain)
	}
	sort.Strings(chains)

	return result, nil
}
