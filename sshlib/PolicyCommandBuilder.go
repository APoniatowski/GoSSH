package sshlib

import "strings"

func (policies *musthavepolicies) policyCommandBuilder(chosenOption string) string {
	policiesCommand := strings.Builder{}
	switch chosenOption {
	case "check":
		policy := strings.Split(policies.polimport, "/")
		policiesCommand.WriteString("semanage export -f /tmp/policy-check.mod;")
		policiesCommand.WriteString("diff /tmp/policy-check.mod /tmp/")
		policiesCommand.WriteString(policy[len(policy)])
	case "apply":
		policiesCommand.WriteString("semanage -i ")
		policy := strings.Split(policies.polimport, "/")
		policiesCommand.WriteString("/tmp/")
		policiesCommand.WriteString(policy[len(policy)])
	default:
		policiesCommand.WriteString("")
	}

	return policiesCommand.String()
}
