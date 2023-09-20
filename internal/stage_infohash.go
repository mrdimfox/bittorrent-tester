package internal

import (
	"fmt"
	"math/rand"
	"os"
	"path"

	tester_utils "github.com/codecrafters-io/tester-utils"
)

func testInfoHash(stageHarness *tester_utils.StageHarness) error {
	initRandom()

	logger := stageHarness.Logger
	executable := stageHarness.Executable

	tempDir, err := os.MkdirTemp("", "torrents")
	if err != nil {
		return err
	}

	shuffled := make([]TestTorrentInfo, len(testTorrents))
	copy(shuffled, testTorrents)
	rand.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })

	for _, torrent := range shuffled {
		if err := copyTorrent(tempDir, torrent.filename); err != nil {
			logger.Errorln("Couldn't copy torrent file")
			return err
		}

		torrentPath := path.Join(tempDir, torrent.filename)

		logger.Infof("Running ./your_bittorrent.sh info %s", torrentPath)
		result, err := executable.Run("info", torrentPath)
		if err != nil {
			return err
		}

		if err = assertExitCode(result, 0); err != nil {
			return err
		}

		expected := fmt.Sprintf("Info Hash: %s", torrent.infohash)

		if err = assertStdoutContains(result, expected); err != nil {
			return err
		}
	}

	return nil
}
