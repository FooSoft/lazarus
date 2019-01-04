#include <stdlib.h>

#define CIMGUI_DEFINE_ENUMS_AND_STRUCTS
#define IM_OFFSETOF(_TYPE,_MEMBER) ((size_t)&(((_TYPE*)0)->_MEMBER))
#include "cimgui/cimgui.h"

// silly trick to get go-vet off our back
inline ImTextureID nativeHandleCast(uintptr_t id) {
    return (ImTextureID)id;
}
