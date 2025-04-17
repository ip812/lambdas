package notifier

type Notifier interface {
	SendSuccessNotification(repositoryName, imageTag, imageDigest string) error
	SendErrorNotification(repositoryName, imageTag, imageDigest string) error
}
