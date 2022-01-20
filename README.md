# dafdardaar 🏢

> You have an office? you want to manage it? you love IoT? we are here for you

## Structure

Each office has rooms and each room has many boards that use their room-id as their id besides their type.
there is an also local server which controls these boards with mqtt and topics like follow:

```
/room-id/ldr

/room-id/door/<open/close>

/room-id/led/<light-level>

/room-id/rfid

/room-id/detect

/room-id/light
```

- The door actuator which is a servo motor can open or close the door.

- The led has light level and can create user requested light level (local server use LDR sensor to set this level with respects to current light level)

- The LDR sensor periodically sends room light level.

- The rfid sends card information in case of usage.

- The is an ultra-sonic sensor in front of each room that detects persons.

- The is a light in front of each door and in case of detection and beging in specific time of day it will be trun on.

As you can see, all of the logic is placed on locl server which is also conntecd over HTTP to a central
server that controls many offices. local server has its database for users and admin but it also exists
on central server too with more information in case of any failure or misuse.
