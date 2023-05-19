package rollout

import (
	"fmt"
	"github.com/cuisongliu/gitops/pkg/client-go/kubernetes"
	"time"
)

type Rollout interface {
	// Rollout restarts the rollout of a resource.
	Rollout(namespace, name string) error
}

func NewDeployment(config string) Rollout {
	ide := kubernetes.NewKubeIdempotency(config)
	if ide != nil {
		return &DeploymentRollout{
			ide: ide,
		}
	}
	return nil
}

func NewStatefulSet(config string) Rollout {
	ide := kubernetes.NewKubeIdempotency(config)
	if ide != nil {
		return &StatefulSetRollout{
			ide: ide,
		}
	}
	return nil
}

func NewDaemonSet(config string) Rollout {
	ide := kubernetes.NewKubeIdempotency(config)
	if ide != nil {
		return &DaemonSetRollout{
			ide: ide,
		}
	}
	return nil
}

func timeString() string {
	return fmt.Sprintf("%s-%s", "timestamp", time.Now().Format("20060102150405"))
}
