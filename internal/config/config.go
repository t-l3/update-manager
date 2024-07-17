package config

type AppConfig struct {
	Apps                []App  `yaml:"apps"`
	TmpDownloadLocation string `yaml:"tmp-download-location" default:"/tmp/update-manager-download"`
}

type App struct {
	Name              string           `yaml:"name"`
	Icon              string           `yaml:"icon"`
	DownloadUrl       string           `yaml:"download-url"`
	InstallDir        InstallDir       `yaml:"install-dir"`
	VersioningChecks  VersioningChecks `yaml:"versioning-checks"`
	PreInstallScript  string           `yaml:"pre-install-script"`
	PostInstallScript string           `yaml:"post-install-script"`
}

type InstallDir struct {
	Path  string `yaml:"path"`
	Owner string `yaml:"owner"`
	Mode  int    `yaml:"mode" default:"0700"`
}

type VersioningChecks struct {
	Installed string `yaml:"installed"`
	Latest    string `yaml:"latest"`
}
