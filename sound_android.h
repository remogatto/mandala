// +build android

#include <assert.h>
#include <jni.h>
#include <string.h>

// for __android_log_print(ANDROID_LOG_INFO, "YourApp", "formatted message");
// #include <android/log.h>

// for native audio
#include <SLES/OpenSLES.h>
#include <SLES/OpenSLES_Android.h>

// for native asset manager
#include <sys/types.h>
#include <android/asset_manager.h>
#include <android/asset_manager_jni.h>

typedef struct buffer_queue_ap
{
  SLObjectItf bqPlayerObject;
  SLPlayItf bqPlayerPlay;
  SLAndroidSimpleBufferQueueItf	bqPlayerBufferQueue;
  SLVolumeItf bqPlayerVolume;
} t_buffer_queue_ap;

extern SLresult initOpenSL();
extern void shutdownOpenSL();
extern SLresult createBufferQueueAudioPlayer(t_buffer_queue_ap *ap);
extern SLresult destroyBufferQueueAudioPlayer(t_buffer_queue_ap *ap);
extern SLresult enqueueBuffer(t_buffer_queue_ap *ap, const void *buffer, SLuint32 size);
extern SLresult getMaxVolumeLevel(t_buffer_queue_ap *ap, SLmillibel *maxLevel);
extern SLresult setVolumeLevel(t_buffer_queue_ap *ap, SLmillibel value);



