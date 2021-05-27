package firewallctl

import (
	log "github.com/sirupsen/logrus"
	"os/exec"
	"strings"
)

func ExecIptables(arg ...string) string {
	cmd := exec.Command("iptables", arg...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("exec iptables %s error %s: %s", strings.Join(arg, " "), err.Error(), output)
	}
	return string(output)
}

func AddChain(name string) {
	_ = ExecIptables("-N", name)
}

func RemoveChain(name string) {
	_ = ExecIptables("-F", name)
	_ = ExecIptables("-X", name)
}

func AddJumpChain(sourceChainName string, toChainName string) {
	_ = ExecIptables("-A", sourceChainName, "-j", toChainName)
}

func RemoveJumpChain(sourceChainName string, toChainName string) {
	_ = ExecIptables("-D", sourceChainName, "-j", toChainName)
}

func AddRule(chainName string, ip string, port string, protocol string, rule string) {
	args := []string{"-A", chainName, "-j", rule}
	doRule(args, ip, port, protocol)
}

func RemoveRule(chainName string, ip string, port string, protocol string, rule string) {
	args := []string{"-D", chainName, "-j", rule}
	doRule(args, ip, port, protocol)
}

func doRule(args []string, ip string, port string, protocol string) {
	if protocol == "" {
		protocol = "tcp"
	}

	args = append(args, "-p", protocol)

	if ip != "" {
		args = append(args, "-s", ip)
	}
	if port != "" {
		args = append(args, "--dport", port)
	}
	_ = ExecIptables(args...)
}
