# Rotel Amplifier Control

Code here is for bridging a Rotel Amplifier to MQTT, for integration into Home Assistant.
To make the docker container,

```
git clone git@github.com:djthorpe/go-rotel.git
cd go-rotel
make docker
```

Then to run the docker container, assuming you have a MQTT broker running on some machine, and
you have a Rotel amplifier connected to the serial port `/dev/ttyUSB1`:

```
docker run --rm --name rotel --device=/dev/ttyUSB1 \
  rotel-linux-aarch64:latest \
  rotel -mqtt mqtt.mutablelogic.com -tty /dev/ttyUSB1
```

The docker container will connect to the MQTT broker and publish messages to the topic `rotel/amp` and