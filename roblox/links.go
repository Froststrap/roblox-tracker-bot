package roblox

import (
	"fmt"
	"net/url"
	"strings"
)

const (
	setupCdnBaseUrl  = "https://setup.rbxcdn.com"
	rddBaseUrl       = "https://rdd.latte.to/"
	apkMirrorBaseUrl = "https://www.apkmirror.com/apk/roblox-corporation/roblox/"
)

type DownloadLink struct {
	Name string
	Url  string
}

func isProductionChannel(channel string) bool {
	return strings.EqualFold(
		normalizeChannel(channel),
		defaultChannel,
	)
}

// ClientSettingsUrl returns the client-settings API URL for a platform and
// channel. Android is not served by client-settings, so it returns "".
func ClientSettingsUrl(platform Platform, channel string) string {
	if platform == PlatformAndroid {
		return ""
	}

	return clientSettingsCdnBaseUrl +
		buildChannelPath(string(platform), channel)
}

// buildRddUrl returns a rdd.latte.to link that assembles a specific
// production deployment into a downloadable zip in the browser.
func buildRddUrl(platform Platform, versionGuid string) string {
	query := url.Values{}
	query.Set("channel", "LIVE")
	query.Set("binaryType", string(platform))
	query.Set("version", versionGuid)

	return rddBaseUrl + "?" + query.Encode()
}

// buildApkMirrorUrl returns the APKMirror release page for an Android
// version such as "2.737.1584". The page only exists once APKMirror has
// published that version.
func buildApkMirrorUrl(version string) string {
	slug := strings.ReplaceAll(version, ".", "-")

	return fmt.Sprintf(
		"%sroblox-%s-release/",
		apkMirrorBaseUrl,
		slug,
	)
}

// DownloadLinks returns download links for a production deployment.
func DownloadLinks(
	platform Platform,
	deployment Deployment,
) []DownloadLink {
	versionGuid := deployment.ClientVersionUpload

	switch platform {
	case PlatformWindowsPlayer:
		if versionGuid == "" {
			return nil
		}

		return []DownloadLink{
			{
				Name: "rdd.latte.to",
				Url:  buildRddUrl(platform, versionGuid),
			},
		}
	case PlatformMacPlayer:
		if versionGuid == "" {
			return nil
		}

		return []DownloadLink{
			{
				Name: "Apple Silicon",
				Url: fmt.Sprintf(
					"%s/mac/arm64/%s-RobloxPlayer.zip",
					setupCdnBaseUrl,
					versionGuid,
				),
			},
			{
				Name: "Intel",
				Url: fmt.Sprintf(
					"%s/mac/%s-RobloxPlayer.zip",
					setupCdnBaseUrl,
					versionGuid,
				),
			},
		}
	case PlatformAndroid:
		if deployment.Version == "" {
			return nil
		}

		return []DownloadLink{
			{
				Name: "APKMirror",
				Url:  buildApkMirrorUrl(deployment.Version),
			},
		}
	default:
		return nil
	}
}
