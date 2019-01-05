#define WINAPI

#define DWORD        unsigned long
#define LPDWORD      unsigned long *
#define LPOVERLAPPED void *
#define TCHAR        char
#define HANDLE       void *
#define LONG         long

#define ERROR_HANDLE_EOF 1002

DWORD GetLastError();

#include <stdlib.h>
#include "native.h"
