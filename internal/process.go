package internal

import (
	"strings"

	"github.com/charmbracelet/log"
	"github.com/shirou/gopsutil/v3/process"
)

// KillClaudeProcesses finds and terminates any running instances of Claude Code
// (e.g., 'claude', 'claude-agent-sdk', 'claude-agent-acp') to prevent them from
// refreshing tokens and corrupting the keychain after an account switch.
func KillClaudeProcesses() {
	procs, err := process.Processes()
	if err != nil {
		log.Warn("unable to list processes to terminate claude", "error", err)
		return
	}

	killedCount := 0
	for _, p := range procs {
		name, err := p.Name()
		if err != nil {
			continue
		}

		// Some processes might be invoked via node, so we can also check the command line
		cmdline, err := p.Cmdline()
		if err != nil {
			continue
		}

		isClaude := false
		if strings.Contains(strings.ToLower(name), "claude") {
			isClaude = true
		} else if strings.Contains(cmdline, "claude-agent") || strings.Contains(cmdline, "claude --") {
			isClaude = true
		}

		if isClaude {
			log.Debug("found running claude process, terminating", "pid", p.Pid, "name", name)
			if err := p.Kill(); err != nil {
				log.Warn("failed to kill claude process", "pid", p.Pid, "error", err)
			} else {
				killedCount++
			}
		}
	}

	if killedCount > 0 {
		log.Info("terminated running claude processes to prevent token corruption", "count", killedCount)
	}
}
