## Install chrome headless
###Uninstall previous chrome
``
sudo apt-get purge chromium-browser
``

### Re-install the new stable chrome headless
``
sudo apt-get update
``

``
sudo apt-get install -y libappindicator1 fonts-liberation
``

``
wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb
``

``
sudo dpkg -i google-chrome*.deb
``

This is going to install dependencies required for chrome to be installed.

In case of error related to unmet dependencies you can do:

``
sudo apt-get -f install
``

``
sudo dpkg --configure -a
``

As last option you can try with:

``
sudo apt-get -u dist-upgrade
``

and then to try to installed it again.

###Other optional dependencies :

``
sudo apt-get -y install dbus-x11 xfonts-base xfonts-100dpi xfonts-75dpi xfonts-cyrillic xfonts-scalable
``

### Test installation version
Verify the version and is it install by
``google-chrome-stable -version``

If you want to test work in headless mode you can do:

``google-chrome-stable --headless --disable-gpu --dump-dom https://www.chromestatus.com/``

and check the output picture

### Create service (ubuntu)

- create file `/etc/systemd/system/redirective.service` with content from `service.example` file
- change `ExecStart` value to executive file path
- create config file `/etc/rsyslog.d/redirective.conf` with content from `service_log.example` file (log file path `/var/log/redirective.log`)
- restart rsyslog by running `sudo service rsyslog restart` command
- start **redirective** service by running `sudo service redirective start` command
