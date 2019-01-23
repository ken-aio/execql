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
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	funk "github.com/thoas/go-funk"
	"golang.org/x/sync/errgroup"
)

// Option is command option
type Option struct {
	Host       string
	Port       int
	User       string
	Password   string
	CQLFile    string `validate:"required"`
	Keyspace   string `validate:"required"`
	Timeout    int
	NumConns   int
	NumThreads int
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = newRootCmd()

func newRootCmd() *cobra.Command {
	o := &Option{}
	cmd := &cobra.Command{
		Use:   "execql",
		Short: "execute cql command",
		Long:  ``,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			errs := validateParams(*o)
			if len(errs) == 0 {
				return nil
			}
			messages := make([]string, len(errs))
			for i, err := range errs {
				text := ""
				switch err.Field() {
				case "CQLFile":
					text = "-f or --file"
				case "Keyspace":
					text = "-k or --keyspace"
				default:
					text = fmt.Sprintf("unknown error field: %s", err.Field())
				}
				messages[i] = validationErrorToText(err, text)
			}
			message := fmt.Sprintf("\n%s\n", strings.Join(messages, "\n"))
			return errors.New(message)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRootCmd(o)
		},
	}

	cmd.Flags().StringVarP(&o.Host, "host", "H", "localhost", "cassandra host. split ',' if many host. e.g.) cassandra01, cassandra02")
	cmd.Flags().IntVarP(&o.Port, "port", "P", 9042, "cassandra port")
	cmd.Flags().StringVarP(&o.User, "user", "u", "", "connection user")
	cmd.Flags().StringVarP(&o.Password, "password", "p", "", "connection password")
	cmd.Flags().StringVarP(&o.CQLFile, "file", "f", "", "cql file path (required)")
	cmd.Flags().StringVarP(&o.Keyspace, "keyspace", "k", "", "exec target keyspace (required)")
	cmd.Flags().IntVar(&o.Timeout, "timeout", 60000, "query timeout(ms)")
	cmd.Flags().IntVarP(&o.NumConns, "num-conns", "n", 10, "connection nums")
	cmd.Flags().IntVarP(&o.NumThreads, "thread", "t", 1, "concurrent query request thread num")

	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
}

func runRootCmd(o *Option) error {
	log.Printf("Reading input cql file... %s\n", o.CQLFile)
	cqls, err := readCQLs(o.CQLFile)
	if err != nil {
		return err
	}
	log.Printf("Complete reading input cql file\n")

	log.Printf("Creating cassandra session...\n")
	sess, err := createSession(o)
	if err != nil {
		return errors.Wrap(err, "create cassandra session error")
	}
	log.Printf("Complete creating cassandra session\n")

	log.Printf("Execute CQL...\n")
	stopCh := make(chan struct{})
	eg := errgroup.Group{}
	chunkedCQLs := funk.Chunk(cqls, (len(cqls)/o.NumThreads)+1).([][]string)
	for i, chunkedCQL := range chunkedCQLs {
		targets := chunkedCQL
		threadNum := i
		eg.Go(func() error {
			return execCQLs(sess, targets, threadNum, stopCh)
		})
	}
	if err := eg.Wait(); err != nil {
		fmt.Printf("err = %+v\n", err)
		close(stopCh)
		return err
	}
	log.Printf("Complete execute CQL\n")
	return nil
}

func readCQLs(path string) ([]string, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "read cql file error: %s", path)
	}
	cql := string(f)
	return strings.Split(cql, ";"), nil
}

func createSession(o *Option) (*gocql.Session, error) {
	cluster := gocql.NewCluster(strings.Split(o.Host, ",")...)
	cluster.Keyspace = o.Keyspace
	cluster.Timeout = time.Duration(o.Timeout) * time.Millisecond
	cluster.NumConns = o.NumConns
	if o.User != "" && o.Password != "" {
		cluster.Authenticator = &gocql.PasswordAuthenticator{Username: o.User, Password: o.Password}
	}
	sess, err := cluster.CreateSession()
	if err != nil {
		return nil, errors.Wrapf(err, "initialize session error")
	}
	return sess, nil
}

func trimCQL(cql string) string {
	cql = strings.Trim(cql, "\n")
	cql = strings.Trim(cql, "\r")
	return cql
}

func execCQLs(sess *gocql.Session, cqls []string, threadNum int, stopCh chan struct{}) error {
	log.Printf("start thread#%d	/ cql num is %d\n", threadNum, len(cqls))
	for _, cql := range cqls {
		cql = trimCQL(cql)
		if cql == "" {
			continue
		}
		if err := sess.Query(cql).Exec(); err != nil {
			return errors.Wrapf(err, "execute cql error. cql: %s", cql)
		}

		select {
		case <-stopCh:
			return errors.New("stop execute cql")
		default:
		}
	}
	log.Printf("Complete thread#%d\n", threadNum)
	return nil
}
