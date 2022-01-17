#!/usr/bin/env bash

SCRIPTS=$(dirname $0)
SERVER_IP="X.X.X.X"
SV_HOME="/home/golangblog"
SV_USER="golangblog"
EXE_NAME="golangblog"

$SCRIPTS/build.sh
if [[ $? -ne 0 ]]
then
	echo "Building error!"
	exit 1
fi

rsync \
	-avz --delete --delete-excluded --exclude=".DS_Store" \
	$SCRIPTS/../files root@$SERVER_IP:$SV_HOME/.

SERVICE=$(cat $SCRIPTS/service/golangblog.service)
FILE=$(find $SCRIPTS/../bin -name $EXE_NAME_\*-linux_amd64)
BASENAME=$(basename $FILE)
ssh root@$SERVER_IP <<-CMD
	service $EXE_NAME stop
	rm -f $SV_HOME/$EXE_NAME_v*
CMD

scp $FILE root@$SERVER_IP:$SV_HOME/.
ssh root@$SERVER_IP <<-CMD
	chmod 700 $SV_HOME/$BASENAME
	chown $SV_USER:$SV_USER $SV_HOME/$BASENAME
	chown -R $SV_USER:$SV_USER $SV_HOME/files
	ln -sf $SV_HOME/$BASENAME $SV_HOME/$EXE_NAME
	chown -h $SV_USER:$SV_USER $SV_HOME/$EXE_NAME
	echo "${SERVICE}" > /etc/systemd/system/$EXE_NAME.service
	systemctl daemon-reload
	service $EXE_NAME start
CMD

exit 0
