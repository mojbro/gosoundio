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

struct SoundIoOutStream* create_outstream(struct SoundIoDevice* device) {
    struct SoundIoOutStream* outstream = soundio_outstream_create(device);
    if (!outstream) {
        return NULL;
    }
    outstream->write_callback = writeCallbackProxy;    
    outstream->format = SoundIoFormatFloat32NE;
    return outstream;
}

void run(struct SoundIoOutStream *outstream) {
    int err;
    if ((err = soundio_outstream_start(outstream))) {
        fprintf(stderr, "unable to start device: %s", soundio_strerror(err));
        return;
    }

    for (;;)
        soundio_wait_events(outstream->device->soundio);
    
    printf("Done waiting for events...\n");
}