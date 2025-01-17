//coverage:ignore
/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var filesUpdatedTotal []string

// ciCmd represents the ci command
var ciCmd = &cobra.Command{
	Use:   "ci",
	Short: "For CI integrations[experimental]",
	Long: `

	Includes features for better CI integrations such as failure when updates available
	for pipelines, allowing users to automatically create PRs when updates are present(custom thresholds) and so on.

Not all those who don't update dependencies are lost.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug().Msgf("Ci stuff... in %s with args count %d", Path, len(args))
		//_, _, directoriesToIgnore, _, _ := getParamsForCheckForUpdatesCMD(cmd.Flags())
		log.Debug().Msg("output format: " + OutputFormat)
		log.Debug().Msgf("Params: Depth=%s, rootDir=%s, Path=%s", strconv.Itoa(Depth), Path, strings.Join(DirectoriesToIgnore, " "))
		rootDir := fixTrailingSlashForPath(Path)
		tf := setupTerraform(rootDir, "1.9.8")
		if tf == nil {
			return
		}
		var modules []map[string]string
		var failureList []map[string]string
		err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
			Check(err, "ci :: command :: ", path)

			Check(err, "checkForUpdates :: command :: ", path)
			isAllowedDir, dirError := directorySearch(rootDir, path, d)
			if errors.Is(dirError, fs.SkipDir) {
				return dirError
			}
			if isAllowedDir {
				err := tf.FormatWrite(context.TODO())
				Check(err, "tf fmt failed")
				filesUpdated := createModuleVersionUpdates(path)
				filesUpdatedTotal = append(filesUpdatedTotal, filesUpdated...)
				filesUpdatedTotal = removeDuplicateStr(filesUpdatedTotal)
				modulesListTotal = append(modulesListTotal, modules...)
				failureListTotal = append(failureListTotal, failureList...)
			}

			return nil
		})
		Check(err, "ci :: command :: unable to walk the directories")
		log.Debug().Msgf("ci :: command :: filesUpdatedTotal :: %s", strings.Join(filesUpdatedTotal, " "))
		//writeCommit(rootDir)

	},
}

func createModuleVersionUpdates(path string) []string {
	files, err := os.ReadDir(fixTrailingSlashForPath(path))
	var filesUpdated []string
	Check(err, "util :: updateTfFiles :: unable to read dir")
	for _, file := range files {
		filesEdited := updateTfFiles(path, file.Name())
		filesUpdated = append(filesUpdated, filesEdited...)
	}
	log.Debug().Msgf("ci :: command :: files :: %s", strings.Join(filesUpdated, " "))
	return filesUpdated
}

func init() {
	cobra.OnInitialize(initConfig)
	checkForUpdatesCmd.AddCommand(ciCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ciCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}
