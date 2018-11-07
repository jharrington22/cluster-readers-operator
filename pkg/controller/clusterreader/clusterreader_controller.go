package clusterreader

import (
	"context"
	"log"
	"reflect"

	clusterreaderv1alpha1 "github.com/jharrington22/cluster-readers/pkg/apis/clusterreader/v1alpha1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new ClusterReader Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileClusterReader{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("clusterreader-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource ClusterReader
	err = c.Watch(&source.Kind{Type: &clusterreaderv1alpha1.ClusterReader{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner ClusterReader
	err = c.Watch(&source.Kind{Type: &rbacv1.ClusterRoleBinding{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &clusterreaderv1alpha1.ClusterReader{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileClusterReader{}

// ReconcileClusterReader reconciles a ClusterReader object
type ReconcileClusterReader struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a ClusterReader object and makes changes based on the state read
// and what is in the ClusterReader.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileClusterReader) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log.Printf("Reconciling ClusterReader %s/%s\n", request.Namespace, request.Name)

	// Fetch the ClusterReader instance
	instance := &clusterreaderv1alpha1.ClusterReader{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// print out my allowed readers
	readers := instance.Spec.Readers

	log.Printf("Readers: %#v\n", readers)

	log.Println("1")

	rbacBindingList := &rbacv1.ClusterRoleBindingList{}

	log.Println("2")

	listOptions := &client.ListOptions{}

	log.Println("3")

	r.client.List(context.TODO(), listOptions, rbacBindingList)

	log.Println("4")

	clusterRoleBindingExists := roleBindingInList(instance.Name, rbacBindingList)

	clusterRoleBinding := createClusterRoleBinding(instance)

	// TODO clean this up clusterRoleBindingExists should just be the existing resource or a new one
	if clusterRoleBindingExists {
		log.Printf("Role Binding %s exist!\n", instance.Name)
		clusterRoleBinding = getClusterRoleBinding(instance.Name, rbacBindingList)
	} else {
		log.Printf("Role Binding %s doesn't exist!\n", instance.Name)
	}

	if err != nil {
		log.Println("Error retriving cluster role bindings")
	}

	log.Println("5")

	if err := controllerutil.SetControllerReference(instance, clusterRoleBinding, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	if !clusterRoleBindingExists {
		err := r.client.Create(context.TODO(), clusterRoleBinding)
		if err != nil {
			log.Printf("Error creating role binding %s\n", clusterRoleBinding.Name)
		}
		log.Printf("Created cluster role binding %s\n", clusterRoleBinding.Name)
	}

	if clusterRoleBindingExists {
		if verifyClusterRoleBindingUsers(instance, clusterRoleBinding) {
			log.Println("Users are the same all good!")
		} else {
			log.Println("Error list has been modified re create!")
		}
	}

	return reconcile.Result{}, nil
}

// createRBACGroup
func createRBACGroup(cr *clusterreaderv1alpha1.ClusterReader) {
}

func roleBindingInList(name string, list *rbacv1.ClusterRoleBindingList) bool {
	for _, binding := range list.Items {
		if name == binding.Name {
			return true
		}
	}
	return false
}

func createSubjects(cr *clusterreaderv1alpha1.ClusterReader) []rbacv1.Subject {
	var subjects []rbacv1.Subject
	for _, name := range cr.Spec.Readers {
		subject := rbacv1.Subject{
			Kind:      "User",
			Namespace: cr.Namespace,
			Name:      name,
			APIGroup:  "rbac.authorization.k8s.io",
		}
		subjects = append(subjects, subject)
	}
	return subjects
}

// createClusterRoleBinding
func createClusterRoleBinding(cr *clusterreaderv1alpha1.ClusterReader) *rbacv1.ClusterRoleBinding {
	subjects := createSubjects(cr)

	labels := map[string]string{
		"app": cr.Name,
	}
	return &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Subjects: subjects,
		RoleRef: rbacv1.RoleRef{
			Kind: "ClusterRole",
			Name: "cluster-reader",
		},
	}
}

// verifyClusterRoleBindingUsers verify users in the CR are the only users in the ClusterRoleBinding
func verifyClusterRoleBindingUsers(cr *clusterreaderv1alpha1.ClusterReader, clusterRoleBinding *rbacv1.ClusterRoleBinding) bool {
	var clusterRoleBindingUsers []string
	for _, subject := range clusterRoleBinding.Subjects {
		clusterRoleBindingUsers = append(clusterRoleBindingUsers, subject.Name)
	}
	if reflect.DeepEqual(cr.Spec.Readers, clusterRoleBindingUsers) {
		log.Printf("clusterRoleBindingUsers: %v", clusterRoleBindingUsers)
		log.Printf("clusterReaderUsers: %v", cr.Spec.Readers)
		return true
	}
	log.Printf("clusterRoleBindingUsers: %v", clusterRoleBindingUsers)
	log.Printf("clusterReaderUsers: %v", cr.Spec.Readers)
	return false
}

// getClusterRoleBinding
func getClusterRoleBinding(name string, clusterRoleBindingList *rbacv1.ClusterRoleBindingList) *rbacv1.ClusterRoleBinding {
	var binding rbacv1.ClusterRoleBinding
	for _, binding := range clusterRoleBindingList.Items {
		if name == binding.Name {
			return &binding
		}
	}
	// TODO return an actual error
	return &binding
}
