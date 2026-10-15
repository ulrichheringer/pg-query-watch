package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type QueryStat struct {
	QueryID string  `json:"queryid"`
	Calls   float64 `json:"calls"`
	TotalMS float64 `json:"total_exec_time"`
	Rows    float64 `json:"rows"`
	Query   string  `json:"query"`
}

func MeanMS(stat QueryStat) float64 {
	if stat.Calls == 0 {
		return 0
	}
	return stat.TotalMS / stat.Calls
}

func ParseStats(r io.Reader) ([]QueryStat, error) {
	var stats []QueryStat
	decoder := json.NewDecoder(r)
	for decoder.More() {
		var stat QueryStat
		if err := decoder.Decode(&stat); err != nil {
			return nil, err
		}
		stats = append(stats, stat)
	}
	return stats, nil
}

func FilterSlow(stats []QueryStat, minMeanMS float64) []QueryStat {
	var slow []QueryStat
	for _, stat := range stats {
		if MeanMS(stat) >= minMeanMS {
			slow = append(slow, stat)
		}
	}
	return slow
}

const statsSQL = `SELECT json_build_object('queryid', queryid::text, 'calls', calls, 'total_exec_time', total_exec_time, 'rows', rows, 'query', query)::text FROM pg_stat_statements ORDER BY total_exec_time / GREATEST(calls, 1) DESC LIMIT %s;`

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("pg-query-watch", flag.ContinueOnError)
	flags.SetOutput(stderr)
	minMS := flags.Float64("min-ms", 0, "minimum mean execution time in milliseconds")
	limit := flags.Int("limit", 25, "number of top statements to fetch (1-500)")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *limit < 1 || *limit > 500 {
		return errors.New("limit must be between 1 and 500")
	}
	if _, err := exec.LookPath("psql"); err != nil {
		return errors.New("psql is required; install PostgreSQL client tools and configure PGHOST/PGDATABASE/PGUSER (and PGPASSWORD if needed)")
	}
	query := fmt.Sprintf(statsSQL, strconv.Itoa(*limit))
	cmd := exec.Command("psql", "-X", "-qAt", "-v", "ON_ERROR_STOP=1", "-c", query)
	cmd.Stderr = stderr
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("query pg_stat_statements: %w (ensure the extension is installed and visible)", err)
	}
	stats, err := ParseStats(strings.NewReader(string(output)))
	if err != nil {
		return err
	}
	slow := FilterSlow(stats, *minMS)
	if len(slow) == 0 {
		_, err := fmt.Fprintln(stdout, "No statements matched the mean latency threshold.")
		return err
	}
	for _, stat := range slow {
		queryText := strings.Join(strings.Fields(stat.Query), " ")
		if len(queryText) > 240 {
			queryText = queryText[:237] + "..."
		}
		if _, err := fmt.Fprintf(stdout, "id=%s calls=%.0f mean=%.2fms total=%.2fms rows=%.0f\n  %s\n", stat.QueryID, stat.Calls, MeanMS(stat), stat.TotalMS, stat.Rows, queryText); err != nil {
			return err
		}
	}
	return nil
}
