package notifier

type Slack struct {
	Channel string
	Webhook string
}

func (s *Slack) SendSuccessNotification(repositoryName, imageTag, imageDigest string) error {
	// Implement the logic to send a success notification to Slack
	return nil
}

func (s *Slack) SendErrorNotification(repositoryName, imageTag, imageDigest string) error {
	// Implement the logic to send an error notification to Slack
	return nil
}
