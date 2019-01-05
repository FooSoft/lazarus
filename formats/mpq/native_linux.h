#define DWORD            unsigned int
#define LPDWORD          unsigned int *
#define LPOVERLAPPED     void *
#define TCHAR            char
#define HANDLE           void *
#define LONG             int
#define ERROR_HANDLE_EOF 1002
#define WINAPI

DWORD GetLastError();

#include <stdlib.h>
#include "native.h"
