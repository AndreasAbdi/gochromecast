package receiver

import (
	"testing"
)

func TestGetSessionByNamespace(t *testing.T) {
	tests := []struct {
		in     []string
		search string
		out    string
	}{
		{[]string{"fakespace"}, "fakespace", "fakespace"},
	}

	for _, test := range tests {
		status := constructReceiverStatus(test.in)
		session := status.GetSessionByNamespace(test.search)
		for _, namespace := range session.Namespaces {
			if namespace.Name == test.out {
				return
			}
		}
		t.Error("Failed to get a session object with the specified namespace")
	}
}

func constructReceiverStatus(namespaces []string) Status {
	sessions := []*ApplicationSession{}
	for _, namespace := range namespaces {
		sessions = append(sessions, &ApplicationSession{
			Namespaces: []*Namespace{
				&Namespace{Name: namespace},
			},
		})
	}
	receiverStatus := Status{
		Applications: sessions,
	}
	return receiverStatus
}
