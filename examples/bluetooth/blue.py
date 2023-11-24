from bluepy.btle import Scanner, DefaultDelegate
import sys

# setting target MAC address
target_mac = "e1:6b:1a:7e:c5:be"

def parseTemperature(scan):
    tempData = scan[1][2][4:]
    firstVal = tempData[-2:]
    secondVal = tempData[-4:-2]
    val = str(firstVal) + str(secondVal)
    val = int(val,16) / 100
    print("Temperature value (Celsius): " + str(val))

class ScanDelegate(DefaultDelegate):
    def __init__(self):
        DefaultDelegate.__init__(self)

    def handleDiscovery(self, dev, isNewDev, isNewData):
        if dev.addr == target_mac:
            print("Received advertisement from", dev.addr)
            print("RSSI:", dev.rssi)
            scan = dev.getScanData()
            print("Advertisement data:", scan)
            parseTemperature(scan)
            sys.exit()

if __name__ == "__main__":
    # scanner to fetch passive data
    scanner = Scanner().withDelegate(ScanDelegate())
    print("Scanning for advertisements...")
    # 10 second scan
    # scan clears, starts, processes and stops ble scan
    devices = scanner.scan(10.0)
    # other devices should be passed
    for dev in devices:
       pass