cd /etc/systemd/system

sudo nano template-backend.service

systemctl daemon-reload

sudo systemctl enable template-backend.service
sudo systemctl start template-backend.service
sudo systemctl status template-backend.service



sudo service template-backend status



cd /lib/systemd/system
sudo nano crfund.service