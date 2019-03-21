package deploy

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/peterq/pan-light/qt/tool-chain/cmd/rcc"
	"github.com/peterq/pan-light/qt/tool-chain/utils"
)

//linux

func linux_sh(target, name string) string {
	bb := new(bytes.Buffer)
	defer bb.Reset()

	fmt.Fprint(bb, "#!/bin/bash\n")
	fmt.Fprint(bb, "appname=`basename $0 | sed s,\\.sh$,,`\n\n")
	fmt.Fprint(bb, "dirname=`dirname $0`\n")
	fmt.Fprint(bb, "tmp=\"${dirname#?}\"\n\n")
	fmt.Fprint(bb, "if [ \"${dirname%$tmp}\" != \"/\" ]; then\n")
	fmt.Fprint(bb, "dirname=$PWD/$dirname\n")
	fmt.Fprint(bb, "fi\n")

	if strings.HasPrefix(target, "rpi") {
		fmt.Fprint(bb, "export DISPLAY=\":0\"\n")
		fmt.Fprint(bb, "export LD_PRELOAD=\"/opt/vc/lib/libGLESv2.so /opt/vc/lib/libEGL.so\"\n")
	}

	if utils.QT_PKG_CONFIG() {
		libDir := strings.TrimSpace(utils.RunCmd(exec.Command("pkg-config", "--variable=libdir", "Qt5Core"), fmt.Sprintf("get lib dir for %v on %v", target, runtime.GOOS)))
		miscDir := utils.QT_MISC_DIR()

		fmt.Fprintf(bb, "export LD_LIBRARY_PATH=\"%v\"\n", libDir)
		fmt.Fprintf(bb, "export QT_PLUGIN_PATH=\"%v\"\n", filepath.Join(miscDir, "plugins"))
		fmt.Fprintf(bb, "export QML_IMPORT_PATH=\"%v\"\n", filepath.Join(miscDir, "qml"))
		fmt.Fprintf(bb, "export QML2_IMPORT_PATH=\"%v\"\n", filepath.Join(miscDir, "qml"))
	} else {
		libDir := "lib"
		if name == libDir {
			libDir = "libs"
		}
		fmt.Fprintf(bb, "export LD_LIBRARY_PATH=\"$dirname/%v\"\n", libDir)
		fmt.Fprint(bb, "export QT_PLUGIN_PATH=\"$dirname/plugins\"\n")
		fmt.Fprint(bb, "export QML_IMPORT_PATH=\"$dirname/qml\"\n")
		fmt.Fprint(bb, "export QML2_IMPORT_PATH=\"$dirname/qml\"\n")
	}
	fmt.Fprint(bb, "$dirname/$appname \"$@\"\n")

	return bb.String()
}

//android

