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

package rollout

import (
	"github.com/cuisongliu/logger"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

var Args struct {
	Namespace  string
	KubeConfig string
}

func Root() {
	if Args.KubeConfig != "" {
		if os.Getenv(clientcmd.RecommendedConfigPathEnvVar) == "" {
			_ = os.Setenv(clientcmd.RecommendedConfigPathEnvVar, Args.KubeConfig)
			defer func() {
				_ = os.Unsetenv(clientcmd.RecommendedConfigPathEnvVar)
			}()
		} else {
			logger.Info("env %s is not empty, env override flag.", clientcmd.RecommendedConfigPathEnvVar)
		}
	}
	configAccess := clientcmd.NewDefaultPathOptions()
	kubeconfigContents, err := configAccess.GetStartingConfig()
	if err != nil {
		logger.Error("get  kubeconfig default cluster: %+v", err)
		os.Exit(-1)
	}
	currentContext := kubeconfigContents.CurrentContext
	if Args.Namespace == "" {
		Args.Namespace = kubeconfigContents.Contexts[currentContext].Namespace
	}
}
