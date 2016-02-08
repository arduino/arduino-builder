## Arduino Builder

A command line tool for compiling Arduino sketches

This tool is able to parse [Arduino Hardware specifications](https://github.com/arduino/Arduino/wiki/Arduino-IDE-1.5-3rd-party-Hardware-specification), properly run `gcc` and produce compiled sketches.

An Arduino sketch differs from a standard C program in that it misses a `main` (provided by the Arduino core), function prototypes are not mandatory, and libraries inclusion is automagic (you just have to `#include` them).
This tool generates function prototypes and gathers library paths, providing `gcc` with all the needed `-I` params.

### Usage

* `-compile` or `-dump-prefs` or `-preprocess` or `-listen:3000`: Optional. If omitted, defaults to `-compile`. `-dump-prefs` will just print all build preferences used, `-compile` will use those preferences to run the actual compiler, `-preprocess` will only print preprocessed code to stdout, `-listen:3000` opens a web server on the specified port. See the section web api below.

* `-hardware`: Mandatory. Folder containing Arduino platforms. An example is the `hardware` folder shipped with the Arduino IDE, or the `packages` folder created by Arduino Boards Manager. Can be specified multiple times. If conflicting hardware definitions are specified, the last one wins.

* `-tools`: Mandatory. Folder containing Arduino tools (`gcc`, `avrdude`...). An example is the `hardware/tools` folder shipped with the Arduino IDE, or the `packages` folder created by Arduino Boards Manager. Can be specified multiple times.

* `-libraries`: Optional. Folder containing Arduino libraries. An example is the `libraries` folder shipped with the Arduino IDE. Can be specified multiple times.

* `-fqbn`: Mandatory. Fully Qualified Board Name, e.g.: arduino:avr:uno

* `-build-path`: Optional. Folder where to save compiled files. If omitted, a folder will be created in the temporary folder specified by your OS.

* `-prefs=key=value`: Optional. It allows to override some build properties.

* `-warnings`: Optional, can be "none", "default", "more" and "all". Defaults to "none". Used to tell `gcc` which warning level to use (`-W` flag).

* `-verbose`: Optional, turns on verbose mode.

* `-quiet`: Optional, supresses almost every output.

* `-debug-level`: Optional, defaults to "5". Used for debugging. Set it to 10 when submitting an issue.

* `-core-api-version`: Optional, defaults to "10600". The version of the Arduino IDE which is using this tool.

* `-logger`: Optional, can be "human" or "machine". Defaults to "human". If "machine", messages emitted will be in a format which the Arduino IDE understands and that it uses for I18N.

* `-version`: if specified, prints version and exits.

* `-build-options-file`: it specifies path to a local `build.options.json` file (see paragraph below), which allows you to omit specifying params such as `-hardware`, `-tools`, `-libraries`, `-fqbn`, `-pref` and `-ide-version`.

* `-vid-pid`: when specified, VID/PID specific build properties are used, if boards supports them.

Final mandatory parameter is the sketch to compile.

### What is and how to use build.options.json file

Every time you run this tool, it will create a `build.options.json` file in build path. It's used to understand if build options (such as hardware folders, fqbn and so on) were changed when compiling the same sketch.
If they changed, the whole build path is wiped out. If they didn't change, previous compiled files will be reused if the corresponding source files didn't change as well.
You can save this file locally and use it instead of specifying `-hardware`, `-tools`, `-libraries`, `-fqbn`, `-pref` and `-ide-version`.

### Using it for continuously verify your libraries or cores

See [Doing continuous integration with arduino builder](https://github.com/arduino/arduino-builder/wiki/Doing-continuous-integration-with-arduino-builder/).

### Building from source

You need [Go 1.4.3](https://golang.org/dl/#go1.4.3).

Repo root contains the script `setup_go_env_vars`. Use it as is or as a template for setting up Go environment variables.

To install `codereview/patch` you have to install [Mercurial](https://www.mercurial-scm.org/) first.

Once done, run the following commands:

```
go get github.com/go-errors/errors
go get github.com/stretchr/testify
go get github.com/jstemmer/go-junit-report
go get golang.org/x/codereview/patch
go get golang.org/x/tools/cmd/vet
go build
```

### Web Api

You can choose to compile the builder with the -api option:

```bash
go build -tags 'api'
```

Then if you launch it with the option `listen=3000` the builder will open a web server listening for requests to the /compile endpoint.

Here's how to request the compilation of the simplest sketch possible:

```
POST /compile HTTP/1.1
Host: localhost:3000
Content-Type: application/json

{
    "fqbn": "arduino:avr:uno",
    "sketch": {
        "main_file": {
            "name": "sketch.ino",
            "source": "void setup() {\n  // initialize digital pin 13 as an output.\n  pinMode(13, OUTPUT);\n}\n// the loop function runs over and over again forever\nvoid loop() {\n  digitalWrite(13, HIGH);   // turn the LED on (HIGH is the voltage level)\n  delay(1000);              // wait for a second\n  digitalWrite(13, LOW);    // turn the LED off by making the voltage LOW\n  delay(1000);              // wait for a second\n}"
        }
    }
}
```

And here's the response (the actual response will be much bigger, but the structure is the same):

```
{
    "binaries": {
        "elf": "f0VMRgEBAQAAAAAAAAAAAAIAUwABAAAAAAAAADQAAACILAAAhQAAAA...",
        "hex": "OjEwMDAwMDAwMEM5NDVDMDAwQzk0NkUwMDBDOTQ2RTAwMEM5NDZFMD..."
    },
    "out": [
        {
            "Level": "warn",
            "Message": "Board Intel:i586:izmir_fd doesn't define a 'build.board' preference. Auto-set to: I586_IZMIR_FD"
        },
        {
            "Level": "info",
            "Message": "\"/opt/tools/avr-gcc/4.8.1-arduino5/bin/avr-g++\" -c -g -Os -w -std=gnu++11 -fno-exceptions -ffunction-sections -fdata-sections -fno-threadsafe-statics  -w -x c++ -E -CC -mmcu=atmega328p -DF_CPU=16000000L -DARDUINO=10608 -DARDUINO_AVR_UNO -DARDUINO_ARCH_AVR   \"-I/opt/cores/arduino/avr/cores/arduino\" \"-I/opt/cores/arduino/avr/variants/standard\" \"/tmp/build/sketch/sketch.ino.cpp\" -o \"/dev/null\""
        },
        {
            "Level": "info",
            "Message": "\"/opt/tools/avr-gcc/4.8.1-arduino5/bin/avr-g++\" -c -g -Os -w -std=gnu++11 -fno-exceptions -ffunction-sections -fdata-sections -fno-threadsafe-statics  -w -x c++ -E -CC -mmcu=atmega328p -DF_CPU=16000000L -DARDUINO=10608 -DARDUINO_AVR_UNO -DARDUINO_ARCH_AVR   \"-I/opt/cores/arduino/avr/cores/arduino\" \"-I/opt/cores/arduino/avr/variants/standard\" \"/tmp/build/sketch/sketch.ino.cpp\" -o \"/dev/null\""
        },
        ...
    ]
}
```

### TDD

In order to run the tests, type:

```
go test -v ./src/arduino.cc/builder/test/...
```

In jenkins, use
```
go test -v ./src/arduino.cc/builder/test/... | bin/go-junit-report > report.xml
```

### License and Copyright

`arduino-builder` is licensed under General Public License version 2, as published by the Free Software Foundation. See [LICENSE.txt](LICENSE.txt).

Copyright (C) 2015 Arduino LLC and contributors

See https://www.arduino.cc/ and https://github.com/arduino/arduino-builder/graphs/contributors
