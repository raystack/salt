package dockertestx

import "runtime"

func DockerHostAddress() string {
	var dockerHostInternal = "host-gateway" // linux by default
	if runtime.GOOS == "darwin" {
		dockerHostInternal = "host.docker.internal"
	}
	return dockerHostInternal
}
