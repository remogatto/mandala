// +build android

#include <android/native_activity.h>

const char *getPackageName(ANativeActivity* activity) {
  JNIEnv* env=0;

  (*activity->vm)->AttachCurrentThread(activity->vm, &env, 0);

  jclass clazz = (*env)->GetObjectClass(env, activity->clazz);
  jmethodID methodID = (*env)->GetMethodID(env, clazz, "getPackageCodePath", "()Ljava/lang/String;");
  jobject result = (*env)->CallObjectMethod(env, activity->clazz, methodID);

  const char* str;
  jboolean isCopy;
  str = (*env)->GetStringUTFChars(env, (jstring)result, &isCopy);

  (*activity->vm)->DetachCurrentThread(activity->vm);
  return str;
}

