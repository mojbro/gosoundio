#include <soundio/soundio.h>
#include <stdio.h>

int gosoundio_connect(void* s);
struct SoundIoOutStream* create_outstream(struct SoundIoDevice* device);
void run(struct SoundIoOutStream *outstream);
void writeCallbackProxy(struct SoundIoOutStream *outstream, int frame_count_min, int frame_count_max);
