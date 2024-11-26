<template>
    <b-card-body>
        <b-card-text class="adanos-help">
            <ul>
                <li>Base Fields:
                    <pre><code>
Action: string
RuleTemplateParsed: string // The alarm summary generated based on the rules of the global template is only available in the action-triggering template section.
PreviewURL: string // Preview address for alarm details
Rule: {
	ID          primitive.ObjectID
	Name        string
	Description string
	Tags        []string

	Interval int64

	Rule            string
	Template        string
	SummaryTemplate string
	Triggers        []Trigger

	Status string

	CreatedAt time.Time
	UpdatedAt time.Time
}
Trigger: {
	ID           primitive.ObjectID
	Name         string
	PreCondition string
	Action       string
	Meta         string
	UserRefs     []primitive.ObjectID
	Status       string
	FailedCount  int
	FailedReason string
},
Group: {
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
                    </code></pre>
                </li>
                <li>支持的对象、函数：
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
                <li>The time format layout such as <code>2006-01-02T15:04:05Z07:00</code> represents <code>RFC3339</code></li>
                <li>OpenFalconIM format is
                    <pre>
<code>type OpenFalconIM struct {
	Priority    int
	Status      string
	Endpoint    string
	Body        string
	CurrentStep int
	FormatTime  string
}
</code>
                    </pre>
                </li>
                <li>jsonutils.KvPair format is
                    <pre>
<code>
type KvPair struct {
	Key   string
	Value string
}
</code>
                    </pre>
                </li>
            </ol>
        </b-card-text>
    </b-card-body>
</template>

<script>
    export default {
        name: "TemplateHelp",
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