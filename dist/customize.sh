#!/system/bin/sh
ui_print "setting permission..."
set_perm  ${MODPATH}/elichika  0  0  0755
set_perm  ${MODPATH}/service.sh  0  0  0755
set_perm  ${MODPATH}/unzip 0 0 0755
if [ -f /sdcard/llas.zip && -f  ]; then	
	ui_print "installing game data.."
	${MODPATH}/unzip -d /sdcard/Android/data/ -o /sdcard/llas.zip
fi
if [ ! -f /sdcard/Android/data/com.klab.lovelive.allstars.global/shared_prefs ]; then
	ui_print "installing profile"
	mkdir /data/data/com.klab.lovelive.allstars.global/shared_prefs
	cp /sdcard/Android/data/com.klab.lovelive.allstars.global.v2.playerprefs.xml /data/data/com.klab.lovelive.allstars.global/shared_prefs/com.klab.lovelive.allstars.global.v2.playerprefs.xml
	set_perm /data/data/com.klab.lovelive.allstars.global/shared_prefs 0 0 0777
	set_perm /data/data/com.klab.lovelive.allstars.global/shared_prefs/com.klab.lovelive.allstars.global.v2.playerprefs.xml 0 0 0777
	ui_print "this module uses elichika(https://github.com/YumeMichi/elichika) as local server with modified assets"
fi
killall elichika
killall niconi
cd ${MODPATH}/
./elichika &
ui_print "enjoy the game!"

