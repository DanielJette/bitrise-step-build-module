package uploaders

import (
	"encoding/json"
	"fmt"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-steplib/bitrise-step-build-router-start/androidartifact"
)

// DeployAPK ...
func DeployAPK(pth string, artifacts []string, buildURL, token, notifyUserGroups, notifyEmails, isEnablePublicPage string) (ArtifactURLs, error) {
	log.Printf("analyzing apk")

	apkInfo, err := androidartifact.GetAPKInfo(pth)
	if err != nil {
		return ArtifactURLs{}, err
	}

	appInfo := map[string]interface{}{
		"app_name":        apkInfo.AppName,
		"package_name":    apkInfo.PackageName,
		"version_code":    apkInfo.VersionCode,
		"version_name":    apkInfo.VersionName,
		"min_sdk_version": apkInfo.MinSDKVersion,
	}

	log.Printf("  apk infos: %v", appInfo)

	if apkInfo.PackageName == "" {
		log.Warnf("Package name is undefined, AndroidManifest.xml package content:\n%s", apkInfo.RawPackageContent)
	}

	if apkInfo.VersionCode == "" {
		log.Warnf("Version code is undefined, AndroidManifest.xml package content:\n%s", apkInfo.RawPackageContent)
	}

	if apkInfo.VersionName == "" {
		log.Warnf("Version name is undefined, AndroidManifest.xml package content:\n%s", apkInfo.RawPackageContent)
	}

	// ---

	fileSize, err := fileSizeInBytes(pth)
	if err != nil {
		return ArtifactURLs{}, fmt.Errorf("failed to get apk size, error: %s", err)
	}

	info := androidartifact.ParseArtifactPath(pth)

	apkInfoMap := map[string]interface{}{
		"file_size_bytes": fmt.Sprintf("%f", fileSize),
		"app_info":        appInfo,
		"module":          info.Module,
		"product_flavour": info.ProductFlavour,
		"build_type":      info.BuildType,
	}
	splitMeta, err := androidartifact.CreateSplitArtifactMeta(pth, artifacts)
	if err != nil {
		log.Errorf("Failed to create split meta, error: %s", err)
	} else {
		apkInfoMap["aab"] = splitMeta.AAB
		apkInfoMap["apk"] = splitMeta.APK
		apkInfoMap["split"] = splitMeta.Split
		apkInfoMap["universal"] = splitMeta.UniversalApk
	}

	artifactInfoBytes, err := json.Marshal(apkInfoMap)
	if err != nil {
		return ArtifactURLs{}, fmt.Errorf("failed to marshal apk infos, error: %s", err)
	}

	// ---

	uploadURL, artifactID, err := createArtifact(buildURL, token, pth, "android-apk")
	if err != nil {
		return ArtifactURLs{}, fmt.Errorf("failed to create apk artifact, error: %s", err)
	}

	if err := uploadArtifact(uploadURL, pth, "application/vnd.android.package-archive"); err != nil {
		return ArtifactURLs{}, fmt.Errorf("failed to upload apk artifact, error: %s", err)
	}

	artifactURLs, err := finishArtifact(buildURL, token, artifactID, string(artifactInfoBytes), notifyUserGroups, notifyEmails, isEnablePublicPage)
	if err != nil {
		return ArtifactURLs{}, fmt.Errorf("failed to finish apk artifact, error: %s", err)
	}

	return artifactURLs, nil
}
