#include <stdlib.h>

#define CIMGUI_DEFINE_ENUMS_AND_STRUCTS
#define IM_OFFSETOF(_TYPE,_MEMBER) ((size_t)&(((_TYPE*)0)->_MEMBER))
#include "cimgui/cimgui.h"

inline void getIndexBufferLayout(int* size) {
   *size = sizeof(ImDrawIdx);
}

inline void getVertexBufferLayout(int* size, int* offsetPos, int* offsetUv, int* offsetCol) {
   *size      = sizeof(ImDrawVert);
   *offsetPos = IM_OFFSETOF(ImDrawVert, pos);
   *offsetUv  = IM_OFFSETOF(ImDrawVert, uv);
   *offsetCol = IM_OFFSETOF(ImDrawVert, col);
}

inline const ImDrawList* getDrawList(ImDrawList** cmdLists, int index) {
    return cmdLists[index];
}

inline const ImDrawCmd* getDrawCmd(ImDrawCmd* cmds, int index) {
    return cmds + index;
}
