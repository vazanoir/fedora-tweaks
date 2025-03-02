package tweaks

type Tweak struct {
	Name              string
	Desc              string
	SelectedByDefault bool
	SupportedVersions []int
	Callback          func() error
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
