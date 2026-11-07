package axmanifest

import (
	"os"
	"path/filepath"
	"testing"
)

func load(t *testing.T, content string) *ApplicationManifestSchema {
	t.Helper()
	path := filepath.Join(t.TempDir(), "manifest.json")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	m, err := LoadManifest(path)
	if err != nil {
		t.Fatal(err)
	}
	return m
}

// The schema 2.2.0 manifest of the consume-scene-metadata example
// (AxisCommunications/acap-native-sdk-examples, device-data-hub).
func TestLoadManifestSchemaV2(t *testing.T) {
	m := load(t, `{
	  "schemaVersion": "2.2.0",
	  "acapPackageConf": {
	    "setup": {
	      "appName": "consume_scene_metadata",
	      "vendor": "Axis Communications",
	      "vendorId": "0123456789",
	      "architecture": "aarch64",
	      "runMode": "never",
	      "version": "1.0.0",
	      "compatibleOsVersions": [{"max": "13"}, {"min": "12.10", "max": "12"}]
	    }
	  },
	  "resources": {
	    "deviceDataHub": {"enabled": true},
	    "deepLearningProcessor": {"enabled": true, "required": false}
	  }
	}`)
	setup := m.ACAPPackageConf.Setup
	if m.SchemaVersion != "2.2.0" || setup.VendorID != "0123456789" || setup.Architecture != "aarch64" {
		t.Fatalf("%+v", m)
	}
	if len(setup.CompatibleOsVersions) != 2 || setup.CompatibleOsVersions[0] != (OsVersionRange{Max: "13"}) ||
		setup.CompatibleOsVersions[1] != (OsVersionRange{Min: "12.10", Max: "12"}) {
		t.Fatalf("compatibleOsVersions %+v", setup.CompatibleOsVersions)
	}
	if m.Resources.DeviceDataHub == nil || !m.Resources.DeviceDataHub.Enabled {
		t.Fatalf("deviceDataHub %+v", m.Resources.DeviceDataHub)
	}
	if m.Resources.DeepLearningProcessor == nil || !m.Resources.DeepLearningProcessor.Enabled || m.Resources.DeepLearningProcessor.Required {
		t.Fatalf("deepLearningProcessor %+v", m.Resources.DeepLearningProcessor)
	}
}

func TestLoadManifestSchemaV1Unchanged(t *testing.T) {
	m := load(t, `{
	  "schemaVersion": "1.7.0",
	  "acapPackageConf": {
	    "setup": {"appName": "app", "vendor": "v", "embeddedSdkVersion": "3.0", "runMode": "respawn", "version": "1.0.0",
	              "user": {"username": "sdk", "group": "sdk"}}
	  }
	}`)
	if m.ACAPPackageConf.Setup.EmbeddedSdkVersion != "3.0" || m.ACAPPackageConf.Setup.VendorID != "" ||
		m.Resources.DeviceDataHub != nil || m.Resources.DeepLearningProcessor != nil {
		t.Fatalf("%+v", m)
	}
}
