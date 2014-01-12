// Gorgasm is a framework for writing Android native applications in
// Go using the Goandroid[1] toolchain. You can develop, test and run
// your application on your desktop and then deploy it to an Android
// device. It encourages the use of idiomatic Go for writing Android
// applications: communication happens through channels, no
// callbacks. The framework is not to be considered as an high-level
// game engine but as a basic layer onto which game engines can be
// build or existing ones can be used. In my opinion, this opens
// interesting scenarios in the developing of native Android
// applications/games in Go. Goandroid's native_activity[2] example
// was the initial source of inspiration for this project.
//
// Please consider that Gorgasm is in a very early stage of
// development: API will change, test coverage is not so good for
// now. Last but not least, Go doesn't officially supports native
// Android development. Regarding this point, I hope that the present
// work could act as a sort of incentive in the direction of an
// official Android support by the Go Team.
//
// Have a nice Gorgasm!
//
//[1] - https://github.com/eliasnaur/goandroid
//
//[2] - https://github.com/eliasnaur/goandroid/tree/master/native-activity
//
package gorgasm
