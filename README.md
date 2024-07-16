# update-manager

A configurable app to manage installs, and updates, for apps that are only released for direct download and not available in package manager repos on all distros.

For example, I use it for Discord, VS Code and Google Chrome which aren't available in Arch repositories and don't auto update.

update-manager will check the locally installed version against the remote latest available version, and initiate a download if there is a difference.

## Config

`update-manager --config /path/to/config.yaml`

### /etc/config-manager/config.yaml (Default)

```yaml
tmp-download-location: [string] File path to use as a temporary download location.

apps:
  - name: [string] A display name for the app, for use in logs.
    download-url: [string] URL to download the latest release of the application
    install-dir:
      path: [string] File path to install app to locally
      owner: [string] Linux file owner string for top level install directory and sub-files (i.e root or root:root)
      mode: [string] Linux file mode representation for top level install directory (i.e 755)
    versioning-checks:
      installed: [string] An executable command to check the locally installed version (i.e jq '.version' .discord/resources/build_info.json)
      latest: [string] An executable command to check the remotely latest available version (i.e curl -s 'https://discord.com/api/download?platform=linux&format=tar.gz' | grep -oE '[0-9]+\.[0-9]+\.[0-9]+' | head -n 1;)
```