func android_config(target, path, depPath string) string {
	jsonStruct := &struct {
		Qt                            string `json:"qt"`
		Sdk                           string `json:"sdk"`
		SdkBuildToolsRevision         string `json:"sdkBuildToolsRevision"`
		Ndk                           string `json:"ndk"`
		Toolchainprefix               string `json:"toolchain-prefix"`
		Toolprefix                    string `json:"tool-prefix"`
		Toolchainversion              string `json:"toolchain-version"`
		Ndkhost                       string `json:"ndk-host"`
		Targetarchitecture            string `json:"target-architecture"`
		AndroidExtraLibs              string `json:"android-extra-libs"`
		AndroidPackageSourceDirectory string `json:"android-package-source-directory"`
		Qmlrootpath                   string `json:"qml-root-path"`
		StdcppPath                    string `json:"stdcpp-path"`
		Applicationbinary             string `json:"application-binary"`
	}{
		Qt:                            filepath.Join(utils.QT_DIR(), utils.QT_VERSION_MAJOR(), "android_armv7"),
		Sdk:                           utils.ANDROID_SDK_DIR(),
		SdkBuildToolsRevision:         "28.0.3",
		Ndk:                           utils.ANDROID_NDK_DIR(),
		Toolchainprefix:               "arm-linux-androideabi",
		Toolprefix:                    "arm-linux-androideabi",
		Toolchainversion:              "4.9",
		Ndkhost:                       runtime.GOOS + "-x86_64",
		Targetarchitecture:            "armeabi-v7a",
		AndroidExtraLibs:              filepath.Join(depPath, "libgo_base.so"),
		AndroidPackageSourceDirectory: filepath.Join(path, target),
		Qmlrootpath:                   path,
		StdcppPath:                    filepath.Join(utils.ANDROID_NDK_DIR(), "sources", "cxx-stl", "llvm-libc++", "libs", "armeabi-v7a", "libc++_shared.so"),
		Applicationbinary:             filepath.Join(depPath, "libgo.so"),
	}

	if target == "android-emulator" {
		jsonStruct.Qt = filepath.Join(utils.QT_DIR(), utils.QT_VERSION_MAJOR(), "android_x86")
		jsonStruct.Toolchainprefix = "x86"
		jsonStruct.Toolprefix = "i686-linux-android"
		jsonStruct.Targetarchitecture = "x86"
		jsonStruct.StdcppPath = filepath.Join(utils.ANDROID_NDK_DIR(), "sources", "cxx-stl", "llvm-libc++", "libs", jsonStruct.Targetarchitecture, "libc++_shared.so")
	}

	if utils.QT_DOCKER() {
		switch target {
		case "android":
			jsonStruct.AndroidExtraLibs += "," + filepath.Join(os.Getenv("HOME"), "openssl-1.0.2q-arm", "libcrypto.so") + "," + filepath.Join(os.Getenv("HOME"), "openssl-1.0.2q-arm", "libssl.so")
		case "android-emulator":
			jsonStruct.AndroidExtraLibs += "," + filepath.Join(os.Getenv("HOME"), "openssl-1.0.2q-x86", "libcrypto.so") + "," + filepath.Join(os.Getenv("HOME"), "openssl-1.0.2q-x86", "libssl.so")
		}
	}

	out, err := json.Marshal(jsonStruct)
	if err != nil {
		utils.Log.WithError(err).Panicf("failed to create json-config file for androiddeployqt on %v", runtime.GOOS)
	}
	return strings.Replace(string(out), `\\`, `/`, -1)
}

//darwin

func darwin_plist(name string) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleExecutable</key>
	<string>%[1]v</string>
	<key>CFBundleGetInfoString</key>
	<string>Created by Qt/QMake</string>
	<key>CFBundleIconFile</key>
	<string>%[1]v.icns</string>
	<key>CFBundleIdentifier</key>
	<string>com.yourcompany.%[1]v</string>
	<key>CFBundleName</key>
	<string>%[1]v</string>
	<key>CFBundlePackageType</key>
	<string>APPL</string>
	<key>CFBundleShortVersionString</key>
	<string>1.0.0</string>
	<key>CFBundleSignature</key>
	<string>????</string>
	<key>LSMinimumSystemVersion</key>
	<string>10.11</string>
	<key>NOTE</key>
	<string>This file was generated by Qt/QMake.</string>
	<key>NSPrincipalClass</key>
	<string>NSApplication</string>
	<key>NSHighResolutionCapable</key>
	<true/>
	<key>NSSupportsAutomaticGraphicsSwitching</key>
	<true/>
</dict>
</plist>
`, name)
}

func darwin_pkginfo() string {
	return "APPL????\n"
}

func darwin_nix_script(name string) string {
	return fmt.Sprintf(`#!/bin/bash
export PATH=$HOME/.nix-profile/bin:$PATH
cd "${0%%/*}"
./%v_bin
`, name)
}

//ios

func ios_c_main_wrapper() string {
	bb := new(bytes.Buffer)
	bb.WriteString("#include \"libgo.h\"\n")
	for _, n := range rcc.ResourceNames {
		fmt.Fprintf(bb, "extern int qInitResources_%v();\n", n)
	}
	bb.WriteString("int main(int argc, char *argv[]) {\n")
	for _, n := range rcc.ResourceNames {
		fmt.Fprintf(bb, "qInitResources_%v();\n", n)
	}
	bb.WriteString("go_main_wrapper(argc, argv);\n}")
	return bb.String()
}

func ios_plist(name string) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleDisplayName</key>
	<string>%[1]v</string>
	<key>CFBundleExecutable</key>
	<string>main</string>
	<key>CFBundleGetInfoString</key>
	<string>Created by Qt/QMake</string>
	<key>CFBundleIdentifier</key>
	<string>%[2]v</string>
	<key>CFBundleName</key>
	<string>%[1]v</string>
	<key>CFBundlePackageType</key>
	<string>APPL</string>
	<key>CFBundleShortVersionString</key>
	<string>1.0.0</string>
	<key>CFBundleSignature</key>
	<string>????</string>
	<key>CFBundleVersion</key>
	<string>1.0.0</string>
	<key>LSRequiresIPhoneOS</key>
	<true/>
	<key>MinimumOSVersion</key>
	<string>${IPHONEOS_DEPLOYMENT_TARGET}</string>
	<key>NOTE</key>
	<string>This file was generated by Qt/QMake.</string>
	<key>UILaunchStoryboardName</key>
	<string>LaunchScreen</string>
	<key>UISupportedInterfaceOrientations</key>
	<array>
		<string>UIInterfaceOrientationPortrait</string>
		<string>UIInterfaceOrientationPortraitUpsideDown</string>
		<string>UIInterfaceOrientationLandscapeLeft</string>
		<string>UIInterfaceOrientationLandscapeRight</string>
	</array>
	<key>QtRunLoopIntegrationDisableSeparateStack</key>
	<true/>
</dict>
</plist>
`, name, strings.Replace(name, "_", "", -1))
}

