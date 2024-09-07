#include <stdint.h>
#include <stdio.h>

#include "bufferpool.h"
#include "diskmanager.h"
#include "page.h"
#ifdef DEBUG
#include <assert.h>
#endif

typedef struct TestRecord TestRecord;
struct TestRecord
{
    uint64_t a;
    uint64_t b;
    uint64_t c;
    uint64_t d;
    uint64_t e;
    uint64_t f;
    uint64_t g;
    uint64_t h;
};

int main()
{
#ifdef DEBUG
    assert(sizeof(Page) == PAGE_SIZE);
#endif

    DiskManager *dm = create_disk_manager("db.txt");
    open_dm(dm);

    BufferPool *bp = create_buffer_pool();

    TestRecord one = {1, 1, 1, 1, 1, 1, 1, 1};
    TestRecord two = {2, 2, 2, 2, 2, 2, 2, 2};
    TestRecord three = {3, 3, 3, 3, 3, 3, 3, 3};

    Page *page0 = create_page(0, sizeof(TestRecord));

#ifdef DEBUG
    assert(page0->page_id == 0 && "failed pageid assertion");
    assert(page0->num_tuples == 0 && "failed num_tuples assertion");
    assert(page0->tuple_size == 64 && "failed tuple_size assertion");
#endif

    add_tuple(page0, (uint8_t *)&one);
    add_tuple(page0, (uint8_t *)&two);

    add_tuple(page0, (uint8_t *)&three);

#ifdef DEBUG
    assert(((TestRecord *)get_tuple(page0, 0))->a == 1);
    assert(((TestRecord *)get_tuple(page0, 1))->a == 2);
    assert(((TestRecord *)get_tuple(page0, 2))->a == 3);
#endif
    write_page(dm, page0);

    Page *page0_from_disk = get_page_from_buffer_pool(bp, dm, 0);

    TestRecord *tuple_one = (TestRecord *)get_tuple(page0_from_disk, 0);
    printf("First: %lu %lu %lu %lu %lu %lu %lu %lu\n", tuple_one->a,
           tuple_one->b, tuple_one->c, tuple_one->d, tuple_one->e, tuple_one->f,
           tuple_one->g, tuple_one->h);

    remove_tuple(page0_from_disk, 1);

    TestRecord *tuple_two = (TestRecord *)get_tuple(page0_from_disk, 1);
    if (tuple_two)
    {
        printf("Second: %lu %lu %lu %lu %lu %lu %lu %lu\n", tuple_two->a,
               tuple_two->b, tuple_two->c, tuple_two->d, tuple_two->e,
               tuple_two->f, tuple_two->g, tuple_two->h);
    }
    TestRecord *tuple_three = (TestRecord *)get_tuple(page0_from_disk, 2);
    printf("Third: %lu %lu %lu %lu %lu %lu %lu %lu\n", tuple_three->a,
           tuple_three->b, tuple_three->c, tuple_three->d, tuple_three->e,
           tuple_three->f, tuple_three->g, tuple_three->h);

    destroy_buffer_pool(bp);
    destroy_disk_manager(dm);
}
