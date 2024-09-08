#pragma once
#include <stddef.h>
#include <stdint.h>

#include "common.h"

typedef struct Page Page;
struct Page
{
  uint32_t page_id;                      // Page identifier
  uint32_t tuple_size;                   // Size of tuples in this page, in bytes TODO: handle
                                         // dynamic sizes; ptrs to overflow?
  uint32_t num_tuples;                   // Number of tuples currently in this page
  uint32_t free_offset;                  // Offset to free space between offset array and tuples
  uint8_t data[PAGE_SIZE - HEADER_SIZE]; // Tuple data and slot array
};
//
//  Page:
//
//  ...| free_offset |                  data                 |
//  ...|             |                                       |
//  ...|     4B      | offsets ->*free_offset      <- tuples |
//
//
Page *create_page(uint32_t page_id, uint32_t tuple_size);
void destroy_page(Page *page);

void add_tuple(Page *page, const uint8_t *tuple);
uint8_t *get_tuple(Page *page, uint8_t tuple_num);
uint32_t *get_offset_array(Page *page);
void remove_tuple(Page *page, uint8_t tuple_num);
