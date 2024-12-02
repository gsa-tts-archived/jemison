# nih performance

To run the analysis, the guestbook and host table are needed as JSON.

sling run --src-conn WORKDB --src-stream 'public.guestbook' --tgt-object 'file://guestbook.json'
sling run --src-conn WORKDB --src-stream 'public.hosts' --tgt-object 'file://hosts.json'

