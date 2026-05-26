#include "Arduino.h"
#include "Ethernet.h"
#include "MQTTClient.h"

byte mac[] = {0xDE, 0xAD, 0xBE, 0xEF, 0xFE, 0xED};

enum state {
  ETH = 1,
  NET = 2,
  HOP = 3,
};

EthernetClient eth;
MQTTClient mqtt;
state s = ETH;

char buf[12];
char *val = &buf[8];

void setup() {
  pinMode(A0, OUTPUT);
  pinMode(A1, OUTPUT);
  pinMode(A2, OUTPUT);

  buf[0] = 0x80 | 1;           // fixmap
  buf[1] = 0xa0 | 5;           // string
  memcpy(&buf[2], "value", 5); // key
  buf[7] = 0xd2;               // i32
}

void setColor(int r, int g, int b) {
  analogWrite(A0, r);
  analogWrite(A1, g);
  analogWrite(A2, b);
}

static inline void write_i32(char* dst, int32_t src) {
  uint32_t s = (uint32_t)src;
  uint8_t* d = (uint8_t*)dst;

  d[0] = (uint8_t)((src >> 24) & 0xFF);
  d[1] = (uint8_t)((src >> 16) & 0xFF);
  d[2] = (uint8_t)((src >>  8) & 0xFF);
  d[3] = (uint8_t)( src        & 0xFF);
}

void loop() {

  switch (s) {
    case ETH:
      setColor(255, 0, 0);

      if (!Ethernet.begin(mac)) {
        delay(100);
        return;
      }

      s = NET;

      break;
    case NET:
      setColor(0, 0, 255);

      mqtt.loop();
      mqtt.begin(MQTT_HOST, MQTT_PORT, eth);

      for (int i = 0; i < 3; ++i) {
        if (mqtt.connect("v2k.hfr-v", "public", "public")) {
          s = HOP;
          return;
        }

        delay(100);
      }

      s = ETH;

      break;
    case HOP:
      setColor(0, 255, 0);

      write_i32(val, analogRead(A3));
      if (!mqtt.publish(MQTT_TOPIC_A3, buf, sizeof(buf), false, 0)) {
        s = NET;
        return;
      }

      write_i32(val, analogRead(A4));
      if (!mqtt.publish(MQTT_TOPIC_A4, buf, sizeof(buf), false, 0)) {
        s = NET;
        return;
      }

      delay(500);

      break;
  }
}