func ios_launchscreen(name string) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="no"?>
	<document type="com.apple.InterfaceBuilder3.CocoaTouch.XIB" version="3.0" toolsVersion="6250" systemVersion="14A343f" targetRuntime="iOS.CocoaTouch" propertyAccessControl="none" useAutolayout="YES" launchScreen="YES" useTraitCollections="YES">
	    <dependencies>
	        <plugIn identifier="com.apple.InterfaceBuilder.IBCocoaTouchPlugin" version="6244"/>
	        <capability name="Constraints with non-1.0 multipliers" minToolsVersion="5.1"/>
	    </dependencies>
	    <objects>
	        <placeholder placeholderIdentifier="IBFilesOwner" id="-1" userLabel="File's Owner"/>
	        <placeholder placeholderIdentifier="IBFirstResponder" id="-2" customClass="UIResponder"/>
	        <view contentMode="scaleToFill" id="iN0-l3-epB">
	            <rect key="frame" x="0.0" y="0.0" width="480" height="480"/>
	            <autoresizingMask key="autoresizingMask" widthSizable="YES" heightSizable="YES"/>
	            <subviews>
	                <label opaque="NO" clipsSubviews="YES" userInteractionEnabled="NO" contentMode="left" horizontalHuggingPriority="251" verticalHuggingPriority="251" text="%v" textAlignment="center" lineBreakMode="middleTruncation" baselineAdjustment="alignBaselines" minimumFontSize="18" translatesAutoresizingMaskIntoConstraints="NO" id="kId-c2-rCX">
	                    <rect key="frame" x="20" y="140" width="441" height="43"/>
	                    <fontDescription key="fontDescription" type="boldSystem" pointSize="36"/>
	                    <color key="textColor" cocoaTouchSystemColor="darkTextColor"/>
	                    <nil key="highlightedColor"/>
	                </label>
	            </subviews>
	            <color key="backgroundColor" white="1" alpha="1" colorSpace="custom" customColorSpace="calibratedWhite"/>
	            <constraints>
	                <constraint firstItem="kId-c2-rCX" firstAttribute="centerY" secondItem="iN0-l3-epB" secondAttribute="bottom" multiplier="1/3" constant="1" id="Kid-kn-2rF"/>
	                <constraint firstAttribute="centerX" secondItem="kId-c2-rCX" secondAttribute="centerX" id="Koa-jz-hwk"/>
	                <constraint firstItem="kId-c2-rCX" firstAttribute="leading" secondItem="iN0-l3-epB" secondAttribute="leading" constant="20" symbolic="YES" id="fvb-Df-36g"/>
	            </constraints>
	            <nil key="simulatedStatusBarMetrics"/>
	            <freeformSimulatedSizeMetrics key="simulatedDestinationMetrics"/>
	            <point key="canvasLocation" x="404" y="445"/>
	        </view>
	    </objects>
	</document>
	`, name)
}

func ios_appicon() string {
	return `{
  "images" : [
    {
      "idiom" : "iphone",
      "size" : "29x29",
      "scale" : "2x"
    },
    {
      "idiom" : "iphone",
      "size" : "29x29",
      "scale" : "3x"
    },
    {
      "idiom" : "iphone",
      "size" : "40x40",
      "scale" : "2x"
    },
    {
      "idiom" : "iphone",
      "size" : "40x40",
      "scale" : "3x"
    },
    {
      "idiom" : "iphone",
      "size" : "60x60",
      "scale" : "2x"
    },
    {
      "idiom" : "iphone",
      "size" : "60x60",
      "scale" : "3x"
    },
    {
      "idiom" : "ipad",
      "size" : "29x29",
      "scale" : "1x"
    },
    {
      "idiom" : "ipad",
      "size" : "29x29",
      "scale" : "2x"
    },
    {
      "idiom" : "ipad",
      "size" : "40x40",
      "scale" : "1x"
    },
    {
      "idiom" : "ipad",
      "size" : "40x40",
      "scale" : "2x"
    },
    {
      "idiom" : "ipad",
      "size" : "76x76",
      "scale" : "1x"
    },
    {
      "idiom" : "ipad",
      "size" : "76x76",
      "scale" : "2x"
    }
  ],
  "info" : {
    "version" : 1,
    "author" : "xcode"
  }
}
`
}

func ios_xcodeproject() string {
	return `// !$*UTF8*$!
{
	archiveVersion = 1;
	classes = {
	};
	objectVersion = 46;
	objects = {

/* Begin PBXBuildFile section */
		254BB84F1B1FD08900C56DE9 /* Images.xcassets in Resources */ = {isa = PBXBuildFile; fileRef = 254BB84E1B1FD08900C56DE9 /* Images.xcassets */; };
		254BB8681B1FD16500C56DE9 /* main in Resources */ = {isa = PBXBuildFile; fileRef = 254BB8671B1FD16500C56DE9 /* main */; };
		25916F411CE65FF600695115 /* LaunchScreen.xib in Resources */ = {isa = PBXBuildFile; fileRef = 25916F401CE65FF600695115 /* LaunchScreen.xib */; };
		25F26AED1CE6675E0045FFBA /* Default-568h@2x.png in Resources */ = {isa = PBXBuildFile; fileRef = 25F26AEC1CE6675E0045FFBA /* Default-568h@2x.png */; };
/* End PBXBuildFile section */

/* Begin PBXFileReference section */
		254BB83E1B1FD08900C56DE9 /* main.app */ = {isa = PBXFileReference; explicitFileType = wrapper.application; includeInIndex = 0; path = main.app; sourceTree = BUILT_PRODUCTS_DIR; };
		254BB8421B1FD08900C56DE9 /* Info.plist */ = {isa = PBXFileReference; lastKnownFileType = text.plist.xml; path = Info.plist; sourceTree = "<group>"; };
		254BB84E1B1FD08900C56DE9 /* Images.xcassets */ = {isa = PBXFileReference; lastKnownFileType = folder.assetcatalog; path = Images.xcassets; sourceTree = "<group>"; };
		254BB8671B1FD16500C56DE9 /* main */ = {isa = PBXFileReference; lastKnownFileType = "compiled.mach-o.executable"; path = main; sourceTree = "<group>"; };
		25916F401CE65FF600695115 /* LaunchScreen.xib */ = {isa = PBXFileReference; fileEncoding = 4; lastKnownFileType = file.xib; path = LaunchScreen.xib; sourceTree = "<group>"; };
		25F26AEC1CE6675E0045FFBA /* Default-568h@2x.png */ = {isa = PBXFileReference; lastKnownFileType = image.png; path = "Default-568h@2x.png"; sourceTree = "<group>"; };
/* End PBXFileReference section */

/* Begin PBXGroup section */
		254BB8351B1FD08900C56DE9 = {
			isa = PBXGroup;
			children = (
				254BB8671B1FD16500C56DE9 /* main */,
				254BB8421B1FD08900C56DE9 /* Info.plist */,
				254BB84E1B1FD08900C56DE9 /* Images.xcassets */,
				25916F401CE65FF600695115 /* LaunchScreen.xib */,
				25F26AEC1CE6675E0045FFBA /* Default-568h@2x.png */,
				254BB83F1B1FD08900C56DE9 /* products */,
			);
			sourceTree = "<group>";
			usesTabs = 0;
		};
		254BB83F1B1FD08900C56DE9 /* products */ = {
			isa = PBXGroup;
			children = (
				254BB83E1B1FD08900C56DE9 /* main.app */,
			);
			name = products;
			sourceTree = "<group>";
		};
/* End PBXGroup section */

/* Begin PBXNativeTarget section */
		254BB83D1B1FD08900C56DE9 /* main */ = {
			isa = PBXNativeTarget;
			buildConfigurationList = 254BB8611B1FD08900C56DE9 /* Build configuration list for PBXNativeTarget "main" */;
			buildPhases = (
				254BB83C1B1FD08900C56DE9 /* Resources */,
			);
			buildRules = (
			);
			dependencies = (
			);
			name = main;
			productName = main;
			productReference = 254BB83E1B1FD08900C56DE9 /* main.app */;
			productType = "com.apple.product-type.application";
		};
/* End PBXNativeTarget section */

/* Begin PBXProject section */
		254BB8361B1FD08900C56DE9 /* Project object */ = {
			isa = PBXProject;
			attributes = {
				LastUpgradeCheck = 0630;
				ORGANIZATIONNAME = Developer;
				TargetAttributes = {
					254BB83D1B1FD08900C56DE9 = {
						CreatedOnToolsVersion = 6.3.1;
					};
				};
			};
			buildConfigurationList = 254BB8391B1FD08900C56DE9 /* Build configuration list for PBXProject "project" */;
			compatibilityVersion = "Xcode 3.2";
			developmentRegion = English;
			hasScannedForEncodings = 0;
			knownRegions = (
				en,
				Base,
			);
			mainGroup = 254BB8351B1FD08900C56DE9;
			productRefGroup = 254BB83F1B1FD08900C56DE9 /* products */;
			projectDirPath = "";
			projectRoot = "";
			targets = (
				254BB83D1B1FD08900C56DE9 /* main */,
			);
		};
/* End PBXProject section */

/* Begin PBXResourcesBuildPhase section */
		254BB83C1B1FD08900C56DE9 /* Resources */ = {
			isa = PBXResourcesBuildPhase;
			buildActionMask = 2147483647;
			files = (
				254BB8681B1FD16500C56DE9 /* main in Resources */,
				25F26AED1CE6675E0045FFBA /* Default-568h@2x.png in Resources */,
				25916F411CE65FF600695115 /* LaunchScreen.xib in Resources */,
				254BB84F1B1FD08900C56DE9 /* Images.xcassets in Resources */,
			);
			runOnlyForDeploymentPostprocessing = 0;
		};
/* End PBXResourcesBuildPhase section */

/* Begin XCBuildConfiguration section */
		254BB8601B1FD08900C56DE9 /* Release */ = {
			isa = XCBuildConfiguration;
			buildSettings = {
				ALWAYS_SEARCH_USER_PATHS = NO;
				CLANG_CXX_LANGUAGE_STANDARD = "gnu++0x";
				CLANG_CXX_LIBRARY = "libc++";
				CLANG_ENABLE_MODULES = YES;
				CLANG_ENABLE_OBJC_ARC = YES;
				CLANG_WARN_BOOL_CONVERSION = YES;
				CLANG_WARN_CONSTANT_CONVERSION = YES;
				CLANG_WARN_DIRECT_OBJC_ISA_USAGE = YES_ERROR;
				CLANG_WARN_EMPTY_BODY = YES;
				CLANG_WARN_ENUM_CONVERSION = YES;
				CLANG_WARN_INT_CONVERSION = YES;
				CLANG_WARN_OBJC_ROOT_CLASS = YES_ERROR;
				CLANG_WARN_UNREACHABLE_CODE = YES;
				CLANG_WARN__DUPLICATE_METHOD_MATCH = YES;
				"CODE_SIGN_IDENTITY[sdk=iphoneos*]" = "iPhone Developer";
				COPY_PHASE_STRIP = NO;
				DEBUG_INFORMATION_FORMAT = "dwarf-with-dsym";
				ENABLE_NS_ASSERTIONS = NO;
				ENABLE_STRICT_OBJC_MSGSEND = YES;
				GCC_C_LANGUAGE_STANDARD = gnu99;
				GCC_NO_COMMON_BLOCKS = YES;
				GCC_WARN_64_TO_32_BIT_CONVERSION = YES;
				GCC_WARN_ABOUT_RETURN_TYPE = YES_ERROR;
				GCC_WARN_UNDECLARED_SELECTOR = YES;
				GCC_WARN_UNINITIALIZED_AUTOS = YES_AGGRESSIVE;
				GCC_WARN_UNUSED_FUNCTION = YES;
				GCC_WARN_UNUSED_VARIABLE = YES;
				IPHONEOS_DEPLOYMENT_TARGET = 10.0;
				MTL_ENABLE_DEBUG_INFO = NO;
				SDKROOT = iphoneos;
				TARGETED_DEVICE_FAMILY = "1,2";
				VALIDATE_PRODUCT = YES;
			};
			name = Release;
		};
		254BB8631B1FD08900C56DE9 /* Release */ = {
			isa = XCBuildConfiguration;
			buildSettings = {
				ASSETCATALOG_COMPILER_APPICON_NAME = AppIcon;
				INFOPLIST_FILE = Info.plist;
				PRODUCT_NAME = "$(TARGET_NAME)";
			};
			name = Release;
		};
/* End XCBuildConfiguration section */

/* Begin XCConfigurationList section */
		254BB8391B1FD08900C56DE9 /* Build configuration list for PBXProject "project" */ = {
			isa = XCConfigurationList;
			buildConfigurations = (
				254BB8601B1FD08900C56DE9 /* Release */,
			);
			defaultConfigurationIsVisible = 0;
			defaultConfigurationName = Release;
		};
		254BB8611B1FD08900C56DE9 /* Build configuration list for PBXNativeTarget "main" */ = {
			isa = XCConfigurationList;
			buildConfigurations = (
				254BB8631B1FD08900C56DE9 /* Release */,
			);
			defaultConfigurationIsVisible = 0;
			defaultConfigurationName = Release;
		};
/* End XCConfigurationList section */
	};
	rootObject = 254BB8361B1FD08900C56DE9 /* Project object */;
}
`
}

//sailfish

func sailfish_spec(name string) string {
	return fmt.Sprintf(`#
