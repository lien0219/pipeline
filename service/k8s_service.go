package service

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// K8sService Kubernetes 服务结构体
type K8sService struct {
	ClientSet *kubernetes.Clientset
}

var K8sSvc *K8sService

// CreateJob 创建 Job
func (s *K8sService) CreateJob(namespace string, job *batchv1.Job) (*batchv1.Job, error) {
	return s.ClientSet.BatchV1().Jobs(namespace).Create(context.TODO(), job, metav1.CreateOptions{})
}

// GetJob 获取 Job
func (s *K8sService) GetJob(namespace, name string) (*batchv1.Job, error) {
	return s.ClientSet.BatchV1().Jobs(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

// UpdateJob 更新 Job
func (s *K8sService) UpdateJob(namespace string, job *batchv1.Job) (*batchv1.Job, error) {
	return s.ClientSet.BatchV1().Jobs(namespace).Update(context.TODO(), job, metav1.UpdateOptions{})
}

// DeleteJob 删除 Job
func (s *K8sService) DeleteJob(namespace, name string) error {
	return s.ClientSet.BatchV1().Jobs(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

// CreateCronJob 创建 CronJob
func (s *K8sService) CreateCronJob(namespace string, cronJob *batchv1beta1.CronJob) (*batchv1beta1.CronJob, error) {
	return s.ClientSet.BatchV1beta1().CronJobs(namespace).Create(context.TODO(), cronJob, metav1.CreateOptions{})
}

// GetCronJob 获取 CronJob
func (s *K8sService) GetCronJob(namespace, name string) (*batchv1beta1.CronJob, error) {
	return s.ClientSet.BatchV1beta1().CronJobs(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

// UpdateCronJob 更新 CronJob
func (s *K8sService) UpdateCronJob(namespace string, cronJob *batchv1beta1.CronJob) (*batchv1beta1.CronJob, error) {
	return s.ClientSet.BatchV1beta1().CronJobs(namespace).Update(context.TODO(), cronJob, metav1.UpdateOptions{})
}

// DeleteCronJob 删除 CronJob
func (s *K8sService) DeleteCronJob(namespace, name string) error {
	return s.ClientSet.BatchV1beta1().CronJobs(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

// CreateDeployment 创建 Deployment
func (s *K8sService) CreateDeployment(namespace string, deployment *appsv1.Deployment) (*appsv1.Deployment, error) {
	return s.ClientSet.AppsV1().Deployments(namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
}

// GetDeployment 获取 Deployment
func (s *K8sService) GetDeployment(namespace, name string) (*appsv1.Deployment, error) {
	return s.ClientSet.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

// UpdateDeployment 更新 Deployment
func (s *K8sService) UpdateDeployment(namespace string, deployment *appsv1.Deployment) (*appsv1.Deployment, error) {
	return s.ClientSet.AppsV1().Deployments(namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
}

// DeleteDeployment 删除 Deployment
func (s *K8sService) DeleteDeployment(namespace, name string) error {
	return s.ClientSet.AppsV1().Deployments(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

// NewK8sService 创建一个新的 K8sService 实例
func NewK8sService(kubeconfigPath string) (*K8sService, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &K8sService{
		ClientSet: clientset,
	}, nil
}
