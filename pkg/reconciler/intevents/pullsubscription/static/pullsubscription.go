/*
Copyright 2019 Google LLC

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

package static

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/reconciler"

	v1 "github.com/google/knative-gcp/pkg/apis/intevents/v1"
	pullsubscriptionreconciler "github.com/google/knative-gcp/pkg/client/injection/reconciler/intevents/v1/pullsubscription"
	psreconciler "github.com/google/knative-gcp/pkg/reconciler/intevents/pullsubscription"
)

// Reconciler implements controller.Reconciler for PullSubscription resources.
type Reconciler struct {
	*psreconciler.Base
}

// Check that our Reconciler implements Interface.
var _ pullsubscriptionreconciler.Interface = (*Reconciler)(nil)

func (r *Reconciler) ReconcileKind(ctx context.Context, ps *v1.PullSubscription) reconciler.Event {
	return r.Base.ReconcileKind(ctx, ps)
}

func (r *Reconciler) ReconcileDeployment(ctx context.Context, ra *appsv1.Deployment, src *v1.PullSubscription) error {
	existing, err := r.Base.GetOrCreateReceiveAdapter(ctx, ra, src)
	if err != nil {
		return err
	}
	if !equality.Semantic.DeepEqual(ra.Spec, existing.Spec) {
		existing.Spec = ra.Spec
		existing, err = r.KubeClientSet.AppsV1().Deployments(src.Namespace).Update(ctx, existing, metav1.UpdateOptions{})
		if err != nil {
			src.Status.MarkDeployedFailed("ReceiveAdapterUpdateFailed", "Error updating the Receive Adapter: %s", err.Error())
			logging.FromContext(ctx).Desugar().Error("Error updating Receive Adapter", zap.Error(err))
			return err
		}
	}

	if minimumReplicasUnavailable := src.Status.PropagateDeploymentAvailability(existing); minimumReplicasUnavailable {
		logging.FromContext(ctx).Desugar().Error("minimumreolicasUnavalibel detected", zap.Error(err))
		podList, _ := r.Base.GetPods(ctx, src, ra)
		if podList != nil {
			for _, pod := range podList.Items {
				logging.FromContext(ctx).Desugar().Info("current pod is : " + pod.Name)
				eventList, _ := r.KubeClientSet.CoreV1().Events(pod.Namespace).List(ctx, metav1.ListOptions{FieldSelector: fmt.Sprintf("metadata.name=%s", pod.Name)})
				for _, event := range eventList.Items {
					logging.FromContext(ctx).Desugar().Info("get the event " + event.Name + event.Type)
					if event.Reason == "FailedMount" {
						src.Status.MarkDeployedFailed("AuthenticationCheckFailed", event.Message)
						return nil
					}
				}
				for _, cs := range pod.Status.ContainerStatuses {
					logging.FromContext(ctx).Desugar().Info("current pod's cs is : " + cs.Name)
					logging.FromContext(ctx).Desugar().Info(fmt.Sprintf("current pod's cs is : %v", cs.State))
					if cs.State.Terminated != nil && strings.Contains(cs.State.Terminated.Message, "auth") {
						src.Status.MarkDeployedFailed("AuthenticationCheckFailed", cs.State.Terminated.Message)
						return nil
					} else if cs.LastTerminationState.Terminated != nil && strings.Contains(cs.LastTerminationState.Terminated.Message, "auth") {
						src.Status.MarkDeployedFailed("AuthenticationCheckFailed", cs.LastTerminationState.Terminated.Message)
					}
				}
			}
		}
	}

	return nil
}

func (r *Reconciler) FinalizeKind(ctx context.Context, ps *v1.PullSubscription) reconciler.Event {
	return r.Base.FinalizeKind(ctx, ps)
}