# Do NOT Edit the Auto-generated Part!
# Generated by: spectacle version 0.27
#

Name:       harbour-%[1]v

# >> macros
# << macros

Summary:    Put your summary here
Version:    0.1
Release:    1
Group:      Qt/Qt
License:    MIT
Source0:    %%{name}-%%{version}.tar.bz2

%%description
Put your description here


%%prep
%%setup -q -n %%{name}-%%{version}

# >> setup
# << setup

%%build
# >> build pre
# << build pre

# >> build post
# << build post

%%install
rm -rf %%{buildroot}
# >> install pre
# << install pre
install -d %%{buildroot}%%{_bindir}
install -p -m 0755 %%(pwd)/%%{name} %%{buildroot}%%{_bindir}/%%{name}
install -d %%{buildroot}%%{_datadir}/applications
install -d %%{buildroot}%%{_datadir}/%%{name}
install -d %%{buildroot}%%{_datadir}/icons/hicolor/86x86/apps
install -m 0444 -t %%{buildroot}%%{_datadir}/icons/hicolor/86x86/apps %%{name}.png
install -p %%(pwd)/%[1]v.desktop %%{buildroot}%%{_datadir}/applications/%%{name}.desktop

# >> install post
# << install post

desktop-file-install --delete-original       \
  --dir %%{buildroot}%%{_datadir}/applications             \
   %%{buildroot}%%{_datadir}/applications/%%{name}.desktop

