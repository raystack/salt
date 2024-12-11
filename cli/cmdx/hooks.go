package cmdx

// AddClientHooks applies all configured hooks to commands annotated with `client:true`.
func (m *Commander) AddClientHooks() {
	for _, cmd := range m.RootCmd.Commands() {
		for _, hook := range m.Hooks {
			if cmd.Annotations["client"] == "true" {
				hook.Behavior(cmd)
			}
		}
	}
}
