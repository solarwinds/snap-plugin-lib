from ctypes import CDLL
import os.path
import platform

# Dependent library
PLUGIN_LIB_EXTENSION = ".dll" if (platform.system() == "Windows") else ".so"
PLUGIN_LIB_FILE = "swisnap-plugin-lib%s" % PLUGIN_LIB_EXTENSION
PLUGIN_LIB_OBJ = CDLL(PLUGIN_LIB_FILE)
