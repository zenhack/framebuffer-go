#include "fb.h"

#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>

int initfb(char *filename, fb_info_t *fbinfo) {
	fbinfo->fd = open(filename, O_RDWR);
	if(fbinfo->fd < 0) {
		return -1;
	}

	if(ioctl(fbinfo->fd, FBIOGET_FSCREENINFO, &fbinfo->fix_info) < 0) {
		return -1;
	}

	if(ioctl(fbinfo->fd, FBIOGET_VSCREENINFO, &fbinfo->var_info) < 0) {
		return -1;
	}
}
