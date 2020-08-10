package cmd

import (
	"errors"
	"github.com/GoodwayGroup/gwsm/s3"
	"github.com/urfave/cli/v2"
)

// S3Get will perform a simple pull of the requested object to the destination path.
func S3Get(c *cli.Context) error {
	if c.NArg() > 2 {
		_ = cli.ShowSubcommandHelp(c)
		return cli.NewExitError(errors.New("ERROR too many arguments passed"), 2)
	}

	src, dest := c.Args().Get(0), c.Args().Get(1)

	err := s3.Get(src, dest)
	if err != nil {
		return cli.NewExitError(err, 2)
	}

	return nil
}
