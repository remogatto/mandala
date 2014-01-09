# What's that?

Gorgasm is a framework for writing Android native applications in
[Go](http://golang.org). You can develop, test and run your application on your desktop
and then deploy it to an Android device. It encourages the use of
idiomatic Go for writing Android applications: communication happens
through channels, without callbacks. The framework is not to be
considered as an high-level game engine but as a basic layer onto
which game engines can be build or existing ones can be used. In my
opinion, this opens interesting scenarios in the developing of native
Android applications/games in Go. Goandroid's [native_activity]
example was the initial source of inspiration for this project.

Please consider that Gorgasm is in a very early stage of development:
API will change, test coverage is not so good for now. Last but not
least, Go doesn't officially supports native Android
development. Regarding this point, I hope that the present work could
act as a sort of incentive in the direction of an official Android
support by the Go Team.

Have a nice Gorgasm!

# Key features

* Code/test/run on your desktop and deploy on the device.
* Build/deploy/run your application using simple shell commands.
* On-device black-box testing.
* Communicate through channels, no callbacks.
* Quick bootstrap using a predefined template.

# Techonologies involved

* [Android NDK](http://developer.android.com/tools/sdk/ndk/index.html)
* [Goandroid](https://github.com/eliasnaur/goandroid)
* [EGL](https://www.khronos.org/egl/)
* [OpenGL ES 2](http://www.khronos.org/opengles/2_X/)
* [GLFW 3](http://www.glfw.org/)
* [Gotask](https://github.com/jingweno/gotask)
* [Loop](http://git.tideland.biz/goas/loop)
* [PrettyTest](https://github.com/remogatto/prettytest)

# How does it work?

Gorgasm use [Goandroid](https://github.com/) toolchain to compile Go
applications for Android. The graphics abstraction between desktop and
device is obtained using a bunch of technologies. In particular

* EGL
* OpenGL ES 2.0
* GLFW 3

The EGL layer is necessary to use an OpenGL ES context on the
desktop. The GLFW library is responsible of managing the rendering
context and the handling of events in a window.

The framework itself provides an event channel from which client code
listen for events happening during program execution. Examples of
events are the interaction with the screen, the creation of the native
rendering context, pausing/resuming/destroying of the application,
etc.

The framework abstracts the Android native events providing a way to
build, run and test the application on the desktop with the promise
that it will behave the same on the device once deployed. Oh well,
this is the long-term aim, at least!

A typical Gorgasm application has two loops: one continously listen to
events, the other is responsible for rendering the scene. In order to
dealing with application resources (images, sounds, configuration
files, etc.), the framework provides an AssetManager object. Client
code sends request to it in order to obtain resources as
<tt>io.Reader</tt> instances. In the desktop application this simply
means opening the file at the given path. In the Android application
the framework will unpack the apk archive on the fly getting the
requested resources from it. However, is the framework responsibility
to deal with the right native method for opening file. From the
client-code point of view the request will be the same.

The bothering steps needed to build, package and deploy the
application on the device are simplified using a set of predefined
[gotask](https://github.com/jingweno/gotask) tasks.

# Prerequisites

* Android NDK
* Goandroid
* EGL
* OpenGL ES 2
* GLFW3
* gotask
* xdotool (to run the tests)

## Android NDK

See here.

## Goandroid

See [here](https://github.com/eliasnaur/goandroid).

## EGL/OpenGL ES 2

On a debian-like system:

<pre>
sudo apt-get install libgles2-mesa-dev libegl1-mesa-dev
</pre>

## GLFW3

Install from source following the instruction [here](http://www.glfw.org/docs/latest/compile.html). Please note that
you have to configure GLFW in order to use EGL and OpenGL ES 2.

## gotask

<pre>
go get github.com/jingweno/gotask
</pre>

## xdotool

On a debian-like system:

<pre>
sudo apt-get install xdotool
</pre>

# Install

Once you have satisfied all the prerequisites:

<pre>
go get https://github.com/remogatto/gorgasm
</pre>

This will install all the remaining dependencies.

# Quick start

To create a basic application simply clone the
<tt>gorgasm-template</tt> repo and run a bunch of tasks:

<pre>
git clone https://github.com/remogatto/gorgasm-template myapp
cd myapp
gotask init MyApp
gotask run android # deploy and run on a connected device
gotask run xorg    # run on a desktop window
</pre>

This will generate a simple Android application showing a colored
triangle on the screen. See
[gorgasm-template](https://github.com/remogatto/gorgasm-template) for
furher info.

# Testing

Setup a testing environment on Android was not straightforward. The
main [issue](https://github.com/eliasnaur/goandroid/issues/20) is
related to the <tt>flag</tt> package. To avoid dependency from it I
had to hack [PrettyTest](https://github.com/remogatto/prettytest) in
order to remove the dependency from <tt>testing</tt> (which in turn
depends on <tt>flag</tt>). So basically, testing a native Android
application is now possible using
[PrettyTest](https://github.com/remogatto/prettytest) but we have to
renounce to the benifits of <tt>testing</tt> (at least for now). See
[test](test/) for further info about testing. To run the tests on your
desktop window you need the <tt>xdotool</tt> (see the Prerequisites
section)

# To do

* Sound support
* More tests
* Write a complete game

# License

See [LICENSE](LICENSE)
