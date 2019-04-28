# Lazarus #

The Lazarus project aims to preserve [Diablo II](https://en.wikipedia.org/wiki/Diablo_II) by reimplementing the game
engine in the Go programming language. Lazarus is a drop-in replacement for the original game executable; the user is
responsible for supplying the game assets (namely the MPQ files) from the official game media.

![](https://foosoft.net/projects/lazarus/img/viewer.png)

## Building ##

It is not currently possible to use `go get` to install all of the packages in the project in one step; some assembly is
required. Follow the instructions below to set up a build environment from a fresh install of your 64-bit operating
system of choice.

### Linux ###

Lazarus is primarily being developed on Fedora, but the required package names are also provided for Ubuntu.

1.  Install the required packages (for [Fedora](https://getfedora.org/) users):
    ```
    # sudo dnf install golang gcc-c++ cmake make git SDL2-devel mesa-libGL-devel zlib-devel bzip2-devel
    ```
    Install the required packages (for [Ubuntu](https://www.ubuntu.com/) users):
    ```
    # sudo apt-get install golang g++ cmake make git libsdl2-dev libgl1-mesa-dev zlib1g-dev libbz2-dev
    ```
2.  Build the [Dear ImGui](https://github.com/ocornut/imgui) wrapper package:
    ```
    $ go get github.com/FooSoft/lazarus/platform/imgui
    ```
    Go will fetch the code, but Cgo will fail to link the [cimgui](https://github.com/cimgui/cimgui) wrapper;
    we need to configure and build it:
    ```
    $ cd $GOPATH/src/github.com/FooSoft/lazarus/platform/imgui/cimgui
    $ cmake -DIMGUI_STATIC="yes" .
    $ make
    ```
    You should now have a `cimgui.a` statically linked library in the `cimgui` directory.
3.  Build the [StormLib](http://zezula.net/en/mpq/stormlib.html) wrapper package:
    ```
    $ go get github.com/FooSoft/lazarus/formats/mpq
    ```
    Go will fetch the code, but Cgo will fail to link the StormLib wrapper;
    we need to configure and build it:
    ```
    $ cd $GOPATH/src/github.com/FooSoft/lazarus/formats/mpq/stormlib
    $ cmake .
    $ make
    ```
    You should now have a `libstorm.a` statically linked library in the `stormlib` directory.

### Windows ###

Lazarus is only tested on Windows 10, but should in theory run on anything newer than Windows XP.

1.  Download the latest 64-bit MSI installer for Go from the [official homepage](https://golang.org/dl/); install to the default directory.
2.  Download and the latest 64-bit EXE installer for MSYS2 from the [official homepage](https://www.msys2.org/); install to the default directory.
3.  Install the required packages (using the MSYS MinGW terminal):
    ```
    $ pacman -S mingw-w64-x86_64-gcc cmake make git mingw-w64-x86_64-SDL2 zlib-devel libbz2-devel
    ```
4.  Add `C:\msys64\usr\bin` and `C:\msys64\mingw64\bin` to your system's `PATH` environment variable.
5.  Build the [Dear ImGui](https://github.com/ocornut/imgui) wrapper package (using the system command prompt):
    ```
    $ go get github.com/FooSoft/lazarus/platform/imgui
    ```
    Go will fetch the code, but Cgo will fail to link the [cimgui](https://github.com/cimgui/cimgui) wrapper;
    we need to configure and build it:
    ```
    $ cd %GOPATH%/src/github.com/FooSoft/lazarus/platform/imgui/cimgui
    $ cmake -DIMGUI_STATIC="yes" .
    $ make
    ```
    You should now have a `cimgui.a` statically linked library in the `cimgui` directory.
6.  Build the [StormLib](http://zezula.net/en/mpq/stormlib.html) wrapper package (using the system command prompt):
    ```
    $ go get github.com/FooSoft/lazarus/formats/mpq
    ```
    Go will fetch the code, but Cgo will fail to link the StormLib wrapper;
    we need to configure and build it:
    ```
    $ cd %GOPATH%/src/github.com/FooSoft/lazarus/formats/mpq/stormlib
    $ cmake .
    $ make
    ```
    You should now have a `libstorm.a` statically linked library in the `stormlib` directory.

## Tools ##

This project includes several tools which are used to demonstrate the capabilities of the engine as well as manipulate
game data for debugging purposes. Make sure to perform the setup steps outlined in the "Building" section before
installing these packages.

### `dc6` ###

Converts the frames of one or more DC6 animations to PNG files, using the provided palette file.

*   Installation:
    ```
    $ go get github.com/FooSoft/lazarus/tools/dc6
    ```
*   Usage:
    ```
    Usage: dc6 [options] palette_file [dc6_files]
    Parameters:

    -target string
            target directory (default ".")
    ```

### `mpq` ###

Extracts the contents of one or more MPQ archives to a target directory, using an optional filter.

*   Installation:
    ```
    $ go get github.com/FooSoft/lazarus/tools/mpq
    ```
*   Usage:
    ```
    Usage: mpq [options] command [files]
    Parameters:

    -filter string
            wildcard file filter (default "**")
    -target string
            target directory (default ".")
    ```

### `viewer` ###

Displays the frames of DC6 animation files, using the provided palette file. A grayscale fallback palette is used if no
palette is provided on the command line.

*   Installation:
    ```
    $ go get github.com/FooSoft/lazarus/tools/viewer
    ```
*   Usage:
    ```
    Usage: viewer [options] file
    Parameters:

    -palette string
            path to palette file
    ```
