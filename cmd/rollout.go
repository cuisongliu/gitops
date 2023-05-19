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
	"github.com/cuisongliu/gitops/pkg/rollout"
	"github.com/cuisongliu/logger"
	"github.com/spf13/cobra"
)

func newRollout() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "rollout",
		Short: "auto rollout some resource for kubernetes",
		PreRun: func(cmd *cobra.Command, args []string) {
			rollout.Root()
		},
	}
	cmd.AddCommand(rolloutDeployCmd)
	cmd.AddCommand(rolloutStsCmd)
	cmd.AddCommand(rolloutDsCmd)
	cmd.PersistentFlags().StringVarP(&rollout.Args.Namespace, "namespace", "n", "", "If present, the namespace scope for this CLI request")
	cmd.PersistentFlags().StringVar(&rollout.Args.KubeConfig, "kubeconfig", "", "Path to the kubeconfig file to use for CLI requests.")
	return cmd
}

var rolloutDeployCmd = &cobra.Command{
	Use:     "deployment",
	Aliases: []string{"deploy"},
	Short:   "auto rollout deployment",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rol := rollout.NewDeployment(rollout.Args.KubeConfig)
		if rol == nil {
			logger.Error("new deployment rollout failed")
			return
		}
		err := rol.Rollout(rollout.Args.Namespace, args[0])
		if err != nil {
			logger.Error("rollout deployment failed: %s", err)
			return
		}
		logger.Info("rollout deployment success")
	},
}

var rolloutStsCmd = &cobra.Command{
	Use:     "statefulset",
	Aliases: []string{"sts"},
	Args:    cobra.ExactArgs(1),
	Short:   "auto rollout statefulsets",
	Run: func(cmd *cobra.Command, args []string) {
		rol := rollout.NewStatefulSet(rollout.Args.KubeConfig)
		if rol == nil {
			logger.Error("new sts rollout failed")
			return
		}
		err := rol.Rollout(rollout.Args.Namespace, args[0])
		if err != nil {
			logger.Error("rollout sts failed: %s", err)
			return
		}
		logger.Info("rollout statefulset success")
	},
}

var rolloutDsCmd = &cobra.Command{
	Use:     "daemonset",
	Aliases: []string{"ds"},
	Args:    cobra.ExactArgs(1),
	Short:   "auto rollout daemonsets",
	Run: func(cmd *cobra.Command, args []string) {
		rol := rollout.NewDaemonSet(rollout.Args.KubeConfig)
		if rol == nil {
			logger.Error("new ds rollout failed")
			return
		}
		err := rol.Rollout(rollout.Args.Namespace, args[0])
		if err != nil {
			logger.Error("rollout ds failed: %s", err)
			return
		}
		logger.Info("rollout daemonset success")
	},
}

func init() {
	rootCmd.AddCommand(newRollout())

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// kubeRolloutCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// kubeRolloutCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
