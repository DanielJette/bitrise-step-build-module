package uploaders

import (
	"encoding/json"
	"fmt"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-steplib/bitrise-step-build-router-start/androidartifact"
	"github.com/bitrise-steplib/bitrise-step-build-router-start/bundletool"
)

// DeployAAB ...
func DeployAAB(pth string, artifacts []string, buildURL, token, bundletoolVersion string) (ArtifactURLs, error) {
	log.Printf("- analyzing aab")

	// get aab manifest dump
	log.Printf("- fetching info")

	r, err := bundletool.New(bundletoolVersion)
	if err != nil {
		return ArtifactURLs{}, err
	}
	cmd := r.Command("dump", "manifest", "--bundle", pth)

	log.Donef("$ %s", cmd.PrintableCommandArgs())

	out, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return ArtifactURLs{}, err
	}

	packageName, versionCode, versionName := androidartifact.ParsePackageInfos(out)

	appInfo := map[string]interface{}{
		"package_name": packageName,
		"version_code": versionCode,
		"version_name": versionName,
	}

	log.Printf("  aab infos: %v", appInfo)

	if packageName == "" {
		log.Warnf("Package name is undefined, AndroidManifest.xml package content:\n%s", out)
	}

	if versionCode == "" {
		log.Warnf("Version code is undefined, AndroidManifest.xml package content:\n%s", out)
	}

	if versionName == "" {
		log.Warnf("Version name is undefined, AndroidManifest.xml package content:\n%s", out)
	}

	// ---

	fileSize, err := fileSizeInBytes(pth)
	if err != nil {
		return ArtifactURLs{}, fmt.Errorf("failed to get apk size, error: %s", err)
	}

	info := androidartifact.ParseArtifactPath(pth)

	aabInfoMap := map[string]interface{}{
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
		aabInfoMap["apk"] = splitMeta.APK
		aabInfoMap["aab"] = splitMeta.AAB
		aabInfoMap["split"] = splitMeta.Split
		aabInfoMap["universal"] = splitMeta.UniversalApk
	}

	artifactInfoBytes, err := json.Marshal(aabInfoMap)
	if err != nil {
		return ArtifactURLs{}, fmt.Errorf("failed to marshal apk infos, error: %s", err)
	}

	// ---

	uploadURL, artifactID, err := createArtifact(buildURL, token, pth, "android-apk")
	if err != nil {
		return ArtifactURLs{}, fmt.Errorf("failed to create apk artifact, error: %s", err)
	}

	if err := uploadArtifact(uploadURL, pth, "application/octet-stream aab"); err != nil {
		return ArtifactURLs{}, fmt.Errorf("failed to upload apk artifact, error: %s", err)
	}
	artifactURLs, err := finishArtifact(buildURL, token, artifactID, string(artifactInfoBytes), "", "", "false")
	if err != nil {
		return ArtifactURLs{}, fmt.Errorf("failed to finish apk artifact, error: %s", err)
	}

	return artifactURLs, nil
}
