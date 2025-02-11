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

type adminBackupOpts struct {
	authDir   string
	configDir string
	logDir    string

	backupDir string
}

func NewAdminBackupCommand() *cobra.Command {
	opts := &adminBackupOpts{}

	cmd := &cobra.Command{
		Use:   "backup",
		Short: "Backup the database and configuration",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Opening auth database...")
			authStorage := auth.NewStorage(
				auth.OptDirectory(opts.authDir),
				auth.OptReadOnly(true),
			)
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
			logStorage := log.NewStorage(
				log.OptDirectory(opts.logDir),
				log.OptReadOnly(true),
			)
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
			configStorage := config.NewStorage(
				config.OptDirectory(opts.configDir),
			)
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

			fmt.Println("Preparing backup directory...")
			err = os.MkdirAll(opts.backupDir, 0700)
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR: Failed to create backup directory:", err)
				exitCode = 1
				return
			}

			fmt.Println("Backuping auth database...")
			authBackupPath := filepath.Join(opts.backupDir, "auth.db")
			authBackupOut, err := os.Create(authBackupPath)
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR: Failed to create backup file:", err)
				exitCode = 1
				return
			}

			defer func() {
				err := authBackupOut.Close()
				if err != nil {
					fmt.Fprintln(os.Stderr, "ERROR: Failed to close backup file:", err)
					exitCode = 1
				}
			}()

			err = authStorage.Backup(context.Background(), authBackupOut)
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR: Failed to backup auth database:", err)
				exitCode = 1
				return
			}

			fmt.Println("Backuping log database...")
			logBackupPath := filepath.Join(opts.backupDir, "log.db")
			logBackupOut, err := os.Create(logBackupPath)
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR: Failed to create backup file:", err)
				exitCode = 1
				return
			}

			defer func() {
				err := logBackupOut.Close()
				if err != nil {
					fmt.Fprintln(os.Stderr, "ERROR: Failed to close backup file:", err)
					exitCode = 1
				}
			}()

			err = logStorage.Backup(context.Background(), logBackupOut)
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR: Failed to backup log database:", err)
				exitCode = 1
				return
			}

			fmt.Println("Backuping config database...")
			configBackupPath := filepath.Join(opts.backupDir, "config.db")
			configBackupOut, err := os.Create(configBackupPath)
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR: Failed to create backup file:", err)
				exitCode = 1
				return
			}

			defer func() {
				err := configBackupOut.Close()
				if err != nil {
					fmt.Fprintln(os.Stderr, "ERROR: Failed to close backup file:", err)
					exitCode = 1
				}
			}()

			err = configStorage.Backup(context.Background(), configBackupOut)
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR: Failed to backup config database:", err)
				exitCode = 1
				return
			}

			fmt.Println("Backup complete.")
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
