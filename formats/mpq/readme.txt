To compile and install stormlib on windows with mingw64:

Copy the stormlib folder to an easy to access directory. I used /mingw64/home/thegtproject/stormlib

(from mingw64 console)

$ cmake -DCMAKE_SYSTEM_NAME="windows" -DCMAKE_INSTALL_PREFIX="/mingw64" .
$ make && make install

Done.
