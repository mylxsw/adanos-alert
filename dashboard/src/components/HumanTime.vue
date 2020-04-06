<template>
    <span>{{ formatted(value) }}</span>
</template>

<script>
    export default {
        name: 'HumanTime',
        props: {
            value: Number
        },
        methods: {
            formatted(time_sec) {
               if (time_sec < 60) {
                    return time_sec + " s"
                }

                if (time_sec < 3600) {
                    return (time_sec / 60).toFixed(1) + " m"
                }

                return (time_sec / 60 / 60).toFixed(1) + " h"
            }
        }
    }
</script>

