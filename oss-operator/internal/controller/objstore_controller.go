/*
Copyright 2025.

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

package controller

import (
	"context"
	"fmt"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	cninfossv1alpha1 "oss-operator/api/v1alpha1"
)

const configMapName = "%s-cm"

// ObjStoreReconciler reconciles a ObjStore object
type ObjStoreReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	OssSvc *oss.Client
}

//+kubebuilder:rbac:groups=cninf-oss.oss.fjj.com,resources=objstores,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cninf-oss.oss.fjj.com,resources=objstores/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cninf-oss.oss.fjj.com,resources=objstores/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resource=configmaps,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ObjStore object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.0/pkg/reconcile
func (r *ObjStoreReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	//创建一个objstore 实例
	instance := &cninfossv1alpha1.ObjStore{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	// TODO(user): your logic here
	//如果state 为空，则辅助为pending
	if instance.Status.State == "" {
		instance.Status.State = cninfossv1alpha1.PENDING_STATE
		r.Status().Update(ctx, instance)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ObjStoreReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cninfossv1alpha1.ObjStore{}).
		Complete(r)
}

func (r *ObjStoreReconciler) createResources(ctx context.Context, objStore *cninfossv1alpha1.ObjStore) error {
	//更新状态
	objStore.Status.State = cninfossv1alpha1.CREATING_STATE
	err := r.Status().Update(ctx, objStore)
	if err != nil {
		return err
	}
	//创建bucket
	r.OssSvc.CreateBucket(objStore.Spec.Name)
	if err != nil {
		return err
	}
	//获取bucket 信息
	bucketInfo, err := r.OssSvc.GetBucketInfo(objStore.Spec.Name)
	if err != nil {
		return err
	}

	//创建configmap
	data := make(map[string]string, 0)
	data["bucketName"] = objStore.Spec.Name
	data["endpoint"] = bucketInfo.BucketInfo.ExtranetEndpoint

	configmap := &v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},

		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf(configMapName, objStore.Name),
			Namespace: objStore.Namespace,
		},
		Data: data,
	}
	err = r.Create(ctx, configmap)
	if err != nil {
		return err
	}
	//更新状态为created
	objStore.Status.State = cninfossv1alpha1.CREATED_STATE
	err = r.Status().Update(ctx, objStore)
	if err != nil {
		return err
	}
	return nil
}

func (r *ObjStoreReconciler) deleteResources(ctx context.Context, objStore *cninfossv1alpha1.ObjStore) error {
	//删除bucket
	err := r.OssSvc.DeleteBucket(objStore.Spec.Name)
	if err != nil {
		return err
	}
	//删除configmap
	configmap := &v1.ConfigMap{}
	err = r.Get(ctx,client, client.ObjectKey{
		Name: fmt.Sprintf(configMapname, objStore.Name),
		Namespace: objStore.Namespace,},configmap)
	if err != nil{
		return err
		}
	err = r.Delete(ctx, configmap)
	if err != nil{
		return err
		}
	return nil
}
