<template>
    <b-card>
        <b-card-text class="adanos-help">
            <ul>
                <li>支持的基本字段如下：
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
	Actions      []Trigger // 事件组关联的所有动作

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
                <li>支持的对象、函数：
                    <ul>
                        <li v-for="(helper, i) in helpers" v-bind:key="i"><code>{{ helper.text }}</code> {{ helper.displayText }}</li>
                    </ul>
                </li>
            </ul>
            <hr />
            <ol>
                <li><b>triggerStatus</b> 状态可选值为 <code>collecting</code>, <code>pending</code>, <code>ok</code>, <code>failed</code>, <code>canceled</code></li>
                <li>时间格式 layout 如 <code>2006-01-02T15:04:05Z07:00</code> 代表了 <code>RFC3339</code></li>
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
}

</style>
