package questions

import (
	"strings"

	"github.com/Squwid/bytegolf/runner"
)

// Check checks to see if the response is the expected response from the question
func (q *Question) Check(resp *runner.CodeResponse) bool {
	if strings.TrimSpace(strings.ToLower(resp.Output)) == strings.TrimSpace(strings.ToLower(q.Answer)) {
		return true
	}
	return false
}