%%files
%%defattr(-,root,root,-)
%%{_bindir}
%%{_datadir}/%%{name}
%%{_datadir}/icons/hicolor/86x86/apps
%%{_datadir}/applications/%%{name}.desktop

# >> files
# << files`, name)
}

func sailfish_desktop(name string) string {
	return fmt.Sprintf(`[Desktop Entry]
Encoding=UTF-8
Version=1.0
Type=Application
X-Nemo-Application-Type=generic
Comment=Put your comment here
Name=%[1]v
Icon=harbour-%[1]v
Exec=harbour-%[1]v`, name)
}

func sailfish_ssh(port, login string, cmd ...string) error {

	typ := "SailfishOS_Emulator"
	if port == "2222" {
		typ = "engine"
	}

	signer, err := ssh.ParsePrivateKey([]byte(utils.Load(filepath.Join(utils.SAILFISH_DIR(), "vmshare", "ssh", "private_keys", typ, login))))
	if err != nil {
		return err
	}

	client, err := ssh.Dial("tcp", "localhost:"+port, &ssh.ClientConfig{User: login, Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)}, HostKeyCallback: ssh.InsecureIgnoreHostKey()})
	if err != nil {
		return err
	}
	defer client.Close()

	sess, err := client.NewSession()
	if err != nil {
		return err
	}

	output, err := sess.CombinedOutput(strings.Join(cmd, " "))
	if err != nil {
		utils.Log.WithField("cmd", strings.Join(cmd, " ")).Debugf("failed to run ssh cmd for %v on %v", typ, runtime.GOOS)
		return errors.New(string(output))
	}

	return nil
}

