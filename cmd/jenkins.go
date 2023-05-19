/*
Copyright 2022 cuisongliu@qq.com.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"github.com/cuisongliu/gitops/pkg/jenkins"
	"github.com/spf13/cobra"
)

// jenkinsCmd represents the jenkins command
var jenkinsCmd = &cobra.Command{
	Use: "jenkins",
	RunE: func(cmd *cobra.Command, args []string) error {
		return jenkins.Do()
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return jenkins.Validate()
	},
}

func init() {
	rootCmd.AddCommand(jenkinsCmd)
	jenkinsCmd.Flags().StringVarP(&jenkins.Args.JobName, "job", "j", "", "jenkins job name")
	jenkinsCmd.Flags().StringVarP(&jenkins.Args.JobParams, "params", "p", "", "jenkins job params")
	jenkinsCmd.Flags().BoolVarP(&jenkins.Args.Wait, "wait", "w", true, "jenkins job wait success")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// jenkinsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// jenkinsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
