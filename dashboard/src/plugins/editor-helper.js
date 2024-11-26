import {CodeMirror} from "vue-codemirror-lite";

let helpers = {
    helpers: [
        {text: "Upper", args: ['KEY'], displayText: "Upper(val string) string | Convert string to uppercase"},
        {text: "Lower", args: ['KEY'], displayText: "Lower(val string) string | Convert string to lowercase"},
        {text: "Now", args: [], displayText: "Now() time.Time | Current time"},
        {text: "ParseTime", args: ['LAYOUT', 'VALUE'], displayText: "ParseTime(layout string, value string) time.Time | time string to time object"},
        {text: "DailyTimeBetween", args: ['START_TIME_STR', 'END_TIME_STR'], displayText: "DailyTimeBetween(startTime, endTime string) bool | Determine whether the current time is between startTime and endTime (every day), the time format is 15: 04"},
        {text: 'SQLFinger', args: ['SQL_STR'], displayText: "SQLFinger(sqlStr string) string | Create SQL fingerprint"},
        {text: 'TrimSuffix', args: ['STR', 'SUFFIX'], displayText: 'TrimSuffix(str, suffix string) string | Remove string suffix'},
        {text: 'TrimPrefix', args: ['STR', 'PREFIX'], displayText: 'TrimPrefix(str, prefix string) string | Remove string prefix'},
        {text: 'CutoffLine', args: ['STR', 'MAX'], displayText: 'CutoffLine(val string, maxLine int) string | String intercept maxLine line'},
        {text: 'MD5', args: ['DATA'], displayText: 'MD5(data interface{}) string | Generate md5 value of data'},
        {text: 'Sha1', args: ['DATA'], displayText: 'Sha1(data interface{}) string | Generate the sha1 value of data'},
        {text: 'Base64', args: ['DATA'], displayText: 'Base64(data interface{}) string | Generate base64 encoded value of data'},
        {text: 'CutOff', args: ['MAXLEN', 'STR'], displayText: 'CutOff(maxLen int, val string) string | Maximum length of string interception'},
        {text: 'Mask', args: ['LEFT', 'STR'], displayText: 'Mask(left int, content string) string | String mask, only keep the left characters on both sides, and replace the rest with *' },
        {text: 'Split', args: ['SEP', 'STR'], displayText: 'Split(sep string, content string) []string | The string is split into a string array using sep'},
        {text: 'FilterEmptyLines', args: ['STR'], displayText: 'FilterEmptyLines(content string) string | Remove empty lines from string'},
        {text: 'Join', args: ['ELEMENTS', 'SEP'], displayText: 'Join(elements interface{}, sep string) string | Arrays are connected to strings using sep'},
        {text: 'Repeat', args: ['COUNT', 'STR'], displayText: 'Repeat(count int, s string) string | Repeat string COUNT times'},
        {text: 'NumberBeauty', args: ['NUMBER'], displayText: 'NumberBeauty(number interface{}) string | Numeric formatting'},
        {text: 'Float', args: ['NUMSTR'], displayText: 'Float(numStr string) float64 | String to float64'},
        {text: 'Int', args: ['NUMSTR'], displayText: 'Int(numStr string) int | string to int'},
        {text: 'Empty', args: ['NUMSTR'], displayText: 'Empty(str string) bool | Check whether the string is empty, `blank`, `0`, `false` in any case will be considered is false, otherwise it is true'},
        {text: 'JQuery', args: ['DATA', 'EXPR', 'SUPPRESS_ERR'], displayText: 'JQuery(data string, expression string, suppressError bool) string | Use JQ expression to extract the value in the json string '},
        {text: 'DOMQueryOne', args: ['SELECTOR', 'INDEX HTML'], displayText: 'DOMQueryOne(selector string, index int, htmlContent string) string | Query the content of the index-th element matching selector from the HTML DOM object '},
        {text: 'DOMQuery', args: ['SELECTOR', 'HTML'], displayText: 'DOMQuery(selector string, htmlContent string) []string | Query all elements matching selector from the HTML DOM object'},
        {text: 'JSONArray', args: ['STR', 'PATH'], displayText: 'JSONArray(content string, path string) []gjson.Result | return array elements from path'},
        {text: 'JSONStrArray', args: ['STR', 'PATH'], displayText: 'JSONStrArray(content string, path string) []string | return string array from json'},
        {text: 'JSONIntArray', args: ['STR', 'PATH'], displayText: 'JSONIntArray(content string, path string) []int64 | return int64 array from json'},
        {text: 'JSONFloatArray', args: ['STR', 'PATH'], displayText: 'JSONFloatArray(content string, path string) []float64 | return float64 array from json'},
        {text: 'JSONBoolArray', args: ['STR', 'PATH'], displayText: 'JSONBoolArray(content string, path string) []bool | return bool array from json'},
        {text: 'JSON', args: ['STR', 'PATH'], displayText: 'JSON(content string, path string) string | return string content from json'},
        {text: 'JSONInt', args: ['STR', 'PATH'], displayText: 'JSONInt(content string, path string) int64 | return int content from json'},
        {text: 'JSONFloat', args: ['STR', 'PATH'], displayText: 'JSONFloat(content string, path string) float64 | return float64 content from json'},
        {text: 'JSONBool', args: ['STR', 'PATH'], displayText: 'JSONBool(content string, path string) bool | return bool content from json'},
        {text: 'String', args: ['DATA'], displayText: 'String(data interface}) string | convert any data to string'},
        {text: 'JSONEncode', args: ['DATA'], displayText: 'JSONEncode(data interface}) string | convert any data to json string'},
        {text: 'HumanDuration', args: ['STR'], displayText: 'HumanDuration(duration string) string | Time period formatting, displayed in human-readable form'}
    ],
    groupMatchRules: [
        {text: 'Content', displayText: 'Content | Type: string | event content, string format'},
        {text: 'Meta[""]', displayText: 'Meta | Type: map[string]interface{} | Field, dictionary type'},
        {text: 'Tags[0]', displayText: 'Tags | Type: []string | Field, array type'},
        {text: 'Origin', displayText: 'Origin | Type: string | event source, string'},
        {text: "JsonGet(KEY, DEFAULT)", displayText: "JsonGet(key string, defaultValue string) string | Parse the event body as json and get the specified key"},
        {text: "IsRecovery()", displayText: "IsRecovery() bool | Determine whether the current event is a recovery event"},
        {text: "IsRecoverable()", displayText: "IsRecoverable() bool | Determine whether the current event is recoverable"},
        {text: "IsPlain()", displayText: "IsPlain() bool | Determine whether the current event is a normal event"},
        {text: "FullJSON()", displayText: "FullJSON() string | Encode the entire event into a unified JSON object and return the string representation"},
    ],
    triggerMatchRules: [
        {text: "Events()", displayText: "Events() []repository.Message | Get all Events in the event group"},
        {text: "EventsMatchRegexCount(REGEX)", displayText: "MessagesMatchRegexCount(regex string) int64 | Get the number of Events that match the specified regular expression"},
        {text: "EventsWithMetaCount(KEY, VALUE)", displayText: "EventsWithMetaCount(key, value string) int64 | Get the number of Events whose meta matches the specified key=value"},
        {text: "EventsWithTagsCount(TAG)", displayText: "EventsWithTagsCount(tags string) int64 | Get the number of Events with the specified tag. Multiple tags are separated by commas"},
        {text: "EventsCount()", displayText: "EventsCount() int64 | Get the number of Events in the event group"},
        {text: "HasEventsMatchRegex(REGEX)", displayText: "HasEventsMatchRegex(regex string) bool | Determine whether there is an event matching the regular expression"},
        {text: "HasEventsMatchRegexs([]REGEX)", displayText: "HasEventsMatchRegexs(regex []string) bool | Determine whether there are events matching regular expressions (multiple or relationships)"},
        {text: "UsersHasProperty(KEY, VALUE, RETURN)", displayText: "UsersHasProperty(key, value string, returnField string) []string | Query users based on fields and values, and return a list of fields specified by the value"},
        {text: "UsersHasPropertyRegex(KEY, VALUE_REGEX, RETURN)", displayText: "UsersHasPropertyRegex(key, valueRegex string, returnField string) []string | Query users based on fields and values. The return value specifies a field list. The field value is a regular expression. match"},
        {text: "UsersIDWithProperty(KEY, VALUE, RETURN)", displayText: "UsersIDWithProperty(key, value, returnField string) []repository.UserIDWithMeta | Query users based on fields and values, return the specified field and id"},
        {text: "UsersIDWithPropertyRegex(KEY, VALUE_REGEX, RETURN)", displayText: "UsersIDWithPropertyRegex(key, valueRegex, returnField string) []repository.UserIDWithMeta | Query users based on fields and values. The return value specifies the field and id. The field value is regular Expression matches "},
        {text: "TriggeredTimesInPeriod(PERIOD_IN_MINUTES, TRIGGER_STATUS)", displayText: "TriggeredTimesInPeriod(periodInMinutes int, triggerStatus string) int64 The number of times the current rule is triggered with status triggerStatus within the specified time range"},
        {text: "LastTriggeredGroup(TRIGGER_STATUS)", displayText: "LastTriggeredGroup(triggerStatus string) repository.MessageGroup The event group that last triggered this rule with a status of triggerStatus"},
        {text: "collecting", displayText: "collecting | TriggerStatus: collecting"},
        {text: "pending", displayText: "pending | TriggerStatus：pending"},
        {text: "ok", displayText: "ok | TriggerStatus：ok"},
        {text: "failed", displayText: "failed | TriggerStatus：failed"},
        {text: "canceled", displayText: "canceled | TriggerStatus：canceled"},
        {text: "2006-01-02T15:04:05Z07:00", displayText: "2006-01-02T15:04:05Z07:00 | Time format"},
        {text: "RFC3339", displayText: "RFC3339 | time format"},

        {text: 'Group', displayText: 'Group | Type: MessageGroup | Belong To: ROOT' },
        {text: 'Trigger', displayText: 'Trigger | Type: Trigger | Belong To: ROOT' },
        {text: 'Group.ID', displayText: 'Group.ID | Type: ObjectID | Belong To: MessageGroup' },
        {text: 'Group.SeqNum', displayText: 'Group.SeqNum | Type: int64 | Belong To: MessageGroup' },
        {text: 'Group.AggregateKey', displayText: 'Group.AggregateKey | Type: string | Belong To: MessageGroup' },
        {text: 'Group.MessageCount', displayText: 'Group.MessageCount | Type: int64 | Belong To: MessageGroup' },
        {text: 'Group.Rule', displayText: 'Group.Rule | Type: MessageGroupRule | Belong To: MessageGroup' },
        {text: 'Group.Actions', displayText: 'Group.Actions | Type: []Trigger | Belong To: MessageGroup' },
        {text: 'Group.Status', displayText: 'Group.Status | Type: string | Belong To: MessageGroup' },
        {text: 'Group.CreatedAt', displayText: 'Group.CreatedAt | Type: Time | Belong To: MessageGroup' },
        {text: 'Group.UpdatedAt', displayText: 'Group.UpdatedAt | Type: Time | Belong To: MessageGroup' },
        {text: 'Group.Rule.ID', displayText: 'Group.Rule.ID | Type: ObjectID | Belong To: MessageGroupRule' },
        {text: 'Group.Rule.Name', displayText: 'Group.Rule.Name | Type: string | Belong To: MessageGroupRule' },
        {text: 'Group.Rule.AggregateKey', displayText: 'Group.Rule.AggregateKey | Type: string | Belong To: MessageGroupRule' },
        {text: 'Group.Rule.ExpectReadyAt', displayText: 'Group.Rule.ExpectReadyAt | Type: Time | Belong To: MessageGroupRule' },
        {text: 'Group.Rule.Rule', displayText: 'Group.Rule.Rule | Type: string | Belong To: MessageGroupRule' },
        {text: 'Group.Rule.Template', displayText: 'Group.Rule.Template | Type: string | Belong To: MessageGroupRule' },
        {text: 'Group.Rule.SummaryTemplate', displayText: 'Group.Rule.SummaryTemplate | Type: string | Belong To: MessageGroupRule' },
        {text: 'Trigger.ID', displayText: 'Trigger.ID | Type: ObjectID | Belong To: Trigger' },
        {text: 'Trigger.Name', displayText: 'Trigger.Name | Type: string | Belong To: Trigger' },
        {text: 'Trigger.PreCondition', displayText: 'Trigger.PreCondition | Type: string | Belong To: Trigger' },
        {text: 'Trigger.Action', displayText: 'Trigger.Action | Type: string | Belong To: Trigger' },
        {text: 'Trigger.Meta', displayText: 'Trigger.Meta | Type: string | Belong To: Trigger' },
        {text: 'Trigger.UserRefs', displayText: 'Trigger.UserRefs | Type: []ObjectID | Belong To: Trigger' },
        {text: 'Trigger.Status', displayText: 'Trigger.Status | Type: string | Belong To: Trigger' },
        {text: 'Trigger.FailedCount', displayText: 'Trigger.FailedCount | Type: int | Belong To: Trigger' },
        {text: 'Trigger.FailedReason', displayText: 'Trigger.FailedReason | Type: string | Belong To: Trigger' },
    ],
    matchRules: [
        {text: "matches", displayText: "\"foo\" matches \"^b.+\" Regular matching"},
        {text: "contains", displayText: "str contains \"xxx\" | string contains"},
        {text: "startsWith", displayText: "str startsWith prefix | prefix matching"},
        {text: "endsWith", displayText: "str endsWith prefix | suffix matching"},
        {text: "in", displayText: "user.Group in [\"human_resources\", \"marketing\"] | contains "},
        {text: "not in", displayText: "user.Group not in [\"human_resources\", \"marketing\"] | not containing"},
        {text: "or", displayText: "or | or"},
        {text: "and", displayText: "and | at the same time"},
        {text: "len", displayText: "len | length of array, map or string"},
        {text: "all", displayText: "all | will return true if all element satisfies the predicate"},
        {text: "none", displayText: "none | will return true if all element does NOT satisfies the predicate"},
        {text: "any", displayText: "any | will return true if any element satisfies the predicate"},
        {text: "any([], {Content contains #})", displayText: 'any(["aa", "bb", "cc"], {Content contains #}) | Determine whether Content contains elements in the array any value'},
        {text: "one", displayText: "one | will return true if exactly ONE element satisfies the predicate"},
        {text: "filter", displayText: "filter | filter array by the predicate"},
        {text: "map", displayText: "map | map all items with the closure"},
        {text: "count", displayText: "count | returns number of elements what satisfies the predicate"},
    ],
    templates: [
        {text: '.Events EVENT_COUNT', displayText: 'Events(limit int64) []repository.Event | Get EVENT_COUNT Events from the event group'},
        {text: ".IsRecovery", displayText: "IsRecovery() bool | Determine whether the event in the current event group is a recovery event"},
        {text: ".IsRecoverable", displayText: "IsRecoverable() bool | Determine whether the event in the current event group is recoverable"},
        {text: ".IsPlain", displayText: "IsPlain() bool | Determine whether the event in the current event group is a plain event"},
        {text: ".EventType", displayText: "EventType() string | Determine the event type in the current event group: recovery/plain/recoverable"},
        {text: ".FirstEvent", displayText: "FirstEvent() repository.Event | Get the first event from the current event group"},
        {text: '{{ }}', displayText: '{{ }} | Golang code block'},
        {text: '{{ range $i, $msg := ARRAY }}\n {{ $i }} {{ $msg }} \n{{ end }}', displayText: '{{ range }} | Golang traverse object'},
        {text: '{{ range $i, $msg := .Messages 4 }} {{ end }}', displayText: '{{ range $i, $msg := .Messages 4 }} {{ end }} | Golang traverses Messages, only taking 4 as summary'},
        {text: '{{ if pipeline }}\n T1 \n{{ else if pipeline }}\n T2 \n{{ else }}\n T3 \n{{ end }}', displayText: '{{ if }} | Golang branch condition'}, {text: '[]()', displayText: 'Markdown connection address'}, {text: 'index MAP_VAR "KEY"', displayText: 'index $msg.Meta "message.message" | Get the value of a Key from the Map'}, {text: 'call FUNC ARGS...', displayText: 'call | Returns the result of calling the first argument, which must be a function, with the remaining arguments as parameters'}, {text: 'html', displayText: 'html | Returns the escaped HTML equivalent of the textual representation of its arguments'}, {text: 'js', displayText: 'js | Returns the escaped JavaScript equivalent of the textual representation of its arguments'}, {text: 'len', displayText: 'len | Returns the integer length of its argument'}, {text: 'urlquery', displayText: 'urlquery | Returns the escaped value of the textual representation of its arguments in a form suitable for embedding in a URL query'}, {text: 'print', displayText: 'print | An alias for fmt.Sprint'}, {text: 'printf', displayText: 'printf | An alias for fmt.Sprintf'}, {text: 'println', displayText: 'println | An alias for fmt.Sprintln'}, {text: 'cutoff MAX_LENGTH STR', displayText: 'cutoff(maxLen int, val string) string | String truncation'}, {text: 'cutoff_line MAX_LINES STR', displayText: 'cutoff_line(maxLines int, val string) string | String truncation intercepts the specified number of lines'}, {text: 'line_filter_include FILTER STR', displayText: 'line_filter_include(filter string, val string) string | Filter the string by line, keep the matching lines'},
        {text: 'line_filter_exclude FILTER STR', displayText: 'line_filter_exclude(filter string, val string) string | Filter the string by line, remove the matching lines'},
        {text: 'line_filter_includes STR FILTERS', displayText: 'line_filter_includes(val string, filters ...string) string | Filter the string by line, keep the matching lines, support multiple filters'},
        {text: 'line_filter_excludes STR FILTERS', displayText: 'line_filter_excludes(val string, filters ...string) string | Filter the string by line, remove the matching lines, support multiple filters'},
        {text: 'implode ELEMENT_ARR ","', displayText: 'implode(elems []string, sep string) string | string array concatenation'},
        {text: 'join "," ELEMENT_ARR', displayText: 'join(sep string, elems []string) string | string array concatenation'},
        {text: 'explode STR ","', displayText: 'explode(s, sep string) []string | string separated into arrays'},
        {text: 'split "," STR', displayText: 'split(sep, s string) []string | string separated into arrays'},
        {text: 'ident "IDENT_STR" STR', displayText: 'ident(ident string, message string) string | multi-line string uniform indentation'},
        {text: 'append_by_lines "APPEND_STR" STR', displayText: 'append_by_lines(appendStr string, message string) string | multi-line string appends string after each line'},
        {text: 'json JSONSTR', displayText: 'json(content string) string | JSON string format'},
        {text: 'json_pretty JSONSTR', displayText: 'json_pretty(content string) string | JSON string format'},
        {text: 'datetime LAYOUT DATETIME', displayText: 'datetime(layout string, datetime time.Time) string | Time format is displayed in 2006-01-02 15:04:05 format, time zone selected Beijing/Chongqing'},
        {text: 'datetime_noloc LAYOUT DATETIME',displayText: 'datetime_noloc(layout string, datetime time.Time) string | Time format is displayed in 2006-01-02 15:04:05 format, default time zone'},
        {text: 'datetime_loc LAYOUT LOC DATETIME',displayText: 'datetime_loc(layout string, locName string, datetime time.Time) string | Format the time to display in 2006-01-02 15:04:05 format, specify the time zone, such as UTC'},
        {text: 'datetime_add_sec DATETIME OFFSET',displayText: 'datetime_add_sec(datetime time.Time, offset int) time.Time | Calculate the time, the offset unit is seconds'},
        {text: 'reformat_datetime_str ORIGINAL_LAYOUT TARGET_LAYOUT DATETIME_STR',displayText: 'reformat_datetime_str(originalLayout, targetLayout string, dt string) string | Reformat the time string'},
        {text: 'parse_datetime_str LAYOUT DATETIME_STR',displayText: 'parse_datetime_str(layout string, dt string) time.Time | Parse the time string into a time object'},
        {text: 'parse_datetime_str_rfc3339 DATETIME_STR',displayText: 'parse_datetime_str_rfc3339(dt string) time.Time | Parse the time string into a time object in RFC3339 format'},
        {text: 'json_get "KEY" "DEFAULT" JSONSTR', displayText: 'json_get(key string, defaultValue string, body string) string | Parse body as json, then get the value of key, return defaultValue if failed'},
        {text: 'json_gets "KEY" "DEFAULT" JSONSTR', displayText: 'json_gets(key string, defaultValue string, body string) string | Parse body as json, then get the value of key (you can use commas to separate multiple keys as an alternative), return defaultValue if failed'},
        {text: 'json_array "KEY" JSONSTR', displayText: 'json_array(key string, body string) []string | Parse body as json, then get the value of key (array value)'},
        {text: 'json_flatten JSONSTR MAX_LEVEL', displayText: 'json_flatten(body string, maxLevel int) []jsonutils.KvPair | Parse body as json, then convert to a flat key-value pair and return'},
        {text: 'json_flatten_str JSONSTR MAX_LEVEL', displayText: 'json_flatten_str(body string, maxLevel int) string | Parse body as json, then convert it into a flat string and return it'},
        {text: 'json_flatten_str_noempty JSONSTR MAX_LEVEL', displayText: 'json_flatten_str_noempty(body string, maxLevel int) string | Parse body as json, then convert it into a flat string and return it, while removing all empty values'},
        {text: 'json_fields_filter JSONSTR KEYS', displayText: 'json_filter(body string, keys ...string) string | Parse body as json, then convert it into a flat string, only keep the specified field list, and empty values ​​will be removed'},
        {text: 'format "FORMAT" VAL', displayText: 'format(format string, a ...interface{}) string | Format display, call fmt.Sprintf'},
        {text: 'number_beauty VAL', displayText: 'number_beauty(num interface{}) string | Format and display numbers, automatically add thousand separators'},
        {text: 'integer STR', displayText: 'integer(str string) int | Convert string to integer'},
        {text: 'float STR', displayText: 'float(str string) float64 | Convert string to floating point'},
        {text: 'mysql_slowlog STR', displayText: 'mysql_slowlog(slowlog string) map[string]string | Parse MySQL slow query log as map'},
        {text: 'sql_finger STR', displayText: 'sql_finger(str string) string | Convert SQL to its fingerprint (replace parameters with placeholders)'},
        {text: 'open_falcon_im STR', displayText: 'open_falcon_im(msg string) OpenFalconIM | Parse OpenFalcon event format'},
        {text: 'string_mask STR LEFT', displayText: 'string_mask(content string, left int) string | Only keep left characters on the left and right sides, and replace all characters in the middle with *'},
        {text: 'string_tags TAG_STR SEPARATOR', displayText: 'string_tags(tags string, sep string) []string | Use sep as a separator to cut the string tags into multiple tags. Empty tags will be excluded'},
        {text: 'remove_empty_line STR', displayText: 'remove_empty_line(content string) string | Remove empty lines from the string'},
        {text: 'serialize VAL', displayText: 'serialize(data interface{}) string | Serialize the object into a string for display'},
        {text: 'json_encode VAL', displayText: 'json_encode(data interface{}) string | Object serialized as a string for display'},
        {text: 'error_notice MSG', displayText: 'error_notice(msg string) string | Display "msg" in red font'},
        {text: 'success_notice MSG', displayText: 'success_notice(msg string) string | Display "msg" in green font'},
        {text: 'error_success_notice IS_SUCCESS MSG', displayText: 'error_success_notice(success bool, msg string) string | Display "msg", if success is true, display green, otherwise display red'},
        {text: 'condition S1 S2 CONDITION', displayText: 'condition(s1, s2 string, condition bool) | Conditional output string, if the condition is met, output s1, otherwise s2'},
        {text: 'recoverable_notice .IsRecovery MSG', displayText: 'recoverable_notice(recovered bool, msg string) string | Display the word "msg". If recovered is true, it will be displayed in green and the word [Recovered] will be automatically added. Otherwise, it will be displayed in red'},

        {text: 'user_metas QUERY_K QUERY_V FIELD', displayText: 'user_metas(queryK, queryV string, field string) []string | Query the user field meta information of queryK=queryV. The query result is a string array'},
        {text: 'events_relation_ids EVENTS', displayText: 'events_relation_ids(events []repository.Event) []primitive.ObjectID | Extract the event association IDs contained in multiple events'},
        {text: 'events_relations RELATION_IDS', displayText: 'events_relations(relationIDs []primitive.ObjectID) []repository.EventRelation | Batch query event associations based on multiple event association IDs'},
        {text: 'event_relation_notes RELATION_ID', displayText: 'event_relation_notes(relationID primitive.ObjectID) []repository.EventRelationNote | Query event-related notes based on event association ID'},

        {text: 'prefix_all_str PREFIX ARR', displayText: 'prefix_all_str(prefix string, arr []string) []string | Add a prefix to each element in the string array'},
        {text: 'suffix_all_str SUFFIX ARR', displayText: 'suffix_all_str(prefix string, arr []string) []string | Add a suffix to each element in the string array'},
        {text: 'json_fields_cutoff LENGTH JSON_STR', displayText: 'json_fields_cutoff(length int, body string) map[string]interface{} | Flatten the JSON string and then cut the specified length for each KV'},
        {text: 'json_fields_cutoff_str LENGTH JSON_STR', displayText: 'json_fields_cutoff_str(length int, body string) string | Flatten the JSON string, then cut the specified length for each KV and output the string'},
        {text: 'map_fields_cutoff LENGTH MAP', displayText: 'map_fields_cutoff(length int, source map[string]interface{}) map[string]interface{} | Cut the specified length for each KV in the Map'},
        {text: 'trim_prefix_map_k PREFIX SOURCE', displayText: 'trim_prefix_map_k(prefix string, source map[string]interface{}) map[string]interface{} | Remove the prefix of all Keys in the Map'},

        {text: 'meta_filter STR FILTER_STR', displayText: 'meta_filter(meta map[string]interface{}, allowKeys ...string) map[string]interface{} | Filter Meta and keep only allowed Key'},
        {text: 'meta_filter_exclude STR FILTER_STR', displayText: 'meta_filter_exclude(meta map[string]interface{}, disableKeys ...string) map[string]interface{} | Filter Meta and remove disallowed Key'},
        {text: 'meta_prefix_filter STR FILTER_PREFIX', displayText: 'meta_prefix_filter(meta map[string]interface{}, allowPrefixes ...string) map[string]interface{} | Filter Meta and keep only Keys containing the specified prefix'},
        {text: 'meta_prefix_filter_exclude STR FILTER_PREFIX', displayText: 'meta_prefix_filter_exclude(meta map[string]interface{}, disablePrefixes ...string) map[string]interface{} | Filter Meta and remove Key that matches prefix'},

        {text: 'starts_with STR "START_STR"', displayText: 'starts_with(haystack string, needles ...string) bool | Check if haystack starts with needles'},
        {text: 'ends_with STR "START_END"', displayText: 'ends_with(haystack string, needles ...string) bool | Check if haystack ends with needles'},
        {text: 'trim STR "CUTSTR"', displayText: 'trim(s string, cutset string) string | Remove the cutset characters on both sides of string s'},
        {text: 'trim_left STR "CUTSTR"', displayText: 'trim_left(s string, cutset string) string | Remove the cutset characters on the left side of string s'},
        {text: 'trim_right STR "CUTSTR"', displayText: 'trim_right(s string, cutset string) string | Remove the cutset characters on the right side of string s'},
        {text: 'trim_space STR', displayText: 'trim_space(s string) string | Remove the spaces on both sides of string s'},

        {text: 'str_upper STR', displayText: 'str_upper(s string) string | Convert string to uppercase'},
        {text: 'str_lower STR', displayText: 'str_lower(s string) string | Convert string to lowercase'},
        {text: 'str_replace STR OLD NEW', displayText: 'str_replace(s string, old string, new string) string | Replace string, replace all old in s with new'},
        {text: 'str_repeat STR COUNT', displayText: 'str_repeat(s string, count int) string | Repeat string s count times'},
        {text: 'str_concat STR1 STR2', displayText: 'str_concat(s ...string) string | Concatenate multiple strings'},

        {text: 'html2md HTML', displayText: 'html2md(html string) string | Convert HTML to Markdown'},
        {text: 'md2html MARKDOWN', displayText: 'md2html(markdown string) string | Convert Markdown to HTML'},
        {text: 'md2confluence MARKDOWN', displayText: 'md2confluence(markdown string) string | Convert Markdown to Confluence (Wiki/Jira, etc.) rich text format'},
        {text: 'html_beauty HTML', displayText: 'html_beauty(html string) string | HTML formatting'},
        {text: 'dom_filter_html selector STR', displayText: 'dom_filter_html(selector string, str string) []string | Extract the content matching the selector selector from the HTML DOM and return it as a string array'},
        {text: 'dom_filter_html_n selector N STR', displayText: 'dom_filter_html_n(selector string, n int, str string) string | extract the content matching the selector selector from the HTML DOM and return the nth one (n starts from 0)'},

        {text: 'md5 DATA', displayText: 'md5(data interface{}) string | Generates the md5 value of data'},
        {text: 'sha1 DATA', displayText: 'sha1(data interface{}) string | Generates the sha1 value of data'},
        {text: 'base64 DATA', displayText: 'base64(data interface{}) string | Generates the base64-encoded value of data'},
        {text: 'base64_encode DATA', displayText: 'base64_encode(data interface{}) string | Generates the base64-encoded value of data'},

        {text: 'build_slack_body', displayText: 'build_slack_body(channelName string, username string, emoji string, text string) string | Generate the body of the Slack message'},

        {text: '.Action', displayText: '.Action | Type: string | Belong To: ROOT' },
        {text: '.RuleTemplateParsed', displayText: '.RuleTemplateParsed | Type: string | Belong To: ROOT' },
        {text: '.PreviewURL', displayText: '.PreviewURL | Type: string | Belong To: ROOT' },
        {text: '.ReportURL', displayText: '.ReportURL | Type: string | Belong To: ROOT' },
        {text: '.Trigger', displayText: '.Trigger | Type: Trigger | Belong To: ROOT' },
        {text: '.Group', displayText: '.Group | Type: MessageGroup | Belong To: ROOT' },
        {text: '.Rule', displayText: '.Rule | Type: Rule | Belong To: ROOT' },
        {text: '.Rule.ID', displayText: '.Rule.ID | Type: ObjectID | Belong To: Rule' },
        {text: '.Rule.Name', displayText: '.Rule.Name | Type: string | Belong To: Rule' },
        {text: '.Rule.Description', displayText: '.Rule.Description | Type: string | Belong To: Rule' },
        {text: '.Rule.Tags', displayText: '.Rule.Tags | Type: []string | Belong To: Rule' },
        {text: '.Rule.AggregateRule', displayText: '.Rule.AggregateRule | Type: string | Belong To: Rule' },
        {text: '.Rule.ReadyType', displayText: '.Rule.ReadyType | Type: string | Belong To: Rule' },
        {text: '.Rule.Interval', displayText: '.Rule.Interval | Type: int64 | Belong To: Rule' },
        {text: '.Rule.DailyTimes', displayText: '.Rule.DailyTimes | Type: []string | Belong To: Rule' },
        {text: '.Rule.Rule', displayText: '.Rule.Rule | Type: string | Belong To: Rule' },
        {text: '.Rule.Template', displayText: '.Rule.Template | Type: string | Belong To: Rule' },
        {text: '.Rule.SummaryTemplate', displayText: '.Rule.SummaryTemplate | Type: string | Belong To: Rule' },
        {text: '.Rule.CreatedAt', displayText: '.Rule.CreatedAt | Type: Time | Belong To: Rule' },
        {text: '.Rule.UpdatedAt', displayText: '.Rule.UpdatedAt | Type: Time | Belong To: Rule' },
        {text: '.Rule.Triggers', displayText: '.Rule.Triggers | Type: []Trigger | Belong To: Rule' },
        {text: '.Rule.Status', displayText: '.Rule.Status | Type: string | Belong To: Rule' },
        {text: '.Rule.TimeRanges', displayText: '.Rule.TimeRanges | Type: []TimeRange | Belong To: Rule' },
        {text: '$timeRange.StartTime', displayText: '$timeRange.StartTime | Type: string | Belong To: TimeRange' },
        {text: '$timeRange.EndTime', displayText: '$timeRange.EndTime | Type: string | Belong To: TimeRange' },
        {text: '$timeRange.Interval', displayText: '$timeRange.Interval | Type: int64 | Belong To: TimeRange' },
        {text: '.Trigger.ID', displayText: '.Trigger.ID | Type: ObjectID | Belong To: Trigger' },
        {text: '.Trigger.Name', displayText: '.Trigger.Name | Type: string | Belong To: Trigger' },
        {text: '.Trigger.PreCondition', displayText: '.Trigger.PreCondition | Type: string | Belong To: Trigger' },
        {text: '.Trigger.Action', displayText: '.Trigger.Action | Type: string | Belong To: Trigger' },
        {text: '.Trigger.Meta', displayText: '.Trigger.Meta | Type: string | Belong To: Trigger' },
        {text: '.Trigger.UserRefs', displayText: '.Trigger.UserRefs | Type: []ObjectID | Belong To: Trigger' },
        {text: '.Trigger.Status', displayText: '.Trigger.Status | Type: string | Belong To: Trigger' },
        {text: '.Trigger.FailedCount', displayText: '.Trigger.FailedCount | Type: int | Belong To: Trigger' },
        {text: '.Trigger.FailedReason', displayText: '.Trigger.FailedReason | Type: string | Belong To: Trigger' },
        {text: '$trigger.ID', displayText: '$trigger.ID | Type: ObjectID | Belong To: Trigger' },
        {text: '$trigger.Name', displayText: '$trigger.Name | Type: string | Belong To: Trigger' },
        {text: '$trigger.PreCondition', displayText: '$trigger.PreCondition | Type: string | Belong To: Trigger' },
        {text: '$trigger.Action', displayText: '$trigger.Action | Type: string | Belong To: Trigger' },
        {text: '$trigger.Meta', displayText: '$trigger.Meta | Type: string | Belong To: Trigger' },
        {text: '$trigger.UserRefs', displayText: '$trigger.UserRefs | Type: []ObjectID | Belong To: Trigger' },
        {text: '$trigger.Status', displayText: '$trigger.Status | Type: string | Belong To: Trigger' },
        {text: '$trigger.FailedCount', displayText: '$trigger.FailedCount | Type: int | Belong To: Trigger' },
        {text: '$trigger.FailedReason', displayText: '$trigger.FailedReason | Type: string | Belong To: Trigger' },
        {text: '$action.ID', displayText: '$action.ID | Type: ObjectID | Belong To: Trigger' },
        {text: '$action.Name', displayText: '$action.Name | Type: string | Belong To: Trigger' },
        {text: '$action.PreCondition', displayText: '$action.PreCondition | Type: string | Belong To: Trigger' },
        {text: '$action.Action', displayText: '$action.Action | Type: string | Belong To: Trigger' },
        {text: '$action.Meta', displayText: '$action.Meta | Type: string | Belong To: Trigger' },
        {text: '$action.UserRefs', displayText: '$action.UserRefs | Type: []ObjectID | Belong To: Trigger' },
        {text: '$action.Status', displayText: '$action.Status | Type: string | Belong To: Trigger' },
        {text: '$action.FailedCount', displayText: '$action.FailedCount | Type: int | Belong To: Trigger' },
        {text: '$action.FailedReason', displayText: '$action.FailedReason | Type: string | Belong To: Trigger' },
        {text: '.Group.ID', displayText: '.Group.ID | Type: ObjectID | Belong To: MessageGroup' },
        {text: '.Group.SeqNum', displayText: '.Group.SeqNum | Type: int64 | Belong To: MessageGroup' },
        {text: '.Group.AggregateKey', displayText: '.Group.AggregateKey | Type: string | Belong To: MessageGroup' },
        {text: '.Group.MessageCount', displayText: '.Group.MessageCount | Type: int64 | Belong To: MessageGroup' },
        {text: '.Group.Rule', displayText: '.Group.Rule | Type: MessageGroupRule | Belong To: MessageGroup' },
        {text: '.Group.Actions', displayText: '.Group.Actions | Type: []Trigger | Belong To: MessageGroup' },
        {text: '.Group.Status', displayText: '.Group.Status | Type: string | Belong To: MessageGroup' },
        {text: '.Group.CreatedAt', displayText: '.Group.CreatedAt | Type: Time | Belong To: MessageGroup' },
        {text: '.Group.UpdatedAt', displayText: '.Group.UpdatedAt | Type: Time | Belong To: MessageGroup' },
        {text: '.Group.Rule.ID', displayText: '.Group.Rule.ID | Type: ObjectID | Belong To: MessageGroupRule' },
        {text: '.Group.Rule.Name', displayText: '.Group.Rule.Name | Type: string | Belong To: MessageGroupRule' },
        {text: '.Group.Rule.AggregateKey', displayText: '.Group.Rule.AggregateKey | Type: string | Belong To: MessageGroupRule' },
        {text: '.Group.Rule.ExpectReadyAt', displayText: '.Group.Rule.ExpectReadyAt | Type: Time | Belong To: MessageGroupRule' },
        {text: '.Group.Rule.Rule', displayText: '.Group.Rule.Rule | Type: string | Belong To: MessageGroupRule' },
        {text: '.Group.Rule.Template', displayText: '.Group.Rule.Template | Type: string | Belong To: MessageGroupRule' },
        {text: '.Group.Rule.SummaryTemplate', displayText: '.Group.Rule.SummaryTemplate | Type: string | Belong To: MessageGroupRule' },
        {text: '$msg.ID', displayText: '$msg.ID | Type: ObjectID | Belong To: Message' },
        {text: '$msg.SeqNum', displayText: '$msg.SeqNum | Type: int64 | Belong To: Message' },
        {text: '$msg.Content', displayText: '$msg.Content | Type: string | Belong To: Message' },
        {text: '$msg.Meta', displayText: '$msg.Meta | Type: map[string]interface{} | Belong To: Message' },
        {text: '$msg.Tags', displayText: '$msg.Tags | Type: []string | Belong To: Message' },
        {text: '$msg.Origin', displayText: '$msg.Origin | Type: string | Belong To: Message' },
        {text: '$msg.GroupID', displayText: '$msg.GroupID | Type: []ObjectID | Belong To: Message' },
        {text: '$msg.Status', displayText: '$msg.Status | Type: string | Belong To: Message' },
        {text: '$msg.CreatedAt', displayText: '$msg.CreatedAt | Type: Time | Belong To: Message' },
    ],
    triggerTemplates: [

    ],
}

