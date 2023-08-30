package internal

import (
	"os"
	"path"

	tester_utils "github.com/codecrafters-io/tester-utils"
)

func testDownloadFile(stageHarness *tester_utils.StageHarness) error {
	initRandom()

	logger := stageHarness.Logger
	executable := stageHarness.Executable

	t := randomTorrent()

	tempDir, err := os.MkdirTemp("", "torrents")
	if err != nil {
		logger.Errorf("Couldn't create temp directory")
		return err
	}

	if err := copyTorrent(tempDir, t.filename); err != nil {
		logger.Errorf("Couldn't copy torrent file")
		return err
	}

	torrentFilePath := path.Join(tempDir, t.filename)
	downloadedFilePath := path.Join(tempDir, t.outputFilename)

	logger.Infof("Running ./your_bittorrent.sh download -o %s %s", downloadedFilePath, torrentFilePath)
	result, err := executable.Run("download", "-o", downloadedFilePath, torrentFilePath)
	if err != nil {
		return err
	}

	if err = assertExitCode(result, 0); err != nil {
		return err
	}

	if err = assertFileSize(downloadedFilePath, t.length); err != nil {
		return err
	}

	if err = assertFileSHA1(downloadedFilePath, t.expectedSha1); err != nil {
		return err
	}

	return nil
}
