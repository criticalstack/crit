package node

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/apimachinery/pkg/util/wait"
	clientset "k8s.io/client-go/kubernetes"
)

func TaintExists(node *v1.Node, taint v1.Taint) bool {
	for _, t := range node.Spec.Taints {
		if t.Key == taint.Key && t.Value == taint.Value {
			return true
		}
	}
	return false
}

func AddTaint(n *v1.Node, t v1.Taint) {
	if TaintExists(n, t) {
		return
	}
	n.Spec.Taints = append(n.Spec.Taints, t)
}

// PatchNodeOnce executes patchFn on the node object found by the node name.
// This is a condition function meant to be used with wait.Poll. false, nil
// implies it is safe to try again, an error indicates no more tries should be
// made and true indicates success.
//
// Copied from github.com/kubernetes/kubernetes/cmd/kubeadm/app/util/apiclient
func PatchNodeOnce(ctx context.Context, client clientset.Interface, nodeName string, patchFn func(*v1.Node)) func() (bool, error) {
	return func() (bool, error) {
		// First get the node object
		n, err := client.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
		if err != nil {
			// TODO this should only be for timeouts
			return false, nil
		}

		// The node may appear to have no labels at first,
		// so we wait for it to get hostname label.
		if _, found := n.ObjectMeta.Labels[v1.LabelHostname]; !found {
			return false, nil
		}

		oldData, err := json.Marshal(n)
		if err != nil {
			return false, errors.Wrapf(err, "failed to marshal unmodified node %q into JSON", n.Name)
		}

		// Execute the mutating function
		patchFn(n)

		newData, err := json.Marshal(n)
		if err != nil {
			return false, errors.Wrapf(err, "failed to marshal modified node %q into JSON", n.Name)
		}

		patchBytes, err := strategicpatch.CreateTwoWayMergePatch(oldData, newData, v1.Node{})
		if err != nil {
			return false, errors.Wrap(err, "failed to create two way merge patch")
		}

		if _, err := client.CoreV1().Nodes().Patch(ctx, n.Name, types.StrategicMergePatchType, patchBytes, metav1.PatchOptions{}); err != nil {
			// TODO also check for timeouts
			if apierrors.IsConflict(err) {
				fmt.Println("Temporarily unable to update node metadata due to conflict (will retry)")
				return false, nil
			}
			return false, errors.Wrapf(err, "error patching node %q through apiserver", n.Name)
		}

		return true, nil
	}
}

const (
	APICallRetryInterval = 500 * time.Millisecond
	PatchNodeTimeout     = 2 * time.Minute
)

// PatchNode tries to patch a node using patchFn for the actual mutating logic.
// Retries are provided by the wait package.
//
// Copied from github.com/kubernetes/kubernetes/cmd/kubeadm/app/util/apiclient
func PatchNode(ctx context.Context, client clientset.Interface, nodeName string, patchFn func(*v1.Node)) error {
	// wait.Poll will rerun the condition function every interval function if
	// the function returns false. If the condition function returns an error
	// then the retries end and the error is returned.
	return wait.Poll(APICallRetryInterval, PatchNodeTimeout, PatchNodeOnce(ctx, client, nodeName, patchFn))
}

func PatchNodeWithContext(ctx context.Context, client clientset.Interface, nodeName string, patchFn func(*v1.Node)) error {
	return wait.PollImmediateUntil(APICallRetryInterval, PatchNodeOnce(ctx, client, nodeName, patchFn), ctx.Done())
}
