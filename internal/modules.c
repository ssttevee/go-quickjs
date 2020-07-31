#include "quickjs/quickjs.h"

JSModuleDef *js_value_get_module_def(JSValue value)
{
    return JS_VALUE_GET_PTR(value);
}
