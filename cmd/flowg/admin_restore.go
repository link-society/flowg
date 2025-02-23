package main

import (
	"fmt"

	"context"
	"time"

	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"link-society.com/flowg/internal/utils/proctree"

	"link-society.com/flowg/internal/storage/auth"
	"link-society.com/flowg/internal/storage/config"
	"link-society.com/flowg/internal/storage/log"
)

type adminRestoreOpts struct {
	authDir   string
	configDir string
	logDir    string

	backupDir string
}

func NewAdminRestoreCommand() *cobra.Command {
	opts := &adminRestoreOpts{}

	cmd := &cobra.Command{
		Use:   "restore",
		Short: "Restore the database and configuration from a backup",
		Run: func(cmd *cobra.Command, args []string) {
			authStorage := auth.NewStorage(auth.OptDirectory(opts.authDir))
			logStorage := log.NewStorage(log.OptDirectory(opts.logDir))
			configStorage := config.NewStorage(config.OptDirectory(opts.configDir))

			p := proctree.NewProcessGroup(
				proctree.DefaultProcessGroupOptions(),
				authStorage,
				logStorage,
				configStorage,
			)

			fmt.Fprintln(os.Stderr, "INFO: Opening databases...")
			p.Start()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			err := p.WaitReady(ctx)
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR: Failed to open databases:", err)
				exitCode = 1
				return
			}

			defer func() {
				fmt.Fprintln(os.Stderr, "INFO: Closing databases...")
				p.Stop()
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				err := p.Join(ctx)
				if err != nil {
					fmt.Fprintln(os.Stderr, "ERROR: Failed to close databases:", err)
					exitCode = 1
				}
			}()

			fmt.Println("Restoring auth database...")
			authBackupPath := filepath.Join(opts.backupDir, "auth.db")
			authBackupIn, err := os.Open(authBackupPath)
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR: Failed to open backup file:", err)
				exitCode = 1
				return
			}

			defer func() {
				err := authBackupIn.Close()
				if err != nil {
					fmt.Fprintln(os.Stderr, "ERROR: Failed to close backup file:", err)
					exitCode = 1
				}
			}()

			err = authStorage.Restore(context.Background(), authBackupIn)
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR: Failed to restore auth database:", err)
				exitCode = 1
				return
			}

			fmt.Println("Restoring log database...")
			logBackupPath := filepath.Join(opts.backupDir, "log.db")
			logBackupIn, err := os.Open(logBackupPath)
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR: Failed to open backup file:", err)
				exitCode = 1
				return
			}

			defer func() {
				err := logBackupIn.Close()
				if err != nil {
					fmt.Fprintln(os.Stderr, "ERROR: Failed to close backup file:", err)
					exitCode = 1
				}
			}()

			err = logStorage.Restore(context.Background(), logBackupIn)
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR: Failed to restore log database:", err)
				exitCode = 1
				return
			}

			fmt.Println("Restoring config database...")
			configBackupPath := filepath.Join(opts.backupDir, "config.db")
			configBackupIn, err := os.Open(configBackupPath)
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR: Failed to open backup file:", err)
				exitCode = 1
				return
			}

			defer func() {
				err := configBackupIn.Close()
				if err != nil {
					fmt.Fprintln(os.Stderr, "ERROR: Failed to close backup file:", err)
					exitCode = 1
				}
			}()

			err = configStorage.Restore(context.Background(), configBackupIn)
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR: Failed to restore config database:", err)
				exitCode = 1
				return
			}

			fmt.Println("Restore complete.")
		},
	}

	cmd.Flags().StringVar(
		&opts.authDir,
		"auth-dir",
		defaultAuthDir,
		"Path to the auth database directory",
	)
	cmd.MarkFlagDirname("auth-dir")

	cmd.Flags().StringVar(
		&opts.configDir,
		"config-dir",
		defaultConfigDir,
		"Path to the config directory",
	)
	cmd.MarkFlagDirname("config-dir")

	cmd.Flags().StringVar(
		&opts.logDir,
		"log-dir",
		defaultLogDir,
		"Path to the log database directory",
	)
	cmd.MarkFlagDirname("log-dir")

	cmd.Flags().StringVar(
		&opts.backupDir,
		"backup-dir",
		defaultBackupDir,
		"Path to the backup directory",
	)
	cmd.MarkFlagFilename("backup-dir")

	return cmd
}
