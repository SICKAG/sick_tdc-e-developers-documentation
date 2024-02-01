# Developed for SICK Mobilisis d.o.o. 27.10.2023
# Python script for enabling MAG / ACC

import DirectTdc as dtdc

enableAcc = "/sys/class/misc/FreescaleAccelerometer/enable"
enableMag = "/sys/class/misc/FreescaleMagnetometer/enable"

if __name__ == "__main__":
    accmag = dtdc.DirectTDC()
    accmag.Echo("echo 1 > /sys/class/misc/FreescaleAccelerometer/enable")
    accmag.Echo("echo 1 > /sys/class/misc/FreescaleMagnetometer/enable")