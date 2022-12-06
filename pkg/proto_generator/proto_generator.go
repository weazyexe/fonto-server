package proto_generator

import (
	"fmt"
	"github.com/weazyexe/fonto-server/pkg/logger"
	"os"
	"os/exec"
	"strings"
)

func Run(configPath string) {
	logger.InitializeLogger()

	config, err := ReadConfig(configPath)
	if err != nil {
		logger.Zap.Fatalf("Error while reading config file: %v\n%v", configPath, err)
	}

	files, err := os.ReadDir(config.SourcesDirectory)
	if err != nil {
		logger.Zap.Fatalf("Error while reading directory %v", config.SourcesDirectory)
	}

	for _, file := range files {
		logger.Zap.Infof("Compiling %v", file.Name())
		nameWithoutExtension := strings.Split(file.Name(), ".")[0]

		createFolderIfNeeds(nameWithoutExtension, config.GenerateDirectory)
		command := composeProtocCommand(
			config.GenerateDirectory,
			config.SourcesDirectory,
			file.Name(),
			nameWithoutExtension,
		)

		if err := command.Run(); err != nil {
			logger.Zap.Fatalf("Error while compiling %v\n%v", file.Name(), err)
		}
	}

	logger.Zap.Info("Proto has been compiled")
}

func createFolderIfNeeds(name, where string) {
	if _, err := os.Stat(where + "/" + name); os.IsNotExist(err) {
		mkdirCommand := exec.Command("mkdir", name)
		mkdirCommand.Dir = where
		if err := mkdirCommand.Run(); err != nil {
			logger.Zap.Fatalf("Can't create new folder: %v\n%v", name, err)
		}
	}
}

func composeProtocCommand(generateDir, protoSourcesDir, fileName, name string) *exec.Cmd {
	outArg := fmt.Sprintf("../%v/%v", generateDir, name)
	optArg := "paths=source_relative"

	command := exec.Command(
		"protoc",
		"--go_out="+outArg,
		"--go_opt="+optArg,
		"--go-grpc_out="+outArg,
		"--go-grpc_opt="+optArg,
		fileName,
	)
	command.Dir = protoSourcesDir
	return command
}
