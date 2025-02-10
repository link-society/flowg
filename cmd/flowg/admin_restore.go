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
			authStorage := auth.NewStorage(auth.OptDirectory(opts.authDir))
			authStorage.Start()
			err := authStorage.WaitStarted()
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR: Failed to open auth database:", err)
				exitCode = 1
				return
			}

			defer func() {
				authStorage.Stop()
				err := authStorage.WaitStopped()
				if err != nil {
					fmt.Fprintln(os.Stderr, "ERROR: Failed to close auth database:", err)
					exitCode = 1
				}
			}()

			logStorage := log.NewStorage(log.OptDirectory(opts.logDir))
			logStorage.Start()
			err = logStorage.WaitStarted()
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR: Failed to open log database:", err)
				exitCode = 1
				return
			}

			defer func() {
				logStorage.Stop()
				err := logStorage.WaitStopped()
				if err != nil {
					fmt.Fprintln(os.Stderr, "ERROR: Failed to close log database:", err)
					exitCode = 1
				}
			}()

			backupConfigStorageDir := filepath.Join(opts.backupDir, "config")
			backupConfigStorage := config.NewStorage(config.OptDirectory(backupConfigStorageDir))
			backupConfigStorage.Start()
			err = backupConfigStorage.WaitStarted()
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR: Failed to open backup config directory:", err)
				exitCode = 1
				return
			}

			defer func() {
				backupConfigStorage.Stop()
				err := backupConfigStorage.WaitStopped()
				if err != nil {
					fmt.Fprintln(os.Stderr, "ERROR: Failed to close backup config directory:", err)
					exitCode = 1
				}
			}()

			configStorage := config.NewStorage(config.OptDirectory(opts.configDir))
			configStorage.Start()
			err = configStorage.WaitStarted()
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR: Failed to open config directory:", err)
				exitCode = 1
				return
			}

			defer func() {
				configStorage.Stop()
				err := configStorage.WaitStopped()
				if err != nil {
					fmt.Fprintln(os.Stderr, "ERROR: Failed to close config directory:", err)
					exitCode = 1
				}
			}()

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

			fmt.Println("Restoring auth database...")
			err = authStorage.Restore(context.Background(), authBackupIn)
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR: Failed to restore auth database:", err)
				exitCode = 1
				return
			}

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

			fmt.Println("Restoring log database...")
			err = logStorage.Restore(context.Background(), logBackupIn)
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR: Failed to restore log database:", err)
				exitCode = 1
				return
			}

			fmt.Println("Restoring configuration...")
			transformers, err := backupConfigStorage.ListTransformers(context.Background())
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR: Failed to list backup transformers:", err)
				exitCode = 1
				return
			}

			pipelines, err := backupConfigStorage.ListPipelines(context.Background())
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR: Failed to list backup pipelines:", err)
				exitCode = 1
				return
			}

			alerts, err := backupConfigStorage.ListAlerts(context.Background())
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR: Failed to list backup alerts:", err)
				exitCode = 1
				return
			}

			for _, transformerName := range transformers {
				transformerContent, err := backupConfigStorage.ReadTransformer(context.Background(), transformerName)
				if err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: Failed to read backup transformer %s: %v\n", transformerName, err)
					exitCode = 1
					return
				}

				err = configStorage.WriteTransformer(context.Background(), transformerName, transformerContent)
				if err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: Failed to restore transformer %s: %v\n", transformerName, err)
					exitCode = 1
					return
				}
			}

			for _, pipelineName := range pipelines {
				pipelineContent, err := backupConfigStorage.ReadPipeline(context.Background(), pipelineName)
				if err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: Failed to read backup pipeline %s: %v\n", pipelineName, err)
					exitCode = 1
					return
				}

				err = configStorage.WritePipeline(context.Background(), pipelineName, pipelineContent)
				if err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: Failed to restore pipeline %s: %v\n", pipelineName, err)
					exitCode = 1
					return
				}
			}

			for _, alertName := range alerts {
				alertContent, err := backupConfigStorage.ReadAlert(context.Background(), alertName)
				if err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: Failed to read backup alert %s: %v\n", alertName, err)
					exitCode = 1
					return
				}

				err = configStorage.WriteAlert(context.Background(), alertName, alertContent)
				if err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: Failed to restore alert %s: %v\n", alertName, err)
					exitCode = 1
					return
				}
			}
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
