package rollout

import (
	"github.com/cuisongliu/gitops/pkg/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

type DeploymentRollout struct {
	ide kubernetes.Idempotency
}

func (d *DeploymentRollout) Rollout(namespace, name string) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		deploy, err := d.ide.GetDeployment(namespace, name)
		if err != nil {
			return err
		}
		if deploy.Spec.Template.Annotations == nil {
			deploy.Spec.Template.Annotations = map[string]string{}
		}
		deploy.Spec.Template.Annotations["tigger-auto-rollout"] = timeString()
		deploy.SetResourceVersion("")
		deploy.SetUID("")
		return d.ide.CreateOrUpdateDeployment(deploy)
	})
}
