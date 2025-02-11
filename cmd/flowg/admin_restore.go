package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

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
			fmt.Println("Opening auth database...")
			authStorage := auth.NewStorage(auth.OptDirectory(opts.authDir))
			authStorage.Start()
			err := authStorage.WaitStarted()
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR: Failed to open auth database:", err)
				exitCode = 1
				return
			}

			defer func() {
				fmt.Println("Closing auth database...")
				authStorage.Stop()
				err := authStorage.WaitStopped()
				if err != nil {
					fmt.Fprintln(os.Stderr, "ERROR: Failed to close auth database:", err)
					exitCode = 1
				}
			}()

			fmt.Println("Opening log database...")
			logStorage := log.NewStorage(log.OptDirectory(opts.logDir))
			logStorage.Start()
			err = logStorage.WaitStarted()
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR: Failed to open log database:", err)
				exitCode = 1
				return
			}

			defer func() {
				fmt.Println("Closing log database...")
				logStorage.Stop()
				err := logStorage.WaitStopped()
				if err != nil {
					fmt.Fprintln(os.Stderr, "ERROR: Failed to close log database:", err)
					exitCode = 1
				}
			}()

			fmt.Println("Opening config database...")
			configStorage := config.NewStorage(config.OptDirectory(opts.configDir))
			configStorage.Start()
			err = configStorage.WaitStarted()
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR: Failed to open config database:", err)
				exitCode = 1
				return
			}

			defer func() {
				fmt.Println("Closing config database...")
				configStorage.Stop()
				err := configStorage.WaitStopped()
				if err != nil {
					fmt.Fprintln(os.Stderr, "ERROR: Failed to close config database:", err)
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
