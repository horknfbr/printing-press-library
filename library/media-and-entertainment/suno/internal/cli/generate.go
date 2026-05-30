// Copyright 2026 horknfbr. Licensed under Apache-2.0. See LICENSE.

package cli

import (
	"github.com/spf13/cobra"
)

// newGenerateCmd is the parent for music/lyrics/video generation. The parent is
// Hidden so it stays out of the top-level --help list (matching the clips and
// billing groups), while every subcommand still resolves normally.
func newGenerateCmd(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "generate",
		Short:  "Music, lyrics, and video generation jobs",
		Hidden: true,
		RunE:   parentNoSubcommandRunE(flags),
	}

	// Hand-authored, captcha-aware generation/transform commands.
	cmd.AddCommand(newSunoGenerateCreateCmd(flags))
	cmd.AddCommand(newSunoDescribeCmd(flags))
	cmd.AddCommand(newSunoExtendCmd(flags))
	cmd.AddCommand(newSunoCoverCmd(flags))
	cmd.AddCommand(newSunoRemasterCmd(flags))

	// Spec endpoint subcommands.
	cmd.AddCommand(newGenerateConcatCmd(flags))
	cmd.AddCommand(newGenerateLyricsCmd(flags))
	cmd.AddCommand(newGenerateLyricsStatusCmd(flags))
	cmd.AddCommand(newGenerateVideoStatusCmd(flags))
	return cmd
}
