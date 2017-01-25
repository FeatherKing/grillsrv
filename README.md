# grillsrv
### Go Web Service / Library for Controlling a Green Mountain Grill Pellet Smoker ###

Early last year I purchased a Pellet Smoker that had an onboard wifi chip. I had the desire to log the temperature of my grill and food over time. The grill has a wifi app to control some of its features. Using a mix of network packet sniffing and android apk decompilation, I was able to piece together most of the commands used to control the grill.


This project has two pieces. A terribly designed (but functional) web ui, and a go library that contains all the grill functions. Early on in the project they were the same file, but a while ago I split out the actual grill commands into a reuseable Go library. I plan to replace the web ui with something else.


The grillsrv.go is written to save its data to a postgres database.

The company that makes these grills has since released a 'Server Mode' that uses AWS to communicate with the grill. Once your grill is set to use server mode, all traffic becomes encrypted and the grill stops responding to commands sent over the LAN. Essentially you have to decide if using their 'Server Mode' is worth losing being able to control the grill directly yourself. I use this app to simulate their server mode by just port forwarding to this web ui.


My typical use is to compile this to run on a raspberry pi and save the results to a postgres container running in docker on CentOS 6.
