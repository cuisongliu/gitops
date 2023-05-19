package rollout

import (
	"github.com/cuisongliu/gitops/pkg/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

type DaemonSetRollout struct {
	ide kubernetes.Idempotency
}

func (d *DaemonSetRollout) Rollout(namespace, name string) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		ds, err := d.ide.GetDaemonSet(namespace, name)
		if err != nil {
			return err
		}
		if ds.Spec.Template.Annotations == nil {
			ds.Spec.Template.Annotations = map[string]string{}
		}
		ds.Spec.Template.Annotations["github.com/cuisongliu/tigger-auto-rollout"] = timeString()
		ds.SetResourceVersion("")
		ds.SetUID("")
		return d.ide.CreateOrUpdateDaemonSet(ds)
	})
}
