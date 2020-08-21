#!/bin/sh

systemctl enable --now disable-swap.service
systemctl daemon-reload