func ubports_desktop(name string) string {
	return fmt.Sprintf(`[Desktop Entry]
Name=%[1]v
Exec=%[1]v
Icon=logo.svg
Terminal=false
Type=Application
X-Ubuntu-Touch=true`, name)
}

func ubports_apparmor() string {
	return `{
    "policy_groups": [],
    "policy_version": 1.3
}`
}

func ubports_manifest(name string) string {
	return fmt.Sprintf(`{
    "name": "%[1]v",
    "description": "description",
    "architecture": "%[2]v",
    "title": "%[1]v",
    "hooks": {
        "%[1]v": {
            "apparmor": "%[1]v.apparmor",
            "desktop":  "%[1]v.desktop"
        }
    },
    "version": "1.0",
    "maintainer": "maintainer_name <maintainer_email>",
    "framework" : "ubuntu-sdk-15.04.6"
}`, name, func() string {
		if utils.QT_UBPORTS_ARCH() == "arm" {
			return "armhf"
		}
		return "amd64"
	}())
}

func relink(env map[string]string, target string) string {
	return fmt.Sprintf(`#!/bin/bash
set -ev

#GO_VERSION: %v
#GO_HOST_OS: %v
#GO_HOST_ARCH: %v
#QT_VERSION: %v

export GOOS=%v
export GOARCH=%v
export GOARM=%v
export CC=%v
export CXX=%v

go tool link -f -o $PWD/relinked -importcfg $PWD/b001/importcfg.link -buildmode=%v -w -extld=%v $PWD/b001/_pkg_.a`,

		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
		utils.QT_VERSION(),

		env["GOOS"],
		env["GOARCH"],
		func() string {
			if env["GOARCH"] == "arm" {
				return env["GOARM"]
			}
			return ""

		}(),
		env["CC"],
		env["CXX"],

		func() string {
			switch target {
			case "ios", "ios-simulator":
				return "c-archive"
			case "android", "android-emulator":
				return "c-shared"
			default:
				return "exe"
			}
		}(),

		func() string {
			switch target {
			case "ios", "ios-simulator", "darwin":
				return "clang++"
			default:
				return "g++"
			}
		}())
}

