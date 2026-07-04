#include <zlog.h>

#include <stdio.h>
#include <stdbool.h>
#include <signal.h>
#include <unistd.h>
#include <errno.h>

bool g_run = true;

void sigint_handler(
    int sig[[maybe_unused]])
{
    g_run = false;
}

int main(
    int    argc[[maybe_unused]],
    char **argv[[maybe_unused]])
{
    int              rc  = 0;
    zlog_category_t *cat = NULL;

    if (SIG_ERR == signal(SIGINT, sigint_handler))
    {
        (void)fprintf(stderr, "signal failed\n");
        rc = errno;
        goto finally;
    }

    if ((rc = zlog_init("/etc/zlog.conf")) != 0)
    {
        (void)fprintf(stderr, "zlog init failed\n");
        goto finally;
    }

    cat = zlog_get_category("bootstrap");
    if (NULL == cat)
    {
        (void)fprintf(stderr, "zlog category init failed\n");
        rc = 1;
        goto finally;
    }

    while (g_run)
    {
        zlog_info(cat, "I'm alive!");
        sleep(3);
    }

finally:
    zlog_fini();

    return rc;
}
