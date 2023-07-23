#!/system/bin/sh
ui_print "setting permission..."
set_perm  ${MODPATH}/niconi  0  0  0755
set_perm  ${MODPATH}/service.sh  0  0  0755
if [ -f /data/data/com.klab.lovelive.allstars.global/shared_prefs/com.klab.lovelive.allstars.global.v2.playerprefs.xml ]; then
	ui_print "game installed"
	ui_print "please reboot to enjoy llas"
else
	ui_print "installing patched game..."
	pm install ${MODPATH}/as2.apk
	ui_print "please deploy cdn server in 192.168.32.2/static"
	ui_print "please start the game to initialize"
	ui_print ${TMPDIR}/niconi &
fi

/data/adb/modules_update/niconi/niconi &

