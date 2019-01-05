#define bool unsigned char

#define FILE_BEGIN   0
#define FILE_CURRENT 1
#define FILE_END     2

#define SFILE_INVALID_SIZE 0xffffffff

bool  WINAPI SFileOpenArchive(const TCHAR * szMpqName, DWORD dwPriority, DWORD dwFlags, HANDLE * phMpq);
bool  WINAPI SFileCloseArchive(HANDLE hMpq);
bool  WINAPI SFileOpenFileEx(HANDLE hMpq, const char * szFileName, DWORD dwSearchScope, HANDLE * phFile);
DWORD WINAPI SFileSetFilePointer(HANDLE hFile, LONG lFilePos, LONG * plFilePosHigh, DWORD dwMoveMethod);
bool  WINAPI SFileReadFile(HANDLE hFile, void * lpBuffer, DWORD dwToRead, LPDWORD pdwRead, LPOVERLAPPED lpOverlapped);
bool  WINAPI SFileCloseFile(HANDLE hFile);
