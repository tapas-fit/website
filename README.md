# Tapas website

Landing page for tapas app with waitlist

## frontend

* icons from https://github.com/gauravghongde/social-icons
* custom checkbox based on https://www.w3schools.com/howto/howto_css_custom_checkbox.asp

### embed video or screenshots?

* can record with scrcpy. test.mkv works better in kdenlive for editing than than mp4
* render to webp and gif is super large.
* mp4 is about half the size of webm, so mp4 it is then..
* but kdenlive always wants to save as full HD , and not sure how to fit the portrait into the webpage anyway. let's go with row of screenshots for now, which i can fit into the layout in a nicer way

### favicon

made with https://favicon.io/favicon-generator/
backgroudn #0AA
font: Leckerli One

## backend

simple go backend to save signups to csv files

## anti-spam

using cloudflare turnstile across frontend and backend

### cases to test

* page submitted after open > 5 min. -> seems to work just fine.
* backend down -> works
* backend 3 different test keys, and frontend all test keys -> all work reasonably
* OK -> works
* submit twice -> weird seems to have worked. maybe cause page open long time. but then another submit failed
