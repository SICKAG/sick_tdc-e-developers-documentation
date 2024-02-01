# Developed for SICK Mobilisis d.o.o. 27.10.2023.
# Main program for fetching and processing acc / mag data

import DirectTdc as dtdc
import math
import time

accPath = "/sys/class/misc/FreescaleAccelerometer/data"
magPath = "/sys/class/misc/FreescaleMagnetometer/data"

# splits text
def splitter(text, character):
    return text.split(character)

# orientation calculation
def parseCompassHeadings(x, y):
    offset_x = 0
    offset_y = 0
    scale_x = 1.0
    scale_y = 1.0

    # Calibrating MAG data
    calibrated_x = (float(x) - offset_x) * scale_x
    calibrated_y = (float(y) - offset_y) * scale_y

    # Calculating yaw - compass heading
    yaw = math.atan2(calibrated_y, calibrated_x)
    yaw = math.degrees(yaw)
    
    # Top of the device (Y-axis) pointing in a direction X° counterclockwise from magnetic North
    if yaw < 0:
        yaw += 360.0

    print("Yaw (Compass Heading): {:.2f} degrees".format(yaw))
    print("\n")

# parses ACC and MAG data
def parseAllData(acc, mag):
    # parsing ACC
    accData = splitter(acc, ",")
    print("ACC X: " + accData[0] + "g")
    print("ACC Y: " + accData[1] + "g")
    print("ACC Z: " + accData[2] + "g")
    print("\n")
    
    # parsing MAG
    magData = splitter(mag, ",")
    print("MAG X: " + magData[0] + "μT")
    print("MAG Y: " + magData[1] + "μT")
    print("MAG Z: " + magData[2] + "μT")
    print("\n")

    parseCompassHeadings(magData[0], magData[1])

if __name__ == "__main__":
    while(1):
        datai = dtdc.DirectTDC()
        acc = datai.ReadData(accPath)
        acc = splitter(acc, "\n")

        print("Accelerometer data: " + acc[0])
        mag = datai.ReadData(magPath)
        mag = splitter(mag, "\n")
        print("Magnetometer data: " + mag[0])

        del datai

        print("\n-----------ACC/MAG-DATA-----------")
        parseAllData(acc[0], mag[0])
        time.sleep(5)
