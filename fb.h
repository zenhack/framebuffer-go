#ifndef GOLANG_FB_H
#define GOLANG_FB_H

#include <linux/fb.h>

typedef struct {
	int fd;
	struct fb_var_screeninfo var_info;
	struct fb_fix_screeninfo fix_info;
} fb_info_t;

int initfb(char *filename, fb_info_t *fbinfo);

#endif
