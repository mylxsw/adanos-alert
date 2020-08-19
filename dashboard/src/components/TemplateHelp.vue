<template>
    <b-card-body>
        <b-card-text class="adanos-help">
            <ul>
                <li>支持的基本字段：
                    <pre><code>
Action:  string
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
	Actions      []Trigger // 分组关联的所有动作

	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
                    </code></pre>
                </li>
                <li>支持的函数：
                    <ul>
                        <li v-for="(helper, i) in helpers" v-bind:key="i"><code>{{ helper.text }}</code> {{ helper.displayText }}</li>
                    </ul>
                </li>
            </ul>
            <hr />
            <ol>
                <li>时间格式 layout 如 <code>2006-01-02T15:04:05Z07:00</code> 代表了 <code>RFC3339</code></li>
                <li>OpenFalconIM 格式为 
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
                <li>jsonutils.KvPair 格式为
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
}

</style>