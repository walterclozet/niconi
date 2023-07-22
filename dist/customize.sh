#!/system/bin/sh
ui_print "setting permission..."
set_perm  ${MODPATH}/niconi  0  0  0755
set_perm  ${MODPATH}/service.sh  0  0  0755
ui_print "installing patched game..."
pm install ${MODPATH}/assets/as2.apk
