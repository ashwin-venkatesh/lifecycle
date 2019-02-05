package restorer

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/buildpack/lifecycle"
	"github.com/buildpack/lifecycle/cmd"
	"github.com/buildpack/lifecycle/image"
	"log"
	"os"
)

var (
	cacheImageTag string
	layersDir     string
	groupPath     string
)

func init() {
	cmd.FlagLayersDir(&layersDir)
	cmd.FlagGroupPath(&groupPath)
}

func main() {
	flag.Parse()
	if flag.NArg() > 1 || flag.Arg(0) == "" {
		args := map[string]interface{}{"narg": flag.NArg(), "layersDir": layersDir}
		cmd.Exit(cmd.FailCode(cmd.CodeInvalidArgs, "parse arguments", fmt.Sprintf("%+v", args)))
	}
	cacheImageTag = flag.Arg(0)
	cmd.Exit(restore())
}

func restore() error {
	var group lifecycle.BuildpackGroup
	if _, err := toml.DecodeFile(groupPath, &group); err != nil {
		return cmd.FailErr(err, "read group")
	}

	restorer := &lifecycle.Restorer{
		LayersDir:  layersDir,
		Buildpacks: group.Buildpacks,
		Out:        log.New(os.Stdout, "", log.LstdFlags),
		Err:        log.New(os.Stderr, "", log.LstdFlags),
	}

	factory, err := image.DefaultFactory()
	if err != nil {
		return err
	}

	cacheImage, err := factory.NewLocal(cacheImageTag, false)
	if err != nil {
		return err
	}

	restorer.Restore(cacheImage)
	return nil
}
