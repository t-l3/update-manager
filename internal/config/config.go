package config

type AppConfig struct {
	Apps                []Apps `yaml:"apps"`
	TmpDownloadLocation string `yaml:"tmp-download-location"`
}

type Apps struct {
	Name             string           `yaml:"name"`
	DownloadUrl      string           `yaml:"download-url"`
	InstallDir       InstallDir       `yaml:"install-dir"`
	VersioningChecks VersioningChecks `yaml:"versioning-checks"`
}

type InstallDir struct {
	Path  string `yaml:"path"`
	Owner string `yaml:"owner"`
	Mode  string `yaml:"mode"`
}

type VersioningChecks struct {
	Installed string `yaml:"installed"`
	Latest    string `yaml:"latest"`
}
