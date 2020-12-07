package database

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/craftcms/nitro/internal/helpers"
	"github.com/craftcms/nitro/labels"
	"github.com/craftcms/nitro/terminal"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

var backupExampleText = `  # backup a database
  nitro db backup`

// backupCommand is the command for backing up an individual database or
func backupCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "backup",
		Short:   "Backup a database",
		Example: backupExampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			env := cmd.Flag("environment").Value.String()
			ctx := cmd.Context()

			// add filters to show only the envrionment and database containers
			filter := filters.NewArgs()
			filter.Add("label", labels.Environment+"="+env)
			filter.Add("label", labels.Type+"=database")

			// get a list of all the databases
			containers, err := docker.ContainerList(cmd.Context(), types.ContainerListOptions{Filters: filter})
			if err != nil {
				return err
			}

			// TODO(jasonmccallister) prompt the user for the container to import
			var containerID, containerCompatability, containerName string
			for _, c := range containers {
				containerID = c.ID
				containerCompatability = c.Labels[labels.DatabaseCompatability]
				containerName = strings.TrimLeft(c.Names[0], "/")
			}

			// create a backup with the current timestamp
			backup := fmt.Sprintf("nitro-backup-%d.sql", time.Now().Unix())

			// create the backup command based on the compatability type
			var backupCmd []string
			switch containerCompatability {
			// TODO(jasonmccallister) add mysql backup
			case "postgres":
				backupCmd = []string{"pg_dump", "-Unitro", "-f", "/tmp/" + backup}
			}

			output.Pending("creating backup", backup)

			// create the command and pass to exec
			exec, err := docker.ContainerExecCreate(ctx, containerID, types.ExecConfig{
				AttachStdout: true,
				AttachStderr: true,
				Tty:          false,
				Cmd:          backupCmd,
			})
			if err != nil {
				return err
			}

			// attach to the container
			stream, err := docker.ContainerExecAttach(ctx, exec.ID, types.ExecConfig{
				AttachStdout: true,
				AttachStderr: true,
				Tty:          false,
				Cmd:          backupCmd,
			})
			if err != nil {
				return err
			}
			defer stream.Close()

			// start the exec
			if err := docker.ContainerExecStart(ctx, exec.ID, types.ExecStartCheck{}); err != nil {
				return fmt.Errorf("unable to start the container, %w", err)
			}

			// wait for the container exec to complete
			waiting := true
			for {
				if waiting {
					resp, err := docker.ContainerExecInspect(ctx, exec.ID)
					if err != nil {
						return err
					}

					waiting = resp.Running
				} else {
					break
				}
			}

			// copy the backup file from the container
			rdr, stat, err := docker.CopyFromContainer(ctx, containerID, "/tmp/"+backup)
			if err != nil || stat.Mode.IsRegular() == false {
				return err
			}
			defer rdr.Close()

			// read the content of the file
			buf := new(bytes.Buffer)
			buf.ReadFrom(rdr)

			// make the backup directory if it does not exist
			backupDir := filepath.Join(home, ".nitro", "backups", env, containerName)
			if err := helpers.MkdirIfNotExists(backupDir); err != nil {
				return err
			}

			// write the file to the backups dir
			if err := ioutil.WriteFile(filepath.Join(backupDir, backup), buf.Bytes(), 0644); err != nil {
				return err
			}

			output.Done()

			output.Info("Backup saved 💾\n =>", filepath.Join(backupDir, backup))

			return nil
		},
	}

	return cmd
}
