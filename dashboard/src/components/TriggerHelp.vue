<template>
    <b-card>
        <b-card-text class="adanos-help">
            <ul>
                <li>Base Fields:
                    <pre><code>
Group   {
	ID     primitive.ObjectID
	SeqNum int64
	MessageCount int64
	AggregateKey string
	Rule         {
		ID              primitive.ObjectID
		Name            string
		ExpectReadyAt   time.Time
		Rule            string
		Template        string
		SummaryTemplate string
	}
	Actions      []Trigger // All actions associated with the event group

	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// 当前触发的动作
Trigger {
	ID           primitive.ObjectID
	Name         string
	PreCondition string
	Action       string
	Meta         string
	UserRefs     []primitive.ObjectID
	Status       string
	FailedCount  int
	FailedReason string
}

UserIDWithMeta {
	UserID string
	Meta   []string
}
                    </code></pre>
                </li>
                <li>Objects and functions:
                    <ul>
                        <li v-for="(helper, i) in helpers" v-bind:key="i">
                            <code>{{ helper.text }}</code>
                            <p>{{ helper.displayText }}</p>
                        </li>
                    </ul>
                </li>
            </ul>
            <hr />
            <ol>
                <li><b>triggerStatus</b> status optional values are <code>collecting</code>, <code>pending</code>, <code>ok</code>, <code>failed</code>, <code>canceled</code></li>
                <li>The time format layout such as <code>2006-01-02T15:04:05Z07:00</code> represents <code>RFC3339</code></li>
            </ol>
        </b-card-text>
    </b-card>
</template>

<script>
    export default {
        name: "TriggerHelp",
        props: {
            helpers: Array,
        }
    }
</script>

<style scoped>

.adanos-help {
    font-size: 80%;
    max-height: 400px;
    overflow-y: auto;
    border: 1px solid #ccc;
    border-radius: 4px;
    padding-top: 12px;
}

</style>
