// +build android

/*
 * Copyright (C) 2010 The Android Open Source Project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

/* This is a JNI example where we use native methods to play sounds
 * using OpenSL ES. See the corresponding Java source file located at:
 *
 *   src/com/example/nativeaudio/NativeAudio/NativeAudio.java
 */

#include <jni.h>
#include <string.h>

// for __android_log_print(ANDROID_LOG_INFO, "YourApp", "formatted message");
// #include <android/log.h>

// for native activity
#include <android/native_activity.h>

// for native audio
#include <SLES/OpenSLES.h>
#include <SLES/OpenSLES_Android.h>

// for native asset manager
#include <sys/types.h>
#include <android/asset_manager.h>
#include <android/asset_manager_jni.h>

#include "sound_android.h"

// engine interfaces
static SLObjectItf engineObject = NULL;
static SLEngineItf engineEngine;

// output mix interfaces
static SLObjectItf outputMixObject = NULL;

extern void playerCallback();

// create the engine and output mix objects
SLresult initOpenSL()
{
    SLresult result;

    // create engine
    result = slCreateEngine(&engineObject, 0, NULL, 0, NULL, NULL);

    if (result != SL_RESULT_SUCCESS) {
      return result;
    }

    // realize the engine
    result = (*engineObject)->Realize(engineObject, SL_BOOLEAN_FALSE);

    if (result != SL_RESULT_SUCCESS) {
      return result;
    }

    // get the engine interface, which is needed in order to create other objects
    result = (*engineObject)->GetInterface(engineObject, SL_IID_ENGINE, &engineEngine);

    if (result != SL_RESULT_SUCCESS) {
      return result;
    }

    // create output mix
    result = (*engineEngine)->CreateOutputMix(engineEngine, &outputMixObject, 0, NULL, NULL);

    if (result != SL_RESULT_SUCCESS) {
      return result;
    }

    // realize the output mix
    result = (*outputMixObject)->Realize(outputMixObject, SL_BOOLEAN_FALSE);

    if (result != SL_RESULT_SUCCESS) {
      return result;
    }

    return SL_RESULT_SUCCESS;
}

// create buffer queue audio player
SLresult createBufferQueueAudioPlayer(t_buffer_queue_ap *ap)
{
    SLresult result;

    // configure audio source
    SLDataLocator_AndroidSimpleBufferQueue loc_bufq = {SL_DATALOCATOR_ANDROIDSIMPLEBUFFERQUEUE, 2};
    SLDataFormat_PCM format_pcm = {SL_DATAFORMAT_PCM, 1, SL_SAMPLINGRATE_44_1,
        SL_PCMSAMPLEFORMAT_FIXED_16, SL_PCMSAMPLEFORMAT_FIXED_16,
        SL_SPEAKER_FRONT_CENTER, SL_BYTEORDER_LITTLEENDIAN};
    SLDataSource audioSrc = {&loc_bufq, &format_pcm};

    // configure audio sink
    SLDataLocator_OutputMix loc_outmix = {SL_DATALOCATOR_OUTPUTMIX, outputMixObject};
    SLDataSink audioSnk = {&loc_outmix, NULL};

    // create audio player
    const SLInterfaceID ids[] = {SL_IID_BUFFERQUEUE, SL_IID_VOLUME};
    const SLboolean req[] = {SL_BOOLEAN_TRUE, SL_BOOLEAN_TRUE};
    result = (*engineEngine)->CreateAudioPlayer(engineEngine, &ap->bqPlayerObject, &audioSrc, &audioSnk, 2, ids, req);
    
    if (result != SL_RESULT_SUCCESS) {
      return result;
    }

    // realize the player
    result = (*ap->bqPlayerObject)->Realize(ap->bqPlayerObject, SL_BOOLEAN_FALSE);

    if (result != SL_RESULT_SUCCESS) {
      return result;
    }

    // get the play interface
    result = (*ap->bqPlayerObject)->GetInterface(ap->bqPlayerObject, SL_IID_PLAY, &ap->bqPlayerPlay);

    if (result != SL_RESULT_SUCCESS) {
      return result;
    }

    // get the buffer queue interface
    result = (*ap->bqPlayerObject)->GetInterface(ap->bqPlayerObject, SL_IID_BUFFERQUEUE,
            &ap->bqPlayerBufferQueue);

    if (result != SL_RESULT_SUCCESS) {
      return result;
    }

    /* // register callback on the buffer queue */
    /* result = (*ap->bqPlayerBufferQueue)->RegisterCallback(ap->bqPlayerBufferQueue, bqPlayerCallback, NULL); */
    /* assert(SL_RESULT_SUCCESS == result); */
    /* (void)result; */

    // get the volume interface
    result = (*ap->bqPlayerObject)->GetInterface(ap->bqPlayerObject, SL_IID_VOLUME, &ap->bqPlayerVolume);

    if (result != SL_RESULT_SUCCESS) {
      return result;
    }

    // set the player's state to playing
    result = (*ap->bqPlayerPlay)->SetPlayState(ap->bqPlayerPlay, SL_PLAYSTATE_PLAYING);

    return result;
}

SLresult getMaxVolumeLevel(t_buffer_queue_ap *ap, SLmillibel *maxLevel)
{
  SLresult result;
  result = (*ap->bqPlayerVolume)->GetVolumeLevel(ap->bqPlayerVolume, maxLevel);
  return result;
}

SLresult setVolumeLevel(t_buffer_queue_ap *ap, SLmillibel value)
{
  SLresult result;
  result = (*ap->bqPlayerVolume)->SetVolumeLevel(ap->bqPlayerVolume, value);
  return result;
}

SLresult enqueueBuffer(t_buffer_queue_ap *ap, const void *buffer, SLuint32 size) 
{
  SLresult result;

  result = (*ap->bqPlayerBufferQueue)->Clear(ap->bqPlayerBufferQueue);

  if (result != SL_RESULT_SUCCESS) {
    return result;
  }
  
  result = (*ap->bqPlayerBufferQueue)->Enqueue(ap->bqPlayerBufferQueue, (short*)buffer, size);

  return result;
}

// shut down the native audio system
void shutdownOpenSL()
{
    // destroy output mix object, and invalidate all associated interfaces
    if (outputMixObject != NULL) {
        (*outputMixObject)->Destroy(outputMixObject);
        outputMixObject = NULL;
    }

    // destroy engine object, and invalidate all associated interfaces
    if (engineObject != NULL) {
        (*engineObject)->Destroy(engineObject);
        engineObject = NULL;
        engineEngine = NULL;
    }

}
