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

typedef struct asset_ap
{
  SLObjectItf fdPlayerObject;
  SLPlayItf fdPlayerPlay;
} t_asset_ap;

extern void createEngine(JNIEnv* env, jclass clazz);
extern jboolean createAssetAudioPlayer(ANativeActivity *act, t_asset_ap *ap, char *filename);
extern void setPlayingAssetAudioPlayer(SLPlayItf fdPlayerPlay, jboolean isPlaying);

