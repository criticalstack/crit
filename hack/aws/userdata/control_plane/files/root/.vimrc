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
