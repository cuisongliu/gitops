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
	"context"
	"fmt"
	"github.com/bndr/gojenkins"
	"github.com/cuisongliu/logger"
	"github.com/spf13/cobra"
	"os"
	"path"
	"strings"
	"time"
)

var jenkinsJobName string
var jenkinsJobParams string
var jenkinsWait bool
var jenkinsEnvs struct {
	jenkinsUrl      string
	jenkinsUser     string
	jenkinsPassword string
}

// jenkinsCmd represents the jenkins command
var jenkinsCmd = &cobra.Command{
	Use: "jenkins",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		jenkins := gojenkins.CreateJenkins(nil, jenkinsEnvs.jenkinsUrl, jenkinsEnvs.jenkinsUser, jenkinsEnvs.jenkinsPassword)
		_, err := jenkins.Init(ctx)
		if err != nil {
			return err
		}
		logger.Info("jenkins init success")
		job, err := jenkins.GetJob(ctx, jenkinsJobName)
		if err != nil {
			return err
		}
		logger.Info("jenkins job get success: %+v", job)
		params := StringToMap(jenkinsJobParams, ",")
		queueid, err := job.InvokeSimple(ctx, params)
		if err != nil {
			return err
		}
		build, err := jenkins.GetBuildFromQueueID(ctx, queueid)
		if err != nil {
			return err
		}
		logger.Info("jenkins job build success")
		pipelineURL := path.Join(build.GetUrl(), "console")
		if jenkinsWait {
			for build.IsRunning(ctx) {
				logger.Info("jenkins job is running,wait 10s to check")
				time.Sleep(10 * time.Second)
				build.Poll(ctx)
			}
			logger.Info("jenkins job is finished")
		} else {
			logger.Warn("jenkins skip wait job build")
		}
		switch build.GetResult() {
		case "SUCCESS":
			logger.Info("jenkins job build success, build number: %d, console url: %s", build.GetBuildNumber(), pipelineURL)
		case "FAILURE":
			logger.Error("jenkins job build failure, build number: %d, console url: %s", build.GetBuildNumber(), pipelineURL)
		default:
			logger.Warn("jenkins job build waiting, build number: %d, console url: %s", build.GetBuildNumber(), pipelineURL)
		}
		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if s, _ := os.LookupEnv("JENKINS_URL"); s == "" {
			return fmt.Errorf("JENKINS_URL is empty")
		} else {
			jenkinsEnvs.jenkinsUrl = s
		}
		if s, _ := os.LookupEnv("JENKINS_USER"); s == "" {
			return fmt.Errorf("JENKINS_USER is empty")
		} else {
			jenkinsEnvs.jenkinsUser = s
		}
		if s, _ := os.LookupEnv("JENKINS_PASSWORD"); s == "" {
			return fmt.Errorf("JENKINS_PASSWORD is empty")
		} else {
			jenkinsEnvs.jenkinsPassword = s
		}
		if jenkinsJobName == "" {
			return fmt.Errorf("jenkins job name is empty")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(jenkinsCmd)
	jenkinsCmd.Flags().StringVarP(&jenkinsJobName, "job", "j", "", "jenkins job name")
	jenkinsCmd.Flags().StringVarP(&jenkinsJobParams, "params", "p", "", "jenkins job params")
	jenkinsCmd.Flags().BoolVarP(&jenkinsWait, "wait", "w", true, "jenkins job wait success")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// jenkinsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// jenkinsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func StringToMap(data string, spilt string) map[string]string {
	list := strings.Split(data, spilt)
	return ListToMap(list)
}

func ListToMap(data []string) map[string]string {
	m := make(map[string]string)
	for _, l := range data {
		if l != "" {
			kv := strings.SplitN(l, "=", 2)
			if len(kv) == 2 {
				m[kv[0]] = kv[1]
			}
		}
	}
	return m
}
