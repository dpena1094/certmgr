package svcmgr

import (
	"os/exec"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	shellBinary    string
	canCheckSyntax bool
)

type commandManager struct {
	command string
}

func (cm commandManager) TakeAction(changeType string, specPath string, caPath string, certPath string, keyPath string) error {
	env := []string{
		"CERTMGR_CA_PATH=" + caPath,
		"CERTMGR_CERT_PATH=" + certPath,
		"CERTMGR_KEY_PATH=" + keyPath,
		"CERTMGR_SPEC_PATH=" + specPath,
		"CERTMGR_CHANGE_TYPE=" + changeType,
	}
	return runEnv(env, shellBinary, "-c", cm.command)
}

func newCommandManager(config *Options) (Manager, error) {
	if config.Service != "" {
		log.Warningf("svcmgr 'command': service '%s' for action '%s' doesn't do anything, ignoring", config.Service, config.Action)
	}
	if canCheckSyntax {
		log.Debugf("svcmgr 'command': validating the action definition %s", config.Action)
		err := runEnv([]string{}, shellBinary, "-n", "-c", config.Action)
		if err != nil {
			return nil, errors.WithMessagef(err, "action %s failed bash -nc parse checks", config.Action)
		}
	} else {
		log.Warningf("svcmgr 'command': skipping parse check for '%s' since bash couldn't be found", config.Action)
	}
	return &commandManager{
		command: config.Action,
	}, nil
}

func init() {
	// prefer bash if we can find it since it allows us to validate
	var err error
	shellBinary, err = exec.LookPath("bash")
	canCheckSyntax = true
	if err != nil {
		log.Infof("svcmgr 'command' couldn't find a bash binary; action statements will not be able to be validated for syntax: err %s", err)
		shellBinary, err = exec.LookPath("sh")
		if err != nil {
			log.Error("svcmgr 'command' is unavailable due to both bash and sh not being found in PATH")
			return
		}
		canCheckSyntax = false
	}
	SupportedBackends["command"] = newCommandManager
}
