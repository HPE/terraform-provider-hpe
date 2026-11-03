#!/usr/bin/env bash
# Rebuilds the skeleton plugin jars used by the hpe_morpheus_plugin acceptance tests.
#
# Two jars are produced from the same class, differing only in the manifest
# Plugin-Version, so the tests can exercise the install path (0.1.0) and the
# re-upload/update path (0.2.0). The jars are committed so the tests need no JDK;
# rerun this script only when SkeletonPlugin.java changes.
#
# Requirements: a JDK (javac/jar) and network access to Maven Central for the
# morpheus-plugin-api jar the class is compiled against.
set -euo pipefail

cd "$(dirname "$0")"

API_VERSION="${MORPHEUS_PLUGIN_API_VERSION:-1.5.0}"
API_JAR="morpheus-plugin-api-${API_VERSION}.jar"
API_URL="https://repo1.maven.org/maven2/com/morpheusdata/morpheus-plugin-api/${API_VERSION}/${API_JAR}"
WORK="$(mktemp -d)"
trap 'rm -rf "$WORK"' EXIT

if [ ! -s "$WORK/$API_JAR" ]; then
  curl -fsSL --retry 3 -o "$WORK/$API_JAR" "$API_URL"
fi

# Java 11 bytecode: the plugin API itself is built for 11 and every supported
# appliance runs a newer JVM.
javac --release 11 -cp "$WORK/$API_JAR" -d "$WORK/classes" SkeletonPlugin.java

for version in 0.1.0 0.2.0; do
  cat > "$WORK/MANIFEST-$version.MF" <<EOF
Manifest-Version: 1.0
Plugin-Class: com.hpe.terraform.skeleton.SkeletonPlugin
Plugin-Version: $version
Morpheus-Name: HPE Terraform Skeleton
Morpheus-Code: hpe-terraform-skeleton-plugin
Morpheus-Organization: hpe
Morpheus-Description: Inert plugin used by the HPE Terraform provider acceptance tests
Morpheus-Labels: Plugin, Test
EOF
  out="hpe-terraform-skeleton-plugin-$version.jar"
  rm -f "$out"
  # The manifest must be the first jar entry: the appliance reads it with
  # JarInputStream.getManifest(), which the jar tool's cfm ordering satisfies.
  jar cfm "$out" "$WORK/MANIFEST-$version.MF" -C "$WORK/classes" .
  printf '%s  %s\n' "$(shasum -a 256 "$out" | cut -d' ' -f1)" "$out"
done
