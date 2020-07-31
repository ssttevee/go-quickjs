#include "quickjs/quickjs.h"

JSValue js_call_func_job(JSContext *ctx, int argc, JSValueConst *argv)
{
    return JS_Call(ctx, argv[0], argv[1], argc-2, &argv[2]);
}
