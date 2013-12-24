// +build android

#include <android/native_activity.h>

int throwException(ANativeActivity *act, const char*err) {
	JNIEnv *env;
	(*(act->vm))->GetEnv(act->vm, (void **)&env, JNI_VERSION_1_6);
	if (env == NULL)
		return 0;
	jclass excClass = (*env)->FindClass(env, "java/lang/Error");
	if (excClass == NULL)
		return 0;
	if ((*env)->ThrowNew(env, excClass, err) < 0)
		return 0;
	return 1;
}

extern void onNativeWindowCreated(ANativeActivity *activity, ANativeWindow* window);
extern void onNativeWindowDestroyed(ANativeActivity *activity, ANativeWindow* window);
extern void onNativeWindowResized(ANativeActivity *activity, ANativeWindow* window);
extern void onInputQueueCreated(ANativeActivity *activity, AInputQueue* queue);
extern void onInputQueueDestroyed(ANativeActivity *activity, AInputQueue* queue);
extern void onCreate(ANativeActivity *activity, void* savedState, size_t savedStateSize);
extern void onDestroy(ANativeActivity *activity);
extern void onResume(ANativeActivity *activity);
extern void onPause(ANativeActivity *activity);
extern void onConfigurationChanged(ANativeActivity *activity);
extern void onWindowFocusChanged(ANativeActivity *activity, int focused);

void ANativeActivity_onCreate(ANativeActivity *activity, void* savedState, size_t savedStateSize) {
	activity->callbacks->onNativeWindowCreated = onNativeWindowCreated;
	activity->callbacks->onNativeWindowDestroyed = onNativeWindowDestroyed;
	activity->callbacks->onInputQueueCreated = onInputQueueCreated;
	activity->callbacks->onInputQueueDestroyed = onInputQueueDestroyed;
	activity->callbacks->onDestroy = onDestroy;
	activity->callbacks->onResume = onResume;
	activity->callbacks->onPause = onPause;
	activity->callbacks->onConfigurationChanged = onConfigurationChanged;
	activity->callbacks->onNativeWindowResized = onNativeWindowResized;
	activity->callbacks->onWindowFocusChanged = onWindowFocusChanged;

	onCreate(activity, savedState, savedStateSize);
}
