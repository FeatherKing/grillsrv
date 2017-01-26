# grillsrv
### Go Web Service / Library for Controlling a Green Mountain Grill Pellet Smoker ###

Early last year I purchased a Pellet Smoker that had an onboard wifi chip. I had the desire to log the temperature of my grill and food over time. The grill manufacturer provides a wifi app to control some of its features. Using a mix of network packet sniffing and android apk decompilation, I was able to piece together most of the commands used to control the grill.


This project has three pieces. A terribly designed (but functional) web ui, a go web service/library that communicates with the grill, and coming soon, an android app to replace the admittedly terrible web UI.


The grillsrv.go also has the ability to persist its data to a postgres database.

---
The company that makes these grills has since released a 'Server Mode' that uses AWS to communicate with the grill. Once your grill is set to use server mode, all traffic becomes encrypted and the grill stops responding to commands sent over the LAN.

Essentially you have to decide if using their 'Server Mode' is worth losing being able to control the grill directly yourself. I use this app to simulate their server mode by just port forwarding to this web ui.

---
My typical use is to compile this to run on a raspberry pi and save the results to a postgres container running in docker on CentOS 6.
