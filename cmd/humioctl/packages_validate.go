// Copyright © 2018 Humio Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/humio/cli/prompt"

	"github.com/spf13/cobra"
)

func validatePackageCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "validate [flags] <repo-or-view-name> <package-dir>",
		Short: "Validate a package's content.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			out := prompt.NewPrompt(cmd.OutOrStdout())

			repoOrViewName := args[0]
			dirPath := args[1]

			if !filepath.IsAbs(dirPath) {
				var err error
				dirPath, err = filepath.Abs(dirPath)
				if err != nil {
					out.Error(fmt.Sprintf("Invalid path: %s", err))
					os.Exit(1)
				}
				dirPath += "/"
			}

			out.Info(fmt.Sprintf("Validating Package in: %s", dirPath))

			// Get the HTTP client
			client := NewApiClient(cmd)

			validationResult, apiErr := client.Packages().Validate(repoOrViewName, dirPath)
			if apiErr != nil {
				out.Error(fmt.Sprintf("Errors validating package: %s", apiErr))
				os.Exit(1)
			}

			if validationResult.IsValid() {
				out.Info("Package is valid")
			} else {
				out.Error("Package is not valid")
				out.Error(out.List(validationResult.InstallationErrors))
				out.Error(out.List(validationResult.ParseErrors))
				os.Exit(1)
			}
		},
	}

	return &cmd
}
