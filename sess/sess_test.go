package sess

import "testing"

func TestRetreiveSession(t *testing.T) {
	sess, err := retreiveSess("abc123")
	if err != nil {
		t.Fatalf("error getting sess: %v", err)
	}
	t.Logf("sess: %v", sess)
}
