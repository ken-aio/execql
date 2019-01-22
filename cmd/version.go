// Copyright © 2019 @ken-aio <suguru.akiho@gmail.com>
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

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version number
	Version string
	// Revision build revision
	Revision string
	// BuildDate is build date
	BuildDate string
	// GoVersion is build go version
	GoVersion string
)

// newVersionCmd represents the version command
func newVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "All software has versions. This is Goenum's.",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			version := `Go enum generator.
 Version:       %s
 Git commit:    %s
 Go version:    %s
 Build date:    %s
`
			fmt.Printf(version, Version, Revision, GoVersion, BuildDate)
		},
	}
	return cmd
}

func init() {
	rootCmd.AddCommand(newVersionCmd())

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
