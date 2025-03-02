package tweaks

type Tweak struct {
	Name              string
	Desc              string
	Callback          func() error
	SelectedByDefault bool
	SupportedVersions []int
}

var Tweaks = []Tweak{
	// dnf
	maxParallelDownloads,
	// flatpak
	removeFedoraRemote,
	addFlathubRemote,
	rpmToFlathub,
	preferedPackagingFormat,
	// system
	i2cKernelModules,
	selinuxAllowExecheap,
	installSystemdContainer,
	vmMaxMapCount,
	nonFreeP7zip,
}
