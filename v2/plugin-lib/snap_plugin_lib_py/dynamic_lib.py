from ctypes import CDLL
import os.path
import platform

# Dependent library
if os.getenv("SNAP_PLUGIN_LIB"):
    PLUGIN_LIB_OBJ = CDLL(os.getenv("SNAP_PLUGIN_LIB"))
else:
    PLUGIN_LIB_EXTENSION = ".dll" if (platform.system() == "Windows") else ".so"
    PLUGIN_LIB_FILE = "snap-plugin-lib%s" % PLUGIN_LIB_EXTENSION
    PLUGIN_LIB_OBJ = CDLL(os.path.join(os.path.dirname(__file__), PLUGIN_LIB_FILE))

