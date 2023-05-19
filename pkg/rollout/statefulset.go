package rollout

import (
	"github.com/cuisongliu/gitops/pkg/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

type StatefulSetRollout struct {
	ide kubernetes.Idempotency
}

func (d *StatefulSetRollout) Rollout(namespace, name string) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		sts, err := d.ide.GetStatefulSet(namespace, name)
		if err != nil {
			return err
		}
		if sts.Spec.Template.Annotations == nil {
			sts.Spec.Template.Annotations = map[string]string{}
		}
		sts.Spec.Template.Annotations["tigger-auto-rollout"] = timeString()
		sts.SetResourceVersion("")
		sts.SetUID("")
		return d.ide.CreateOrUpdateStatefulSet(sts)
	})
}
