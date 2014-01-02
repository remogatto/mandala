// +build gotask

package test

import (
	"fmt"
	"github.com/jingweno/gotask/tasking"
	"os"
	"path/filepath"
)

var (
	// The project name.
	ProjectName = "runner"

	// The path for the ARM binary. The binary is then copied on
	// each of SharedLibraryPaths.
	ARMBinaryPath = "bin/linux_arm"

	// The path for shared libraries.
	SharedLibraryPaths = []string{
		"android/obj/local/armeabi-v7a/",
		"android/libs/armeabi-v7a/",
	}

	// Android path
	AndroidPath = "android"

	buildFun = map[string]func(*tasking.T, bool){
		"xorg":    buildXorg,
		"android": buildAndroid,
	}

	runFun = map[string]func(*tasking.T, string){
		"xorg":    runXorg,
		"android": runAndroid,
	}

	labelFAIL = red("F")
	labelPASS = green("OK")
)

// NAME
//    build - Build the tests
//
// DESCRIPTION
//    Build the tests for the given platforms (xorg/android).
//
// OPTIONS
//    --all, -a
//        force rebuilding of packages that are already up-to-date
//    --verbose, -v
//        run in verbose mode
func TaskBuild(t *tasking.T) {
	if len(t.Args) == 0 {
		t.Error("At least a platform name must be specified!")
	}
	for _, platform := range t.Args {
		buildFun[platform](t, t.Flags.Bool("a"))
	}
	if t.Failed() {
		t.Fatalf("%-20s %s\n", status(t.Failed()), "Build the tests for the given platforms.")
	}
	t.Logf("%-20s %s\n", status(t.Failed()), "Build the tests for the given platforms.")
}

// NAME
//    test - Run the tests
//
// DESCRIPTION
//    Build and run the tests on the given platforms.
//
// OPTIONS
//    --flags=""
//        pass the given flags to the executable
//    --logcat
//        show logcat output (android only)
//    --verbose, -v
//        run in verbose mode
func TaskTest(t *tasking.T) {
	TaskBuild(t)
	for _, platform := range t.Args {
		if f, ok := runFun[platform]; ok {
			f(t, t.Flags.String("flags"))
		}
	}
	if t.Failed() {
		t.Fatalf("%-20s %s\n", status(t.Failed()), "Run the example on the given platforms.")
	}
	t.Logf("%-20s %s\n", status(t.Failed()), "Run the example on the given platforms.")
}

// NAME
//    deploy - Deploy the application
//
// DESCRIPTION
//    Build and deploy the application on the device via ant.
//
// OPTIONS
//    --verbose, -v
//        run in verbose mode
func TaskDeploy(t *tasking.T) {
	deployAndroid(t)
	if t.Failed() {
		t.Fatalf("%-20s %s\n", status(t.Failed()), "Build and deploy the application on the device via ant.")
	}
	t.Logf("%-20s %s\n", status(t.Failed()), "Build and deploy the application on the device via ant.")
}

// NAME
//    clean - Clean all generated files
//
// DESCRIPTION
//    Clean all generated files and paths.
//
// OPTIONS
//    --verbose, -v
//        run in verbose mode
func TaskClean(t *tasking.T) {
	var paths []string

	// Remove shared libraries
	paths = append(paths, SharedLibraryPaths...)

	// Remove ARM binaries
	paths = append(paths, ARMBinaryPath, filepath.Join("bin", ProjectName))

	// Remove APK files
	apkFiles, _ := filepath.Glob(filepath.Join(AndroidPath, "bin/*.apk"))
	paths = append(paths, apkFiles...)

	// Actually remove files using rm
	for _, path := range paths {
		rm_rf(t, path)
	}

	if t.Failed() {
		t.Fatalf("%-20s %s\n", status(t.Failed()), "Clean all generated files and paths.")
	}
	t.Logf("%-20s %s\n", status(t.Failed()), "Clean all generated files and paths.")
}

func buildXorg(t *tasking.T, buildAll bool) {
	allFlagString := ""
	if buildAll {
		allFlagString = "-a"
	}
	err := t.Exec(
		`sh -c "`,
		"GOPATH=`pwd`:$GOPATH",
		`go build`, allFlagString,
		ProjectName, `"`,
	)
	if err != nil {
		t.Error(err)
	}
}

func buildAndroid(t *tasking.T, buildAll bool) {
	allFlagString := ""

	os.MkdirAll("android/libs/armeabi-v7a", 0777)
	os.MkdirAll("android/obj/local/armeabi-v7a", 0777)

	if buildAll {
		allFlagString = "-a"
	}

	err := t.Exec(`sh -c "`,
		`CC="$NDK_ROOT/bin/arm-linux-androideabi-gcc"`,
		"GOPATH=`pwd`:$GOPATH",
		`GOROOT=""`,
		"GOOS=linux",
		"GOARCH=arm",
		"GOARM=7",
		"CGO_ENABLED=1",
		"$GOANDROID/go install", allFlagString,
		"$GOFLAGS",
		`-ldflags=\"-android -shared -extld $NDK_ROOT/bin/arm-linux-androideabi-gcc -extldflags '-march=armv7-a -mfloat-abi=softfp -mfpu=vfpv3-d16'\"`,
		"-tags android",
		ProjectName, `"`,
	)

	if err != nil {
		t.Error(err)
	}

	os.MkdirAll(ARMBinaryPath, 0777)

	for _, path := range SharedLibraryPaths {
		err := t.Exec(
			"cp",
			filepath.Join(ARMBinaryPath, ProjectName),
			filepath.Join(path, "lib"+ProjectName+".so"),
		)

		if err != nil {
			t.Error(err)
		}
	}
}

func runXorg(t *tasking.T, flags string) {
	buildXorg(t, false)
	err := t.Exec(
		fmt.Sprintf("./%s", ProjectName),
		flags,
	)
	if err != nil {
		t.Error(err)
	}
}

func runAndroid(t *tasking.T, flags string) {
	buildAndroid(t, false)
	deployAndroid(t)
	err := t.Exec(
		fmt.Sprintf(
			"adb shell am start -a android.intent.action.MAIN -n net.gorgasm.%s/android.app.NativeActivity",
			ProjectName,
		))
	if err != nil {
		t.Error(err)
	}
	if t.Flags.Bool("logcat") {
		err = t.Exec("adb shell logcat")
		if err != nil {
			t.Error(err)
		}
	}
}

func deployAndroid(t *tasking.T) {
	buildAndroid(t, false)
	err := t.Exec("ant -f android/build.xml clean debug")
	if err != nil {
		t.Error(err)
	}
	uploadAndroid(t)
}

func uploadAndroid(t *tasking.T) {
	err := t.Exec(fmt.Sprintf("adb install -r android/bin/%s-debug.apk", ProjectName))
	if err != nil {
		t.Error(err)
	}
}

func cp(t *tasking.T, src, dest string) error {
	return t.Exec("cp", src, dest)
}

func mv(t *tasking.T, src, dest string) error {
	return t.Exec("mv", src, dest)
}

func rm_rf(t *tasking.T, path string) error {
	return t.Exec("rm -rf", path)
}

func green(text string) string {
	return "\033[32m" + text + "\033[0m"
}

func red(text string) string {
	return "\033[31m" + text + "\033[0m"
}

func yellow(text string) string {
	return "\033[33m" + text + "\033[0m"
}

func status(failed bool) string {
	if failed {
		return red("FAIL")
	}
	return green("OK")
}
