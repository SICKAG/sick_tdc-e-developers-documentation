import subprocess
import random
import time

# executes echo command
def execute(cmd):
    try:
        output = subprocess.check_output(["sh", "-c", cmd], stderr=subprocess.STDOUT, universal_newlines=True)
        print(output)
    except subprocess.CalledProcessError as e:
        print(e.output)

# ON, OFF; LED is on A (496)
commands = ["echo 1 > /sys/class/gpio/gpio496/value", "echo 0 > /sys/class/gpio/gpio496/value"]

while True:
    # executes a random command from the list
    execute(random.choice(commands))
    time.sleep(5)