//js/wasm

func js_c_main_wrapper(target string) string {
	bb := new(bytes.Buffer)
	bb.WriteString("#include <emscripten.h>\n")
	for _, n := range rcc.ResourceNames {
		fmt.Fprintf(bb, "extern int qInitResources_%v();\n", n)
	}
	bb.WriteString("int main(int argc, char *argv[]) {\n")
	for _, n := range rcc.ResourceNames {
		fmt.Fprintf(bb, "qInitResources_%v();\n", n)
	}

	//TODO: use emscripten_sync_run_in_main_runtime_thread once thread support is there ?
	bb.WriteString("emscripten_run_script(\"Module._goMain()\");\n")

	bb.WriteString("return 0;\n")
	bb.WriteString("}")
	return bb.String()
}

func wasm_js() string {
	return `

	if (!WebAssembly.instantiateStreaming) { // polyfill 
		WebAssembly.instantiateStreaming = async (resp, importObject) => { 
			const source = await (await resp).arrayBuffer(); 
			return await WebAssembly.instantiate(source, importObject); 
		};
	} 

	let go = new Go(); 
	let instance;

	let fetchPromise = fetch("go.wasm");
	WebAssembly.instantiateStreaming(fetchPromise, go.importObject).then((result) => { 
		instance = result.instance;
	}).catch((err) => { 
		//console.log(err); 

		//fallback for wrong MIME type
		fetchPromise.then((response) =>
			response.arrayBuffer()
		).then((bytes) =>
			WebAssembly.instantiate(bytes, go.importObject)
		).then((result) =>
			instance = result.instance
		);
	});

	Module._goMain = function() {
		go.run(instance);
	};
})();`
}
