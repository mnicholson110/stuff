#pragma once
#include <stdint.h>

#include "common.h"
#include "diskmanager.h"
#include "page.h"

#define NUM_BUFFERS 16

// Buffer Pool struct
typedef struct BufferPool BufferPool;
struct BufferPool
{
  Page page_buffers[NUM_BUFFERS];
  uint8_t in_use[NUM_BUFFERS];    // 1 if the buffer is in use, 0 otherwise
  uint32_t page_dir[NUM_BUFFERS]; // Keeps track of which page is stored in
                                  // each buffer
};

BufferPool *create_buffer_pool();

Page *get_page_from_buffer_pool(BufferPool *bp, DiskManager *dm,
                                uint32_t page_id);

void set_page_in_buffer_pool(BufferPool *bp, uint32_t page_id, Page *page);

void remove_page_from_buffer_pool(BufferPool *bp, uint32_t page_id);

void destroy_buffer_pool(BufferPool *bp);
