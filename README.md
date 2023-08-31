# niconi
Magisk module of elichika for simpily install in rooted Android phone.
## manual install
* install patched as.apk
* install niconi module
* reboot phone
* connect to pc, set pc ip to 192.168.31.2, or you can modify /data/adb/modules/niconi/config.json to you pc ip
* in pc run cdn server. use any tool is ok, put assets to static/
* start game it will download assets for pc, takes about 10 minutes and 17GB space.
## restore existing game
* rename /sdcard/Android/data/com.klab.lovelive.allstars.global(jp version path is similar) to something else
* copy /data/data/com.klab.lovelive.allstars.global/shared_prefs/com.klab.lovelive.allstars.v2.playerprefs.xml to /sdcard
* install patched as.apk
* install niconi reboot phone
* start game onece, it will show error. force close it
* clear app data
* start game again, close it
* edit copy /data/data/com.klab.lovelive.allstars.global/shared_prefs/com.klab.lovelive.allstars.v2.playerprefs.xml <string name="SQ">wgY7bcfzRmHWOxt0FWfwu2qxz85CPS142N7T8leTY7g%3D</string> change key between SQ"> and </string> to value of your backup up file.
* remove /sdcard/Android/data/com.klab.lovelive.allstars.global, restore your origianl folder
* start the game, it should work. if not, try manual install.

## fast install
* uninstall game if already installed
* copy as.apk and llas.apk to /sdcard
* install as.apk, start game, it will show error
* clear app data
* install niconi module, it will take about 7-9 minutes
* reboot
* start the game, it should work.
> llas.zip contains com.klab.lovelive.allstars.global folder for /sdcard/Android/data and a matching com.klab.lovelive.allstars.v2.playerprefs.xml file.

# elichika

Local Server for Love Live! All Stars Japanese / Global

Clone this repository.
```
# clone
git clone https://github.com/YumeMichi/elichika

# or update
git pull
```
Edit cdn server (Usually local ip:port) in config.json

Build executable for Android.
```
./b.sh
```
> windows user can copy commands in b.sh to build, zip command is needed.

copy niconi.zip to your phone.

### Assets

Put databases and assets into `static/2d61e7b4e89961c7` (Global) or `static/b66ec2295e9a00aa` (Japanese).

You can download assets from [ll-sifas-cdn-data](https://archive.org/download/ll-sifas-cdn-data).

File list:
| File name                                | Description       |
| :--------------------------------------- | :---------------- |
| sifas-jp-cdn-assets-b66ec2295e9a00aa.tar | assets (Japanese) |
| sifas-gl-cdn-assets-2d61e7b4e89961c7.tar | assets (Global)   |

### Clients
Use [3.12.0 clients](https://selenachina-my.sharepoint.com/:u:/p/walter/ER6SWMf1vKBMsm39VVXq9AABH4SCwxx7pnB4aWJanR116A?e=snCBhP 作者：walterclozet https://www.bilibili.com/read/cv25706340/?spm_id_from=333.788.0.0 出处：bilibili). For fast install you need to download [llas.zip](https://selenachina-my.sharepoint.com/:u:/p/walter/EUkzHP6PAMNEjRzCjBZCiEYBC2tI5NxwnrrCyAEGUhFI4g?e=E0L15K).