helpers.matchRules.push(...helpers.helpers.map(item => {
    return {
        text: item.text + "(" + item.args.join(", ") + ")",
        displayText: item.displayText,
    }
}));
helpers.templates.push(...helpers.helpers.map(item => {
    return {
        text: 'helpers.' + item.text + " " + item.args.join(" "),
        displayText: 'helpers.' + item.displayText,
    };
}))

let hintHandler = function (editor) {
    let sources = [];
    switch (editor.options.hintOptions.adanosType) {
        case 'GroupMatchRule':
            sources.push(...helpers.groupMatchRules);
            sources.push(...helpers.matchRules);
            break;
        case 'TriggerMatchRule':
            sources.push(...helpers.triggerMatchRules);
            sources.push(...helpers.matchRules);
            break;
        case 'Template':
            sources.push(...helpers.templates);
            break;
        case 'DingTemplate':
            sources.push(...helpers.templates);
            sources.push(...helpers.triggerTemplates);
            break;
        case 'AllMatchRule':
            sources.push(...helpers.groupMatchRules);
            sources.push(...helpers.triggerMatchRules);
            sources.push(...helpers.matchRules);
            break;
        default:
    }

    let cur = editor.getCursor();
    let token = editor.getTokenAt(cur), start, end, search;
    if (token.end > cur.ch) {
        token.end = cur.ch;
        token.string = token.string.slice(0, cur.ch - token.start);
    }

    if (token.string.match(/^[.`"\w@][.\w$#]*$/g)) {
        search = token.string;
        start = token.start;
        end = token.end;
    } else {
        start = end = cur.ch;
        search = "";
    }

    search = search.toLowerCase()

    let list = [];
    if (search.trim() === '') {
        list = sources;
    } else {
        if (search.charAt(0) === '"' || search.charAt(0) === '.' || search.charAt(0) === "'") {
            search = search.substring(1);
        }

        for (let s in sources) {
            let str = sources[s];
            if (typeof str !== "string") {
                str = str.text;
            }
            if (str.toLowerCase().indexOf(search) >= 0) {
                list.push(sources[s]);
            }
        }
    }

    return {list: list, from: CodeMirror.Pos(cur.line, start), to: CodeMirror.Pos(cur.line, end)};
};

export {helpers, hintHandler}