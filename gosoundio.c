#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdint.h>
#include <math.h>

#include "gosoundio.h"

int gosoundio_connect(void* s) {
    struct SoundIo* snd = (struct SoundIo*)s;
    int res = soundio_connect(snd);
    soundio_flush_events(snd);
    return res;
}

const char* gosoundio_device_name(struct SoundIoDevice* device) {
    return device->name;
}