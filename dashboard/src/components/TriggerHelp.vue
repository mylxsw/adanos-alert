<template>
    <b-card>
        <b-card-text>
            <li>Upper(val string) string</li>
            <li>Lower(val string) string</li>
            <li>Now() time.Time</li>
            <li>ParseTime(layout string, value string) time.Time/li>
            <li>Messages() []repository.Message</li>
            <li>TriggeredTimesInPeriod(periodInMinutes int, triggerStatus string) int64</li>
            <li>LastTriggeredGroup(triggerStatus string) repository.MessageGroup</li>
        </b-card-text>
    </b-card>
</template>

<script>
    export default {
        name: "TriggerHelp"
    }
</script>

<style scoped>

</style>