# Python script for direct access to data files on TDC-E
# Created for SICK Mobilsis d.o.o. 27.10.2023.

import subprocess
import os

#Reads data directly from TDC-E
class DirectTDC:
    def ReadData(self, path):
        try:
            output = subprocess.check_output(['cat', path], universal_newlines=True)
            return output
        except subprocess.CalledProcessError as e:
            print("Error: {e}")
            return e
    def Echo(self, command):
        try:
            os.system(command)
        except Exception as e:
            print(f"Error: {e}")