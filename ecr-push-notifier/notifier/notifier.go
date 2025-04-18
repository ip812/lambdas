package notifier

type Notifier interface {
	SendSuccessNotification(repositoryName, imageTag string) error
	SendErrorNotification(repositoryName, imageTag string) error
}
