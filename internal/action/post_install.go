package action

import (
	"fmt"

	"github.com/craftcms/nitro/validate"
)

func ConfigurePHPMemoryLimit(name, php, limit string) (*Action, error) {
	if err := validate.MachineName(name); err != nil {
		return nil, err
	}
	if err := validate.PHPVersion(php); err != nil {
		return nil, err
	}

	cmd := fmt.Sprintf("s|memory_limit = 128M|memory_limit = %s|g", limit)

	return &Action{
		Type:       "exec",
		UseSyscall: false,
		Args:       []string{"exec", name, "--", "sudo", "sed", "-i", cmd, "/etc/php/" + php + "/fpm/php.ini"},
	}, nil
}

func ConfigurePHPExecutionTimeLimit(name, php, limit string) (*Action, error) {
	if err := validate.MachineName(name); err != nil {
		return nil, err
	}
	if err := validate.PHPVersion(php); err != nil {
		return nil, err
	}

	cmd := fmt.Sprintf("s|max_execution_time = 30|max_execution_time = %s|g", limit)

	return &Action{
		Type:       "exec",
		UseSyscall: false,
		Args:       []string{"exec", name, "--", "sudo", "sed", "-i", cmd, "/etc/php/" + php + "/fpm/php.ini"},
	}, nil
}

func ConfigureXdebug(name, php string) (*Action, error) {
	if err := validate.MachineName(name); err != nil {
		return nil, err
	}
	if err := validate.PHPVersion(php); err != nil {
		return nil, err
	}

	return &Action{
		Type:       "exec",
		UseSyscall: false,
		Args:       []string{"exec", name, "--", "printf", "xdebug.remote_enable=1\nxdebug.remote_connect_back=0\nxdebug.remote_host=localhost\nxdebug.remote_port=9000\nxdebug.remote_log=/var/log/nginx/xdebug.log\n", "|", "sudo", "tee", "-a", "/etc/php/" + php + "/mods-available/xdebug.ini"},
	}, nil
}
