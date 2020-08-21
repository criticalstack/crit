#!/bin/bash

dnf install -y bind-utils vim bash-completion nc go git jq nload tcpdump

# sensible vimrc
cat << EOF > /root/.vimrc
set encoding=utf-8

" Exit insert mode
imap jk <Esc>

" Easier menu access and remap repeat motion
nnoremap ; :
nnoremap m ;
nnoremap M ,

" Start/End Line Movement
nnoremap H ^
nnoremap L $
vnoremap H ^
vnoremap L $

" Split line
nnoremap S i<Enter><Esc>^

set pastetoggle=<F11>
EOF

# completions
kubectl completion bash > /usr/share/bash-completion/completions/kubectl
kubectl completion bash | sed 's/kubectl/k/g' > /usr/share/bash-completion/completions/k

# aliases

cat << EOF >> /etc/bashrc
alias k='kubectl'
alias gp='kubectl get pods --all-namespaces'
alias gn='kubectl get nodes'
alias gnw='kubectl get nodes -o wide'
alias si='sudo -i'
alias dnstools='kubectl run -it --rm --restart=Never --image=infoblox/dnstools:latest dnstools'
alias ll='ls -lhaF'
alias setns='kubectl config set-context --current --namespace'
alias jc='journalctl --no-tail -xeu'
EOF
