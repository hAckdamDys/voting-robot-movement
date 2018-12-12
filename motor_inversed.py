#!/usr/bin/python3

from ev3dev.ev3 import *
from time import sleep
import urllib.request as ur

mB, mA = LargeMotor('outA'), LargeMotor('outB')

btn = Button()

def setSpeed(spA, spB):
	spA = max(-500, min(500, spA))
	spB = max(-500, min(500, spB))
	mA.run_forever(speed_sp=spA)
	mB.run_forever(speed_sp=spB)

# main
while not btn.any():
	
	mystr = ""

	# Try reading url
	try:
		mystr = ur.urlopen("http://192.168.44.1:8080/getCommand").read().decode("utf8")
		spA, spB = mystr.split('|')
		setSpeed(int(spA), int(spB))

	except:
		mA.stop()
		mB.stop()
		print("Cannot open 192.168.44.1:8080/getCommand")
	

	sleep(1)

mA.stop()
mB.stop()

