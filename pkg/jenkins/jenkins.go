/*
Copyright 2023 cuisongliu@qq.com.

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

package jenkins

import (
	"context"
	"fmt"
	"github.com/bndr/gojenkins"
	"github.com/cuisongliu/gitops/pkg/utils"
	"github.com/cuisongliu/logger"
	"os"
	"path"
	"time"
)

var jenkinsEnvs struct {
	jenkinsUrl      string
	jenkinsUser     string
	jenkinsPassword string
}

var Args struct {
	JobName   string
	JobParams string
	Wait      bool
}

func Validate() error {
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
	if Args.JobName == "" {
		return fmt.Errorf("jenkins job name is empty")
	}
	return nil
}

func Do() error {
	ctx := context.Background()
	jenkins := gojenkins.CreateJenkins(nil, jenkinsEnvs.jenkinsUrl, jenkinsEnvs.jenkinsUser, jenkinsEnvs.jenkinsPassword)
	_, err := jenkins.Init(ctx)
	if err != nil {
		return err
	}
	logger.Info("jenkins init success")
	job, err := jenkins.GetJob(ctx, Args.JobName)
	if err != nil {
		return err
	}
	logger.Info("jenkins job get success: %+v", job)
	params := utils.StringToMap(Args.JobParams, ",")
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
	if Args.Wait {
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
}
