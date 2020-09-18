package cmd

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger lumberjack.Logger

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "lumberjack-logger [flags] <filename>",
	Short: "Log to rolling files from stdin.",
	Long: `
A small command line wrapper around lumberjack (a rolling log library for go
by Nate Finch).

* Directories in the path given for <filename> are created if necessary.
* Anything that looks (by its filename) like the log file or one of its
  backups will be treated that way. For example, if the path specifies a
  directory, lumberjack will attempt to move the directory in the same maner
  as it would a full log file.
* If it's not possible to write to the given location, an error message will
  be printed to stderr for every write attempt that fails.

If you're planning to start a command and then leave it running after
disconnecting the terminal, you will probably want to run it with nohup, e.g.

	nohup yes "yes" < /dev/null 2> /dev/null | lumberjack-logger -s 1 -b 2 yeslog.txt 2>&1 > /dev/null & disown

For more about that, see
* & vs nohup vs disown: https://unix.stackexchange.com/a/148698
* more about nohup: https://stackoverflow.com/a/10408906
* piping stderr and stdout: https://stackoverflow.com/a/16497456
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("no filename")
		}
		var err error
		logger.Filename, err = filepath.Abs(args[0])
		if err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			_, err := logger.Write(append([]byte(scanner.Text()), '\n'))
			if err != nil {
				logrus.Error(err)
			}
		}
		if err := scanner.Err(); err != nil {
			logrus.Error(err)
		}
	},
}

func init() {
	// local flags
	rootCmd.Flags().IntVarP(&logger.MaxSize, "maxsize", "s", 0, "max size (MB) of the log file before it gets rotated (default: 0 == 100)")
	rootCmd.Flags().IntVarP(&logger.MaxAge, "maxage", "a", 0, "max number of 24 hour days to retain old log files after rotation (default: 0 == keep all)")
	rootCmd.Flags().IntVarP(&logger.MaxBackups, "maxbackups", "b", 0, "max number of old log files to retain (default: 0 == keep all)")
	rootCmd.Flags().BoolVarP(&logger.LocalTime, "localtime", "l", false, "if true use local time for timestamps (default: false == use UTC)")
	rootCmd.Flags().BoolVarP(&logger.Compress, "compress", "c", false, "if true compress rotated log files (default: false)")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
