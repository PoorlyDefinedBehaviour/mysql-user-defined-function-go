#include <m_string.h>
#include <my_dbug.h>
#include <mysql/plugin.h>
#include <mysql/plugin_audit.h>
#include <mysql_version.h>
#include <mysqld_error.h>
#include <string.h>

static int temporal_plugin_init(void *p);
static int temporal_plugin_deinit(void *p);
static void temporal_event_handler(MYSQL_THD, mysql_event_class_t,
                                   const void *);

static struct st_mysql_audit temporal_descriptor = {
    MYSQL_AUDIT_INTERFACE_VERSION,
    NULL,
    NULL,
    NULL,
    NULL,
    NULL,
    temporal_event_handler,
    NULL,
    NULL,
    NULL,
    NULL,
    NULL,
    NULL,
    NULL,
    NULL,
    NULL};

mysql_declare_plugin(temporal){MYSQL_AUDIT_PLUGIN,
                               &temporal_descriptor,
                               "temporal",
                               "TODO",
                               "Temporal Tables Plugin",
                               PLUGIN_LICENSE_GPL,
                               temporal_plugin_init,
                               temporal_plugin_deinit,
                               0x0100,
                               NULL,
                               NULL,
                               NULL} mysql_declare_plugin_end;

static int temporal_plugin_init(void *p) { return 0; }

static int temporal_plugin_deinit(void *p) { return 0; }

static void temporal_event_handler(MYSQL_THD thd,
                                   mysql_event_class_t event_class,
                                   const void *event) {
  if (event_class == MYSQL_AUDIT_GENERAL_CLASS) {
    const struct mysql_event_general *event_general =
        (const struct mysql_event_general *)event;

    if (strcmp(event_general->command, "select") == 0) {
      handle_temporal_select(thd, event_general);
    }
  }
}

void handle_temporal_select(MYSQL_THD thd,
                            const struct mysql_event_general *event) {
  const char *query = event->query;
  char rewritten_query[2048];

  // Check if the query requests historical data
  if (strstr(query, "FOR SYSTEM_TIME") != NULL) {
    // Parse the query to extract table name and rewrite it to select from the
    // history table
    char *table_name = extract_table_name(query);
    if (table_name) {
      snprintf(rewritten_query, sizeof(rewritten_query),
               "SELECT * FROM %s_history WHERE %s", table_name,
               query + strlen("SELECT * FROM ") + strlen(table_name));

      // Execute the rewritten query
      if (mysql_query(thd, rewritten_query)) {
        fprintf(stderr, "Failed to execute rewritten query: %s\n",
                mysql_error(thd));
      } else {
        // Handle the result set
        MYSQL_RES *result = mysql_store_result(thd);
        if (result) {
          // Process the result
          mysql_free_result(result);
        }
      }

      free(table_name);
    }
  }
}

char *extract_table_name(const char *query) {
  const char *from_str = strstr(query, " FROM ");
  if (!from_str) return NULL;

  from_str += 6;  // Move past " FROM "
  const char *end = strpbrk(from_str, " \t\n\r");

  size_t table_name_len = end ? (size_t)(end - from_str) : strlen(from_str);
  char *table_name = (char *)malloc(table_name_len + 1);
  if (!table_name) return NULL;

  strncpy(table_name, from_str, table_name_len);
  table_name[table_name_len] = '\0';

  return table_name;
}
